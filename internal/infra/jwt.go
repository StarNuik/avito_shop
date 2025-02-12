package infra

import (
	"fmt"
	"github.com/avito_shop/internal/dto"
	"github.com/gin-gonic/gin"
)

const IdentityKey = "identity"

var (
	ErrJwtPayload error = fmt.Errorf("jwt payload error")
)

func JwtPayload(ctx *gin.Context) (*dto.JwtPayload, error) {
	any, ok := ctx.Get(IdentityKey)
	if !ok {
		return nil, ErrJwtPayload
	}

	dto, ok := any.(*dto.JwtPayload)
	if !ok {
		return nil, ErrJwtPayload
	}

	return dto, nil
}
