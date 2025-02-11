package domain

import (
	"context"
	"fmt"
	"github.com/avito_shop/internal/dto"
)

func Auth(ctx context.Context, repo ShopRepo, req dto.AuthRequest) (*dto.JwtPayload, error) {
	user, err := repo.User(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", dto.ErrInternalServer, err)
	}

	// TODO: password hashing
	if req.Username == user.Username && req.Password == user.PasswordHash {
		return &dto.JwtPayload{
			UserId: user.Id,
		}, nil
	}

	return nil, dto.ErrUnauthorized
}
