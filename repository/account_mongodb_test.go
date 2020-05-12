package repository

import (
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database/stub"
	"gopkg.in/mgo.v2"
)

func TestAccountMongoDB_Store(t *testing.T) {
	type args struct {
		account domain.Account
	}

	tests := []struct {
		name          string
		args          args
		expected      domain.Account
		expectedError interface{}
		repository    AccountMongoDB
		expectedErr   bool
	}{
		{
			name:       "Success to create account",
			args:       args{account: domain.Account{}},
			expected:   domain.Account{},
			repository: NewAccountMongoDB(stub.MongoHandlerSuccessStub{}),
		},
		{
			name:        "Error to create account",
			args:        args{account: domain.Account{}},
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{}),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repository.Store(tt.args.account)

			if (err != nil) != tt.expectedErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | ExpectedErr: '%v'", tt.name, err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestAccountMongoDB_UpdateBalance(t *testing.T) {
	type args struct {
		ID      string
		balance float64
	}

	tests := []struct {
		name        string
		repository  AccountMongoDB
		args        args
		expectedErr bool
	}{
		{
			name:       "Success to account update",
			repository: NewAccountMongoDB(stub.MongoHandlerSuccessStub{}),
			args: args{
				ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				balance: 100.00,
			},
		},
		{
			name:       "Error to account update",
			repository: NewAccountMongoDB(stub.MongoHandlerErrorStub{}),
			args: args{
				ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				balance: 1.00,
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.repository.UpdateBalance(tt.args.ID, tt.args.balance); (err != nil) != tt.expectedErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | ExpectedErr: '%v'", tt.name, err, tt.expectedErr)
			}
		})
	}
}

func TestAccountMongoDB_FindAll(t *testing.T) {
	tests := []struct {
		name        string
		repository  AccountMongoDB
		expected    []domain.Account
		expectedErr bool
	}{
		{
			name:       "Success in finding all accounts",
			repository: NewAccountMongoDB(stub.MongoHandlerSuccessStub{}),

			expected: []domain.Account{},
		},
		{
			name:        "Error in finding all accounts",
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{}),
			expected:    []domain.Account{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repository.FindAll()

			if (err != nil) != tt.expectedErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | ExpectedErr: '%v'", tt.name, err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestAccountMongoDB_FindByID(t *testing.T) {
	type args struct {
		ID string
	}

	tests := []struct {
		name        string
		repository  AccountMongoDB
		args        args
		expected    *domain.Account
		expectedErr bool
	}{
		{
			name:       "Success to find account",
			repository: NewAccountMongoDB(stub.MongoHandlerSuccessStub{}),
			args:       args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:   &domain.Account{},
		},
		{
			name:        "Error to find account",
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{}),
			args:        args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:    &domain.Account{},
			expectedErr: true,
		},
		{
			name:        "Account not found",
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{TypeErr: mgo.ErrNotFound}),
			args:        args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:    &domain.Account{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repository.FindByID(tt.args.ID)

			if (err != nil) != tt.expectedErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | ExpectedErr: '%v'", tt.name, err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestAccountMongoDB_FindBalance(t *testing.T) {
	type args struct {
		ID string
	}

	tests := []struct {
		name        string
		repository  AccountMongoDB
		args        args
		expected    domain.Account
		expectedErr bool
	}{
		{
			name:       "Success to find account balance",
			repository: NewAccountMongoDB(stub.MongoHandlerSuccessStub{}),
			args:       args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:   domain.Account{},
		},
		{
			name:        "Error to find account balance",
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{}),
			args:        args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:    domain.Account{},
			expectedErr: true,
		},
		{
			name:        "Account balance not found",
			repository:  NewAccountMongoDB(stub.MongoHandlerErrorStub{TypeErr: mgo.ErrNotFound}),
			args:        args{ID: "3c096a40-ccba-4b58-93ed-57379ab04680"},
			expected:    domain.Account{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repository.FindBalance(tt.args.ID)

			if (err != nil) != tt.expectedErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | ExpectedErr: '%v'", tt.name, err, tt.expectedErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}
