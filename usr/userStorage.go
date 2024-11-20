package usr

import (
	"fmt"

	"github.com/fridrock/users/api"
	"github.com/fridrock/users/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStorage interface {
	CheckUser(userDto api.AuthUserDto) (api.User, error)
	SaveUser(userDto api.UserDto) (api.User, error)
	GetProfiles([]uuid.UUID) ([]api.GetUserDto, error)
	FindUsers(string) ([]api.GetUserDto, error)
	GetUsers() ([]api.GetUserDto, error)
	FilterNotInFriends([]api.GetUserDto, uuid.UUID) ([]api.GetUserDto, error)
}

type UserStorageImpl struct {
	db     *sqlx.DB
	hasher utils.PasswordHasher
}

func (us *UserStorageImpl) GetProfiles(ids []uuid.UUID) ([]api.GetUserDto, error) {
	var users []api.GetUserDto
	q := `SELECT id, username, email, name, surname FROM users WHERE id IN(?)`
	q, args, err := sqlx.In(q, ids)
	if err != nil {
		fmt.Println("error building query in string" + err.Error())
	}
	q = us.db.Rebind(q)
	err = us.db.Select(&users, q, args...)
	return users, err
}

func (us *UserStorageImpl) FindUsers(username string) ([]api.GetUserDto, error) {
	var userDtos []api.GetUserDto
	q := `SELECT id, username, email, name, surname FROM users WHERE username LIKE $1`
	err := us.db.Select(&userDtos, q, "%"+username+"%")
	return userDtos, err
}

func (us *UserStorageImpl) FilterNotInFriends(users []api.GetUserDto, userId uuid.UUID) ([]api.GetUserDto, error) {
	q := `SELECT fr2id FROM friends where fr1id = $1`
	var friendsIds []uuid.UUID
	var usersNotFriends []api.GetUserDto
	err := us.db.Select(&friendsIds, q, userId)
	if err != nil {
		return usersNotFriends, err
	}
	friendsIds = append(friendsIds, userId)
	for _, user := range users {
		flag := true
		for _, id := range friendsIds {
			if id == user.Id {
				flag = false
				break
			}
		}
		if flag {
			usersNotFriends = append(usersNotFriends, user)
		}
	}
	return usersNotFriends, nil
}

func (us *UserStorageImpl) GetUsers() ([]api.GetUserDto, error) {
	var userDtos []api.GetUserDto
	q := `SELECT id, username, email, name, surname FROM users`
	err := us.db.Select(&userDtos, q)
	return userDtos, err
}

func (us *UserStorageImpl) CheckUser(authUserDto api.AuthUserDto) (api.User, error) {
	var user api.User
	q := `SELECT * FROM users WHERE username=$1`

	err := us.db.QueryRowx(q, authUserDto.Username).StructScan(&user)
	if err != nil {
		return user, err
	}
	if !us.hasher.CheckPassword(authUserDto.Password, user.HashedPassword) {
		return user, fmt.Errorf("wrong password for user")
	}
	return user, nil
}

func (us *UserStorageImpl) SaveUser(userDto api.UserDto) (api.User, error) {
	var user api.User
	if us.checkIfUserExist(userDto) {
		return user, fmt.Errorf("such user already exist")
	}
	q := `INSERT INTO users(id, username, email, name, surname, hashed_password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	hashedPassword, err := us.hasher.HashPassword(userDto.Password)
	if err != nil {
		return user, err
	}
	var id uuid.UUID
	err = us.db.QueryRow(
		q,
		uuid.New().String(),
		userDto.Username,
		userDto.Email,
		userDto.Name,
		userDto.Surname,
		hashedPassword).Scan(&id)
	if err != nil {
		return user, err
	}
	user = us.fillUserFromUserDto(userDto, id, hashedPassword)
	return user, nil
}

func (us *UserStorageImpl) fillUserFromUserDto(userDto api.UserDto, id uuid.UUID, hashedPassword string) api.User {
	return api.User{
		Id:             id,
		Username:       userDto.Username,
		Name:           userDto.Name,
		Surname:        userDto.Surname,
		Email:          userDto.Email,
		HashedPassword: hashedPassword,
	}
}

func (us *UserStorageImpl) checkIfUserExist(u api.UserDto) bool {
	var user api.User
	q := `SELECT * FROM users WHERE username=$1 OR email=$2`
	row := us.db.QueryRowx(q, u.Username, u.Email)
	err := row.StructScan(&user)
	return err == nil && user.Username != ""
}

func NewUserStorage(db *sqlx.DB) UserStorage {
	return &UserStorageImpl{
		db:     db,
		hasher: utils.NewPasswordHasher(),
	}
}
