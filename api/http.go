package api

import (
	"log"
	"net/http"

	"github.com/gsabadini/go-stone/api/action"

	"github.com/gorilla/mux"
)

type HTTPServer struct{}

func (s HTTPServer) Listen() {
	var (
		router = mux.NewRouter()
	)

	router.HandleFunc("/health-check", action.HealthCheck)

	log.Printf("Iniciando servidor HTTP na porta %d", 3001)
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatalln("Erro ao iniciar API HTTP", err)
	}
}
