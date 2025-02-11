package infra

import (
	"context"
	"github.com/avito_shop/internal/domain"
)

type mockRepository struct{}

var _ domain.ShopRepo = (*mockRepository)(nil)

func NewMockRepo() *mockRepository {
	return &mockRepository{}
}

func (repo *mockRepository) User(_ context.Context, username string) (*domain.User, error) {
	return &domain.User{
		Id:           -1,
		Username:     username,
		PasswordHash: "password", // TODO: hash passwords
	}, nil
}
