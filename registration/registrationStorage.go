package registration

import (
	"github.com/fridrock/users/api"
	"github.com/fridrock/users/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegistrationStorage interface {
	SaveUser(userDto api.UserDto) (api.User, error)
}

type RegistrationStorageImpl struct {
	db     *sqlx.DB
	hasher utils.PasswordHasher
}

func (rs RegistrationStorageImpl) SaveUser(userDto api.UserDto) (api.User, error) {
	var user api.User
	q := `INSERT INTO users(id, username, email, name, surname, hashed_password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	hashedPassword, err := rs.hasher.HashPassword(userDto.Password)
	if err != nil {
		return user, err
	}
	var id uuid.UUID
	err = rs.db.QueryRow(
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
	user = rs.fillUserFromUserDto(userDto, id, hashedPassword)
	return user, nil
}

func (rs RegistrationStorageImpl) fillUserFromUserDto(userDto api.UserDto, id uuid.UUID, hashedPassword string) api.User {
	return api.User{
		Id:             id,
		Username:       userDto.Username,
		Name:           userDto.Name,
		Surname:        userDto.Surname,
		Email:          userDto.Email,
		HashedPassword: hashedPassword,
	}
}
func NewRegistrationStorage(db *sqlx.DB) RegistrationStorage {
	return &RegistrationStorageImpl{
		db:     db,
		hasher: utils.NewPasswordHasher(),
	}
}
