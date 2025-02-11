package handler

import (
	"fmt"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
)

func LoginBadRequest(ctx *gin.Context, chain gin.HandlerFunc) {
	var req dto.AuthRequest
	err := ctx.BindJSON(&req)
	if err != nil || len(req.Username) == 0 || len(req.Password) == 0 {
		ctx.AbortWithStatusJSON(400, dto.ErrorResponse{
			Errors: "incorrect AuthRequest json",
		})
		return
	}
	chain(ctx)
}

func Authenticator(ctx *gin.Context, repo domain.ShopRepo, log infra.Logger) (interface{}, error) {
	var req dto.AuthRequest
	_ = ctx.BindJSON(&req)

	resp, err := domain.Auth(ctx, repo, req)
	if err != nil {
		log.LogError(err)
		ctx.Status(500)
		return nil, fmt.Errorf("internal server error")
	}

	if resp == nil {
		return nil, fmt.Errorf("incorrect username or password")
	}

	return resp, nil
}

func Unauthorized(c *gin.Context, code int, message string) {
	resp := dto.ErrorResponse{
		Errors: message,
	}
	c.JSON(code, resp)
}
