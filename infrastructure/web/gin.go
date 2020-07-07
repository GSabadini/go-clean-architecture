package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/repository/mongodb"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gin-gonic/gin"
)

type ginServer struct {
	router    *gin.Engine
	log       logger.Logger
	db        repository.NoSQLHandler
	validator validator.Validator
	port      Port
}

//NewGinServer constrói um ginServer com todas as suas dependências
func NewGinServer(
	log logger.Logger,
	db repository.NoSQLHandler,
	validator validator.Validator,
	port Port,
) ginServer {
	return ginServer{
		router:    gin.New(),
		log:       log,
		db:        db,
		validator: validator,
		port:      port,
	}
}

func (g ginServer) Listen() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()

	g.setAppHandlers(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.router,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Errors starting HTTP server")
	}
}

func (g ginServer) setAppHandlers(router *gin.Engine) {
	router.POST("/api/transfers", g.buildActionStoreTransfer())
	router.GET("/api/transfers", g.buildActionIndexTransfer())

	router.GET("/api/accounts/:account_id/balance", g.buildActionFindBalanceAccount())
	router.POST("/api/accounts", g.buildActionStoreAccount())
	router.GET("/api/accounts", g.buildActionIndexAccount())

	router.GET("/api/healthcheck", g.healthcheck())
}

func (g ginServer) buildActionStoreTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferRepository = mongodb.NewTransferRepository(g.db)
			accountRepository  = mongodb.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
		)

		var transferAction = action.NewTransfer(transferUseCase, g.log, g.validator)

		transferAction.Store(c.Writer, c.Request)
	}
}

func (g ginServer) buildActionIndexTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferRepository = mongodb.NewTransferRepository(g.db)
			accountRepository  = mongodb.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
			transferAction     = action.NewTransfer(transferUseCase, g.log, g.validator)
		)

		transferAction.Index(c.Writer, c.Request)
	}
}

func (g ginServer) buildActionStoreAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.Store(c.Writer, c.Request)
	}
}

func (g ginServer) buildActionIndexAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.Index(c.Writer, c.Request)
	}
}

func (g ginServer) buildActionFindBalanceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountRepository = mongodb.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log, g.validator)
		)

		q := c.Request.URL.Query()
		q.Add("account_id", c.Param("account_id"))
		c.Request.URL.RawQuery = q.Encode()

		accountAction.FindBalance(c.Writer, c.Request)
	}
}

func (g ginServer) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		action.HealthCheck(c.Writer, c.Request)
	}
}
