package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type mockFindBalanceAccount struct {
	result output.AccountBalance
	err    error
}

func (m mockFindBalanceAccount) Execute(_ context.Context, _ domain.AccountID) (output.AccountBalance, error) {
	return m.result, m.err
}

func TestFindBalanceAccountAction_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		accountID string
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.FindBalanceAccount
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "FindBalanceAccountAction success",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: output.AccountBalance{
					Balance: 10,
				},
				err: nil,
			},
			expectedBody:       []byte(`{"balance":10}`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindBalanceAccountAction generic error",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: output.AccountBalance{},
				err:    errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "FindBalanceAccountAction error parameter invalid",
			args: args{
				accountID: "error",
			},
			ucMock: mockFindBalanceAccount{
				result: output.AccountBalance{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["parameter invalid"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "FindBalanceAccountAction error fetching account",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: output.AccountBalance{},
				err:    domain.ErrNotFound,
			},
			expectedBody:       []byte(`{"errors":["not found"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/accounts/%s/balance", tt.args.accountID)
			req, _ := http.NewRequest(http.MethodGet, uri, nil)

			q := req.URL.Query()
			q.Add("account_id", tt.args.accountID)
			req.URL.RawQuery = q.Encode()

			var (
				w      = httptest.NewRecorder()
				action = NewFindBalanceAccountAction(tt.ucMock, log.LoggerMock{})
			)

			action.Execute(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = bytes.TrimSpace(w.Body.Bytes())
			if !bytes.Equal(result, tt.expectedBody) {
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}
