package handler

import (
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
)

func Info(ctx *gin.Context, repo domain.ShopRepo, log infra.Logger) {
	jwt, err := infra.JwtPayload(ctx)
	if err != nil {
		ctx.JSON(401, dto.ErrorResponse{"incorrect jwt"})
		return
	}

	out, err := domain.Info(ctx, repo, jwt.UserId)
	if domain.IsDomainError(err) {
		ctx.JSON(400, dto.ErrorResponse{err.Error()})
		return
	} else if err != nil {
		// TODO: trace id
		log.LogError(err)
		ctx.JSON(500, dto.ErrorResponse{"server error"})
		return
	}

	ctx.JSON(200, out)
}
