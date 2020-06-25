package usecase

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/mock"
)

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		input TransferInput
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
				input: TransferInput{
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               20,
				},
			},
			usecase: NewTransfer(mock.TransferRepositoryStubSuccess{}, mock.AccountRepositoryStubSuccess{}),
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
				input: TransferInput{
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               20,
				},
			},
			usecase:       NewTransfer(mock.TransferRepositoryStubError{}, mock.AccountRepositoryStubSuccess{}),
			expectedError: "Error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer amount not have sufficient",
			args: args{
				input: TransferInput{
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               200,
				},
			},
			usecase:       NewTransfer(mock.TransferRepositoryStubSuccess{}, mock.AccountRepositoryStubSuccess{}),
			expectedError: "origin account does not have sufficient balance",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error find account",
			args: args{
				input: TransferInput{
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
					Amount:               20,
				},
			},
			usecase:       NewTransfer(mock.TransferRepositoryStubSuccess{}, mock.AccountRepositoryStubError{}),
			expectedError: "Error",
			expected:      TransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.usecase.Store(tt.args.input)

			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}
			fmt.Printf("\nGOT - %+v\n", got)
			fmt.Printf("\nEXP - %+v\n", tt.expected)

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
			usecase: NewTransfer(mock.TransferRepositoryStubSuccess{}, mock.AccountRepositoryStubSuccess{}),
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
			usecase:       NewTransfer(mock.TransferRepositoryStubError{}, mock.AccountRepositoryStubSuccess{}),
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
