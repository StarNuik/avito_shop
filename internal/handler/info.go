package handler

import (
	"github.com/avito_shop/internal/dto"
	"github.com/gin-gonic/gin"
)

func Info(ctx *gin.Context) {
	ctx.JSON(500, dto.ErrorResponse{
		Errors: "not implemented",
	})
}
