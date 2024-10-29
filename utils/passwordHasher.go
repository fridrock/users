package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(string) (string, error)
	CheckPassword(string, string) bool
}
type PasswordHasherImpl struct{}

func (ph PasswordHasherImpl) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func (ph PasswordHasherImpl) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func NewPasswordHasher() PasswordHasher {
	return &PasswordHasherImpl{}
}
