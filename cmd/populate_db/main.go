package main

import (
	"context"
	"fmt"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/shoptest"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	dbUrl := "postgres://postgres:password@localhost:5432/shop"
	db, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Panic(err)
	}
	repo := shoptest.NewShopRepo(db)
	hasher := infra.BcryptHasher{}

	repo.Clear("Purchases")
	repo.Clear("Transfers")

	for _, item := range shoptest.Inventory {
		_, err := repo.InsertInventory(item)
		if err != nil {
			log.Panic(err)
		}
	}
	fmt.Println("Inventory added")

	for idx := range 100_000 {
		username := fmt.Sprintf("user#%d", idx)
		password := fmt.Sprintf("password#%d", idx)

		passwordHash, err := hasher.HashFast(password)
		if err != nil {
			log.Panic(err)
		}

		user := domain.User{
			Username:     username,
			PasswordHash: passwordHash,
		}
		err = repo.InsertUserFast(user, shoptest.DefaultBalance)
		if err != nil {
			log.Panic(err)
		}

		if idx%1000 == 0 {
			fmt.Printf("Insert users: %d%%\n", idx/1000)
		}
	}
	fmt.Println("Users added")

}
