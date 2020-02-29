package usecase

import (
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

func TestCreate(t *testing.T) {
	type args struct {
		repository repository.AccountRepository
		account    domain.Account
	}

	tests := []struct {
		name          string
		args          args
		expected      domain.Account
		expectedError interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				repository: repository.AccountRepositoryMockSuccess{},
			},
			expected: domain.Account{
				ID:      "5e570851adcef50116aa7a5c",
				Name:    "Test",
				CPF:     "028.155.170-78",
				Balance: 100,
			},
		},
		{
			name: "Create account error",
			args: args{
				repository: repository.AccountRepositoryMockError{},
			},
			expectedError: "Error",
			expected:      domain.Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreAccount(tt.args.repository, tt.args.account)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(got, tt.expected) || got != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	type args struct {
		repository repository.AccountRepository
	}

	tests := []struct {
		name          string
		args          args
		expected      []domain.Account
		expectedError interface{}
	}{
		{
			name: "Success when returning the account list",
			args: args{
				repository: repository.AccountRepositoryMockSuccess{},
			},
			expected: []domain.Account{
				{
					ID:      "5e570851adcef50116aa7a5c",
					Name:    "Test-0",
					CPF:     "028.155.170-78",
					Balance: 0,
				},
				{
					ID:      "5e570854adcef50116aa7a5d",
					Name:    "Test-1",
					CPF:     "028.155.170-78",
					Balance: 50.25,
				},
			},
		},
		{
			name: "Error when returning the list of accounts",
			args: args{
				repository: repository.AccountRepositoryMockError{},
			},
			expectedError: "Error",
			expected:      []domain.Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := FindAllAccount(tt.args.repository)
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestFindBalanceAccount(t *testing.T) {
	type args struct {
		repository repository.AccountRepository
		id         string
	}

	tests := []struct {
		name          string
		args          args
		expected      domain.Account
		expectedError interface{}
	}{
		{
			name: "Success when returning the account balance",
			args: args{
				repository: repository.AccountRepositoryMockSuccess{},
				id:         "5e519055ba39bfc244dc4625",
			},
			expected: domain.Account{
				Balance: 100.00,
			},
		},
		{
			name: "Error returning account balance",
			args: args{
				repository: repository.AccountRepositoryMockError{},
				id:         "5e519055ba39bfc244dc4625",
			},
			expectedError: "Error",
			expected:      domain.Account{},
		},
	}

	for _, tt := range tests {
		got, err := FindBalanceAccount(tt.args.repository, tt.args.id)

		if (err != nil) && (err.Error() != tt.expectedError) {
			t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			return
		}

		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
		}
	}
}
