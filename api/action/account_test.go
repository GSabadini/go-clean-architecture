package action

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
)

/* TODO remover dependecia do gorilla mux */
func TestAccount_Store(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		accountAction      Account
		args               args
	}{
		{
			name:               "Store action success",
			expectedStatusCode: http.StatusCreated,
			accountAction:      NewAccount(usecase.AccountUseCaseStubSuccess{}, logger.LoggerMock{}, validator),
			args: args{
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
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
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
			name:          "Store action invalid balance",
			accountAction: NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": -1 
					}`,
				),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:          "Store action invalid name fields",
			accountAction: NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
				rawPayload: []byte(
					`{
						"name123": "test",
						"cpf1231": "44451598087", 
						"balance12312": 1 
					}`,
				),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Store action invalid JSON",
			expectedStatusCode: http.StatusBadRequest,
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
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

			r.HandleFunc("/accounts", tt.accountAction.Store).Methods(http.MethodPost)
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

func TestAccount_Index(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	tests := []struct {
		name               string
		expectedStatusCode int
		accountAction      Account
	}{
		{
			name:               "Index handler success",
			expectedStatusCode: http.StatusOK,
			accountAction:      NewAccount(usecase.AccountUseCaseStubSuccess{}, logger.LoggerMock{}, validator),
		},
		{
			name:               "Index handler error",
			expectedStatusCode: http.StatusInternalServerError,
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
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

			r.HandleFunc("/accounts", tt.accountAction.Index).Methods(http.MethodGet)
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

func TestAccount_FindBalance(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	type args struct {
		accountID string
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		accountAction      Account
		args               args
	}{
		{
			name:               "FindBalance action success",
			expectedStatusCode: http.StatusOK,
			accountAction:      NewAccount(usecase.AccountUseCaseStubSuccess{}, logger.LoggerMock{}, validator),
			args: args{
				accountID: "59e09306b5174ba2986a7ce36aa2afd9",
			},
		},
		{
			name:               "FindBalance action error",
			expectedStatusCode: http.StatusInternalServerError,
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
		},
		{
			name:               "FindBalance action parameter invalid",
			expectedStatusCode: http.StatusBadRequest,
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{}, logger.LoggerMock{}, validator),
			args: args{
				accountID: "error",
			},
		},
		{
			name:               "FindBalance action error fetching account",
			expectedStatusCode: http.StatusBadRequest,
			accountAction:      NewAccount(usecase.AccountUseCaseStubError{TypeErr: domain.ErrNotFound}, logger.LoggerMock{}, validator),
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
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

			q := req.URL.Query()
			q.Add("account_id", tt.args.accountID)
			req.URL.RawQuery = q.Encode()

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/accounts/{account_id}/balance", tt.accountAction.FindBalance).Methods(http.MethodGet)
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

func TestAccount_Find(t *testing.T) {
	r := domain.IsValidUUID("59e09306b5174ba2986a7ce36aa2afd9")
	t.Log(r)
}
