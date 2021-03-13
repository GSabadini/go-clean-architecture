package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsabadini/go-bank-transfer/adapter/api/action"
	"github.com/gsabadini/go-bank-transfer/adapter/logger"
	"github.com/gsabadini/go-bank-transfer/adapter/presenter"
	"github.com/gsabadini/go-bank-transfer/adapter/repository"
	"github.com/gsabadini/go-bank-transfer/adapter/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type ginEngine struct {
	router     *gin.Engine
	log        logger.Logger
	db         repository.NoSQL
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newGinServer(
	log logger.Logger,
	db repository.NoSQL,
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
		if err := server.ListenAndServe(); err != nil {
			g.log.WithError(err).Fatalln("Error starting HTTP server")
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		g.log.WithError(err).Fatalln("Server Shutdown Failed")
	}

	g.log.Infof("Service down")
}

/* TODO ADD MIDDLEWARE */
func (g ginEngine) setAppHandlers(router *gin.Engine) {
	router.POST("/v1/transfers", g.buildCreateTransferAction())
	router.GET("/v1/transfers", g.buildFindAllTransferAction())

	router.GET("/v1/accounts/:account_id/balance", g.buildFindBalanceAccountAction())
	router.POST("/v1/accounts", g.buildCreateAccountAction())
	router.GET("/v1/accounts", g.buildFindAllAccountAction())

	router.GET("/v1/health", g.healthcheck())
}

func (g ginEngine) buildCreateTransferAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateTransferInteractor(
				repository.NewTransferNoSQL(g.db),
				repository.NewAccountNoSQL(g.db),
				presenter.NewCreateTransferPresenter(),
				g.ctxTimeout,
			)

			act = action.NewCreateTransferAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindAllTransferAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindAllTransferInteractor(
				repository.NewTransferNoSQL(g.db),
				presenter.NewFindAllTransferPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllTransferAction(uc, g.log)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildCreateAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewCreateAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewCreateAccountAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindAllAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindAllAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewFindAllAccountPresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAllAccountAction(uc, g.log)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildFindBalanceAccountAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindBalanceAccountInteractor(
				repository.NewAccountNoSQL(g.db),
				presenter.NewFindAccountBalancePresenter(),
				g.ctxTimeout,
			)
			act = action.NewFindAccountBalanceAction(uc, g.log)
		)

		q := c.Request.URL.Query()
		q.Add("account_id", c.Param("account_id"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		action.HealthCheck(c.Writer, c.Request)
	}
}
