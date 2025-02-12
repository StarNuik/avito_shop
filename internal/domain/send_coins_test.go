package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendCoins_IncorrectUserId_ErrNotFound(t *testing.T) {
	panic("not implemented")
}

func TestSendCoins_TransferSumLteZero_ErrNotAllowed(t *testing.T) {
	// Arrange
	require := require.New(t)

	// Act
	transferSum := int64(-100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, nil, -1, -1, transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotAllowed)
}

func TestSendCoins_TargetDoesntExist_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username:     "username",
		PasswordHash: "password",
	}

	repo := infra.NewInmemRepo()
	repo.InsertUser(userFrom)
	userFrom = repo.Users[0]

	// Act
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, -1, 100)

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestSendCoins_LowBalance_ErrNotEnough(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username:     "username1",
		PasswordHash: "password",
	}
	userTo := domain.User{
		Username:     "username2",
		PasswordHash: "password",
	}
	balanceOp := domain.BalanceOperation{
		User:   userFrom.Id,
		Delta:  50,
		Result: 50,
	}

	repo := infra.NewInmemRepo()
	repo.InsertUser(userFrom)
	repo.InsertUser(userTo)
	repo.InsertBalanceOperation(balanceOp)
	userFrom = repo.Users[0]
	userTo = repo.Users[1]
	balanceOp = repo.Operations[0]

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, userTo.Id, transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotEnough)
}

func TestSendCoins_HappyPath_TransferAdded(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username:     "username1",
		PasswordHash: "password",
	}
	userTo := domain.User{
		Username:     "username2",
		PasswordHash: "password",
	}
	balanceOp := domain.BalanceOperation{
		User:   userFrom.Id,
		Delta:  1000,
		Result: 1000,
	}

	repo := infra.NewInmemRepo()
	repo.InsertUser(userFrom)
	repo.InsertUser(userTo)
	repo.InsertBalanceOperation(balanceOp)
	userFrom = repo.Users[0]
	userTo = repo.Users[1]
	balanceOp = repo.Operations[0]

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, userTo.Id, transferSum)

	// Assert
	require.NoError(err)
	require.Len(repo.Operations, 3) // +1 for the balanceOp
	require.Len(repo.Transfers, 1)

	// require.Contains didn't work :(
	srcOp := repo.Operations[repo.Transfers[0].SourceOp]
	require.Equal(srcOp.User, userFrom.Id)
	require.Equal(srcOp.Delta, -transferSum)
	require.Equal(srcOp.Result, balanceOp.Result-transferSum)

	destOp := repo.Operations[repo.Transfers[0].TargetOp]
	require.Equal(destOp.User, userTo.Id)
	require.Equal(destOp.Delta, transferSum)
	require.Equal(destOp.Result, transferSum)
}

func TestSendCoins_MultipleSends_CorrectResult(t *testing.T) {
	panic("not implemented")
}
