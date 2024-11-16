package friend

import (
	"encoding/json"
	h "net/http"

	"github.com/fridrock/users/api"
	"github.com/fridrock/users/utils"
	"github.com/google/uuid"
)

type FriendParser interface {
	ParseFriendDto(r *h.Request) (api.FriendDto, error)
	GetUserId(r *h.Request) uuid.UUID
}

type FriendParserImpl struct{}

func (fp FriendParserImpl) ParseFriendDto(r *h.Request) (api.FriendDto, error) {
	var friendDto api.FriendDto
	err := json.NewDecoder(r.Body).Decode(&friendDto)
	if err != nil {
		return friendDto, err
	}
	friendDto.UserId = fp.GetUserId(r)
	return friendDto, nil
}

func (fp FriendParserImpl) GetUserId(r *h.Request) uuid.UUID {
	return utils.UserFromContext(r.Context())
}
func newFriendParser() FriendParser {
	return FriendParserImpl{}
}
