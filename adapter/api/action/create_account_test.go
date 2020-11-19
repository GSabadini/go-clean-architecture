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

	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validation"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type mockAccountCreateAccount struct {
	result usecase.CreateAccountOutput
	err    error
}

func (m mockAccountCreateAccount) Execute(_ context.Context, _ usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
	return m.result, m.err
}

func TestCreateAccountAction_Execute(t *testing.T) {
	t.Parallel()

	validator, _ := validation.NewValidatorFactory(validation.InstanceGoPlayground)

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.CreateAccountUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "CreateAccountAction success",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": 10050
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "07094564964",
					Balance:   10.5,
					CreatedAt: time.Time{}.String(),
				},
				err: nil,
			},
			expectedBody:       `{"id":"3c096a40-ccba-4b58-93ed-57379ab04680","name":"Test","cpf":"07094564964","balance":10.5,"created_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "CreateAccountAction success",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087", 
						"balance": 100000
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "07094564964",
					Balance:   10000,
					CreatedAt: time.Time{}.String(),
				},
				err: nil,
			},
			expectedBody:       `{"id":"3c096a40-ccba-4b58-93ed-57379ab04680","name":"Test","cpf":"07094564964","balance":10000,"created_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "CreateAccountAction generic error",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087",
						"balance": 10
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{},
				err:    errors.New("error"),
			},
			expectedBody:       `{"errors":["error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "CreateAccountAction error invalid balance",
			args: args{
				rawPayload: []byte(
					`{
						"name": "test",
						"cpf": "44451598087",
						"balance": -1
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["Balance must be greater than 0"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateAccountAction error invalid fields",
			args: args{
				rawPayload: []byte(
					`{
						"name123": "test",
						"cpf1231": "44451598087",
						"balance12312": 1
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["Name is a required field","CPF is a required field","Balance must be greater than 0"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateAccountAction error invalid JSON",
			args: args{
				rawPayload: []byte(
					`{
						"name":
					}`,
				),
			},
			ucMock: mockAccountCreateAccount{
				result: usecase.CreateAccountOutput{},
				err:    nil,
			},
			expectedBody:       `{"errors":["invalid character '}' looking for beginning of value"]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateAccountAction(tt.ucMock, log.LoggerMock{}, validator)
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
