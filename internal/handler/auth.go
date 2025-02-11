package handler

import (
	"errors"
	"fmt"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
)

func Authenticator(ctx *gin.Context, repo domain.ShopRepo, log infra.Logger) (interface{}, error) {
	var req dto.AuthRequest
	err := ctx.BindJSON(&req)
	if err != nil || len(req.Username) == 0 || len(req.Password) == 0 {
		return nil, fmt.Errorf("%w: incorrect AuthRequest json", dto.ErrBadRequest)
	}

	resp, err := domain.Auth(ctx, repo, req)
	if errors.Is(err, dto.ErrInternalServer) {
		log.LogError(err)
		return nil, dto.ErrInternalServer
	} else if err != nil {
		return nil, err
	}
	return resp, nil
}
