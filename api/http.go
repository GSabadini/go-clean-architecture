package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/api/middleware"
	"github.com/gsabadini/go-bank-transfer/config"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

//HTTPServer armazena as dependências do servidor HTTP
type HTTPServer struct {
	appConfig          config.Config
	databaseConnection database.DbHandler
	log                *logrus.Logger
}

//NewHTTPServer constrói um HTTPServer com suas dependências
func NewHTTPServer(config config.Config) HTTPServer {
	return HTTPServer{
		appConfig:          config,
		databaseConnection: config.DatabaseConnection,
		log:                config.Logger,
	}
}

//Listen inicia o servidor HTTP
func (s HTTPServer) Listen() {
	var (
		router         = mux.NewRouter()
		negroniHandler = negroni.New()
		address        = fmt.Sprintf(":%d", s.appConfig.APIPort)
	)

	s.setAppHandlers(router)
	negroniHandler.UseHandler(router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         address,
		Handler:      negroniHandler,
	}

	s.log.Infoln("Starting HTTP server on the port", s.appConfig.APIPort)
	if err := server.ListenAndServe(); err != nil {
		s.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (s HTTPServer) setAppHandlers(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/transfers", s.buildActionStoreTransfer()).Methods(http.MethodPost)
	api.Handle("/transfers", s.buildActionIndexTransfer()).Methods(http.MethodGet)

	api.Handle("/accounts/{account_id}/balance", s.buildActionFindBalanceAccount()).Methods(http.MethodGet)
	api.Handle("/accounts", s.buildActionStoreAccount()).Methods(http.MethodPost)
	api.Handle("/accounts", s.buildActionIndexAccount()).Methods(http.MethodGet)

	api.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
}

func (s HTTPServer) buildActionStoreTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferRepository = repository.NewTransfer(s.databaseConnection)
			accountRepository  = repository.NewAccount(s.databaseConnection)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
		)

		var transferAction = action.NewTransfer(transferUseCase, s.log)

		transferAction.Store(res, req)
	}

	var (
		logging  = middleware.NewLogger(s.log).Execute
		validate = middleware.NewValidateTransfer(s.log).Execute
	)

	return negroni.New(
		negroni.HandlerFunc(logging),
		negroni.HandlerFunc(validate),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionIndexTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			transferRepository = repository.NewTransfer(s.databaseConnection)
			accountRepository  = repository.NewAccount(s.databaseConnection)
			transferUseCase    = usecase.NewTransfer(transferRepository, accountRepository)
		)

		var transferAction = action.NewTransfer(transferUseCase, s.log)

		transferAction.Index(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionStoreAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepostory = repository.NewAccount(s.databaseConnection)
			accountUseCase   = usecase.NewAccount(accountRepostory)
			accountAction    = action.NewAccount(accountUseCase, s.log)
		)

		accountAction.Store(res, req)
	}

	var (
		logging  = middleware.NewLogger(s.log).Execute
		validate = middleware.NewValidateAccount(s.log).Execute
	)

	return negroni.New(
		negroni.HandlerFunc(logging),
		negroni.HandlerFunc(validate),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionIndexAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepostory = repository.NewAccount(s.databaseConnection)
			accountUseCase   = usecase.NewAccount(accountRepostory)
			accountAction    = action.NewAccount(accountUseCase, s.log)
		)

		accountAction.Index(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionFindBalanceAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			accountRepostory = repository.NewAccount(s.databaseConnection)
			accountUseCase   = usecase.NewAccount(accountRepostory)
			accountAction    = action.NewAccount(accountUseCase, s.log)
		)

		accountAction.FindBalance(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}
