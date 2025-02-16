package repository

import (
	"context"
	"errors"
	"github.com/avito_shop/internal/domain"
	"github.com/jackc/pgx/v5"
)

type shopRepoPostgres struct {
	*pgx.Conn
}

type shopTxPostgres struct {
	pgx.Tx
	ctx context.Context
}

var _ (domain.ShopRepo) = &shopRepoPostgres{}
var _ (domain.ShopTx) = &shopTxPostgres{}

func NewShopPostgres(db *pgx.Conn) domain.ShopRepo {
	return &shopRepoPostgres{
		Conn: db,
	}
}

func (repo *shopRepoPostgres) User(ctx context.Context, username string) (domain.User, error) {
	row := repo.QueryRow(ctx, `
		select Id, Username, PasswordHash
		from Users
		where Username = $1
	`, username)

	out := domain.User{}
	err := row.Scan(&out.Id, &out.Username, &out.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, domain.ErrNotFound
	}
	return out, err
}

func (repo *shopRepoPostgres) InventoryItem(ctx context.Context, itemName string) (domain.InventoryItem, error) {
	row := repo.QueryRow(ctx, `
		select Id, Name, Price
		from Inventory
		where Name = $1
	`, itemName)

	out := domain.InventoryItem{}
	err := row.Scan(&out.Id, &out.Name, &out.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, domain.ErrNotFound
	}
	return out, err
}

func (repo *shopRepoPostgres) Begin(ctx context.Context) (domain.ShopTx, error) {
	tx, err := repo.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &shopTxPostgres{
		Tx:  tx,
		ctx: ctx,
	}, nil
}

func (tx *shopTxPostgres) Commit() error {
	return tx.Tx.Commit(tx.ctx)
}

func (tx *shopTxPostgres) Rollback() error {
	return tx.Tx.Rollback(tx.ctx)
}

func (tx *shopTxPostgres) UserBalanceLock(userId int64) (int64, error) {
	row := tx.QueryRow(tx.ctx, `
		select Coins
		from Users
		where Id = $1
		for update
	`, userId)

	var out int64
	err := row.Scan(&out)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, domain.ErrNotFound
	}
	return out, err
}

func (tx *shopTxPostgres) UserPairBalanceLock(fromUser int64, toUser int64) (int64, int64, error) {
	rows, err := tx.Query(tx.ctx, `
		select Id, Coins
		from Users
		where Id = $1 or Id = $2
		for update
	`, fromUser, toUser)
	if err != nil {
		return 0, 0, err
	}

	cache := make(map[int64]int64)
	for rows.Next() {
		userId, balance := int64(0), int64(0)

		err := rows.Scan(&userId, &balance)
		if err != nil {
			return 0, 0, err
		}

		cache[userId] = balance
	}

	fromBalance, ok := cache[fromUser]
	if !ok {
		return 0, 0, domain.ErrNotFound
	}

	toBalance, ok := cache[toUser]
	if !ok {
		return 0, 0, domain.ErrNotFound
	}

	return fromBalance, toBalance, nil
}

func (tx *shopTxPostgres) UpdateBalance(userId int64, balance int64) error {
	_, err := tx.Exec(tx.ctx, `
		update Users
		set Coins = $2
		where Id = $1
	`, userId, balance)

	return err
	//if tag.RowsAffected() == 0 {
	//	return domain.ErrNotFound
	//}
}

func (tx *shopTxPostgres) InsertTransfer(transfer domain.Transfer) (int64, error) {
	row := tx.QueryRow(tx.ctx, `
		insert into Transfers(Delta, FromUser, ToUser)
		values ($1, $2, $3)
		returning Id
	`, transfer.Delta, transfer.FromUser, transfer.ToUser)

	var transferId int64
	err := row.Scan(&transferId)
	return transferId, err
}

func (tx *shopTxPostgres) InsertPurchase(purchase domain.Purchase) (int64, error) {
	row := tx.QueryRow(tx.ctx, `
		insert into Purchases(Price, Item, UserId)
		values ($1, $2, $3)
		returning Id
	`, purchase.Price, purchase.Item, purchase.UserId)

	var purchaseId int64
	err := row.Scan(&purchaseId)
	return purchaseId, err
}

func (tx *shopTxPostgres) UserTransfers(userId int64) ([]domain.TransferInfo, error) {
	rows, err := tx.Query(tx.ctx, `
		select Delta, FromUser, ToUser, u1.Username FromUsername, u2.Username ToUsername
		from Transfers t
		    join Users u1
		    on t.FromUser = u1.Id
		    join Users u2
		    on t.ToUser = u2.Id
		where t.FromUser = $1 or t.ToUser = $1
	`, userId)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.TransferInfo])
}

func (tx *shopTxPostgres) InventoryInfo(userId int64) ([]domain.InventoryInfo, error) {
	rows, err := tx.Query(tx.ctx, `
		select i.Name, count(p.Id) Quantity
		from Purchases p
		    join Inventory i
		    on p.Item = i.Id
		where p.UserId = $1
		group by i.Id
	`, userId)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.InventoryInfo])
}
