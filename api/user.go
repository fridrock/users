package api

import "github.com/google/uuid"

type UserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type User struct {
	Id             uuid.UUID
	Username       string
	Email          string
	Name           string
	Surname        string
	HashedPassword string
}
