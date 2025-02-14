package shoptest

import (
	"github.com/avito_shop/internal/domain"
	"log"
)

type shopRepoBuilder struct {
	*inmemRepository
}

func NewShopRepoBuilder() *shopRepoBuilder {
	inmem := NewInmemRepo()
	return &shopRepoBuilder{inmem}
}

func (repo *shopRepoBuilder) AddStagingValues(hash domain.PasswordHasher) {
	err := repo.AddStagingUsers(hash, DefaultBalance)
	if err != nil {
		log.Panic(err)
	}

	repo.AddStagingInventory()
}

func (repo *shopRepoBuilder) AddStagingUsers(hash domain.PasswordHasher, startingBalance int64) error {
	for _, user := range Users {
		passwordHash, err := hash.Hash(user.Password)
		if err != nil {
			return err
		}
		repo.InsertUser(domain.User{
			Username:     user.Password,
			PasswordHash: passwordHash,
		}, startingBalance)
	}
	return nil
}

func (repo *shopRepoBuilder) AddStagingInventory() {
	for _, item := range Inventory {
		repo.InsertInventory(item)
	}
}
