package usecase

import (
	"context"
	"errors"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type mockTransferRepoFindAll struct {
	domain.TransferRepository

	result []domain.Transfer
	err    error
}

func (m mockTransferRepoFindAll) FindAll(_ context.Context) ([]domain.Transfer, error) {
	return m.result, m.err
}

type mockTransferPresenterFindAll struct {
	output.TransferPresenter

	result []output.TransferOutput
}

func (m mockTransferPresenterFindAll) OutputList(_ []domain.Transfer) []output.TransferOutput {
	return m.result
}

func TestTransferFindAllInteractor_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expected      []output.TransferOutput
		transferRepo  domain.TransferRepository
		presenter     output.TransferPresenter
		expectedError string
	}{
		{
			name: "Success when returning the transfer list",
			transferRepo: mockTransferRepoFindAll{
				result: []domain.Transfer{
					domain.NewTransfer(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						100,
						time.Time{},
					),
					domain.NewTransfer(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						500,
						time.Time{},
					),
				},
				err: nil,
			},
			presenter: mockTransferPresenterFindAll{
				result: []output.TransferOutput{
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
						Amount:               1,
						CreatedAt:            time.Time{},
					},
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
						Amount:               5,
						CreatedAt:            time.Time{},
					},
				},
			},
			expected: []output.TransferOutput{
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               1,
					CreatedAt:            time.Time{},
				},
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               5,
					CreatedAt:            time.Time{},
				},
			},
		},
		{
			name: "Success when returning the empty transfer list",
			transferRepo: mockTransferRepoFindAll{
				result: []domain.Transfer{},
				err:    nil,
			},
			presenter: mockTransferPresenterFindAll{
				result: []output.TransferOutput{},
			},
			expected: []output.TransferOutput{},
		},
		{
			name: "Error when returning the transfer list",
			transferRepo: mockTransferRepoFindAll{
				result: []domain.Transfer{},
				err:    errors.New("error"),
			},
			presenter: mockTransferPresenterFindAll{
				result: []output.TransferOutput{},
			},
			expected:      []output.TransferOutput{},
			expectedError: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewFindAllTransferInteractor(tt.transferRepo, tt.presenter, time.Second)

			result, err := uc.Execute(context.Background())
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
