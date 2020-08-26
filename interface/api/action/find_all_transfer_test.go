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

type mockFindAllTransfer struct {
	result []output.Transfer
	err    error
}

func (m mockFindAllTransfer) Execute(_ context.Context) ([]output.Transfer, error) {
	return m.result, m.err
}

func TestTransfer_Index(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		ucMock             usecase.FindAllTransfer
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "FindAllTransferAction success one transfer",
			ucMock: mockFindAllTransfer{
				result: []output.Transfer{
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04679",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
						Amount:               10,
						CreatedAt:            time.Time{},
					},
				},
				err: nil,
			},
			expectedBody:       []byte(`[{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","account_origin_id":"3c096a40-ccba-4b58-93ed-57379ab04680","account_destination_id":"3c096a40-ccba-4b58-93ed-57379ab04681","amount":10,"created_at":"0001-01-01T00:00:00Z"}]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTransferAction success empty",
			ucMock: mockFindAllTransfer{
				result: []output.Transfer{},
				err:    nil,
			},
			expectedBody:       []byte(`[]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "FindAllTransferAction generic error",
			ucMock: mockFindAllTransfer{
				err: errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
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
