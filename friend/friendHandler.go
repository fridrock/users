package friend

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type FriendHandler interface {
	AddFriend(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteFriend(w http.ResponseWriter, r *http.Request) (int, error)
	GetFriends(w http.ResponseWriter, r *http.Request) (int, error)
}

type FriendHandlerImpl struct {
	storage FriendStorage
	parser  FriendParser
}

func (fh *FriendHandlerImpl) AddFriend(w http.ResponseWriter, r *http.Request) (int, error) {
	// parse request and user info
	friendDto, err := fh.parser.ParseFriendDto(r)
	if err != nil {
		slog.Debug("error parsing friend dto" + err.Error())
		return http.StatusBadRequest, err
	}
	err = fh.storage.AddFriend(friendDto.UserId, friendDto.FriendId)
	if err != nil {
		slog.Debug("error adding friend in storage" + err.Error())
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func (fh *FriendHandlerImpl) DeleteFriend(w http.ResponseWriter, r *http.Request) (int, error) {
	friendDto, err := fh.parser.ParseFriendDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = fh.storage.DeleteFriend(friendDto.UserId, friendDto.FriendId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
func (fh *FriendHandlerImpl) GetFriends(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := fh.parser.GetUserId(r)
	friends, err := fh.storage.GetFriends(userId)
	if err != nil {
		return http.StatusNotFound, err
	}
	w.Header().Set("Content-Type", "application/json")
	// Маршалим срез в JSON
	response, err := json.MarshalIndent(friends, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Записываем JSON-данные в ответ
	w.Write(response)

	return http.StatusOK, nil
}

func NewFriendHandler(storage FriendStorage) FriendHandler {
	return &FriendHandlerImpl{
		storage: storage,
		parser:  newFriendParser(),
	}
}
