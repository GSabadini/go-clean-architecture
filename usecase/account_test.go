package usecase

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
)

func TestCreate(t *testing.T) {
	type args struct {
		repository repository.Account
		account    *domain.Account
	}

	//var timeNow = time.Now()

	tests := []struct {
		name          string
		args          args
		expected      interface{}
		expectedError interface{}
	}{
		//{
		//	name: "Create account successful",
		//	args: args{
		//		repository: repository.NewAccount(database.MongoHandlerSuccessMock{}),
		//		account: domain.Account{
		//			Name:      "Test",
		//			Cpf:       "44451598087",
		//			Balance:  10.12,
		//			CreatedAt: timeNow,
		//		},
		//	},
		//	expected: domain.Account{
		//		Name:      "Test",
		//		Cpf:       "44451598087",
		//		Balance:  10.12,
		//		CreatedAt: timeNow,
		//	},
		//},
		//{
		//	name: "Create account error",
		//	args: args{
		//		repository: repository.NewAccount(database.MongoHandlerErrorMock{}),
		//		account: &domain.Account{
		//			Name:      "Test",
		//			Cpf:       "44451598087",
		//			Balance:  0,
		//			CreatedAt: time.Now(),
		//		},
		//	},
		//	expectedError: "Error",
		//	expected:      nil,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreAccount(tt.args.repository, tt.args.account)

			fmt.Println(got, tt.expected)

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
	//var timeNow = time.Now()

	type args struct {
		repository repository.AccountRepository
	}

	tests := []struct {
		name          string
		args          args
		expected      []domain.Account
		expectedError interface{}
	}{
		//{
		//	name: "Success return list accounts",
		//	args: args{
		//		repository: repository.NewAccount(database.MongoHandlerSuccessMock{}),
		//	},
		//	expected: []domain.Account{
		//		{
		//			Id:        "0",
		//			Name:      "Test-0",
		//			Cpf:       "",
		//			Balance:  0,
		//			CreatedAt: timeNow,
		//		},
		//		{
		//			Id:        "1",
		//			Name:      "Test-1",
		//			Cpf:       "",
		//			Balance:  120,
		//			CreatedAt: timeNow,
		//		},
		//	},
		//},
		{
			name: "Empty return list accounts",
			args: args{
				repository: repository.NewAccount(database.MongoHandlerSuccessMock{}),
			},
			expected: []domain.Account{},
		},
		{
			name: "Error return list accounts",
			args: args{
				repository: repository.NewAccount(database.MongoHandlerErrorMock{}),
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
		expected      *AccountBalance
		expectedError interface{}
	}{
		//{
		//	name: "Success return balance account",
		//	args: args{
		//		repository: repository.NewAccount(database.MongoHandlerSuccessMock{}),
		//		id:         "5e519055ba39bfc244dc4625",
		//	},
		//	expected: domain.Account{
		//		Balance: 100.00,
		//	},
		//},
		{
			name: "Error return balance account",
			args: args{
				repository: repository.NewAccount(database.MongoHandlerErrorMock{}),
				id:         "5e519055ba39bfc244dc4625",
			},
			expectedError: "Error",
			expected:      &AccountBalance{},
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
