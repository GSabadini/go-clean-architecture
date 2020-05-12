package repository

import (
	"reflect"
	"testing"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database/stub"
)

func TestTransferMongoDB_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		transfer domain.Transfer
	}

	tests := []struct {
		name        string
		repository  TransferMongoDB
		args        args
		expected    domain.Transfer
		expectedErr bool
	}{
		{
			name:       "Success to create transfer",
			args:       args{transfer: domain.Transfer{}},
			repository: NewTransferMongoDB(stub.MongoHandlerSuccessStub{}),
			expected:   domain.Transfer{},
		},
		{
			name:        "Error to create transfer",
			args:        args{transfer: domain.Transfer{}},
			repository:  NewTransferMongoDB(stub.MongoHandlerErrorStub{}),
			expected:    domain.Transfer{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repository.Store(tt.args.transfer)

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

func TestTransferMongoDB_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		repository  TransferMongoDB
		expected    []domain.Transfer
		expectedErr bool
	}{
		{
			name:       "Success to find all the transfers",
			repository: NewTransferMongoDB(stub.MongoHandlerSuccessStub{}),
			expected:   []domain.Transfer{},
		},
		{
			name:        "Error to find all the transfers",
			repository:  NewTransferMongoDB(stub.MongoHandlerErrorStub{}),
			expected:    []domain.Transfer{},
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
