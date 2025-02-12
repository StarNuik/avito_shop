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
		IdentityKey:   "identity",
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: func() time.Time {
			return time.Now().UTC()
		},

		PayloadFunc:     handler.PackClaims,
		IdentityHandler: handler.UnpackClaims,
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

func addRoutes(engine *gin.Engine, auth *jwt.GinJWTMiddleware, repo domain.ShopRepo, log infra.Logger) {
	engine.POST("/api/auth", auth.LoginHandler)

	authRequired := engine.Group("/api", auth.MiddlewareFunc())
	authRequired.GET("/info", func(ctx *gin.Context) {
		handler.Info(ctx, repo, log)
	})
	authRequired.GET("/buy/:item", func(ctx *gin.Context) {
		handler.BuyItem(ctx, repo, log)
	})
	authRequired.POST("/sendCoin", func(ctx *gin.Context) {
		handler.SendCoins(ctx, repo, log)
	})
}

func TestRepo() domain.ShopRepo {
	repo := infra.NewInmemRepo()
	// TODO: remove this
	u1 := repo.InsertUser(domain.User{
		Id:           -1,
		Username:     "admin",
		PasswordHash: "admin",
	})
	u2 := repo.InsertUser(domain.User{
		Id:           -2,
		Username:     "test",
		PasswordHash: "test",
	})
	repo.InsertBalanceOperation(domain.BalanceOperation{
		User:   u1.Id,
		Delta:  1000,
		Result: 1000,
	})
	repo.InsertBalanceOperation(domain.BalanceOperation{
		User:   u2.Id,
		Delta:  1000,
		Result: 1000,
	})
	repo.InsertInventory(domain.InventoryEntry{Name: "hoodie", Price: 100})
	repo.InsertInventory(domain.InventoryEntry{Name: "keychain", Price: 10})
	return repo
}

func Run() {
	engine := gin.Default()
	log := new(infra.FmtLogger)
	repo := TestRepo()

	auth := newJwt(repo, log)
	engine.Use(handlerMiddleware(auth))

	addRoutes(engine, auth, repo, log)

	_ = engine.Run()
}

func main() {
	Run()
}
