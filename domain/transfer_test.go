package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTransfer(t *testing.T) {
	type args struct {
		ID                   TransferID
		accountOriginID      AccountID
		accountDestinationID AccountID
		amount               Money
		createdAt            time.Time
	}

	tests := []struct {
		name     string
		args     args
		expected Transfer
	}{
		{
			name: "",
			args: args{
				ID:                   "",
				accountOriginID:      "",
				accountDestinationID: "",
				amount:               0,
				createdAt:            time.Time{},
			},
			expected: Transfer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTransfer(
				tt.args.ID,
				tt.args.accountOriginID,
				tt.args.accountDestinationID,
				tt.args.amount,
				tt.args.createdAt,
			)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}

func TestTransfer_Fields(t *testing.T) {
	type fields struct {
		id                   TransferID
		accountOriginID      AccountID
		accountDestinationID AccountID
		amount               Money
		createdAt            time.Time
	}

	tests := []struct {
		name   string
		fields fields
		want   AccountID
	}{
		{
			name: "Get fields",
			fields: fields{
				id:                   "",
				accountOriginID:      "",
				accountDestinationID: "",
				amount:               0,
				createdAt:            time.Time{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := NewTransfer(
				tt.fields.id,
				tt.fields.accountOriginID,
				tt.fields.accountDestinationID,
				tt.fields.amount,
				tt.fields.createdAt,
			)

			if got := transfer.ID(); got != tt.fields.id {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.id)
			}

			if got := transfer.AccountOriginID(); got != tt.fields.accountOriginID {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.accountOriginID)
			}

			if got := transfer.AccountDestinationID(); got != tt.fields.accountDestinationID {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.accountDestinationID)
			}

			if got := transfer.Amount(); got != tt.fields.amount {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.amount)
			}

			if got := transfer.CreatedAt(); got != tt.fields.createdAt {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.createdAt)
			}

		})
	}
}
