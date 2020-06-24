package web

import (
	"fmt"
	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/repository/mongodb"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/** TODO verificar middlewares */
type Gin struct {
	log  logger.Logger
	db   database.NoSQLHandler
	port Port
}

func NewGin(
	log logger.Logger,
	db database.NoSQLHandler,
	port Port,
) Gin {
	return Gin{
		log:  log,
		db:   db,
		port: port,
	}
}

func (g Gin) Listen() {
	var router = gin.New()
	gin.SetMode(gin.ReleaseMode)

	g.setAppHandlers(router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      router,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (g Gin) setAppHandlers(router *gin.Engine) {
	router.POST("/api/transfers", g.buildActionStoreTransfer())
	router.GET("/api/transfers", g.buildActionIndexTransfer())

	router.GET("/api/accounts/:account_id/balance", g.buildActionFindBalanceAccount())
	router.POST("/api/accounts", g.buildActionStoreAccount())
	router.GET("/api/accounts", g.buildActionIndexAccount())

	router.GET("/api/healthcheck", g.healthcheck())
}

func (g Gin) buildActionStoreTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferRepository = mongodb.NewTransferRepository(g.db)
			accountRepository  = mongodb.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
		)

		var transferAction = action.NewTransfer(transferUseCase, g.log)

		transferAction.Store(c.Writer, c.Request)
	}
}

func (g Gin) buildActionIndexTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferRepository = mongodb.NewTransferRepository(g.db)
			accountRepository  = mongodb.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
			transferAction     = action.NewTransfer(transferUseCase, g.log)
		)

		transferAction.Index(c.Writer, c.Request)
	}
}

func (g Gin) buildActionStoreAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		accountAction.Store(c.Writer, c.Request)
	}
}

func (g Gin) buildActionIndexAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		accountAction.Index(c.Writer, c.Request)
	}
}

func (g Gin) buildActionFindBalanceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		q := c.Request.URL.Query()
		q.Add("account_id", c.Param("account_id"))
		c.Request.URL.RawQuery = q.Encode()

		accountAction.FindBalance(c.Writer, c.Request)
	}
}

func (g Gin) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		action.HealthCheck(c.Writer, c.Request)
	}
}
