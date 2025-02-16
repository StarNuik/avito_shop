package handler

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
)

func Authenticator(ctx *gin.Context, repo domain.ShopRepo, log infra.Logger, hash domain.PasswordHasher) (interface{}, error) {
	var req dto.AuthRequest
	err := ctx.BindJSON(&req)
	if err != nil || len(req.Username) == 0 || len(req.Password) == 0 {
		ctx.AbortWithStatus(400)
		return nil, fmt.Errorf("incorrect AuthRequest json")
	}

	resp, err := domain.Auth(ctx, repo, hash, req)
	if domain.IsDomainError(err) {
		return nil, fmt.Errorf("incorrect username or password")
	} else if err != nil {
		// TODO: trace id
		log.LogError(err)
		ctx.AbortWithStatus(500)
		return nil, fmt.Errorf("server error")
	}

	return &resp, nil
}

func PackClaims(in interface{}) jwt.MapClaims {
	payload, ok := in.(*dto.JwtPayload)
	if !ok {
		return jwt.MapClaims{}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return jwt.MapClaims{}
	}

	return jwt.MapClaims{
		infra.IdentityKey: payload.UserId,
		"_payload":        string(bytes),
	}
}

func UnpackClaims(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	obj, ok := claims["_payload"]
	if !ok {
		return nil
	}

	str, ok := obj.(string)
	if !ok {
		return nil
	}

	var payload dto.JwtPayload
	err := json.Unmarshal([]byte(str), &payload)
	if err != nil {
		return nil
	}

	return &payload
}

func Unauthorized(c *gin.Context, code int, message string) {
	resp := dto.ErrorResponse{
		Errors: message,
	}
	c.JSON(code, resp)
}
