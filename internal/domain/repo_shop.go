package domain

import (
	"context"
)

type ShopRepo interface {
	User(ctx context.Context, username string) (User, error)
	InventoryItem(ctx context.Context, itemName string) (InventoryItem, error)

	Begin(ctx context.Context) (ShopTx, error)
}

type ShopTx interface {
	// UserBalanceLock locks the Users table row
	UserBalanceLock(userId int64) (int64, error)
	// UserPairBalanceLock locks 2 rows in the Users table
	UserPairBalanceLock(fromUser int64, toUser int64) (int64, int64, error)

	InventoryInfo(userId int64) ([]InventoryInfo, error)
	UserTransfers(userId int64) ([]TransferInfo, error)

	UpdateBalance(userId int64, balance int64) error

	InsertTransfer(t Transfer) (int64, error)
	InsertPurchase(p Purchase) (int64, error)

	Commit() error
	Rollback() error
}
