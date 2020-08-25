package web

import (
	"fmt"
	mongodb2 "github.com/gsabadini/go-bank-transfer/interface/repository/account/mongodb"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/interface/presenter"
	"github.com/gsabadini/go-bank-transfer/interface/repository"
	"github.com/gsabadini/go-bank-transfer/interface/repository/transfer/mongodb"
)

type ginEngine struct {
	router     *gin.Engine
	log        logger.Logger
	db         repository.NoSQLHandler
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newGinServer(
	log logger.Logger,
	db repository.NoSQLHandler,
	validator validator.Validator,
	port Port,
	t time.Duration,
) *ginEngine {
	return &ginEngine{
		router:     gin.New(),
		log:        log,
		db:         db,
		validator:  validator,
		port:       port,
		ctxTimeout: t,
	}
}

func (g ginEngine) Listen() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()

	g.setAppHandlers(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.router,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (g ginEngine) setAppHandlers(router *gin.Engine) {
	router.POST("/v1/transfers", g.buildActionStoreTransfer())
	router.GET("/v1/transfers", g.buildActionFindAllTransfer())

	router.GET("/v1/accounts/:account_id/balance", g.buildActionFindBalanceAccount())
	router.POST("/v1/accounts", g.buildActionStoreAccount())
	router.GET("/v1/accounts", g.buildActionFindAllAccount())

	router.GET("/v1/healthcheck", g.healthcheck())
}

func (g ginEngine) buildActionStoreTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferUseCase = transfer.NewTransfer(
				mongodb.NewTransferRepository(g.db),
				mongodb2.NewAccountRepository(g.db),
				presenter.NewTransferPresenter(),
				g.ctxTimeout,
			)

			transferAction = transfer.NewTransfer(transferUseCase, g.log, g.validator)
		)

		transferAction.Store(c.Writer, c.Request)
	}
}

func (g ginEngine) buildActionFindAllTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			transferUseCase = transfer.NewTransfer(
				mongodb.NewTransferRepository(g.db),
				mongodb2.NewAccountRepository(g.db),
				presenter.NewTransferPresenter(),
				g.ctxTimeout,
			)
			transferAction = transfer.NewTransfer(transferUseCase, g.log, g.validator)
		)

		transferAction.FindAll(c.Writer, c.Request)
	}
}

func (g ginEngine) buildActionStoreAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountUseCase = account.NewAccount(
				mongodb2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.Store(c.Writer, c.Request)
	}
}

func (g ginEngine) buildActionFindAllAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountUseCase = account.NewAccount(
				mongodb2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.FindAll(c.Writer, c.Request)
	}
}

func (g ginEngine) buildActionFindBalanceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accountUseCase = account.NewAccount(
				mongodb2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		q := c.Request.URL.Query()
		q.Add("account_id", c.Param("account_id"))
		c.Request.URL.RawQuery = q.Encode()

		accountAction.FindBalance(c.Writer, c.Request)
	}
}

func (g ginEngine) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		healtcheck.HealthCheck(c.Writer, c.Request)
	}
}
