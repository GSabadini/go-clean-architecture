package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/repository"
)

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		accountOriginID      string
		accountDestinationID string
		amount               float64
	}

	tests := []struct {
		name          string
		args          args
		usecase       TransferUseCase
		expected      TransferOutput
		expectedError string
	}{
		{
			name: "Create transfer successful",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				amount:               20,
			},
			usecase: NewTransfer(repository.TransferRepositoryStubSuccess{}, repository.AccountRepositoryStubSuccess{}),
			expected: TransferOutput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               20,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "Create transfer error",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               20,
			},
			usecase:       NewTransfer(repository.TransferRepositoryStubError{}, repository.AccountRepositoryStubSuccess{}),
			expectedError: "Error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer amount not have sufficient",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               200,
			},
			usecase:       NewTransfer(repository.TransferRepositoryStubSuccess{}, repository.AccountRepositoryStubSuccess{}),
			expectedError: "origin account does not have sufficient balance",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error find account",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               200,
			},
			usecase:       NewTransfer(repository.TransferRepositoryStubSuccess{}, repository.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      TransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.usecase.Store(tt.args.accountOriginID, tt.args.accountDestinationID, tt.args.amount)

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

func TestTransfer_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expected      []TransferOutput
		usecase       TransferUseCase
		expectedError string
	}{
		{
			name:    "Success when returning the transfer list",
			usecase: NewTransfer(repository.TransferRepositoryStubSuccess{}, repository.AccountRepositoryStubSuccess{}),
			expected: []TransferOutput{
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               100,
					CreatedAt:            time.Time{},
				},
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               500,
					CreatedAt:            time.Time{},
				},
			},
		},
		{
			name:          "Error when returning the transfer list",
			usecase:       NewTransfer(repository.TransferRepositoryStubError{}, repository.AccountRepositoryStubSuccess{}),
			expectedError: "Error",
			expected:      []TransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.usecase.FindAll()

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}
