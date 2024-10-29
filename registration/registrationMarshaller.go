package registration

import (
	"encoding/json"
	"net/http"

	"github.com/fridrock/users/api"
)

type RegistrationParser interface {
	GetDto(r *http.Request) (api.UserDto, error)
}

type RegistrationParserImpl struct {
}

func (rm RegistrationParserImpl) GetDto(r *http.Request) (api.UserDto, error) {
	var userDto api.UserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return userDto, err
	}
	return userDto, nil
}
