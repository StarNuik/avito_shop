package infra

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"math"
)

type inmemRepository struct {
	Users      map[int64]domain.User
	Operations map[int64]domain.BalanceOperation
	Transfers  map[int64]domain.Transfer
	Inventory  map[int64]domain.InventoryEntry
	Purchases  map[int64]domain.Purchase
}

var _ domain.ShopRepo = (*inmemRepository)(nil)

func NewInmemRepo() *inmemRepository {
	return &inmemRepository{
		Users:      make(map[int64]domain.User),
		Operations: make(map[int64]domain.BalanceOperation),
		Transfers:  make(map[int64]domain.Transfer),
		Inventory:  make(map[int64]domain.InventoryEntry),
		Purchases:  make(map[int64]domain.Purchase),
	}
}

func (repo *inmemRepository) InsertUser(user domain.User) {
	id := int64(len(repo.Users))
	user.Id = id
	repo.Users[id] = user
}

func (repo *inmemRepository) InsertInventory(item domain.InventoryEntry) {
	id := int64(len(repo.Inventory))
	item.Id = id
	repo.Inventory[id] = item
}

func (repo *inmemRepository) InsertBalanceOperation(op domain.BalanceOperation) {
	tx, _ := repo.Begin(context.Background())
	_, _ = tx.InsertBalanceOperation(op)
}

func (repo *inmemRepository) User(_ context.Context, username string) (domain.User, error) {
	for _, user := range repo.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrNotFound
}

func (repo *inmemRepository) UserBalance(ctx context.Context, userId int64) (int64, error) {
	best := domain.BalanceOperation{Id: math.MinInt64}

	for _, op := range repo.Operations {
		if op.Id != userId {
			continue
		}
		if op.Id > best.Id {
			best = op
		}
	}

	return best.Result, nil
}

func (repo *inmemRepository) InventoryItem(ctx context.Context, itemName string) (domain.InventoryEntry, error) {
	for _, item := range repo.Inventory {
		if item.Name == itemName {
			return item, nil
		}
	}

	return domain.InventoryEntry{}, domain.ErrNotFound
}

func (repo *inmemRepository) Begin(ctx context.Context) (domain.ShopTx, error) {
	return &inmemTx{repo, false}, nil
}

type inmemTx struct {
	*inmemRepository
	commit bool
}

var _ domain.ShopTx = (*inmemTx)(nil)

func (tx *inmemTx) InsertBalanceOperation(op domain.BalanceOperation) (int64, error) {
	id := int64(len(tx.Operations))
	op.Id = id
	tx.Operations[id] = op
	return id, nil
}

func (tx *inmemTx) InsertTransfer(t domain.Transfer) (int64, error) {
	id := int64(len(tx.Transfers))
	t.Id = id
	tx.Transfers[id] = t
	return id, nil
}

func (tx *inmemTx) InsertPurchase(p domain.Purchase) (int64, error) {
	id := int64(len(tx.Purchases))
	p.Id = id
	tx.Purchases[id] = p
	return id, nil
}

func (tx *inmemTx) Commit() error {
	tx.commit = true
	return nil
}

func (tx *inmemTx) Rollback() error {
	if !tx.commit {
		panic("not implemented")
	}
	return nil
}
