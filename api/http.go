package api

import (
	"log"
	"net/http"

	"github.com/gsabadini/go-stone/api/action"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type HTTPServer struct{}

func (s HTTPServer) Listen() {
	var (
		router = mux.NewRouter()
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

	router.HandleFunc("/health-check", action.HealthCheck).Methods(http.MethodGet)
}

func (s HTTPServer) buildActionCreateAccount() *negroni.Negroni {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var accountAction = action.NewAccountAction("connection database")

		accountAction.CreateAccount(res, req)
	}

	return negroni.New(negroni.Wrap(handler))
}
