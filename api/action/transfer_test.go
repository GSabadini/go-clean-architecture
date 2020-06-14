package action

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/mock"

	"github.com/gorilla/mux"
)

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		transferAction     Transfer
		args               args
	}{
		{
			name:           "Store action success",
			transferAction: NewTransfer(mock.TransferUseCaseStubSuccess{}, mock.LoggerMock{}),
			args: args{
				rawPayload: []byte(`{
					"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
					"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
					"amount": 10 
				}`),
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:           "Store action error",
			transferAction: NewTransfer(mock.TransferUseCaseStubError{}, mock.LoggerMock{}),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Store action insufficient balance",
			transferAction: NewTransfer(
				mock.TransferUseCaseStubError{TypeErr: domain.ErrInsufficientBalance},
				mock.LoggerMock{},
			),
			args: args{
				rawPayload: []byte(
					`{
						"account_destination_id": "3c096a40-ccba-4b58-93ed-57379ab04680",
						"account_origin_id": "3c096a40-ccba-4b58-93ed-57379ab04681",
						"amount": 10
					}`,
				),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		//{
		//	name:           "Store action invalid UUID",
		//	transferAction: NewTransfer(stubTransfer.TransferUseCaseStubError{}, mock.LoggerMock{),
		//	args: args{
		//		rawPayload: []byte(
		//			`{
		//				"account_destination_id": "test",
		//				"account_origin_id": "test",
		//				"amount": 10
		//			}`,
		//		),
		//	},
		//	expectedStatusCode: http.StatusBadRequest,
		//},
		{
			name:           "Store action invalid JSON",
			transferAction: NewTransfer(mock.TransferUseCaseStubError{}, mock.LoggerMock{}),
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

func TestTransfer_Index(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		expectedStatusCode int
		transferAction     Transfer
	}{
		{
			name:               "Index action success",
			expectedStatusCode: http.StatusOK,
			transferAction:     NewTransfer(mock.TransferUseCaseStubSuccess{}, mock.LoggerMock{}),
		},
		{
			name:               "Index action error",
			expectedStatusCode: http.StatusInternalServerError,
			transferAction:     NewTransfer(mock.TransferUseCaseStubError{}, mock.LoggerMock{}),
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
