package domain

import (
	"context"
	"github.com/avito_shop/internal/dto"
)

// Auth *may* return a nil dto. That means the user isn't authorized
func Auth(ctx context.Context, repo ShopRepo, req dto.AuthRequest) (*dto.JwtPayload, error) {
	user, err := repo.User(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// TODO: password hashing
	if req.Username != user.Username || req.Password != user.PasswordHash {
		return nil, nil
	}

	return &dto.JwtPayload{
		UserId: user.Id,
	}, nil
}
