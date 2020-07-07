package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTransfer(t *testing.T) {
	type args struct {
		ID                   string
		accountOriginID      string
		accountDestinationID string
		amount               float64
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
