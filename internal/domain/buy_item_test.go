package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuyItem_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()

	// Act
	ctx := context.Background()
	err := domain.BuyItem(ctx, repo, -1, "doesnt exist")

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestBuyItem_ItemDoesntExist_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username: "username",
	}

	repo := shoptest.NewInmemRepo()
	user = repo.InsertUser(user, shoptest.DefaultBalance)

	// Act
	ctx := context.Background()
	err := domain.BuyItem(ctx, repo, user.Id, "doesnt exist")

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
	require.Empty(repo.Purchases)
}

func TestBuyItem_NotEnoughCoins_ErrNotEnough(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username: "username",
	}
	item := domain.InventoryItem{
		Name:  "buy-me",
		Price: shoptest.DefaultBalance + 100,
	}

	repo := shoptest.NewInmemRepo()

	user = repo.InsertUser(user, shoptest.DefaultBalance)
	item = repo.InsertInventory(item)

	// Act
	ctx := context.Background()
	err := domain.BuyItem(ctx, repo, user.Id, item.Name)

	// Assert
	require.ErrorIs(err, domain.ErrNotEnough)
	require.Empty(repo.Purchases)
}

func TestBuyItem_HappyPath_PurchaseAdded(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username: "username",
	}
	item := domain.InventoryItem{
		Name:  "buy-me",
		Price: 100,
	}

	repo := shoptest.NewInmemRepo()
	user = repo.InsertUser(user, shoptest.DefaultBalance)
	item = repo.InsertInventory(item)

	// Act
	ctx := context.Background()
	err := domain.BuyItem(ctx, repo, user.Id, item.Name)

	// Assert
	require.NoError(err)
	require.Equal(repo.Coins[user.Id], shoptest.DefaultBalance-item.Price)

	require.Len(repo.Purchases, 1)
	require.Equal(repo.Purchases[0].UserId, user.Id)
	require.Equal(repo.Purchases[0].Item, item.Id)
}
