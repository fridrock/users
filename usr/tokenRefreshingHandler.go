package usr

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fridrock/users/token"
)

type TokenRefreshHandler interface {
	HandleRefreshToken(w http.ResponseWriter, r *http.Request) (int, error)
}

type TokenRefreshHandlerImpl struct {
	tokenService token.TokenService
	parser       token.TokenDtoParser
}

// TODO tests
func (rs *TokenRefreshHandlerImpl) HandleRefreshToken(w http.ResponseWriter, r *http.Request) (int, error) {
	tokenDto, err := rs.parser.GetDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	receivedTokenPair, err := rs.tokenService.RefreshToken(tokenDto)
	if err != nil {
		slog.Debug(err.Error())
		return http.StatusBadRequest, err
	}
	responseText, err := json.MarshalIndent(receivedTokenPair, "", " ")
	if err != nil {
		return http.StatusBadRequest, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseText)
	return http.StatusCreated, nil
}
func NewTokenRefreshHandler() TokenRefreshHandler {
	return &TokenRefreshHandlerImpl{
		parser:       token.NewTokenDtoParser(),
		tokenService: token.NewTokenService(),
	}
}
