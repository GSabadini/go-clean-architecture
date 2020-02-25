package api

import (
	"fmt"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/api/middleware"
	"github.com/gsabadini/go-bank-transfer/config"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

//HTTPServer armazena as dependências do servidor HTTP
type HTTPServer struct {
	appConfig          config.Config
	databaseConnection database.NoSQLDBHandler
	log                *logrus.Logger
}

//NewHTTPServer constrói um HTTPServer com suas dependências
func NewHTTPServer(config config.Config) HTTPServer {
	return HTTPServer{
		appConfig:          config,
		databaseConnection: createDatabaseConnection(config),
	}
}

//Listen inicia o servidor HTTP
func (s HTTPServer) Listen() {
	var (
		router         = mux.NewRouter()
		negroniHandler = negroni.New()
		address        = fmt.Sprintf(":%d", s.appConfig.ApiPort)
	)

	//@TODO REVER START DO LOGGER
	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	s.log = log

	s.setAppHandlers(router)
	negroniHandler.UseHandler(router)

	s.log.Infoln("Starting HTTP server on the port", s.appConfig.ApiPort)
	if err := http.ListenAndServe(address, negroniHandler); err != nil {
		s.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (s HTTPServer) setAppHandlers(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/transfers", s.buildActionStoreTransfer()).Methods(http.MethodPost)

	api.Handle("/accounts/{account_id}/balance", s.buildActionShowBalanceAccount()).Methods(http.MethodGet)
	api.Handle("/accounts", s.buildActionStoreAccount()).Methods(http.MethodPost)
	api.Handle("/accounts", s.buildActionIndexAccount()).Methods(http.MethodGet)

	api.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
}

func (s HTTPServer) buildActionStoreTransfer() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewTransfer(s.databaseConnection, s.log)

		accountAction.Store(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Logging),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionStoreAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection, s.log)

		accountAction.Store(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Logging),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionIndexAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection, s.log)

		accountAction.Index(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Logging),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func (s HTTPServer) buildActionShowBalanceAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection, s.log)

		accountAction.ShowBalance(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Logging),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func createDatabaseConnection(config config.Config) *database.MongoHandler {
	handler, err := database.NewMongoHandler(config.DatabaseHost, config.DatabaseName)
	if err != nil {
		logrus.Infoln("Could not make a connection to the database")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	logrus.Infoln("Successfully connected to the database")

	return handler
}
