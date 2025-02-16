package domain

import (
	"context"
	"github.com/avito_shop/internal/dto"
	"slices"
)

func Info(ctx context.Context, repo ShopRepo, userId int64) (dto.InfoResponse, error) {
	out := dto.InfoResponse{}

	tx, err := repo.Begin(ctx)
	if err != nil {
		return dto.InfoResponse{}, err
	}
	defer tx.Rollback()

	out.Coins, err = tx.UserBalanceLock(userId)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	inventory, err := tx.InventoryInfo(userId)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	out.Inventory = make([]dto.InventoryInfo, 0, len(inventory))
	for _, item := range inventory {
		dto := dto.InventoryInfo{
			Type:     item.Name,
			Quantity: item.Quantity,
		}
		out.Inventory = append(out.Inventory, dto)
	}

	transfers, err := tx.UserTransfers(userId)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	out.CoinHistory.Received = make([]dto.BalanceDebitInfo, 0)
	out.CoinHistory.Sent = make([]dto.BalanceCreditInfo, 0)
	for _, transfer := range transfers {
		if transfer.FromUser == userId {
			out.CoinHistory.Sent = append(out.CoinHistory.Sent, dto.BalanceCreditInfo{
				ToUser: transfer.ToUsername,
				Amount: transfer.Delta,
			})
		} else if transfer.ToUser == userId {
			out.CoinHistory.Received = append(out.CoinHistory.Received, dto.BalanceDebitInfo{
				FromUser: transfer.FromUsername,
				Amount:   transfer.Delta,
			})
		}
	}
	slices.Reverse(out.CoinHistory.Received)
	slices.Reverse(out.CoinHistory.Sent)

	return out, tx.Commit()
}
