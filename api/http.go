package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gsabadini/go-stone/api/action"
	"github.com/gsabadini/go-stone/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type HTTPServer struct {
	databaseConnection database.NoSQLDBHandler
}

func NewHTTPServer() HTTPServer {
	return HTTPServer{
		databaseConnection: createDatabaseConnection(getDatabaseURI()),
	}
}

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
	router.PathPrefix("/account").Handler(s.buildActionCreateAccount()).Methods(http.MethodPost)

	router.HandleFunc("/healthcheck", action.HealthCheck).Methods(http.MethodGet)
}

func (s HTTPServer) buildActionCreateAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccount(s.databaseConnection)

		accountAction.Create(res, req)
	}

	return negroni.New(negroni.Wrap(handler))
}

func createDatabaseConnection(uri string) *database.MongoHandler {
	handler, err := database.NewMongoHandler(uri)

	if err != nil {
		log.Println("Não foi possível realizar a conexão com o banco de dados")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	log.Println("Conexão com o banco de dados realizada com sucesso")

	return handler
}

func getDatabaseURI() string {
	if uri := os.Getenv("MONGO_DB_URI"); uri != "" {
		return uri
	}

	panic("Variável de ambiente 'MONGO_DB_URI' não foi definida")
}
