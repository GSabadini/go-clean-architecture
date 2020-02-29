package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
	"reflect"
	"testing"
	"time"
)

func TestStoreTransfer(t *testing.T) {
	type args struct {
		transferRepository repository.TransferRepository
		accountRepository  repository.AccountRepository
		transfer           domain.Transfer
	}

	tests := []struct {
		name          string
		args          args
		expected      domain.Transfer
		expectedError string
	}{
		{
			name: "Create transfer successful",
			args: args{
				transferRepository: repository.TransferRepositoryMockSuccess{},
				accountRepository:  repository.AccountRepositoryMockSuccess{},
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
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
				transferRepository: repository.TransferRepositoryMockError{},
				accountRepository:  repository.AccountRepositoryMockSuccess{},
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			expectedError: "Error",
			expected:      domain.Transfer{},
		},
		{
			name: "Create transfer error find account",
			args: args{
				transferRepository: repository.TransferRepositoryMockError{},
				accountRepository:  repository.AccountRepositoryMockError{},
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			expectedError: "error fetching account: Error",
			expected:      domain.Transfer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreTransfer(tt.args.transferRepository, tt.args.accountRepository, tt.args.transfer)

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
	type args struct {
		repository repository.TransferRepository
	}

	tests := []struct {
		name          string
		args          args
		expected      []domain.Transfer
		expectedError string
	}{
		{
			name: "Success when returning the transfer list",
			args: args{
				repository: repository.TransferRepositoryMockSuccess{},
			},
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
			name: "Error when returning the transfer list",
			args: args{
				repository: repository.TransferRepositoryMockError{},
			},
			expectedError: "Error",
			expected:      []domain.Transfer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindAllTransfer(tt.args.repository)

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
