package action

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsabadini/go-stone/infrastructure/database"
)

func TestAccountCreate(t *testing.T) {
	type args struct {
		accountAction Account
		rawPayload    []byte
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		args               args
	}{
		{
			name:               "Create handler success",
			expectedStatusCode: http.StatusCreated,
			args: args{
				accountAction: NewAccount(database.MongoHandlerSuccessMock{}),
				rawPayload:    []byte(`{"name": "test","cpf": "44451598087", "ballance": 10 }`),
			},
		},
		{
			name:               "Create handler database error",
			expectedStatusCode: http.StatusInternalServerError,
			args: args{
				accountAction: NewAccount(database.MongoHandlerErrorMock{}),
				rawPayload:    []byte(``),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := callHandler(t, tt.args.rawPayload, tt.args.accountAction)

			if err != nil {
				t.Fatal(err)
			}

			if statusCode != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] O handler retornou um HTTP status code inesperado: retornado '%v' esperado '%v'",
					tt.name,
					statusCode,
					tt.expectedStatusCode,
				)
			}
		})
	}

}

func callHandler(t *testing.T, rawPayload []byte, action Account) (int, error) {
	var body io.Reader = bytes.NewReader(rawPayload)

	req, err := http.NewRequest(http.MethodPost, "/account", body)
	if err != nil {
		t.Fatal(err)
	}

	var (
		rr      = httptest.NewRecorder()
		handler = http.HandlerFunc(action.Create)
	)

	handler.ServeHTTP(rr, req)

	return rr.Code, nil
}
