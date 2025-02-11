package domain

import (
	"context"
	"github.com/avito_shop/internal/dto"
)

// Auth returns `ErrNotFound` if it couldn't authorize a user
func Auth(ctx context.Context, repo ShopRepo, req dto.AuthRequest) (*dto.JwtPayload, error) {
	user, err := repo.User(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// TODO: password hashing
	if req.Username != user.Username || req.Password != user.PasswordHash {
		return nil, ErrNotFound
	}

	return &dto.JwtPayload{
		UserId: user.Id,
	}, nil
}
