package usecase

import (
	"errors"
	"github.com/gsabadini/go-bank-transfer/domain"
	"reflect"
	"testing"
	"time"
)

type mockAccountRepoStore struct {
	domain.AccountRepository

	result domain.Account
	err    error
}

func (m mockAccountRepoStore) Store(_ domain.Account) (domain.Account, error) {
	return m.result, m.err
}

func TestAccount_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		name, CPF string
		balance   float64
	}

	tests := []struct {
		name          string
		args          args
		repository    domain.AccountRepository
		expected      AccountOutput
		expectedError interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				name:    "Test",
				CPF:     "02815517078",
				balance: 100,
			},
			repository: mockAccountRepoStore{
				result: domain.Account{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   100,
					CreatedAt: time.Time{},
				},
				err: nil,
			},
			expected: AccountOutput{
				ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
				Name:      "Test",
				CPF:       "02815517078",
				Balance:   100,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "Create account generic error",
			args: args{
				name:    "",
				CPF:     "",
				balance: 0,
			},
			repository: mockAccountRepoStore{
				result: domain.Account{},
				err:    errors.New("error"),
			},
			expectedError: "error",
			expected:      AccountOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewAccount(tt.repository)
			result, err := uc.Store(tt.args.name, tt.args.CPF, tt.args.balance)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

type mockAccountRepoFindAll struct {
	domain.AccountRepository

	result []domain.Account
	err    error
}

func (m mockAccountRepoFindAll) FindAll() ([]domain.Account, error) {
	return m.result, m.err
}

func TestAccount_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		repository    domain.AccountRepository
		expected      []AccountOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the account list",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
						Name:      "Test",
						CPF:       "02815517078",
						Balance:   100,
						CreatedAt: time.Time{},
					},
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "02815517071",
						Balance:   1000,
						CreatedAt: time.Time{},
					},
				},
				err: nil,
			},
			expected: []AccountOutput{
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   100,
					CreatedAt: time.Time{},
				},
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
					Name:      "Test",
					CPF:       "02815517071",
					Balance:   1000,
					CreatedAt: time.Time{},
				},
			},
		},
		{
			name: "Success when returning the empty account list",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    nil,
			},
			expected: []AccountOutput{},
		},
		{
			name: "Error when returning the list of accounts",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    errors.New("error"),
			},
			expectedError: "error",
			expected:      []AccountOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewAccount(tt.repository)
			result, err := uc.FindAll()

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

type mockAccountRepoFindBalance struct {
	domain.AccountRepository

	result domain.Account
	err    error
}

func (m mockAccountRepoFindBalance) FindBalance(_ domain.AccountID) (domain.Account, error) {
	return m.result, m.err
}

func TestAccount_FindBalance(t *testing.T) {
	t.Parallel()

	type args struct {
		ID domain.AccountID
	}

	tests := []struct {
		name          string
		args          args
		repository    domain.AccountRepository
		expected      AccountBalanceOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			repository: mockAccountRepoFindBalance{
				result: domain.Account{
					Balance: 100,
				},
				err: nil,
			},
			expected: AccountBalanceOutput{
				Balance: 100.00,
			},
		},
		{
			name: "Error returning account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			repository: mockAccountRepoFindBalance{
				result: domain.Account{},
				err:    errors.New("error"),
			},
			expectedError: "error",
			expected:      AccountBalanceOutput{},
		},
	}

	for _, tt := range tests {
		var uc = NewAccount(tt.repository)
		result, err := uc.FindBalance(tt.args.ID)

		if (err != nil) && (err.Error() != tt.expectedError) {
			t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			return
		}

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
		}
	}
}
