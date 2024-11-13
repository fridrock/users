package token

import (
	"encoding/json"
	"net/http"
)

type TokenDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type TokenDtoParser interface {
	GetDto(*http.Request) (TokenDto, error)
}
type TokenDtoParserImpl struct {
}

func NewTokenDtoParser() TokenDtoParser {
	return TokenDtoParserImpl{}
}

func (tdp TokenDtoParserImpl) GetDto(r *http.Request) (TokenDto, error) {
	var tokenDto TokenDto
	err := json.NewDecoder(r.Body).Decode(&tokenDto)
	return tokenDto, err
}
