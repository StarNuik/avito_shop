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

	balance, err := repo.UserBalance(ctx, userId)
	if err != nil {
		return err
	}

	if item.Price > balance {
		return ErrNotEnough
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	op := BalanceOperation{
		User:   userId,
		Delta:  -item.Price,
		Result: balance - item.Price,
	}
	opId, err := tx.InsertBalanceOperation(op)
	if err != nil {
		return err
	}

	purchase := Purchase{
		Item:      item.Id,
		User:      userId,
		Operation: opId,
	}
	_, err = tx.InsertPurchase(purchase)
	if err != nil {
		return err
	}

	return tx.Commit()
}
