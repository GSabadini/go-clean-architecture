package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository/stub"
)

func TestStoreTransfer(t *testing.T) {
	t.Parallel()

	type args struct {
		transfer domain.Transfer
	}

	tests := []struct {
		name          string
		args          args
		usecase       TransferUseCase
		expected      domain.Transfer
		expectedError string
	}{
		{
			name: "Create transfer successful",
			args: args{
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			usecase: NewTransferService(stub.TransferRepositoryStubSuccess{}, stub.AccountRepositoryStubSuccess{}),
			expected: domain.Transfer{
				ID:                   "5e570851adcef50116aa7a5a",
				AccountOriginID:      "5e570851adcef50116aa7a5d",
				AccountDestinationID: "5e570851adcef50116aa7a5c",
				Amount:               20,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "Create transfer error",
			args: args{
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			usecase:       NewTransferService(stub.TransferRepositoryStubError{}, stub.AccountRepositoryStubSuccess{}),
			expectedError: "Error",
			expected:      domain.Transfer{},
		},
		{
			name: "Create transfer amount not have sufficient",
			args: args{
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               200,
				},
			},
			usecase:       NewTransferService(stub.TransferRepositoryStubSuccess{}, stub.AccountRepositoryStubSuccess{}),
			expectedError: "Origin account does not have sufficient balance",
			expected:      domain.Transfer{},
		},
		{
			name: "Create transfer error find account",
			args: args{
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			usecase:       NewTransferService(stub.TransferRepositoryStubSuccess{}, stub.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      domain.Transfer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.usecase.StoreTransfer(tt.args.transfer)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) || got != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

func TestFindAllTransfer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expected      []domain.Transfer
		usecase       TransferUseCase
		expectedError string
	}{
		{
			name:    "Success when returning the transfer list",
			usecase: NewTransferService(stub.TransferRepositoryStubSuccess{}, stub.AccountRepositoryStubSuccess{}),
			expected: []domain.Transfer{
				{
					ID:                   "5e570851adcef50116aa7a5a",
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               100,
					CreatedAt:            time.Time{},
				},
				{
					ID:                   "5e570851adcef50116aa7a5b",
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               500,
					CreatedAt:            time.Time{},
				},
			},
		},
		{
			name:          "Error when returning the transfer list",
			usecase:       NewTransferService(stub.TransferRepositoryStubError{}, stub.AccountRepositoryStubSuccess{}),
			expectedError: "Error",
			expected:      []domain.Transfer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.usecase.FindAllTransfer()

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
