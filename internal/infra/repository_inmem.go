package infra

import (
	"context"
	"github.com/avito_shop/internal/domain"
)

type inmemRepository struct {
	users map[int64]domain.User
}

var _ domain.ShopRepo = (*inmemRepository)(nil)

func NewInmemRepo() *inmemRepository {
	return &inmemRepository{
		users: make(map[int64]domain.User),
	}
}

func (repo *inmemRepository) InsertUser(user domain.User) {
	repo.users[user.Id] = user
}

func (repo *inmemRepository) User(_ context.Context, username string) (*domain.User, error) {
	for _, user := range repo.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, domain.ErrNotFound
}
