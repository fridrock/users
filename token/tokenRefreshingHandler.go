package token

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type TokenRefreshHandler interface {
	HandleRefreshToken(w http.ResponseWriter, r *http.Request) (int, error)
}

type TokenRefreshHandlerImpl struct {
	tokenService TokenService
	parser       TokenDtoParser
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
	w.Write(responseText)
	return http.StatusCreated, nil
}
func NewTokenRefreshHandler() TokenRefreshHandler {
	return &TokenRefreshHandlerImpl{
		parser:       NewTokenDtoParser(),
		tokenService: NewTokenService(),
	}
}