package main

import (
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/handler"
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func newJwt(repo domain.ShopRepo, log infra.Logger) *jwt.GinJWTMiddleware {
	params := jwt.GinJWTMiddleware{
		Realm:         "const: realm",
		Key:           []byte("const: secret-key"),
		Timeout:       0,
		MaxRefresh:    0,
		IdentityKey:   "const: identity-key",
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: func() time.Time {
			return time.Now().UTC()
		},

		//IdentityHandler: nil,
		//PayloadFunc:     nil,
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			return handler.Authenticator(ctx, repo, log)
		},
		//Authorizator:    nil,
		Unauthorized: handler.Unauthorized,
	}
	auth, err := jwt.New(&params)
	if err != nil {
		panic(err)
	}
	return auth
}

func handlerMiddleware(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func addRoutes(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	engine.POST("/api/auth", auth.LoginHandler)

	authRequired := engine.Group("/api", auth.MiddlewareFunc())
	authRequired.GET("/api/info")
	authRequired.GET("/api/buy/{item}")
	authRequired.POST("/api/sendCoin")
}

func Run() {
	engine := gin.Default()
	repo := infra.NewMockRepo()
	log := new(infra.FmtLogger)

	auth := newJwt(repo, log)
	engine.Use(handlerMiddleware(auth))

	addRoutes(engine, auth)

	_ = engine.Run()
}

func main() {
	Run()
}
