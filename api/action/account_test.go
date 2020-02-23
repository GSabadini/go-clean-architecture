package action

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
)

func TestAccountCreate(t *testing.T) {
	type args struct {
		accountAction AccountAction
		rawPayload    []byte
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Create handler success",
			expectedStatusCode: http.StatusCreated,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerSuccessMock{}),
				rawPayload:    []byte(`{"name": "test","cpf": "44451598087", "ballance": 10 }`),
			},
		},
		{
			name:               "Create handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerErrorMock{}),
				rawPayload:    []byte(``),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.args.rawPayload)
			req, err := http.NewRequest(http.MethodPost, "/accounts", body)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts", tt.args.accountAction.Create).Methods(http.MethodPost)
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

func TestAccountIndex(t *testing.T) {
	type args struct {
		accountAction AccountAction
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Index handler success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerSuccessMock{}),
			},
		},
		{
			name:               "Index handler error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerErrorMock{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts", tt.args.accountAction.Index).Methods(http.MethodGet)
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

//@TODO Corrigir testes de handler
func TestAccountShow(t *testing.T) {
	type args struct {
		accountAction AccountAction
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Show handler success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerSuccessMock{}),
			},
		},
		{
			name:               "Show handler error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccountAction(database.MongoHandlerErrorMock{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/accounts/5e5282beba39bfc244dc4c4b/ballance", nil)
			req = mux.SetURLVars(req, map[string]string{
				"param1": "param1",
				"param2": "param2",
			})
			if err != nil {
				t.Fatal(err)
			}

			ctx := req.Context()
			ctx = context.WithValue(ctx, "request-uuid", "myvalue")
			req = req.WithContext(ctx)
			req = mux.SetURLVars(req, map[string]string{"account_id": "5e5282beba39bfc244dc4c4b"})

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts/5e5282beba39bfc244dc4c4b/ballance", tt.args.accountAction.Show).Methods(http.MethodGet)
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

//func callHandler(t *testing.T, rawPayload []byte, action Account, httpMethod string) (int, error) {
//	var body io.Reader
//	if httpMethod != http.MethodGet {
//		body = bytes.NewReader(rawPayload)
//	}
//
//	req, err := http.NewRequest(httpMethod, "/account", body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var (
//		rr = httptest.NewRecorder()
//		r  = mux.NewRouter()
//	)
//
//	r.HandleFunc("/account", action.Create).Methods(httpMethod)
//	r.ServeHTTP(rr, req)
//
//	return rr.Code, nil
//}
