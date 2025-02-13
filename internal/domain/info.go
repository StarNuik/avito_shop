package domain

import (
	"context"
	"github.com/avito_shop/internal/dto"
	"slices"
)

func Info(ctx context.Context, repo ShopRepo, userId int64) (dto.InfoResponse, error) {
	out := dto.InfoResponse{}

	var err error
	out.Coins, err = repo.UserBalance(ctx, userId)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	inventory, err := repo.InventoryInfo(ctx, userId)
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

	balance, err := repo.BalanceInfo(ctx, userId)
	if err != nil {
		return dto.InfoResponse{}, err
	}

	out.CoinHistory.Received = make([]dto.BalanceDebitInfo, 0)
	out.CoinHistory.Sent = make([]dto.BalanceCreditInfo, 0)
	for _, op := range balance {
		if op.Delta >= 0 {
			dto := dto.BalanceDebitInfo{
				FromUser: op.ForeignUsername,
				Amount:   op.Delta,
			}
			out.CoinHistory.Received = append(out.CoinHistory.Received, dto)
		} else {
			dto := dto.BalanceCreditInfo{
				ToUser: op.ForeignUsername,
				Amount: -op.Delta,
			}
			out.CoinHistory.Sent = append(out.CoinHistory.Sent, dto)
		}
	}
	slices.Reverse(out.CoinHistory.Received)
	slices.Reverse(out.CoinHistory.Sent)

	return out, nil
}
