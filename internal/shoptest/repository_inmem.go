package shoptest

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

func (repo *inmemRepository) InsertUser(user domain.User) domain.User {
	id := int64(len(repo.Users))
	user.Id = id
	repo.Users[id] = user
	return user
}

func (repo *inmemRepository) InsertInventory(item domain.InventoryEntry) domain.InventoryEntry {
	id := int64(len(repo.Inventory))
	item.Id = id
	repo.Inventory[id] = item
	return item
}

func (repo *inmemRepository) InsertBalanceOperation(op domain.BalanceOperation) domain.BalanceOperation {
	tx, _ := repo.Begin(context.Background())
	id, _ := tx.InsertBalanceOperation(op)
	op.Id = id
	return op
}

func (repo *inmemRepository) InsertTransfer(transfer domain.Transfer) domain.Transfer {
	tx, _ := repo.Begin(context.Background())
	id, _ := tx.InsertTransfer(transfer)
	transfer.Id = id
	return transfer
}

func (repo *inmemRepository) InsertPurchase(purchase domain.Purchase) domain.Purchase {
	tx, _ := repo.Begin(context.Background())
	id, _ := tx.InsertPurchase(purchase)
	purchase.Id = id
	return purchase
}

func (repo *inmemRepository) User(_ context.Context, username string) (domain.User, error) {
	for _, user := range repo.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrNotFound
}

// user exists, but has no operation should return (0, nil)
func (repo *inmemRepository) UserBalance(_ context.Context, userId int64) (int64, error) {
	if _, ok := repo.Users[userId]; !ok {
		return 0, domain.ErrNotFound
	}

	last := domain.BalanceOperation{Id: math.MinInt64}

	for idx := range len(repo.Operations) {
		op := repo.Operations[int64(idx)]
		if op.User != userId {
			continue
		}
		last = op
	}

	return last.Result, nil
}

func (repo *inmemRepository) InventoryItem(_ context.Context, itemName string) (domain.InventoryEntry, error) {
	for _, item := range repo.Inventory {
		if item.Name == itemName {
			return item, nil
		}
	}

	return domain.InventoryEntry{}, domain.ErrNotFound
}

// select group sum()
func (repo *inmemRepository) InventoryInfo(_ context.Context, userId int64) ([]domain.InventoryInfo, error) {
	if _, ok := repo.Users[userId]; !ok {
		return nil, domain.ErrNotFound
	}

	tmp := make(map[string]int64)

	for _, item := range repo.Purchases {
		if item.User != userId {
			continue
		}
		inventory := repo.Inventory[item.Item]
		tmp[inventory.Name]++
	}

	out := make([]domain.InventoryInfo, 0, len(tmp))
	for name, quantity := range tmp {
		item := domain.InventoryInfo{
			Name:     name,
			Quantity: quantity,
		}
		out = append(out, item)
	}

	return out, nil
}

func (repo *inmemRepository) BalanceInfo(_ context.Context, userId int64) ([]domain.BalanceInfo, error) {
	if _, ok := repo.Users[userId]; !ok {
		return nil, domain.ErrNotFound
	}

	out := []domain.BalanceInfo{}
	// TODO: is this SQL-able? split this into two functions?
	for idx := range len(repo.Transfers) {
		transfer := repo.Transfers[int64(idx)]
		srcOp := repo.Operations[transfer.SourceOp]
		dstOp := repo.Operations[transfer.TargetOp]

		srcUser := repo.Users[srcOp.User]
		dstUser := repo.Users[dstOp.User]

		if srcUser.Id != userId && dstUser.Id != userId {
			continue
		}

		item := domain.BalanceInfo{}
		switch userId {
		case dstUser.Id:
			item = domain.BalanceInfo{
				ForeignUsername: srcUser.Username,
				Delta:           -srcOp.Delta,
			}
		case srcUser.Id:
			item = domain.BalanceInfo{
				ForeignUsername: dstUser.Username,
				Delta:           -dstOp.Delta,
			}
		}
		out = append(out, item)
	}

	return out, nil
}

func (repo *inmemRepository) Begin(_ context.Context) (domain.ShopTx, error) {
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
