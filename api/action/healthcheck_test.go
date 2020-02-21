package action

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	var (
		rr = httptest.NewRecorder()
		r  = mux.NewRouter()
	)

	r.HandleFunc("/healthcheck", HealthCheck).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
			status,
			http.StatusOK,
		)
	}
}
