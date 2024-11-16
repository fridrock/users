package usr

import (
	"encoding/json"
	"net/http"

	"github.com/fridrock/users/api"
	"github.com/go-playground/validator/v10"
)

type UserParser interface {
	GetAuthUserDto(*http.Request) (api.AuthUserDto, error)
	GetUserDto(*http.Request) (api.UserDto, error)
	GetUsername(*http.Request) (string, error)
}

type UserParserImpl struct {
	validate *validator.Validate
}

func newUserParser() UserParser {
	return UserParserImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (up UserParserImpl) GetAuthUserDto(r *http.Request) (api.AuthUserDto, error) {
	var authDto api.AuthUserDto
	err := json.NewDecoder(r.Body).Decode(&authDto)
	return authDto, err
}

func (up UserParserImpl) GetUserDto(r *http.Request) (api.UserDto, error) {
	var userDto api.UserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return userDto, err
	}
	err = up.validate.Struct(userDto)
	return userDto, err
}

func (up UserParserImpl) GetUsername(r *http.Request) (string, error) {
	var userDto api.GetUserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	return userDto.Username, err
}
