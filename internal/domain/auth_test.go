package domain_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuth_NoUser_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	hash := shoptest.NewNoopHash()
	repo := shoptest.NewInmemRepo()

	// Act
	ctx := context.Background()
	req := dto.AuthRequest{
		Username: "doesnt matter",
		Password: "doesnt matter",
	}
	_, err := domain.Auth(ctx, repo, hash, req)

	// Assert
	require.Error(err)
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestAuth_IncorrectPassword_ErrNotFound(t *testing.T) {
	// Arrange
	require := require.New(t)

	hash := shoptest.NewNoopHash()
	repo := shoptest.NewInmemRepo()
	user := repo.InsertUser(domain.User{
		Username:     "username",
		PasswordHash: "password",
	})

	// Act
	ctx := context.Background()
	req := dto.AuthRequest{
		Username: user.Username,
		Password: "incorrect password",
	}
	_, err := domain.Auth(ctx, repo, hash, req)

	// Assert
	require.ErrorIs(err, domain.ErrNotFound)
}

func TestAuth_UserExists_ReturnsUserId(t *testing.T) {
	// Arrange
	require := require.New(t)

	hash := shoptest.NewNoopHash()
	repo := shoptest.NewInmemRepo()
	user := repo.InsertUser(domain.User{
		Username:     "username",
		PasswordHash: "password",
	})

	// Act
	ctx := context.Background()
	req := dto.AuthRequest{
		Username: user.Username,
		Password: user.PasswordHash,
	}
	jwt, err := domain.Auth(ctx, repo, hash, req)

	// Assert
	require.NoError(err)
	require.Equal(jwt.UserId, user.Id)
}
