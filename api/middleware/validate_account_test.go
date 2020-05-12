package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestValidateAccount_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		rawPayload         []byte
		expectedStatusCode int
	}{
		{
			name: "Valid account",
			rawPayload: []byte(
				`{
					"name": "Test",
					"cpf": "44451598087",
					"balance": 10 
				}`,
			),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid JSON",
			rawPayload: []byte(`
				{
					"name":
					"cpf": 
					"balance": 
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid name",
			rawPayload: []byte(
				`{
					"name": "",
					"cpf": "44451598087",
					"balance": 10 
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid CPF",
			rawPayload: []byte(
				`{
					"name": "Test",
					"cpf": "0", 
					"balance": 10 
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid balance",
			rawPayload: []byte(
				`{
					"name": "Test",
					"cpf": "44451598087", 
					"balance": -10.00 
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.rawPayload)
			req, err := http.NewRequest(http.MethodPost, "/accounts", body)
			if err != nil {
				t.Fatal(err)
			}

			var loggerMock, _ = test.NewNullLogger()

			// transformando middleware em um http.Handler
			middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
				next := func(w http.ResponseWriter, r *http.Request) {}
				middleware := NewValidateAccount(loggerMock)
				middleware.Execute(w, r, next)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts", middlewareHandler).Methods(http.MethodPost)
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					rr.Code,
					tt.expectedStatusCode,
				)
			}
		})
	}
}
