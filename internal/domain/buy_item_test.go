package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuyItem_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()

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

	repo := infra.NewInmemRepo()
	user = repo.InsertUser(user)

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
	item := domain.InventoryEntry{
		Name:  "buy-me",
		Price: 100,
	}

	repo := infra.NewInmemRepo()

	user = repo.InsertUser(user)
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
	item := domain.InventoryEntry{
		Name:  "buy-me",
		Price: 100,
	}
	balanceOp := domain.BalanceOperation{
		User:   user.Id,
		Delta:  1000,
		Result: 1000,
	}

	repo := infra.NewInmemRepo()
	user = repo.InsertUser(user)
	item = repo.InsertInventory(item)
	balanceOp = repo.InsertBalanceOperation(balanceOp)

	// Act
	ctx := context.Background()
	err := domain.BuyItem(ctx, repo, user.Id, item.Name)

	// Assert
	require.NoError(err)
	require.Len(repo.Operations, 2) // +1 for the balanceOp
	require.Equal(repo.Operations[1].Delta, -item.Price)
	require.Equal(repo.Operations[1].Result, balanceOp.Result-item.Price)
	require.Len(repo.Purchases, 1)
	require.Equal(repo.Purchases[0].User, user.Id)
	require.Equal(repo.Purchases[0].Item, item.Id)
	require.Equal(repo.Purchases[0].Operation, repo.Operations[1].Id)
}
