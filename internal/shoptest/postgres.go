package shoptest

import (
	"context"
	"fmt"
	"github.com/avito_shop/internal/domain"
	"github.com/jackc/pgx/v5"
)

type shopRepoPostgres struct {
	*pgx.Conn
}

func NewShopRepo(db *pgx.Conn) *shopRepoPostgres {
	return &shopRepoPostgres{
		Conn: db,
	}
}

func (repo *shopRepoPostgres) Clear(table string) error {
	ctx := context.Background()
	_, err := repo.Exec(ctx, fmt.Sprintf("delete from %s", table))
	return err
}

func (repo *shopRepoPostgres) InsertUser(user domain.User, balance int64) (int64, error) {
	ctx := context.Background()
	row := repo.QueryRow(ctx, `
        insert into Users (Username, PasswordHash, Coins)
        values ($1, $2, $3)
        returning Id
        `, user.Username, user.PasswordHash, balance)

	var userId int64
	err := row.Scan(&userId)
	return userId, err
}

func (repo *shopRepoPostgres) InsertInventory(item domain.InventoryItem) (int64, error) {
	ctx := context.Background()
	row := repo.QueryRow(ctx, `
        insert into Inventory (Name, Price)
        values ($1, $2)
        returning Id
    `, item.Name, item.Price)

	var userId int64
	err := row.Scan(&userId)
	return userId, err
}
