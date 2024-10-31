package registration

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fridrock/users/token"
)

type RegistrationHandler interface {
	HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error)
}

type RegistrationHandlerImpl struct {
	storage      RegistrationStorage
	tokenService token.TokenService
	parser       RegistrationParser
}

// TODO make validation
// TODO tests
func (rs *RegistrationHandlerImpl) HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error) {
	userDto, err := rs.parser.GetDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	user, err := rs.storage.SaveUser(userDto)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	responseDto, err := rs.tokenService.GenerateToken(user)
	if err != nil {
		slog.Info("Error processing token %s", err.Error())
		return http.StatusBadRequest, err
	}
	responseText, err := json.MarshalIndent(responseDto, "", " ")
	if err != nil {
		return http.StatusBadRequest, err
	}
	w.Write(responseText)
	return http.StatusCreated, nil
}
func NewRegistrationHandler(storage RegistrationStorage) RegistrationHandler {
	return &RegistrationHandlerImpl{
		storage:      storage,
		parser:       &RegistrationParserImpl{},
		tokenService: token.NewTokenService(),
	}
}
