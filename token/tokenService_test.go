package token

import (
	"os"
	"testing"

	"github.com/fridrock/users/api"
	"github.com/google/uuid"
)

var tokenService TokenService

func TestMain(m *testing.M) {
	//setup
	os.Setenv("SECRET_KEY", "SECRET_FOR_TEST")
	os.Setenv("REFRESH_KEY", "SECRET_FOR_TEST")
	tokenService = NewTokenService()
	//running test
	m.Run()
}

func TestParseTokenSuccess(t *testing.T) {
	id := uuid.New()
	tokenDto, _ := tokenService.GenerateToken(api.User{
		Id: id,
	})
	parsedUser, err := tokenService.ValidateToken(tokenDto.AccessToken)
	if err != nil {
		t.Errorf(err.Error())
	}
	if parsedUser.Id.String() != id.String() {
		t.Errorf("Created with : %v, received with :%v", id, parsedUser.Id)
	}
}
func TestParseExpiredToken(t *testing.T) {
	expiredTokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzAzOTk2OTMsImlkIjoiZDUyN2Y5YzUtNGFmOS00NzdlLThlZDItOWNkMDAxMzk5NjkyIn0.i6C2Zaua3_gqYS_D1oxPDiuJDtcGbXTqwmiB-RkT8rs"
	_, err := tokenService.ValidateToken(expiredTokenString)
	if err == nil {
		t.Errorf("Without error on expired token")
	}
}
