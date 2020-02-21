package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gsabadini/go-bank-transfer/api/action"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//HTTPServer armazena as dependências do servidor HTTP
type HTTPServer struct {
	databaseConnection database.NoSQLDBHandler
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

	s.setAppHandlers(router)
	negroniHandler.UseHandler(router)

	log.Printf("Iniciando servidor HTTP na porta %d", 3001)
	if err := http.ListenAndServe(":3001", negroniHandler); err != nil {
		log.Fatalln("Erro ao iniciar API HTTP", err)
	}
}

func (s HTTPServer) setAppHandlers(router *mux.Router) {
	router.PathPrefix("/accounts/{account_id}/ballance").Handler(s.buildActionShowAccount()).Methods(http.MethodGet)

	router.PathPrefix("/accounts").Handler(s.buildActionCreateAccount()).Methods(http.MethodPost)
	router.PathPrefix("/accounts").Handler(s.buildActionIndexAccount()).Methods(http.MethodGet)

	router.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
}

func (s HTTPServer) buildActionCreateAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection)

		accountAction.Create(res, req)
	}

	return negroni.New(negroni.Wrap(handler))
}

func (s HTTPServer) buildActionIndexAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection)

		accountAction.Index(res, req)
	}

	return negroni.New(negroni.Wrap(handler))
}

func (s HTTPServer) buildActionShowAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection)

		accountAction.Show(res, req)
	}

	return negroni.New(negroni.Wrap(handler))
}

func createDatabaseConnection(host, databaseName string) *database.MongoHandler {
	handler, err := database.NewMongoHandler(host, databaseName)

	if err != nil {
		log.Println("Não foi possível realizar a conexão com o banco de dados")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	log.Println("Conexão com o banco de dados realizada com sucesso")

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
