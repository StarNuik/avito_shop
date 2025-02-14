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
	err := repo.AddStagingUsers(hash)
	if err != nil {
		log.Panic(err)
	}

	repo.InitBalances()
	repo.AddStagingInventory()
}

func (repo *shopRepoBuilder) AddStagingUsers(hash domain.PasswordHasher) error {
	for _, user := range Users {
		passwordHash, err := hash.Hash(user.Password)
		if err != nil {
			return err
		}
		repo.InsertUser(domain.User{
			Username:     user.Password,
			PasswordHash: passwordHash,
		})
	}
	return nil
}

func (repo *shopRepoBuilder) AddStagingInventory() {
	for _, item := range Inventory {
		repo.InsertInventory(item)
	}
}

func (repo *shopRepoBuilder) InitBalances() {
	for userId := range repo.Users {
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   userId,
			Delta:  DefaultBalance,
			Result: DefaultBalance,
		})
	}
}
