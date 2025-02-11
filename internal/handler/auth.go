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
		ctx.AbortWithStatus(400)
		return nil, fmt.Errorf("incorrect AuthRequest json")
	}

	resp, err := domain.Auth(ctx, repo, req)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("incorrect username or password")
	} else if err != nil {
		log.LogError(err)
		ctx.AbortWithStatus(500)
		return nil, fmt.Errorf("internal server error")
	}

	return resp, nil
}

func Unauthorized(c *gin.Context, code int, message string) {
	resp := dto.ErrorResponse{
		Errors: message,
	}
	c.JSON(code, resp)
}
