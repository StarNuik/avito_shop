package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendCoins_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, -1, "", transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestSendCoins_TransferSumLteZero_ErrNotAllowed(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()

	// Act
	transferSum := int64(-100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, -1, "", transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotAllowed)
}

func TestSendCoins_TargetDoesntExist_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username: "username",
	}

	repo := infra.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom)

	// Act
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, "", 100)

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestSendCoins_TargetIsUser_ErrNotAllowed(t *testing.T) {
	// Arrange
	require := require.New(t)

	user := domain.User{
		Username: "username",
	}
	balanceOp := domain.BalanceOperation{
		User:   user.Id,
		Delta:  500,
		Result: 500,
	}

	repo := infra.NewInmemRepo()
	user = repo.InsertUser(user)
	balanceOp = repo.InsertBalanceOperation(balanceOp)

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, user.Id, user.Username, transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotAllowed)
}

func TestSendCoins_LowBalance_ErrNotEnough(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username: "username1",
	}
	userTo := domain.User{
		Username: "username2",
	}
	balanceOp := domain.BalanceOperation{
		User:   userFrom.Id,
		Delta:  50,
		Result: 50,
	}

	repo := infra.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom)
	userTo = repo.InsertUser(userTo)
	balanceOp = repo.InsertBalanceOperation(balanceOp)

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, userTo.Username, transferSum)

	// Assert
	require.ErrorIs(err, domain.ErrNotEnough)
}

func TestSendCoins_HappyPath_TransferAdded(t *testing.T) {
	// Arrange
	require := require.New(t)

	userFrom := domain.User{
		Username: "username1",
	}
	userTo := domain.User{
		Username: "username2",
	}
	balanceOp := domain.BalanceOperation{
		User:   userFrom.Id,
		Delta:  1000,
		Result: 1000,
	}

	repo := infra.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom)
	userTo = repo.InsertUser(userTo)
	balanceOp = repo.InsertBalanceOperation(balanceOp)

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, userTo.Username, transferSum)

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
	// Arrange
	require := require.New(t)

	repo := infra.NewInmemRepo()
	users := []domain.User{
		repo.InsertUser(domain.User{Username: "user1"}),
		repo.InsertUser(domain.User{Username: "user2"}),
		repo.InsertUser(domain.User{Username: "user3"}),
	}
	repo.InsertBalanceOperation(domain.BalanceOperation{
		User:   users[0].Id,
		Delta:  1000,
		Result: 1000,
	})
	repo.InsertBalanceOperation(domain.BalanceOperation{
		User:   users[1].Id,
		Delta:  1000,
		Result: 1000,
	})
	repo.InsertBalanceOperation(domain.BalanceOperation{
		User:   users[2].Id,
		Delta:  1000,
		Result: 1000,
	})

	// Act
	ctx := context.Background()

	_ = domain.SendCoins(ctx, repo, users[0].Id, users[1].Username, 50)
	_ = domain.SendCoins(ctx, repo, users[0].Id, users[2].Username, 50)
	_ = domain.SendCoins(ctx, repo, users[1].Id, users[0].Username, 100)
	_ = domain.SendCoins(ctx, repo, users[1].Id, users[2].Username, 100)
	_ = domain.SendCoins(ctx, repo, users[2].Id, users[0].Username, 500)
	_ = domain.SendCoins(ctx, repo, users[2].Id, users[1].Username, 500)

	// Assert
	require.Len(repo.Transfers, 6)

	balance0, _ := repo.UserBalance(ctx, users[0].Id)
	require.Equal(balance0, int64(1000-2*50+100+500))

	balance1, _ := repo.UserBalance(ctx, users[1].Id)
	require.Equal(balance1, int64(1000-2*100+50+500))

	balance2, _ := repo.UserBalance(ctx, users[2].Id)
	require.Equal(balance2, int64(1000-2*500+50+100))
}
