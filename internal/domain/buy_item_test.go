package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuyItem_IncorrectUserId_ErrNotFound(t *testing.T) {
	panic("not implemented")
}

func TestBuyItem_ItemDoesntExist_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username:     "username",
		PasswordHash: "password",
	}

	repo := infra.NewInmemRepo()
	repo.InsertUser(user)
	user = repo.Users[0]

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
		Username:     "username",
		PasswordHash: "password",
	}
	item := domain.InventoryEntry{
		Name:  "buy-me",
		Price: 100,
	}

	repo := infra.NewInmemRepo()
	repo.InsertUser(user)
	repo.InsertInventory(item)
	user = repo.Users[0]
	item = repo.Inventory[0]

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
		Username:     "username",
		PasswordHash: "password",
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
	repo.InsertUser(user)
	repo.InsertInventory(item)
	repo.InsertBalanceOperation(balanceOp)
	user = repo.Users[0]
	item = repo.Inventory[0]
	balanceOp = repo.Operations[0]

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
