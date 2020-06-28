package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/api/middleware"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/repository/postgres"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type GorillaMux struct {
	log        logger.Logger
	db         database.SQLHandler
	validation validator.Validator
	port       Port
}

func NewGorillaMux(
	log logger.Logger,
	dbConnSQL database.SQLHandler,
	validation validator.Validator,
	port Port,
) GorillaMux {
	return GorillaMux{
		log:        log,
		db:         dbConnSQL,
		validation: validation,
		port:       port,
	}
}

//Listen inicia o servidor HTTP
func (g GorillaMux) Listen() {
	var (
		router         = mux.NewRouter()
		negroniHandler = negroni.New()
	)

	g.setAppHandlers(router)
	negroniHandler.UseHandler(router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      negroniHandler,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (g GorillaMux) setAppHandlers(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/transfers", g.buildActionStoreTransfer()).Methods(http.MethodPost)
	api.Handle("/transfers", g.buildActionIndexTransfer()).Methods(http.MethodGet)

	api.Handle("/accounts/{account_id}/balance", g.buildActionFindBalanceAccount()).Methods(http.MethodGet)
	api.Handle("/accounts", g.buildActionStoreAccount()).Methods(http.MethodPost)
	api.Handle("/accounts", g.buildActionIndexAccount()).Methods(http.MethodGet)

	api.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
}

func (g GorillaMux) buildActionStoreTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferRepository = postgres.NewTransferRepository(g.db)
			accountRepository  = postgres.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
		)

		var transferAction = action.NewTransfer(transferUseCase, g.log, g.validation)

		transferAction.Store(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g GorillaMux) buildActionIndexTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferRepository = postgres.NewTransferRepository(g.db)
			accountRepository  = postgres.NewAccountRepository(g.db)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
			transferAction     = action.NewTransfer(transferUseCase, g.log, g.validation)
		)

		transferAction.Index(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g GorillaMux) buildActionStoreAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepository = postgres.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		accountAction.Store(res, req)
	}

	var (
		logging  = middleware.NewLogger(g.log).Execute
		validate = middleware.NewValidateAccount(g.log).Execute
	)

	return negroni.New(
		negroni.HandlerFunc(logging),
		negroni.HandlerFunc(validate),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g GorillaMux) buildActionIndexAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepository = postgres.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		accountAction.Index(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (g GorillaMux) buildActionFindBalanceAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepository = postgres.NewAccountRepository(g.db)
			accountUseCase    = usecase.NewAccount(accountRepository)
			accountAction     = action.NewAccount(accountUseCase, g.log)
		)

		var (
			vars = mux.Vars(req)
			q    = req.URL.Query()
		)

		accountID, _ := vars["account_id"]
		q.Add("account_id", accountID)
		req.URL.RawQuery = q.Encode()

		accountAction.FindBalance(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(g.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}
