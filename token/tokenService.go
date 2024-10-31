package token

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/fridrock/users/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(api.User) (api.TokenDto, error)
	ParseToken(string) (api.User, error)
	RefreshToken(api.TokenDto) (api.TokenDto, error)
}

type TokenServiceImpl struct {
	SECRET_KEY []byte
}

func (ts *TokenServiceImpl) GenerateToken(user api.User) (api.TokenDto, error) {
	var dto api.TokenDto
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
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	return accessToken.SignedString(ts.SECRET_KEY)
}

func (ts *TokenServiceImpl) ParseToken(tokenString string) (api.User, error) {
	var dto api.User
	tokenObj, err := ts.validateToken(tokenString)
	if err != nil {
		return dto, err
	}
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		id, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
		if err != nil {
			return dto, err
		}
		exp, _ := claims.GetExpirationTime()

		slog.Info(fmt.Sprintf("%v", exp))
		if time.Now().After(exp.Time) {
			return dto, fmt.Errorf("Expired token %v", exp.Time)
		}
		dto.Id = id
	}
	return dto, nil
}
func (ts *TokenServiceImpl) validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ts.SECRET_KEY, nil
	})
}

func (ts *TokenServiceImpl) RefreshToken(api.TokenDto) (api.TokenDto, error) {
	var dto api.TokenDto
	return dto, nil
}

func NewTokenService() TokenService {
	varName := "SECRET_KEY"
	secret, exists := os.LookupEnv(varName)
	if !exists {
		log.Fatalf("Can't load env variable: %v", varName)
	}
	return &TokenServiceImpl{
		SECRET_KEY: []byte(secret),
	}
}
