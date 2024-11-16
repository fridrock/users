package api

import "github.com/google/uuid"

type FriendDto struct {
	UserId   uuid.UUID
	FriendId uuid.UUID `json:"friendId"`
}
