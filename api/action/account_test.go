package action

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
)

func TestAccountCreate(t *testing.T) {
	type args struct {
		accountAction Account
		rawPayload    []byte
		httpMethod    string
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
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}),
				rawPayload:    []byte(`{"name": "test","cpf": "44451598087", "ballance": 10 }`),
				httpMethod:    http.MethodPost,
			},
		},
		{
			name:               "Create handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(database.MongoHandlerErrorMock{}),
				rawPayload:    []byte(``),
				httpMethod:    http.MethodPost,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := callHandler(t, tt.args.rawPayload, tt.args.accountAction, tt.args.httpMethod)

			if err != nil {
				t.Fatal(err)
			}

			if statusCode != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					statusCode,
					tt.expectedStatusCode,
				)
			}
		})
	}
}

//@TODO Corrigir testes de handler
//func TestAccountIndex(t *testing.T) {
//	type args struct {
//		accountAction Account
//		httpMethod    string
//	}
//
//	tests := []struct {
//		name               string
//		expectedStatusCode int
//		args               args
//	}{
//		{
//			name:               "Index handler success",
//			expectedStatusCode: http.StatusOK,
//			args: args{
//				accountAction: NewAccount(database.MongoHandlerSuccessMock{}),
//				httpMethod:    http.MethodGet,
//			},
//		},
//		{
//			name:               "Index handler error",
//			expectedStatusCode: http.StatusInternalServerError,
//			args: args{
//				accountAction: NewAccount(database.MongoHandlerErrorMock{}),
//				httpMethod:    http.MethodGet,
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			statusCode, err := callHandler(t, nil, tt.args.accountAction, tt.args.httpMethod)
//
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			if statusCode != tt.expectedStatusCode {
//				t.Errorf(
//					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
//					tt.name,
//					statusCode,
//					tt.expectedStatusCode,
//				)
//			}
//		})
//	}
//}

func callHandler(t *testing.T, rawPayload []byte, action Account, httpMethod string) (int, error) {
	var body io.Reader
	if httpMethod != http.MethodGet {
		body = bytes.NewReader(rawPayload)
	}

	req, err := http.NewRequest(httpMethod, "/account", body)
	if err != nil {
		t.Fatal(err)
	}

	var (
		rr = httptest.NewRecorder()
		r  = mux.NewRouter()
	)

	r.HandleFunc("/account", action.Create).Methods(httpMethod)
	r.ServeHTTP(rr, req)

	return rr.Code, nil
}
