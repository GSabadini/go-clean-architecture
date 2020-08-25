package action

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	var (
		rr      = httptest.NewRecorder()
		handler = http.NewServeMux()
	)

	handler.HandleFunc("/healthcheck", HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
			status,
			http.StatusOK,
		)
	}
}
