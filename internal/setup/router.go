package setup

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/handler"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/repository"
	"github.com/avito_shop/internal/shoptest"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func jwtMiddleware(repo domain.ShopRepo, logger infra.Logger, hash domain.PasswordHasher) *jwt.GinJWTMiddleware {
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

func connectDb(env env) (domain.ShopRepo, func() error) {
	db, err := pgx.Connect(context.Background(), env.dbAddress)
	if err != nil {
		log.Panic(err)
	}

	repo := repository.NewShopPostgres(db)

	// TODO: remove this
	shoptest.AddStagingValues(db, new(infra.BcryptHasher))

	return repo, func() error { return db.Close(context.Background()) }
}

// Router may panic if it couldn't initialize any of router's internal components
func Router() *gin.Engine {
	env := GetEnv()

	// todo defer close()
	repo, _ := connectDb(env)

	router := gin.Default()

	log := new(infra.FmtLogger)
	hash := new(infra.BcryptHasher)

	auth := jwtMiddleware(repo, log, hash)
	router.Use(initJwtMiddleware(auth))

	addRoutes(router, auth, repo, log)

	return router
}
