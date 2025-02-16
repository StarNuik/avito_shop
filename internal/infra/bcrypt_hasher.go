package infra

import (
	"github.com/avito_shop/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

var _ domain.PasswordHasher = (*BcryptHasher)(nil)

func (*BcryptHasher) Hash(password string) (string, error) {
	return hash(password, bcrypt.DefaultCost)
}

func (*BcryptHasher) Same(inPassword string, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inPassword))
	return err == nil
}

func (*BcryptHasher) HashFast(password string) (string, error) {
	return hash(password, bcrypt.MinCost)
}

func hash(password string, cost int) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
