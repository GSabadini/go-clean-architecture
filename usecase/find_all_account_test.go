package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type mockAccountRepoFindAll struct {
	domain.AccountRepository

	result []domain.Account
	err    error
}

func (m mockAccountRepoFindAll) FindAll(_ context.Context) ([]domain.Account, error) {
	return m.result, m.err
}

type mockFindAllAccountPresenter struct {
	result []FindAllAccountOutput
}

func (m mockFindAllAccountPresenter) Output(_ []domain.Account) []FindAllAccountOutput {
	return m.result
}

func TestFindAllAccountInteractor_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		repository    domain.AccountRepository
		presenter     FindAllAccountPresenter
		expected      []FindAllAccountOutput
		expectedError interface{}
	}{
		{
			name: "Success when returning the account list",
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
			presenter: mockFindAllAccountPresenter{
				result: []FindAllAccountOutput{
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
						Name:      "Test",
						CPF:       "02815517078",
						Balance:   1.25,
						CreatedAt: time.Time{}.String(),
					},
					{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "02815517071",
						Balance:   999.99,
						CreatedAt: time.Time{}.String(),
					},
				},
			},
			expected: []FindAllAccountOutput{
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Test",
					CPF:       "02815517078",
					Balance:   1.25,
					CreatedAt: time.Time{}.String(),
				},
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
					Name:      "Test",
					CPF:       "02815517071",
					Balance:   999.99,
					CreatedAt: time.Time{}.String(),
				},
			},
		},
		{
			name: "Success when returning the empty account list",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    nil,
			},
			presenter: mockFindAllAccountPresenter{
				result: []FindAllAccountOutput{},
			},
			expected: []FindAllAccountOutput{},
		},
		{
			name: "Error when returning the list of accounts",
			repository: mockAccountRepoFindAll{
				result: []domain.Account{},
				err:    errors.New("error"),
			},
			presenter: mockFindAllAccountPresenter{
				result: []FindAllAccountOutput{},
			},
			expectedError: "error",
			expected:      []FindAllAccountOutput{},
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
