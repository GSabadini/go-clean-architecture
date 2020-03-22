package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestValidateTransfer(t *testing.T) {
	tests := []struct {
		name               string
		rawPayload         []byte
		expectedStatusCode int
	}{
		{
			name: "Valid transfer",
			rawPayload: []byte(
				`{
					"account_destination_id": "5e5282beba39bfc244dc4c4b" ,
					"account_origin_id": "5e5282beba39bfc244dc4c4a",
					"amount": 1.00
				}`,
			),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid JSON",
			rawPayload: []byte(
				`{
					"account_destination_id": ,
					"account_origin_id": ,
					"amount": 1.00
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid amount",
			rawPayload: []byte(
				`{
					"account_destination_id": "5e5282beba39bfc244dc4c4b",
					"account_origin_id": "5e5282beba39bfc244dc4c4a",
					"amount": -1.00
				}`,
			),
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.rawPayload)

			req, err := http.NewRequest(http.MethodPost, "/", body)
			if err != nil {
				t.Fatal(err)
			}

			var loggerMock, _ = test.NewNullLogger()

			// transformando middleware em um http.Handler
			middlewareHandler := func(w http.ResponseWriter, r *http.Request) {
				next := func(w http.ResponseWriter, r *http.Request) {}
				middleware := NewValidateTransfer(loggerMock)
				middleware.Validate(w, r, next)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/", middlewareHandler).Methods(http.MethodPost)
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
