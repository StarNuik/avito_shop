package domain

import (
	"context"
)

// BuyItem may return these errors: NotEnough, NotFound
func BuyItem(ctx context.Context, repo ShopRepo, userId int64, itemName string) error {
	item, err := repo.InventoryItem(ctx, itemName)
	if err != nil {
		return err
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balance, err := tx.UserBalanceLock(userId)
	if err != nil {
		return err
	}

	if item.Price > balance {
		return ErrNotEnough
	}

	purchase := Purchase{
		Item:   item.Id,
		UserId: userId,
		Price:  item.Price,
	}
	_, err = tx.InsertPurchase(purchase)
	if err != nil {
		return err
	}

	err = tx.UpdateBalance(userId, balance-item.Price)

	return tx.Commit()
}
