package repository

import (
	"context"
	"github.com/avito_shop/internal/domain"
)

type shopRepoPostgres struct{}

type shopTxPostgres struct{}

var _ (domain.ShopRepo) = &shopRepoPostgres{}
var _ (domain.ShopTx) = &shopTxPostgres{}

func NewShopPostgres(url string) domain.ShopRepo {
	//
}

func (s shopRepoPostgres) User(ctx context.Context, username string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopRepoPostgres) UserBalance(ctx context.Context, userId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopRepoPostgres) InventoryItem(ctx context.Context, itemName string) (domain.InventoryEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopRepoPostgres) InventoryInfo(ctx context.Context, userId int64) ([]domain.InventoryInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopRepoPostgres) BalanceInfo(ctx context.Context, userId int64) ([]domain.BalanceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopRepoPostgres) Begin(ctx context.Context) (domain.ShopTx, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopTxPostgres) InsertBalanceOperation(op domain.BalanceOperation) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopTxPostgres) InsertTransfer(t domain.Transfer) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopTxPostgres) InsertPurchase(p domain.Purchase) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s shopTxPostgres) Commit() error {
	//TODO implement me
	panic("implement me")
}

func (s shopTxPostgres) Rollback() error {
	//TODO implement me
	panic("implement me")
}
