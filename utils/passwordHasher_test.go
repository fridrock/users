package utils

import (
	"testing"
)

func Test_checkPassword(t *testing.T) {
	hasher := PasswordHasherImpl{}
	password := "really long and strong password"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	checkResult := hasher.CheckPassword(password, hash)
	if !checkResult {
		t.Error("password and hash didn't pass check")
	}
}
