package action

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockFindBalanceAccount struct {
	result usecase.FindAccountBalanceOutput
	err    error
}

func (m mockFindBalanceAccount) Execute(_ context.Context, _ domain.AccountID) (usecase.FindAccountBalanceOutput, error) {
	return m.result, m.err
}

func TestFindAccountBalanceAction_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		accountID string
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.FindAccountBalanceUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "FindAccountBalanceAction success",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: usecase.FindAccountBalanceOutput{
					Balance: 10,
				},
				err: nil,
			},
			expectedBody:       `{"balance":10}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAccountBalanceAction generic error",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: usecase.FindAccountBalanceOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "FindAccountBalanceAction error parameter invalid",
			args: args{
				accountID: "error",
			},
			ucMock: mockFindBalanceAccount{
				result: usecase.FindAccountBalanceOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["parameter invalid"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "FindAccountBalanceAction error fetching account",
			args: args{
				accountID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			ucMock: mockFindBalanceAccount{
				result: usecase.FindAccountBalanceOutput{},
				err:    domain.ErrAccountNotFound,
			},
			expectedBody:       `{"errors":["account not found"]}`,
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
				action = NewFindAccountBalanceAction(tt.ucMock, log.LoggerMock{})
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

			var result = strings.TrimSpace(w.Body.String())
			if !strings.EqualFold(result, tt.expectedBody) {
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
