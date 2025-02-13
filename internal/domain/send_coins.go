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

	if destUser.Id == userIdFrom {
		return ErrNotAllowed
	}

	destBalance, err := repo.UserBalance(ctx, destUser.Id)
	if err != nil {
		return err
	}

	fromBalance, err := repo.UserBalance(ctx, userIdFrom)
	if err != nil {
		return err
	}

	if transferSum > fromBalance {
		return ErrNotEnough
	}

	tx, err := repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	srcOp := BalanceOperation{
		User:   userIdFrom,
		Delta:  -transferSum,
		Result: fromBalance - transferSum,
	}
	srcOpId, err := tx.InsertBalanceOperation(srcOp)
	if err != nil {
		return err
	}

	destOp := BalanceOperation{
		User:   destUser.Id,
		Delta:  transferSum,
		Result: destBalance + transferSum,
	}
	destOpId, err := tx.InsertBalanceOperation(destOp)
	if err != nil {
		return err
	}

	transfer := Transfer{
		SourceOp: srcOpId,
		TargetOp: destOpId,
	}
	_, err = tx.InsertTransfer(transfer)
	if err != nil {
		return err
	}

	return tx.Commit()
}
