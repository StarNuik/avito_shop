package shoptest

import (
	"context"
	"github.com/avito_shop/internal/domain"
)

type inmemRepository struct {
	Users     map[int64]domain.User
	Transfers map[int64]domain.Transfer
	Inventory map[int64]domain.InventoryItem
	Purchases map[int64]domain.Purchase
	Coins     map[int64]int64
}

var _ domain.ShopRepo = (*inmemRepository)(nil)

func NewInmemRepo() *inmemRepository {
	return &inmemRepository{
		Users:     make(map[int64]domain.User),
		Coins:     make(map[int64]int64),
		Transfers: make(map[int64]domain.Transfer),
		Inventory: make(map[int64]domain.InventoryItem),
		Purchases: make(map[int64]domain.Purchase),
	}
}

func (repo *inmemRepository) InsertUser(user domain.User, coins int64) domain.User {
	id := int64(len(repo.Users))
	user.Id = id
	repo.Users[id] = user
	repo.Coins[id] = coins
	return user
}

func (repo *inmemRepository) InsertInventory(item domain.InventoryItem) domain.InventoryItem {
	id := int64(len(repo.Inventory))
	item.Id = id
	repo.Inventory[id] = item
	return item
}

func (repo *inmemRepository) InsertTransfer(transfer domain.Transfer) domain.Transfer {
	tx, _ := repo.Begin(context.Background())
	id, _ := tx.InsertTransfer(transfer)
	_ = tx.Commit()
	transfer.Id = id
	return transfer
}

func (repo *inmemRepository) InsertPurchase(purchase domain.Purchase) domain.Purchase {
	tx, _ := repo.Begin(context.Background())
	id, _ := tx.InsertPurchase(purchase)
	_ = tx.Commit()
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

func (repo *inmemRepository) InventoryItem(_ context.Context, itemName string) (domain.InventoryItem, error) {
	for _, item := range repo.Inventory {
		if item.Name == itemName {
			return item, nil
		}
	}

	return domain.InventoryItem{}, domain.ErrNotFound
}

func (repo *inmemRepository) Begin(_ context.Context) (domain.ShopTx, error) {
	return &inmemTx{repo, false}, nil
}

type inmemTx struct {
	*inmemRepository
	commit bool
}

var _ domain.ShopTx = (*inmemTx)(nil)

// inmemRepository does not implement locking
func (tx *inmemTx) UserBalanceLock(userId int64) (int64, error) {
	if _, ok := tx.Users[userId]; !ok {
		return 0, domain.ErrNotFound
	}

	return tx.Coins[userId], nil
}

// inmemRepository does not implement locking
func (tx *inmemTx) UserPairBalanceLock(fromUser int64, toUser int64) (int64, int64, error) {
	if _, ok := tx.Users[fromUser]; !ok {
		return 0, 0, domain.ErrNotFound
	}
	if _, ok := tx.Users[toUser]; !ok {
		return 0, 0, domain.ErrNotFound
	}

	return tx.Coins[fromUser], tx.Coins[toUser], nil
}

// select group sum()
func (tx *inmemTx) InventoryInfo(userId int64) ([]domain.InventoryInfo, error) {
	if _, ok := tx.Users[userId]; !ok {
		return nil, domain.ErrNotFound
	}

	tmp := make(map[string]int64)

	for _, item := range tx.Purchases {
		if item.UserId != userId {
			continue
		}
		inventory := tx.Inventory[item.Item]
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

func (tx *inmemTx) UserTransfers(userId int64) ([]domain.TransferInfo, error) {
	if _, ok := tx.Users[userId]; !ok {
		return nil, domain.ErrNotFound
	}

	out := []domain.TransferInfo{}
	for id := range len(tx.Transfers) {
		transfer := tx.Transfers[int64(id)]
		if transfer.FromUser != userId &&
			transfer.ToUser != userId {
			continue
		}
		out = append(out, domain.TransferInfo{
			Transfer:     transfer,
			FromUsername: tx.Users[transfer.FromUser].Username,
			ToUsername:   tx.Users[transfer.ToUser].Username,
		})
	}
	return out, nil
}

func (tx *inmemTx) UpdateBalance(userId int64, amount int64) error {
	tx.Coins[userId] = amount
	return nil
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
	//tx.commit = true
	return nil
}

func (tx *inmemTx) Rollback() error {
	//if !tx.commit {
	//	panic("not implemented")
	//}
	return nil
}
