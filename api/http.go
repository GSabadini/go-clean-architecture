package api

import (
	"net/http"
	"os"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/api/middleware"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

//HTTPServer armazena as dependências do servidor HTTP
type HTTPServer struct {
	databaseConnection database.NoSQLDBHandler
	log                *logrus.Logger
}

//NewHTTPServer constrói um HTTPServer as suas dependências
func NewHTTPServer() HTTPServer {
	return HTTPServer{
		databaseConnection: createDatabaseConnection(getDatabaseHost(), getDatabaseName()),
	}
}

//Listen inicia o servidor HTTP
func (s HTTPServer) Listen() {
	var (
		router         = mux.NewRouter()
		negroniHandler = negroni.New()
	)

	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	s.log = log

	s.setAppHandlers(router)
	negroniHandler.UseHandler(router)

	log.Infoln("Iniciando servidor HTTP na porta", 3001)
	if err := http.ListenAndServe(":3001", negroniHandler); err != nil {
		log.WithError(err).Fatalln("Erro ao iniciar servidor HTTP")
	}
}

func (s HTTPServer) setAppHandlers(router *mux.Router) {
	router.PathPrefix("/accounts/{account_id}/ballance").Handler(s.buildActionShowBallanceAccount()).Methods(http.MethodGet)

	router.PathPrefix("/accounts").Handler(s.buildActionStoreAccount()).Methods(http.MethodPost)
	router.PathPrefix("/accounts").Handler(s.buildActionIndexAccount()).Methods(http.MethodGet)

	router.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
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

func (s HTTPServer) buildActionShowBallanceAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection, s.log)

		accountAction.ShowBallance(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewLogger(s.log).Logging),
		negroni.NewRecovery(),
		negroni.Wrap(handler),
	)
}

func createDatabaseConnection(host, databaseName string) *database.MongoHandler {
	handler, err := database.NewMongoHandler(host, databaseName)

	if err != nil {
		logrus.Infoln("Não foi possível realizar a conexão com o banco de dados")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	logrus.Infoln("Conexão com o banco de dados realizada com sucesso")

	return handler
}

func getDatabaseHost() string {
	if uri := os.Getenv("MONGODB_HOST"); uri != "" {
		return uri
	}

	panic("Variável de ambiente 'MONGODB_HOST' não foi definida")
}

func getDatabaseName() string {
	if uri := os.Getenv("MONGODB_DATABASE"); uri != "" {
		return uri
	}

	panic("Variável de ambiente 'MONGODB_DATABASE' não foi definida")
}
