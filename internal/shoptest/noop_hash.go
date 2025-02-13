package shoptest

import "github.com/avito_shop/internal/domain"

type noopPasswordHash struct{}

var _ domain.PasswordHash = (*noopPasswordHash)(nil)

func NewNoopHash() domain.PasswordHash {
	return &noopPasswordHash{}
}

func (_ *noopPasswordHash) Hash(pass string) (string, error) {
	return pass, nil
}

func (_ *noopPasswordHash) Same(lhs string, rhs string) bool {
	return lhs == rhs
}
