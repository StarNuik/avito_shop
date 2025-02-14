package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendCoins_IncorrectUserId_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()

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

	repo := shoptest.NewInmemRepo()

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

	repo := shoptest.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom, shoptest.DefaultBalance)

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

	repo := shoptest.NewInmemRepo()
	user = repo.InsertUser(user, shoptest.DefaultBalance)

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

	repo := shoptest.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom, 50)
	userTo = repo.InsertUser(userTo, 50)

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

	repo := shoptest.NewInmemRepo()
	userFrom = repo.InsertUser(userFrom, shoptest.DefaultBalance)
	userTo = repo.InsertUser(userTo, shoptest.DefaultBalance)

	// Act
	transferSum := int64(100)
	ctx := context.Background()
	err := domain.SendCoins(ctx, repo, userFrom.Id, userTo.Username, transferSum)

	// Assert
	require.NoError(err)

	require.Equal(repo.Coins[userFrom.Id], shoptest.DefaultBalance-transferSum)
	require.Equal(repo.Coins[userTo.Id], shoptest.DefaultBalance+transferSum)

	require.Len(repo.Transfers, 1)
	require.Equal(repo.Transfers[0].FromUser, userFrom.Id)
	require.Equal(repo.Transfers[0].ToUser, userTo.Id)
	require.Equal(repo.Transfers[0].Delta, transferSum)
}

func TestSendCoins_MultipleSends_CorrectResult(t *testing.T) {
	// Arrange
	require := require.New(t)

	repo := shoptest.NewInmemRepo()
	users := []domain.User{
		repo.InsertUser(domain.User{Username: "user1"}, shoptest.DefaultBalance),
		repo.InsertUser(domain.User{Username: "user2"}, shoptest.DefaultBalance),
		repo.InsertUser(domain.User{Username: "user3"}, shoptest.DefaultBalance),
	}

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

	balance0, _ := repo.Coins[users[0].Id]
	require.Equal(balance0, int64(shoptest.DefaultBalance-2*50+100+500))

	balance1, _ := repo.Coins[users[1].Id]
	require.Equal(balance1, int64(shoptest.DefaultBalance-2*100+50+500))

	balance2, _ := repo.Coins[users[2].Id]
	require.Equal(balance2, int64(shoptest.DefaultBalance-2*500+50+100))
}
