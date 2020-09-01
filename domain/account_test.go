package domain

import (
	"reflect"
	"testing"
	"time"
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

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		ID        AccountID
		name      string
		CPF       string
		balance   Money
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

func TestAccount_Fields(t *testing.T) {
	type fields struct {
		id        AccountID
		name      string
		cpf       string
		balance   Money
		createdAt time.Time
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Get fields",
			fields: fields{
				id:        "",
				name:      "",
				cpf:       "",
				balance:   0,
				createdAt: time.Time{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount(
				tt.fields.id,
				tt.fields.name,
				tt.fields.cpf,
				tt.fields.balance,
				tt.fields.createdAt,
			)

			if got := account.ID(); got != tt.fields.id {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.id)
			}

			if got := account.Name(); got != tt.fields.name {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.name)
			}

			if got := account.CPF(); got != tt.fields.cpf {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.cpf)
			}

			if got := account.Balance(); got != tt.fields.balance {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.balance)
			}

			if got := account.CreatedAt(); got != tt.fields.createdAt {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.fields.createdAt)
			}
		})
	}
}
