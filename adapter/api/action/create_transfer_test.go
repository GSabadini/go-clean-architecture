package action

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validation"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockCreateTransfer struct {
	result usecase.CreateTransferOutput
	err    error
}

func (m mockCreateTransfer) Execute(_ context.Context, _ usecase.CreateTransferInput) (usecase.CreateTransferOutput, error) {
	return m.result, m.err
}

func TestCreateTransferAction_Execute(t *testing.T) {
	t.Parallel()

	validator, _ := validation.NewValidatorFactory(validation.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.CreateTransferUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "CreateTransferAction success",
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 10
				}`),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04679",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               10,
					CreatedAt:            time.Time{}.String(),
				},
				err: nil,
			},
			expectedBody:       `{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","account_origin_id":"3c096a40-ccba-4b58-93ed-57379ab04680","account_destination_id":"3c096a40-ccba-4b58-93ed-57379ab04681","amount":10,"created_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "CreateTransferAction generic error",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "CreateTransferAction error insufficient balance",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    domain.ErrInsufficientBalance,
			},
			expectedBody:       `{"errors":["origin account does not have sufficient balance"]}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "CreateTransferAction error not found account origin",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    domain.ErrAccountOriginNotFound,
			},
			expectedBody:       `{"errors":["account origin not found"]}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "CreateTransferAction error not found account destination",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    domain.ErrAccountDestinationNotFound,
			},
			expectedBody:       `{"errors":["account destination not found"]}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "CreateTransferAction error account origin equals account destination",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"amount": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["account origin equals destination account"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateTransferAction error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": ,
						"account_origin_id": ,
						"amount":
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["invalid character ',' looking for beginning of value"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateTransferAction error invalid amount",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": -1
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["Amount must be greater than 0"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateTransferAction error invalid fields",
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id123": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id123": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount123": 10
					}`,
				),
			},
			ucMock: mockCreateTransfer{
				result: usecase.CreateTransferOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["AccountOriginID is a required field","AccountDestinationID is a required field","Amount must be greater than 0"]}`,
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
				action = NewCreateTransferAction(tt.ucMock, log.LoggerMock{}, validator)
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
