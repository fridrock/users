package usr

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fridrock/users/token"
	"github.com/fridrock/users/utils"
	"github.com/google/uuid"
)

type UserHandler interface {
	HandleRegistration(w http.ResponseWriter, r *http.Request) (int, error)
	HandleAuth(w http.ResponseWriter, r *http.Request) (int, error)
	FindUser(w http.ResponseWriter, r *http.Request) (int, error)
	GetUsers(w http.ResponseWriter, r *http.Request) (int, error)
	GetProfiles(w http.ResponseWriter, r *http.Request) (int, error)
}

type UserHandlerImpl struct {
	storage      UserStorage
	tokenService token.TokenService
	parser       UserParser
}

func (uh *UserHandlerImpl) GetProfiles(w http.ResponseWriter, r *http.Request) (int, error) {
	r.ParseForm()
	idsParam := r.FormValue("ids")
	idsString := strings.Split(idsParam, ",")
	var ids []uuid.UUID
	for _, idString := range idsString {
		id, _ := uuid.Parse(idString)
		ids = append(ids, id)
	}
	users, err := uh.storage.GetProfiles(ids)
	if err != nil {
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func (uh *UserHandlerImpl) FindUser(w http.ResponseWriter, r *http.Request) (int, error) {
	username, err := uh.parser.GetUsername(r)
	if err != nil {
		slog.Error("error parsing username" + err.Error())
		return http.StatusBadRequest, err
	}
	usersFound, err := uh.storage.FindUsers(username)
	if err != nil {
		slog.Debug(fmt.Sprintf("user with name %v not found: %v", username, err.Error()))
		return http.StatusNotFound, err
	}
	filtered, err := uh.storage.FilterNotInFriends(usersFound, utils.UserFromContext(r.Context()))
	if err != nil {
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(filtered, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}
func (uh *UserHandlerImpl) GetUsers(w http.ResponseWriter, r *http.Request) (int, error) {
	usersFound, err := uh.storage.GetUsers()
	if err != nil {
		return http.StatusNotFound, err
	}
	filtered, err := uh.storage.FilterNotInFriends(usersFound, utils.UserFromContext(r.Context()))
	if err != nil {
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(filtered, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusCreated, nil
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
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
