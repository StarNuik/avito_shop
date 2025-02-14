package shoptest

import "github.com/avito_shop/internal/domain"

type noopHasher struct{}

var _ domain.PasswordHasher = (*noopHasher)(nil)

func NewNoopHash() domain.PasswordHasher {
	return &noopHasher{}
}

func (_ *noopHasher) Hash(password string) (string, error) {
	return password, nil
}

func (_ *noopHasher) Same(inPassword string, storedHash string) bool {
	return inPassword == storedHash
}
