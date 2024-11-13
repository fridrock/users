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

type User struct {
	Id             uuid.UUID `db:"id"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	Name           string    `db:"name"`
	Surname        string    `db:"surname"`
	HashedPassword string    `db:"hashed_password"`
}
