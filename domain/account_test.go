package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestAccount_Deposit(t *testing.T) {
	t.Parallel()

	type args struct {
		amount float64
	}

	tests := []struct {
		name     string
		account  Account
		args     args
		expected float64
	}{
		{
			name: "Successful depositing balance",
			account: Account{
				Balance: 0,
			},
			args: args{
				amount: 10,
			},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.account.Deposit(tt.args.amount)

			if tt.account.Balance != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					tt.account.Balance,
					tt.expected,
				)
			}
		})
	}
}

func TestAccount_Withdraw(t *testing.T) {
	t.Parallel()

	type args struct {
		amount float64
	}

	tests := []struct {
		name        string
		account     Account
		args        args
		expected    float64
		expectedErr error
	}{
		{
			name: "Success in withdrawing balance",
			account: Account{
				Balance: 10,
			},
			args: args{
				amount: 10,
			},
			expected: 0,
		},
		{
			name: "error when withdrawing account balance without sufficient balance",
			account: Account{
				Balance: 0,
			},
			args: args{
				amount: 10,
			},
			expectedErr: ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.account.Withdraw(tt.args.amount); (err != nil) && (err.Error() != tt.expectedErr.Error()) {
				t.Errorf("[TestCase '%s'] ResultError: '%v' | ExpectedError: '%v'",
					tt.name,
					err,
					tt.expectedErr.Error(),
				)
				return
			}

			if tt.account.Balance != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					tt.account.Balance,
					tt.expected,
				)
			}
		})
	}
}

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		ID        string
		name      string
		CPF       string
		balance   float64
		createdAt time.Time
	}

	tests := []struct {
		name     string
		args     args
		expected Account
	}{
		{
			name: "Create Account instance",
			args: args{
				ID:        "",
				name:      "",
				CPF:       "",
				balance:   0,
				createdAt: time.Time{},
			},
			expected: Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAccount(tt.args.ID, tt.args.name, tt.args.CPF, tt.args.balance, tt.args.createdAt)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}
