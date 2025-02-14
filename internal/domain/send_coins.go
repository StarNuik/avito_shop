package domain

import (
	"context"
)

// SendCoins may send these errors: ErrNotAllowed, ErrNotEnough, ErrNotFound
func SendCoins(ctx context.Context, repo ShopRepo, userIdFrom int64, usernameTo string, transferSum int64) error {
	if transferSum <= 0 {
		return ErrNotAllowed
	}

	destUser, err := repo.User(ctx, usernameTo)
	if err != nil {
		return err
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if destUser.Id == userIdFrom {
		return ErrNotAllowed
	}

	fromBalance, destBalance, err := tx.UserPairBalanceLock(userIdFrom, destUser.Id)
	if err != nil {
		return err
	}

	if transferSum > fromBalance {
		return ErrNotEnough
	}

	transfer := Transfer{
		FromUser: userIdFrom,
		ToUser:   destUser.Id,
		Delta:    transferSum,
	}
	_, err = tx.InsertTransfer(transfer)
	if err != nil {
		return err
	}

	err = tx.UpdateBalance(destUser.Id, destBalance+transferSum)
	if err != nil {
		return err
	}

	err = tx.UpdateBalance(userIdFrom, fromBalance-transferSum)
	if err != nil {
		return err
	}

	return tx.Commit()
}
