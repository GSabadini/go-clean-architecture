package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type mockAccountRepoFindBalance struct {
	domain.AccountRepository

	result domain.Account
	err    error
}

func (m mockAccountRepoFindBalance) FindBalance(_ context.Context, _ domain.AccountID) (domain.Account, error) {
	return m.result, m.err
}

type mockAccountPresenterFindBalance struct {
	output.AccountPresenter

	result output.AccountBalance
}

func (m mockAccountPresenterFindBalance) OutputBalance(_ domain.Money) output.AccountBalance {
	return m.result
}

func TestFindBalanceAccountInteractor_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		ID domain.AccountID
	}

	tests := []struct {
		name          string
		args          args
		repository    domain.AccountRepository
		presenter     output.AccountPresenter
		expected      output.AccountBalance
		expectedError interface{}
	}{
		{
			name: "Success when returning the account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			repository: mockAccountRepoFindBalance{
				result: domain.NewAccountBalance(100),
				err:    nil,
			},
			presenter: mockAccountPresenterFindBalance{
				result: output.AccountBalance{Balance: 1},
			},
			expected: output.AccountBalance{Balance: 1},
		},
		{
			name: "Success when returning the account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			repository: mockAccountRepoFindBalance{
				result: domain.NewAccountBalance(20050),
				err:    nil,
			},
			presenter: mockAccountPresenterFindBalance{
				result: output.AccountBalance{Balance: 200.5},
			},
			expected: output.AccountBalance{Balance: 200.5},
		},
		{
			name: "Error returning account balance",
			args: args{
				ID: "3c096a40-ccba-4b58-93ed-57379ab04680",
			},
			repository: mockAccountRepoFindBalance{
				result: domain.Account{},
				err:    errors.New("error"),
			},
			presenter: mockAccountPresenterFindBalance{
				result: output.AccountBalance{},
			},
			expectedError: "error",
			expected:      output.AccountBalance{},
		},
	}

	for _, tt := range tests {
		var uc = NewFindBalanceAccountInteractor(tt.repository, tt.presenter, time.Second)

		result, err := uc.Execute(context.Background(), tt.args.ID)
		if (err != nil) && (err.Error() != tt.expectedError) {
			t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			return
		}

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
		}
	}
}
