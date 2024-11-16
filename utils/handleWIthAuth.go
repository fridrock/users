package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fridrock/users/api"
	"github.com/fridrock/users/token"
	"github.com/google/uuid"
)

type UserContextKey string

const key UserContextKey = "user"

func UserFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(key).(uuid.UUID)
}

type AuthManager interface {
	HandleWithAuth(h HandlerWithError) HandlerWithError
}
type AuthManagerImpl struct {
	tokenService token.TokenService
}

func NewAuthManager() AuthManager {
	return &AuthManagerImpl{
		tokenService: token.NewTokenService(),
	}
}

func (am AuthManagerImpl) getUserFromToken(r *http.Request) (api.User, error) {
	var user api.User
	// Извлечение заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return user, fmt.Errorf("empty auth header")
	}

	// Проверка и разбиение заголовка
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return user, fmt.Errorf("wrong format of token")
	}

	user, err := am.tokenService.ValidateToken(parts[1])
	if err != nil {
		return user, fmt.Errorf("token invalidated")
	}
	return user, nil
}

func (am AuthManagerImpl) HandleWithAuth(h HandlerWithError) HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		user, err := am.getUserFromToken(r)
		if err != nil {
			return http.StatusUnauthorized, err
		}
		ctx := context.WithValue(r.Context(), key, user.Id)
		return h(w, r.WithContext(ctx))
	}
}
