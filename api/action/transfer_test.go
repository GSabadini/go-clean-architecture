package action

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockTransferStore struct {
	usecase.TransferUseCase

	result usecase.TransferOutput
	err    error
}

func (m mockTransferStore) Store(
	_ context.Context,
	_ domain.AccountID,
	_ domain.AccountID,
	_ float64,
) (usecase.TransferOutput, error) {
	return m.result, m.err
}

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.TransferUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "Store action success",
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 10
				}`),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04679",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               10,
					CreatedAt:            time.Time{},
				},
				err: nil,
			},
			expectedBody:       []byte(`{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","account_origin_id":"3c096a40-ccba-4b58-93ed-57379ab04680","account_destination_id":"3c096a40-ccba-4b58-93ed-57379ab04681","amount":10,"created_at":"0001-01-01T00:00:00Z"}`),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Store action generic error",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       []byte(`{"errors":["error"]}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Store action error insufficient balance",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    domain.ErrInsufficientBalance,
			},
			expectedBody:       []byte(`{"errors":["origin account does not have sufficient balance"]}`),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Store action error account origin equals account destination",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"amount": 10
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["account origin equals destination account"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store action error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": ,
						"account_origin_id": ,
						"amount":
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["invalid character ',' looking for beginning of value"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store action error invalid amount",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": -1
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["Amount must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store action error invalid fields",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id123": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id123": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount123": 10
					}`,
				),
			},
			ucMock: mockTransferStore{
				result: usecase.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`{"errors":["AccountOriginID is a required field","AccountDestinationID is a required field","Amount must be greater than 0"]}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/transfers",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewTransfer(tt.ucMock, logger.LoggerMock{}, validator)
			)

			action.Store(w, req)

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

type mockTransferFindAll struct {
	usecase.TransferUseCase

	result []usecase.TransferOutput
	err    error
}

func (m mockTransferFindAll) FindAll(_ context.Context) ([]usecase.TransferOutput, error) {
	return m.result, m.err
}

func TestTransfer_Index(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground)

	tests := []struct {
		name               string
		ucMock             usecase.TransferUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "Index handler success one transfer",
			ucMock: mockTransferFindAll{
				result: []usecase.TransferOutput{
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
			name: "Index handler success empty",
			ucMock: mockTransferFindAll{
				result: []usecase.TransferOutput{},
				err:    nil,
			},
			expectedBody:       []byte(`[]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Index handler generic error",
			ucMock: mockTransferFindAll{
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
				action = NewTransfer(tt.ucMock, logger.LoggerMock{}, validator)
			)

			action.Index(w, req)

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
