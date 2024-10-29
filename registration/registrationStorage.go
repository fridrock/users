package registration

import (
	"github.com/fridrock/users/api"
	"github.com/fridrock/users/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegistrationStorage interface {
	SaveUser(userDto api.UserDto) (uuid.UUID, error)
}

type RegistrationStorageImpl struct {
	db     *sqlx.DB
	hasher utils.PasswordHasher
}

func (rs RegistrationStorageImpl) SaveUser(userDto api.UserDto) (uuid.UUID, error) {
	var id uuid.UUID
	q := `INSERT INTO users(id, username, email, name, surname, hashed_password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	hashedPassword, err := rs.hasher.HashPassword(userDto.Password)
	if err != nil {
		return id, err
	}
	err = rs.db.QueryRow(
		q,
		uuid.New().String(),
		userDto.Username,
		userDto.Email,
		userDto.Name,
		userDto.Surname,
		hashedPassword).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func NewRegistrationStorage(db *sqlx.DB) RegistrationStorage {
	return &RegistrationStorageImpl{
		db:     db,
		hasher: utils.NewPasswordHasher(),
	}
}
