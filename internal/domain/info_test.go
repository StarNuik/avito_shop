package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInfo_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()

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
	repo := shoptest.NewInmemRepo()
	user = repo.InsertUser(user, shoptest.DefaultBalance)

	// Act
	ctx := context.Background()
	have, err := domain.Info(ctx, repo, user.Id)

	// Assert
	require.NoError(err)
	require.Equal(have.Coins, shoptest.DefaultBalance)
	require.Len(have.Inventory, 0)
	require.Len(have.CoinHistory.Sent, 0)
	require.Len(have.CoinHistory.Received, 0)
}

func TestInfo_HappyPath_CorrectFields(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()
	user := repo.InsertUser(domain.User{Username: "username"}, shoptest.DefaultBalance)

	usersForeign := []domain.User{
		repo.InsertUser(domain.User{Username: "foreign1"}, shoptest.DefaultBalance),
		repo.InsertUser(domain.User{Username: "foreign2"}, shoptest.DefaultBalance),
	}
	inventory := []domain.InventoryItem{
		repo.InsertInventory(domain.InventoryItem{Name: "merch1", Price: 10}),
		repo.InsertInventory(domain.InventoryItem{Name: "merch2", Price: 100}),
	}

	_ = []domain.Transfer{
		repo.InsertTransfer(domain.Transfer{
			FromUser: user.Id,
			ToUser:   usersForeign[0].Id,
			Delta:    500,
		}),
		repo.InsertTransfer(domain.Transfer{
			FromUser: usersForeign[1].Id,
			ToUser:   user.Id,
			Delta:    200,
		}),
	}

	_ = []domain.Purchase{
		repo.InsertPurchase(domain.Purchase{
			Item:   inventory[1].Id,
			UserId: user.Id,
			Price:  inventory[1].Price,
		}),
		repo.InsertPurchase(domain.Purchase{
			Item:   inventory[0].Id,
			UserId: user.Id,
			Price:  inventory[0].Price,
		}),
		repo.InsertPurchase(domain.Purchase{
			Item:   inventory[0].Id,
			UserId: user.Id,
			Price:  inventory[0].Price,
		}),
	}

	// Act
	ctx := context.Background()
	have, err := domain.Info(ctx, repo, user.Id)

	// Assert
	require.NoError(err)
	require.Equal(have.Coins, shoptest.DefaultBalance)

	require.Len(have.CoinHistory.Sent, 1)
	require.Equal(have.CoinHistory.Sent[0].Amount, int64(500))
	require.Equal(have.CoinHistory.Sent[0].ToUser, usersForeign[0].Username)

	require.Len(have.CoinHistory.Received, 1)
	require.Equal(have.CoinHistory.Received[0].Amount, int64(200))
	require.Equal(have.CoinHistory.Received[0].FromUser, usersForeign[1].Username)

	require.Len(have.Inventory, 2)
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
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()
	user := repo.InsertUser(domain.User{Username: "username"}, shoptest.DefaultBalance)

	usersForeign := []domain.User{
		repo.InsertUser(domain.User{Username: "foreign1"}, shoptest.DefaultBalance),
		repo.InsertUser(domain.User{Username: "foreign2"}, shoptest.DefaultBalance),
	}

	_ = []domain.Transfer{
		repo.InsertTransfer(domain.Transfer{
			FromUser: user.Id,
			ToUser:   usersForeign[0].Id,
			Delta:    500,
		}),
		repo.InsertTransfer(domain.Transfer{
			FromUser: usersForeign[1].Id,
			ToUser:   user.Id,
			Delta:    200,
		}),
		repo.InsertTransfer(domain.Transfer{
			FromUser: user.Id,
			ToUser:   usersForeign[0].Id,
			Delta:    100,
		}),
		repo.InsertTransfer(domain.Transfer{
			FromUser: usersForeign[1].Id,
			ToUser:   user.Id,
			Delta:    500,
		}),
	}

	// Act
	ctx := context.Background()
	have, err := domain.Info(ctx, repo, user.Id)

	// Assert
	require.NoError(err)
	require.Equal(have.Coins, shoptest.DefaultBalance)

	require.Len(have.CoinHistory.Sent, 2)
	require.Equal(have.CoinHistory.Sent[0].Amount, int64(100))
	require.Equal(have.CoinHistory.Sent[0].ToUser, usersForeign[0].Username)
	require.Equal(have.CoinHistory.Sent[1].Amount, int64(500))
	require.Equal(have.CoinHistory.Sent[1].ToUser, usersForeign[0].Username)

	require.Len(have.CoinHistory.Received, 2)
	require.Equal(have.CoinHistory.Received[0].Amount, int64(500))
	require.Equal(have.CoinHistory.Received[0].FromUser, usersForeign[1].Username)
	require.Equal(have.CoinHistory.Received[1].Amount, int64(200))
	require.Equal(have.CoinHistory.Received[1].FromUser, usersForeign[1].Username)
}
