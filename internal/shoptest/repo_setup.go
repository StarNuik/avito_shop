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

func (repo *shopRepoBuilder) AddStagingValues(hash domain.PasswordHash) {
	err := repo.AddStagingUsers(hash)
	if err != nil {
		log.Panic(err)
	}

	repo.InitBalances()

	repo.InsertInventory(domain.InventoryEntry{Name: "hoodie", Price: 100})
	repo.InsertInventory(domain.InventoryEntry{Name: "keychain", Price: 10})
}

func (repo *shopRepoBuilder) AddStagingUsers(hash domain.PasswordHash) error {
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

func (repo *shopRepoBuilder) InitBalances() {
	for userId := range repo.Users {
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   userId,
			Delta:  DefaultBalance,
			Result: DefaultBalance,
		})
	}
}
