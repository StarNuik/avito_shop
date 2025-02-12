package domain

import (
	"context"
)

type ShopRepo interface {
	User(ctx context.Context, username string) (User, error)
	UserBalance(ctx context.Context, userId int64) (int64, error)
	InventoryItem(ctx context.Context, itemName string) (InventoryEntry, error)
	InventoryInfo(ctx context.Context, userId int64) ([]InventoryInfo, error)
	BalanceInfo(ctx context.Context, userId int64) ([]BalanceInfo, error)

	Begin(ctx context.Context) (ShopTx, error)
}

type ShopTx interface {
	InsertBalanceOperation(op BalanceOperation) (int64, error)
	InsertTransfer(t Transfer) (int64, error)
	InsertPurchase(p Purchase) (int64, error)

	Commit() error
	Rollback() error
}
