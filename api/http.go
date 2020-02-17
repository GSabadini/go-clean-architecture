package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {}

func (s HTTPServer) Listen() {
	var (
		router         = mux.NewRouter()
	)

	router.HandleFunc("/", HealthCheck)

	log.Printf("Iniciando servidor HTTP na porta %d", 3001)
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatalln("Erro ao iniciar API HTTP", err)
	}
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Gorilla!\n"))
	w.WriteHeader(http.StatusOK)
}