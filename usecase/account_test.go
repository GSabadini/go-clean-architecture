package usecase

import (
	"github.com/gsabadini/go-bank-transfer/repository/stub"
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
)

func TestStoreAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		account domain.Account
	}

	tests := []struct {
		name          string
		args          args
		usecase       AccountUseCase
		expected      domain.Account
		expectedError interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				account: domain.Account{},
			},
			usecase: NewAccount(stub.AccountRepositoryStubSuccess{}),
			expected: domain.Account{
				ID:      "5e570851adcef50116aa7a5c",
				Name:    "Test",
				CPF:     "02815517078",
				Balance: 100,
			},
		},
		{
			name:          "Create account error",
			args:          args{account: domain.Account{}},
			usecase:       NewAccount(stub.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      domain.Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.usecase.Store(tt.args.account)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

func TestFindAllAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		usecase       AccountUseCase
		expected      []domain.Account
		expectedError interface{}
	}{
		{
			name:    "Success when returning the account list",
			usecase: NewAccount(stub.AccountRepositoryStubSuccess{}),
			expected: []domain.Account{
				{
					ID:      "5e570851adcef50116aa7a5c",
					Name:    "Test-0",
					CPF:     "02815517078",
					Balance: 0,
				},
				{
					ID:      "5e570854adcef50116aa7a5d",
					Name:    "Test-1",
					CPF:     "02815517078",
					Balance: 50.25,
				},
			},
		},
		{
			name:          "Error when returning the list of accounts",
			usecase:       NewAccount(stub.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      []domain.Account{},
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

func TestFindBalanceAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		ID string
	}

	tests := []struct {
		name          string
		args          args
		usecase       AccountUseCase
		expected      domain.Account
		expectedError interface{}
	}{
		{
			name: "Success when returning the account balance",
			args: args{
				ID: "5e519055ba39bfc244dc4625",
			},
			usecase: NewAccount(stub.AccountRepositoryStubSuccess{}),
			expected: domain.Account{
				Balance: 100.00,
			},
		},
		{
			name: "Error returning account balance",
			args: args{
				ID: "5e519055ba39bfc244dc4625",
			},
			usecase:       NewAccount(stub.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      domain.Account{},
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
