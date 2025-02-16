package setup

import (
	"context"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/handler"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

var (
	jwtSecret = "jG>ke8$%t[I|%-Sa5+O*3+];3ZX4}_WIAeld+(NWA2NPM~U*4t-3mWIRg>CEd'"
)

func jwtMiddleware(repo domain.ShopRepo, logger infra.Logger, hash domain.PasswordHasher) *jwt.GinJWTMiddleware {
	nowUtc := func() time.Time {
		return time.Now().UTC()
	}
	authenticator := func(ctx *gin.Context) (interface{}, error) {
		return handler.Authenticator(ctx, repo, logger, hash)
	}

	params := jwt.GinJWTMiddleware{
		Realm:         "avito-shop",
		Key:           []byte(jwtSecret),
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

type engine struct {
	*gin.Engine
	db   *pgx.Conn
	port int
}

// Router may panic if it couldn't initialize any of router's internal components
func Router() *engine {
	env := GetEnv()

	db, err := pgx.Connect(context.Background(), env.DatabaseUrl)
	if err != nil {
		log.Panic(err)
	}

	repo := repository.NewShopPostgres(db)

	router := gin.Default()

	log := new(infra.FmtLogger)
	hash := new(infra.BcryptHasher)

	auth := jwtMiddleware(repo, log, hash)
	router.Use(initJwtMiddleware(auth))

	addRoutes(router, auth, repo, log)

	return &engine{
		Engine: router,
		db:     db,
		port:   env.ServerPort,
	}
}

func (r *engine) Run() error {
	addr := fmt.Sprintf(":%d", r.port)
	return r.Engine.Run(addr)
}

func (r *engine) Close() error {
	return r.db.Close(context.Background())
}
