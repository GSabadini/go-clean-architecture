package action

import (
	"bytes"
	"fmt"
	"github.com/gsabadini/go-bank-transfer/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/gsabadini/go-bank-transfer/usecase/stub"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestAccountStore(t *testing.T) {
	t.Parallel()

	type args struct {
		accountAction Account
		rawPayload    []byte
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Store action success",
			expectedStatusCode: http.StatusCreated,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubSuccess{}, loggerMock),
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": 10 
					}`,
				),
			},
		},
		{
			name:               "Store action error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{}, loggerMock),
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": 10 
					}`,
				),
			},
		},
		{
			name:               "Store action invalid JSON",
			expectedStatusCode: http.StatusBadRequest,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{}, loggerMock),
				rawPayload: []byte(
					`{
						"name": 
					}`,
				),
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

			r.HandleFunc("/accounts", tt.args.accountAction.Store).Methods(http.MethodPost)
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
	t.Parallel()

	type args struct {
		accountAction Account
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Index handler success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubSuccess{}, loggerMock),
			},
		},
		{
			name:               "Index handler error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{}, loggerMock),
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

func TestAccountFindBalance(t *testing.T) {
	t.Parallel()

	type args struct {
		accountAction Account
		accountID     string
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "FindBalance action success",
			expectedStatusCode: http.StatusOK,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubSuccess{}, loggerMock),
				accountID:     "5e5282beba39bfc244dc4c4b",
			},
		},
		{
			name:               "FindBalance action error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{}, loggerMock),
				accountID:     "5e5282beba39bfc244dc4c4b",
			},
		},
		{
			name:               "FindBalance action parameter invalid",
			expectedStatusCode: http.StatusBadRequest,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{}, loggerMock),
				accountID:     "1",
			},
		},
		{
			name:               "FindBalance action error fetching account",
			expectedStatusCode: http.StatusBadRequest,
			args: args{
				accountAction: NewAccount(mock.AccountUseCaseStubError{TypeErr: domain.ErrNotFound}, loggerMock),
				accountID:     "5e5282beba39bfc244dc4c4b",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/accounts/%s/balance", tt.args.accountID)
			req, err := http.NewRequest(http.MethodGet, uri, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"account_id": "5e5282beba39bfc244dc4c4b"})

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts/{account_id}/balance", tt.args.accountAction.FindBalance).Methods(http.MethodGet)
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
