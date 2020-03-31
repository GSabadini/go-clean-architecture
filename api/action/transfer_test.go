package action

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	stubTransfer "github.com/gsabadini/go-bank-transfer/usecase/stub"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestTransferStore(t *testing.T) {
	type args struct {
		rawPayload []byte
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		transferAction     Transfer
		args               args
	}{
		{
			name:           "Store action success",
			transferAction: NewTransfer(stubTransfer.TransferUseCaseStubSuccess{}, loggerMock),
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "5e570851adcef50116aa7a5c",
					"account_origin_id": "5e570851adcef50116aa7a5a",
					"amount": 10 
				}`),
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:           "Store action error",
			transferAction: NewTransfer(stubTransfer.TransferUseCaseStubError{}, loggerMock),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "5e551c2c5bb0cb0107b058e1",
						"account_origin_id": "5e551c315bb0cb0107b058e2",
						"amount": 10
					}`,
				),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Store action insufficient balance",
			transferAction: NewTransfer(
				stubTransfer.TransferUseCaseStubError{TypeErr: domain.ErrInsufficientBalance},
				loggerMock,
			),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "5e551c2c5bb0cb0107b058e1",
						"account_origin_id": "5e551c315bb0cb0107b058e2",
						"amount": 10
					}`,
				),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "Store action invalid ObjectID",
			transferAction: NewTransfer(stubTransfer.TransferUseCaseStubError{}, loggerMock),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "test",
						"account_origin_id": "test", 
						"amount": 10 
					}`,
				),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Store action invalid JSON",
			transferAction: NewTransfer(stubTransfer.TransferUseCaseStubError{}, loggerMock),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": ,
						"account_origin_id": ,
						"amount": 
					}`,
				),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = bytes.NewReader(tt.args.rawPayload)
			req, err := http.NewRequest(http.MethodPost, "/transfers", body)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/transfers", tt.transferAction.Store).Methods(http.MethodPost)
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

func TestTransferIndex(t *testing.T) {
	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		transferAction     Transfer
	}{
		{
			name:               "Index action success",
			expectedStatusCode: http.StatusOK,
			transferAction:     NewTransfer(stubTransfer.TransferUseCaseStubSuccess{}, loggerMock),
		},
		{
			name:               "Index action error",
			expectedStatusCode: http.StatusInternalServerError,
			transferAction:     NewTransfer(stubTransfer.TransferUseCaseStubError{}, loggerMock),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/transfers", nil)
			if err != nil {
				t.Fatal(err)
			}

			var (
				rr = httptest.NewRecorder()
				r  = mux.NewRouter()
			)

			r.HandleFunc("/transfers", tt.transferAction.Index).Methods(http.MethodGet)
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
