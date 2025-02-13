package domain

import (
	"context"
	"github.com/avito_shop/internal/dto"
)

// Auth may return ErrNotFound
func Auth(ctx context.Context, repo ShopRepo, hash PasswordHash, req dto.AuthRequest) (dto.JwtPayload, error) {
	user, err := repo.User(ctx, req.Username)
	if err != nil {
		return dto.JwtPayload{}, err
	}

	if req.Username != user.Username {
		return dto.JwtPayload{}, ErrNotFound
	}

	currentHash, err := hash.Hash(req.Password)
	if err != nil {
		return dto.JwtPayload{}, err
	}

	if !hash.Same(user.PasswordHash, currentHash) {
		return dto.JwtPayload{}, ErrNotFound
	}

	return dto.JwtPayload{
		UserId: user.Id,
	}, nil
}
