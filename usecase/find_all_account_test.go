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

type mockAccountRepoFindAll struct {
	domain.AccountRepository

	result []domain.Account
	err    error
}

func (m mockAccountRepoFindAll) FindAll(_ context.Context) ([]domain.Account, error) {
	return m.result, m.err
}

type mockAccountPresenterFindAll struct {
	output.AccountPresenter

	result []output.AccountOutput
}

func (m mockAccountPresenterFindAll) OutputList(_ []domain.Account) []output.AccountOutput {
	return m.result
}

func TestFindAllAccountInteractor_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		repository    domain.AccountRepository
		presenter     output.AccountPresenter
		expected      []output.AccountOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the validator list",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{
					domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"Test",
						"02815517078",
						125,
						time.Time{},
					),
					domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"02815517071",
						99999,
						time.Time{},
					),
				},
				err: nil,
			},
			presenter: mockAccountPresenterFindAll{
				result: []output.AccountOutput{
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
						Name:      "Test",
						CPF:       "02815517078",
						Balance:   1.25,
						CreatedAt: time.Time{},
					},
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "02815517071",
						Balance:   999.99,
						CreatedAt: time.Time{},
					},
				},
			},
			expected: []output.AccountOutput{
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   1.25,
					CreatedAt: time.Time{},
				},
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
					Name:      "Test",
					CPF:       "02815517071",
					Balance:   999.99,
					CreatedAt: time.Time{},
				},
			},
		},
		{
			name: "Success when returning the empty validator list",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    nil,
			},
			presenter: mockAccountPresenterFindAll{
				result: []output.AccountOutput{},
			},
			expected: []output.AccountOutput{},
		},
		{
			name: "Error when returning the list of accounts",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    errors.New("error"),
			},
			presenter: mockAccountPresenterFindAll{
				result: []output.AccountOutput{},
			},
			expectedError: "error",
			expected:      []output.AccountOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewFindAllAccountInteractor(tt.repository, tt.presenter, time.Second)

			result, err := uc.Execute(context.Background())
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
			}
		})
	}
}
