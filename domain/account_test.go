package domain

import (
	"testing"
)

func TestAccount_Deposit(t *testing.T) {
	t.Parallel()

	type args struct {
		amount Money
	}

	tests := []struct {
		name     string
		account  Account
		args     args
		expected Money
	}{
		{
			name: "Successful depositing balance",
			args: args{
				amount: 10,
			},
			account:  NewAccountBalance(0),
			expected: 10,
		},
		{
			name: "Successful depositing balance",
			args: args{
				amount: 102098,
			},
			account:  NewAccountBalance(0),
			expected: 102098,
		},
		{
			name: "Successful depositing balance",
			args: args{
				amount: 4498,
			},
			account:  NewAccountBalance(98),
			expected: 4596,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.account.Deposit(tt.args.amount)

			if tt.account.Balance() != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					tt.account.Balance(),
					tt.expected,
				)
			}
		})
	}
}

func TestAccount_Withdraw(t *testing.T) {
	t.Parallel()

	type args struct {
		amount Money
	}

	tests := []struct {
		name        string
		account     Account
		args        args
		expected    Money
		expectedErr error
	}{
		{
			name: "Success in withdrawing balance",
			args: args{
				amount: 10,
			},
			account:  NewAccountBalance(10),
			expected: 0,
		},
		{
			name: "Success in withdrawing balance",
			args: args{
				amount: 10012,
			},
			account:  NewAccountBalance(10013),
			expected: 1,
		},
		{
			name: "Success in withdrawing balance",
			args: args{
				amount: 25,
			},
			account:  NewAccountBalance(125),
			expected: 100,
		},
		{
			name: "error when withdrawing account balance without sufficient balance",
			args: args{
				amount: 564,
			},
			account:     NewAccountBalance(62),
			expectedErr: ErrInsufficientBalance,
		},
		{
			name: "error when withdrawing account balance without sufficient balance",
			args: args{
				amount: 5,
			},
			account:     NewAccountBalance(1),
			expectedErr: ErrInsufficientBalance,
		},
		{
			name: "error when withdrawing account balance without sufficient balance",
			args: args{
				amount: 10,
			},
			account:     NewAccountBalance(0),
			expectedErr: ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.account.Withdraw(tt.args.amount); (err != nil) && (err.Error() != tt.expectedErr.Error()) {
				t.Errorf("[TestCase '%s'] ResultError: '%v' | ExpectedError: '%v'",
					tt.name,
					err,
					tt.expectedErr.Error(),
				)
				return
			}

			if tt.expectedErr == nil && tt.account.Balance() != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					tt.account.Balance(),
					tt.expected,
				)
			}
		})
	}
}
