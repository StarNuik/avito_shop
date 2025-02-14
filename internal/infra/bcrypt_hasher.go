package infra

import (
	"github.com/avito_shop/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

var _ domain.PasswordHasher = (*BcryptHasher)(nil)

func (*BcryptHasher) Hash(password string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (*BcryptHasher) Same(inPassword string, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inPassword))
	return err == nil
}
