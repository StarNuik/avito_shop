package shoptest

import (
	"context"
	"fmt"
	"github.com/avito_shop/internal/setup"
	"github.com/jackc/pgx/v5"
)

func ClearRepo() error {
	env := setup.GetEnv()
	db, err := pgx.Connect(context.Background(), env.DatabaseUrl)
	if err != nil {
		return err
	}

	repo := NewShopRepo(db)
	repo.Clear("Purchases")
	repo.Clear("Transfers")

	tag, err := db.Exec(context.Background(), `
        update Users
        set Coins = $1
    `, DefaultBalance)

	if tag.RowsAffected() != UserCount {
		return fmt.Errorf("RowsAffected != UserCount")
	}
	return nil
}
