package api

import "github.com/google/uuid"

type UserDto struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUserDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetUserDto struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Name     string    `json:"name" db:"name"`
	Surname  string    `json:"surname" db:"surname"`
}

type User struct {
	Id             uuid.UUID `db:"id"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	Name           string    `db:"name"`
	Surname        string    `db:"surname"`
	HashedPassword string    `db:"hashed_password"`
}
