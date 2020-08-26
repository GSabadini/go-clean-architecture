package action

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type mockFindAllAccount struct {
	result []output.Account
	err    error
}

func (m mockFindAllAccount) Execute(_ context.Context) ([]output.Account, error) {
	return m.result, m.err
}

func TestFindAllAccountAction_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		ucMock             usecase.FindAllAccount
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "FindAllAccountAction success one account",
			ucMock: mockFindAllAccount{
				result: []output.Account{
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
						Name:      "Test",
						CPF:       "07094564964",
						Balance:   10,
						CreatedAt: time.Time{},
					},
				},
				err: nil,
			},
			expectedBody:       []byte(`[{"id":"3c096a40-ccba-4b58-93ed-57379ab04680","name":"Test","cpf":"07094564964","balance":10,"created_at":"0001-01-01T00:00:00Z"}]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllAccountAction success empty",
			ucMock: mockFindAllAccount{
				result: []output.Account{},
				err:    nil,
			},
			expectedBody:       []byte(`[]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllAccountAction generic error",
			ucMock: mockFindAllAccount{
				err: errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)

			var (
				w      = httptest.NewRecorder()
				action = NewFindAllAccountAction(tt.ucMock, log.LoggerMock{})
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
