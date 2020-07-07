package usecase

import (
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/repository"
)

func TestAccount_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		name, CPF string
		balance   float64
	}

	tests := []struct {
		name          string
		args          args
		usecase       AccountUseCase
		expected      AccountOutput
		expectedError interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				name:    "",
				CPF:     "",
				balance: 0,
			},
			usecase: NewAccount(repository.AccountRepositoryStubSuccess{}),
			expected: AccountOutput{
				ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				Name:    "Test",
				CPF:     "02815517078",
				Balance: 100,
			},
		},
		{
			name: "Create account error",
			args: args{
				name:    "",
				CPF:     "",
				balance: 0,
			},
			usecase:       NewAccount(repository.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      AccountOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.usecase.Store(tt.args.name, tt.args.CPF, tt.args.balance)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

func TestAccount_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		usecase       AccountUseCase
		expected      []AccountOutput
		expectedError interface{}
	}{
		{
			name:    "Success when returning the account list",
			usecase: NewAccount(repository.AccountRepositoryStubSuccess{}),
			expected: []AccountOutput{
				{
					ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:    "Test-0",
					CPF:     "02815517078",
					Balance: 0,
				},
				{
					ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:    "Test-1",
					CPF:     "02815517078",
					Balance: 50.25,
				},
			},
		},
		{
			name:          "Error when returning the list of accounts",
			usecase:       NewAccount(repository.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      []AccountOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.usecase.FindAll()

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

func TestAccount_FindBalance(t *testing.T) {
	t.Parallel()

	type args struct {
		ID string
	}

	tests := []struct {
		name          string
		args          args
		usecase       AccountUseCase
		expected      AccountBalanceOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			usecase: NewAccount(repository.AccountRepositoryStubSuccess{}),
			expected: AccountBalanceOutput{
				Balance: 100.00,
			},
		},
		{
			name: "Error returning account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			usecase:       NewAccount(repository.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      AccountBalanceOutput{},
		},
	}

	for _, tt := range tests {
		result, err := tt.usecase.FindBalance(tt.args.ID)

		if (err != nil) && (err.Error() != tt.expectedError) {
			t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			return
		}

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
		}
	}
}
