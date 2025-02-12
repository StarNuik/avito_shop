package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInfo_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()

	// Act
	ctx := context.Background()
	_, err := domain.Info(ctx, repo, -1)

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestInfo_NewUser_NoErrors(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username:     "username",
		PasswordHash: "password",
	}
	repo := infra.NewInmemRepo()
	user = repo.InsertUser(user)

	// Act
	ctx := context.Background()
	have, err := domain.Info(ctx, repo, user.Id)

	// Assert
	require.NoError(err)
	require.Equal(have.Coins, int64(0))
	require.Len(have.Inventory, 0)
	require.Len(have.CoinHistory.Sent, 0)
	require.Len(have.CoinHistory.Received, 0)
}

func TestInfo_HappyPath_CorrectFields(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()
	user := repo.InsertUser(domain.User{Username: "username"})

	usersForeign := []domain.User{
		repo.InsertUser(domain.User{Username: "foreign1"}),
		repo.InsertUser(domain.User{Username: "foreign2"}),
	}
	inventory := []domain.InventoryEntry{
		repo.InsertInventory(domain.InventoryEntry{Name: "merch1", Price: 10}),
		repo.InsertInventory(domain.InventoryEntry{Name: "merch2", Price: 100}),
	}

	balanceOps := []domain.BalanceOperation{
		// init
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  1000,
			Result: 1000,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   usersForeign[0].Id,
			Delta:  1000,
			Result: 1000,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   usersForeign[1].Id,
			Delta:  1000,
			Result: 1000,
		}),
		// transfer 1
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  -500,
			Result: 500,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   usersForeign[0].Id,
			Delta:  500,
			Result: 1500,
		}),
		// transfer 2
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   usersForeign[1].Id,
			Delta:  -200,
			Result: 800,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  200,
			Result: 700,
		}),
		// purchases
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  -100,
			Result: 600,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  -10,
			Result: 590,
		}),
		repo.InsertBalanceOperation(domain.BalanceOperation{
			User:   user.Id,
			Delta:  -10,
			Result: 580,
		}),
	}

	_ = []domain.Transfer{
		repo.InsertTransfer(domain.Transfer{
			SourceOp: balanceOps[3].Id, TargetOp: balanceOps[4].Id,
		}),
		repo.InsertTransfer(domain.Transfer{
			SourceOp: balanceOps[5].Id, TargetOp: balanceOps[6].Id,
		}),
	}

	_ = []domain.Purchase{
		repo.InsertPurchase(domain.Purchase{
			Item:      inventory[1].Id,
			User:      user.Id,
			Operation: balanceOps[7].Id,
		}),
		repo.InsertPurchase(domain.Purchase{
			Item:      inventory[0].Id,
			User:      user.Id,
			Operation: balanceOps[8].Id,
		}),
		repo.InsertPurchase(domain.Purchase{
			Item:      inventory[0].Id,
			User:      user.Id,
			Operation: balanceOps[9].Id,
		}),
	}

	// Act
	ctx := context.Background()
	have, err := domain.Info(ctx, repo, user.Id)

	// Assert
	require.NoError(err)
	require.Equal(have.Coins, int64(580))
	require.Len(have.Inventory, 2)
	require.Len(have.CoinHistory.Sent, 1)
	require.Equal(have.CoinHistory.Sent[0].Amount, int64(500))
	require.Equal(have.CoinHistory.Sent[0].ToUser, usersForeign[0].Username)
	require.Len(have.CoinHistory.Received, 1)
	require.Equal(have.CoinHistory.Received[0].Amount, int64(200))
	require.Equal(have.CoinHistory.Received[0].FromUser, usersForeign[1].Username)

	require.Contains(have.Inventory, dto.InventoryInfo{
		Type:     inventory[0].Name,
		Quantity: 2,
	})
	require.Contains(have.Inventory, dto.InventoryInfo{
		Type:     inventory[1].Name,
		Quantity: 1,
	})
}

func TestInfo_CoinHistoryOrder_NewFirst(t *testing.T) {
	panic("not implemented")
}
