package action

import (
	"bytes"
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

func (m mockTransferStore) Store(_, _ string, _ float64) (usecase.TransferOutput, error) {
	return m.result, m.err
}

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		transferAction     Transfer
		args               args
		usecaseMock        usecase.TransferUseCase
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name: "Store action success",
			transferAction: NewTransfer(
				usecase.TransferUseCaseStubSuccess{},
				logger.LoggerMock{},
				validator,
			),
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 10
				}`),
			},
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
			usecaseMock: mockTransferStore{
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
				action = NewTransfer(tt.usecaseMock, logger.LoggerMock{}, validator)
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

func TestTransfer_Index(t *testing.T) {
	t.Parallel()

	validator, _ := validator.NewValidatorFactory(validator.InstanceGoPlayground, logger.LoggerMock{})

	tests := []struct {
		name               string
		expectedStatusCode int
		transferAction     Transfer
	}{
		{
			name:               "Index action success",
			expectedStatusCode: http.StatusOK,
			transferAction: NewTransfer(
				usecase.TransferUseCaseStubSuccess{},
				logger.LoggerMock{},
				validator,
			),
		},
		{
			name:               "Index action error",
			expectedStatusCode: http.StatusInternalServerError,
			transferAction: NewTransfer(
				usecase.TransferUseCaseStubError{},
				logger.LoggerMock{},
				validator,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/transfers", nil)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr      = httptest.NewRecorder()
				handler = http.NewServeMux()
			)

			handler.HandleFunc("/transfers", tt.transferAction.Index)
			handler.ServeHTTP(rr, req)

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
