package shoptest

import (
	"github.com/avito_shop/internal/domain"
)

type shopRepoBuilder struct {
	*inmemRepository
}

func NewShopRepoBuilder() *shopRepoBuilder {
	inmem := NewInmemRepo()
	return &shopRepoBuilder{inmem}
}

func (repo *shopRepoBuilder) AddStagingValues() {
	repo.InsertUser(domain.User{
		Id:           -1,
		Username:     "admin",
		PasswordHash: "admin",
	})
	repo.InsertUser(domain.User{
		Id:           -2,
		Username:     "test",
		PasswordHash: "test",
	})

	repo.InitBalances()

	repo.InsertInventory(domain.InventoryEntry{Name: "hoodie", Price: 100})
	repo.InsertInventory(domain.InventoryEntry{Name: "keychain", Price: 10})
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
