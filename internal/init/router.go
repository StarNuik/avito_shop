package init

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/handler"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/shoptest"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func jwtMiddleware(repo domain.ShopRepo, logger infra.Logger, hash domain.PasswordHash) *jwt.GinJWTMiddleware {
	nowUtc := func() time.Time {
		return time.Now().UTC()
	}
	authenticator := func(ctx *gin.Context) (interface{}, error) {
		return handler.Authenticator(ctx, repo, logger, hash)
	}

	params := jwt.GinJWTMiddleware{
		Realm: "avito-shop",
		// TODO
		Key:           []byte("const: secret-key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   "identity",
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      nowUtc,

		PayloadFunc:     handler.PackClaims,
		IdentityHandler: handler.UnpackClaims,
		Authenticator:   authenticator,
		//Authorizator:    nil,
		Unauthorized: handler.Unauthorized,
	}

	auth, err := jwt.New(&params)
	if err != nil {
		log.Panic(err)
	}
	return auth
}

func initJwtMiddleware(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Panic("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
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

// Router may panic if it couldn't initialize any of router's internal components
func Router() *gin.Engine {
	router := gin.Default()

	log := new(infra.FmtLogger)
	hash := shoptest.NewNoopHash()

	// TODO: change this to pg repo
	repo := shoptest.NewShopRepoBuilder()
	repo.AddStagingValues()

	auth := jwtMiddleware(repo, log, hash)
	router.Use(initJwtMiddleware(auth))

	addRoutes(router, auth, repo, log)

	return router
}
