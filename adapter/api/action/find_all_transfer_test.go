package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockFindAllTransfer struct {
	result []usecase.FindAllTransferOutput
	err    error
}

func (m mockFindAllTransfer) Execute(_ context.Context) ([]usecase.FindAllTransferOutput, error) {
	return m.result, m.err
}

func TestTransfer_Index(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		ucMock             usecase.FindAllTransferUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "FindAllTransferAction success one transfer",
			ucMock: mockFindAllTransfer{
				result: []usecase.FindAllTransferOutput{
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04679",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
						Amount:               10,
						CreatedAt:            time.Time{}.String(),
					},
				},
				err: nil,
			},
			expectedBody:       `[{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","account_origin_id":"3c096a40-ccba-4b58-93ed-57379ab04680","account_destination_id":"3c096a40-ccba-4b58-93ed-57379ab04681","amount":10,"created_at":"0001-01-01 00:00:00 +0000 UTC"}]`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTransferAction success empty",
			ucMock: mockFindAllTransfer{
				result: []usecase.FindAllTransferOutput{},
				err:    nil,
			},
			expectedBody:       `[]`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTransferAction generic error",
			ucMock: mockFindAllTransfer{
				err: errors.New("error"),
			},
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/transfers", nil)

			var (
				w      = httptest.NewRecorder()
				action = NewFindAllTransferAction(tt.ucMock, log.LoggerMock{})
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
