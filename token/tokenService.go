package token

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fridrock/users/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(api.User) (TokenDto, error)
	ValidateToken(string) (api.User, error)
	RefreshToken(TokenDto) (TokenDto, error)
}

type TokenServiceImpl struct {
	SECRET_KEY  []byte
	REFRESH_KEY string
}

func (ts *TokenServiceImpl) GenerateToken(user api.User) (TokenDto, error) {
	var dto TokenDto
	accessTokenString, err := ts.generateAccess(user)
	if err != nil {
		return dto, err
	}
	dto.AccessToken = accessTokenString
	refreshTokenString, err := ts.generateRefresh(user)
	if err != nil {
		return dto, err
	}
	dto.RefreshToken = refreshTokenString
	return dto, nil
}

func (ts *TokenServiceImpl) generateAccess(user api.User) (string, error) {
	//When we need some additional info we can add it and rewrite logic of refresh token generation
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"exp": jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	})
	return accessToken.SignedString(ts.SECRET_KEY)
}

func (ts *TokenServiceImpl) generateRefresh(user api.User) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            user.Id,
		"exp":           jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		"refreshSecret": ts.REFRESH_KEY,
	})
	return refreshToken.SignedString(ts.SECRET_KEY)
}

func (ts *TokenServiceImpl) ValidateToken(tokenString string) (api.User, error) {
	var dto api.User
	parsed, err := ts.parseToken(tokenString)
	if err != nil {
		return dto, err
	}
	if time.Now().After(parsed.exp.Time) {
		return dto, fmt.Errorf("expired token %v", parsed.exp.Time)
	}
	dto.Id = parsed.id
	return dto, nil
}

type tokenParsed struct {
	id            uuid.UUID
	exp           *jwt.NumericDate
	refreshSecret string
}

func (ts *TokenServiceImpl) parseToken(tokenString string) (tokenParsed, error) {
	var dto tokenParsed
	tokenObj, err := ts.checkSigning(tokenString)
	if err != nil {
		return dto, err
	}
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		id, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
		if err != nil {
			return dto, err
		}
		dto.id = id
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return dto, err
		}
		dto.exp = exp
		dto.refreshSecret = fmt.Sprintf("%v", claims["refreshSecret"])

	}
	return dto, nil
}

func (ts *TokenServiceImpl) checkSigning(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.SECRET_KEY, nil
	})
}

func (ts *TokenServiceImpl) RefreshToken(incomingToken TokenDto) (TokenDto, error) {
	var dto TokenDto
	parsed, err := ts.parseToken(incomingToken.RefreshToken)
	if err != nil {
		return dto, err
	}
	if parsed.refreshSecret != ts.REFRESH_KEY {
		return dto, fmt.Errorf("wrong refresh key")
	}
	user := api.User{
		Id: parsed.id,
	}
	dto, err = ts.GenerateToken(user)
	return dto, err
}

func NewTokenService() TokenService {
	varName := "SECRET_KEY"
	secret, exists := os.LookupEnv(varName)
	if !exists {
		log.Fatalf("Can't load env variable: %v", varName)
	}
	varName = "REFRESH_KEY"
	refresh, exists := os.LookupEnv(varName)
	if !exists {
		log.Fatalf("Can't load env variable: %v", varName)
	}
	return &TokenServiceImpl{
		SECRET_KEY:  []byte(secret),
		REFRESH_KEY: refresh,
	}
}
