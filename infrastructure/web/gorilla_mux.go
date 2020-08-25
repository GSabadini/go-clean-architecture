package web

import (
	"fmt"
	postgres2 "github.com/gsabadini/go-bank-transfer/interface/repository/account/postgres"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/interface/api/middleware"
	"github.com/gsabadini/go-bank-transfer/interface/presenter"
	"github.com/gsabadini/go-bank-transfer/interface/repository"
	"github.com/gsabadini/go-bank-transfer/interface/repository/transfer/postgres"
	"github.com/urfave/negroni"
)

type gorillaMux struct {
	router     *mux.Router
	middleware *negroni.Negroni
	log        logger.Logger
	db         repository.SQLHandler
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newGorillaMux(
	log logger.Logger,
	db repository.SQLHandler,
	validator validator.Validator,
	port Port,
	t time.Duration,
) *gorillaMux {
	return &gorillaMux{
		router:     mux.NewRouter(),
		middleware: negroni.New(),
		log:        log,
		db:         db,
		validator:  validator,
		port:       port,
		ctxTimeout: t,
	}
}

//Listen inicia o servidor HTTP
func (g gorillaMux) Listen() {
	g.setAppHandlers(g.router)
	g.middleware.UseHandler(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.middleware,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (g gorillaMux) setAppHandlers(router *mux.Router) {
	api := router.PathPrefix("/v1").Subrouter()

	api.Handle("/transfers", g.buildActionStoreTransfer()).Methods(http.MethodPost)
	api.Handle("/transfers", g.buildActionIndexTransfer()).Methods(http.MethodGet)

	api.Handle("/accounts/{account_id}/balance", g.buildActionFindBalanceAccount()).Methods(http.MethodGet)
	api.Handle("/accounts", g.buildActionStoreAccount()).Methods(http.MethodPost)
	api.Handle("/accounts", g.buildActionFindAllAccount()).Methods(http.MethodGet)

	api.HandleFunc("/healthcheck", healtcheck.HealthCheck).Methods(http.MethodGet)
}

func (g gorillaMux) buildActionStoreTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferUseCase = transfer.NewTransfer(
				postgres.NewTransferRepository(g.db),
				postgres2.NewAccountRepository(g.db),
				presenter.NewTransferPresenter(),
				g.ctxTimeout,
			)

			transferAction = transfer.NewTransfer(transferUseCase, g.log, g.validator)
		)

		transferAction.Store(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildActionIndexTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferUseCase = transfer.NewTransfer(
				postgres.NewTransferRepository(g.db),
				postgres2.NewAccountRepository(g.db),
				presenter.NewTransferPresenter(),
				g.ctxTimeout,
			)
			transferAction = transfer.NewTransfer(transferUseCase, g.log, g.validator)
		)

		transferAction.FindAll(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildActionStoreAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountUseCase = account.NewAccount(
				postgres2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.Store(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildActionFindAllAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountUseCase = account.NewAccount(
				postgres2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		accountAction.FindAll(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g gorillaMux) buildActionFindBalanceAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountUseCase = account.NewAccount(
				postgres2.NewAccountRepository(g.db),
				presenter.NewAccountPresenter(),
				g.ctxTimeout,
			)
			accountAction = account.NewAccount(accountUseCase, g.log, g.validator)
		)

		var (
			vars = mux.Vars(req)
			q    = req.URL.Query()
		)

		q.Add("account_id", vars["account_id"])
		req.URL.RawQuery = q.Encode()

		accountAction.FindBalance(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}
