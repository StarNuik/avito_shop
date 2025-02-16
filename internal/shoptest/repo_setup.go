package shoptest

import (
	"github.com/avito_shop/internal/domain"
	"github.com/jackc/pgx/v5"
	"log"
)

type shopRepoBuilder struct {
	*inmemRepository
}

func NewShopRepoBuilder() *shopRepoBuilder {
	inmem := NewInmemRepo()
	return &shopRepoBuilder{inmem}
}

func AddStagingValues(db *pgx.Conn, hash domain.PasswordHasher) {
	repo := NewShopRepo(db)
	for range 4 {
		repo.Clear("Users")
		repo.Clear("Purchases")
		repo.Clear("Inventory")
		repo.Clear("Transfers")
	}
	err := AddStagingUsers(repo, hash, DefaultBalance)
	if err != nil {
		log.Panic(err)
	}

	AddStagingInventory(repo)
}

func AddStagingUsers(repo *shopRepoPostgres, hash domain.PasswordHasher, startingBalance int64) error {
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

func AddStagingInventory(repo *shopRepoPostgres) {
	for _, item := range Inventory {
		repo.InsertInventory(item)
	}
}
