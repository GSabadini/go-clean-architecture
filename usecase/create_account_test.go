package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type mockAccountRepoStore struct {
	domain.AccountRepository

	result domain.Account
	err    error
}

func (m mockAccountRepoStore) Create(_ context.Context, _ domain.Account) (domain.Account, error) {
	return m.result, m.err
}

type mockAccountPresenterStore struct {
	output.AccountPresenter

	result output.Account
}

func (m mockAccountPresenterStore) Output(_ domain.Account) output.Account {
	return m.result
}

func TestCreateAccountInteractor_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		input input.Account
	}

	tests := []struct {
		name          string
		args          args
		repository    domain.AccountRepository
		presenter     output.AccountPresenter
		expected      output.Account
		expectedError interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				input: input.Account{
					Name:    "Test",
					CPF:     "02815517078",
					Balance: 19944,
				},
			},
			repository: mockAccountRepoStore{
				result: domain.NewAccount(
					"3c096a40-ccba-4b58-93ed-57379ab04680",
					"Test",
					"02815517078",
					19944,
					time.Time{},
				),
				err: nil,
			},
			presenter: mockAccountPresenterStore{
				result: output.Account{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   199.44,
					CreatedAt: time.Time{},
				},
			},
			expected: output.Account{
				ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
				Name:      "Test",
				CPF:       "02815517078",
				Balance:   199.44,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "Create account successful",
			args: args{
				input: input.Account{
					Name:    "Test",
					CPF:     "02815517078",
					Balance: 2350,
				},
			},
			repository: mockAccountRepoStore{
				result: domain.NewAccount(
					"3c096a40-ccba-4b58-93ed-57379ab04680",
					"Test",
					"02815517078",
					2350,
					time.Time{},
				),
				err: nil,
			},
			presenter: mockAccountPresenterStore{
				result: output.Account{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   23.5,
					CreatedAt: time.Time{},
				},
			},
			expected: output.Account{
				ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
				Name:      "Test",
				CPF:       "02815517078",
				Balance:   23.5,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "Create account generic error",
			args: args{
				input: input.Account{
					Name:    "",
					CPF:     "",
					Balance: 0,
				},
			},
			repository: mockAccountRepoStore{
				result: domain.Account{},
				err:    errors.New("error"),
			},
			presenter: mockAccountPresenterStore{
				result: output.Account{},
			},
			expectedError: "error",
			expected:      output.Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateAccountInteractor(tt.repository, tt.presenter, time.Second)

			result, err := uc.Execute(context.TODO(), tt.args.input)
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}
