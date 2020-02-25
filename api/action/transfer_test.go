package action

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

//TODO REVER TESTES
func TestTransferStore(t *testing.T) {
	type args struct {
		transferAction Transfer
		rawPayload     []byte
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		//{
		//	name: "Store handler success",
		//	args: args{
		//		transferAction: NewTransfer(database.MongoHandlerSuccessMock{}, loggerMock),
		//		rawPayload: []byte(
		//			`{
		//				"account_destination_id": "5e551c2c5bb0cb0107b058e1",
		//				"account_origin_id": "5e551c315bb0cb0107b058e2",
		//				"amount": 10
		//			}`,
		//		),
		//	},
		//	expectedStatusCode: http.StatusNoContent,
		//},
		{
			name: "Store handler invalid ObjectID",
			args: args{
				transferAction: NewTransfer(database.MongoHandlerSuccessMock{}, loggerMock),
				rawPayload:     []byte(`{"account_destination_id": "test","account_origin_id": "test", "amount": 10 }`),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Store handler database error",
			args: args{
				transferAction: NewTransfer(database.MongoHandlerErrorMock{}, loggerMock),
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

			r.HandleFunc("/transfers", tt.args.transferAction.Store).Methods(http.MethodPost)
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
	type args struct {
		transferAction Transfer
	}

	var loggerMock, _ = test.NewNullLogger()

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Index handler database success",
			expectedStatusCode: http.StatusOK,
			args: args{
				transferAction: NewTransfer(database.MongoHandlerSuccessMock{}, loggerMock),
			},
		},
		{
			name:               "Index handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				transferAction: NewTransfer(database.MongoHandlerErrorMock{}, loggerMock),
			},
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

			r.HandleFunc("/transfers", tt.args.transferAction.Index).Methods(http.MethodGet)
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
