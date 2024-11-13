package usr

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fridrock/users/token"
)

type UserHandler interface {
	HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error)
	HandleAuth(w http.ResponseWriter, r *http.Request) (int, error)
}

type UserHandlerImpl struct {
	storage      UserStorage
	tokenService token.TokenService
	parser       UserParser
}

// TODO tests
func (uh *UserHandlerImpl) HandleAuth(w http.ResponseWriter, r *http.Request) (int, error) {
	authUserDto, err := uh.parser.GetAuthUserDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	user, err := uh.storage.CheckUser(authUserDto)

	if err != nil {
		return http.StatusUnauthorized, err
	}

	responseDto, err := uh.tokenService.GenerateToken(user)
	if err != nil {
		slog.Debug(err.Error())
		return http.StatusBadRequest, err
	}
	responseText, err := json.MarshalIndent(responseDto, "", " ")
	if err != nil {
		return http.StatusBadRequest, err
	}
	w.Write(responseText)
	return http.StatusCreated, nil
}

// TODO tests
func (uh *UserHandlerImpl) HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error) {
	userDto, err := uh.parser.GetUserDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	user, err := uh.storage.SaveUser(userDto)
	if err != nil {
		return http.StatusConflict, err
	}
	responseDto, err := uh.tokenService.GenerateToken(user)
	if err != nil {
		slog.Debug(err.Error())
		return http.StatusBadRequest, err
	}
	responseText, err := json.MarshalIndent(responseDto, "", " ")
	if err != nil {
		return http.StatusBadRequest, err
	}
	w.Write(responseText)
	return http.StatusCreated, nil
}

func NewUserHandler(storage UserStorage) UserHandler {
	return &UserHandlerImpl{
		storage:      storage,
		parser:       newUserParser(),
		tokenService: token.NewTokenService(),
	}
}
