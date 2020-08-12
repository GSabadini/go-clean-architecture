package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type mockTransferRepoStore struct {
	domain.TransferRepository

	result domain.Transfer
	err    error
}

func (m mockTransferRepoStore) Store(_ context.Context, _ domain.Transfer) (domain.Transfer, error) {
	return m.result, m.err
}

type invoked struct {
	call bool
}

type mockAccountRepo struct {
	domain.AccountRepository

	updateBalanceOriginFake      func() error
	updateBalanceDestinationFake func() error
	invokedUpdate                *invoked

	findByIDOriginFake      func() (domain.Account, error)
	findByIDDestinationFake func() (domain.Account, error)
	invokedFind             *invoked
}

func (m mockAccountRepo) UpdateBalance(_ context.Context, _ domain.AccountID, _ domain.Money) error {
	if m.invokedUpdate != nil && m.invokedUpdate.call {
		return m.updateBalanceDestinationFake()
	}

	if m.invokedUpdate != nil {
		m.invokedUpdate.call = true
	}
	return m.updateBalanceOriginFake()
}

func (m mockAccountRepo) FindByID(_ context.Context, _ domain.AccountID) (domain.Account, error) {
	if m.invokedFind != nil && m.invokedFind.call {
		return m.findByIDDestinationFake()
	}

	if m.invokedFind != nil {
		m.invokedFind.call = true
	}
	return m.findByIDOriginFake()
}

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		accountOriginID      domain.AccountID
		accountDestinationID domain.AccountID
		amount               domain.Money
	}

	tests := []struct {
		name          string
		args          args
		transferRepo  domain.TransferRepository
		accountRepo   domain.AccountRepository
		expected      TransferOutput
		expectedError string
	}{
		{
			name: "Create transfer successful",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				amount:               2999,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.NewTransfer(
					"3c096a40-ccba-4b58-93ed-57379ab04680",
					"3c096a40-ccba-4b58-93ed-57379ab04681",
					"3c096a40-ccba-4b58-93ed-57379ab04682",
					2999,
					time.Time{},
				),
				err: nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						5000,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Test2",
						"13098565491",
						3000,
						time.Time{},
					), nil
				},
			},
			expected: TransferOutput{
				ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               29.99,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "Create transfer generic error transfer repository",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               200,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    errors.New("error"),
			},
			accountRepo: mockAccountRepo{
				AccountRepository: nil,
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						1000,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Test2",
						"13098565491",
						3000,
						time.Time{},
					), nil
				},
			},
			expectedError: "error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error find origin account",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               1999,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.Account{}, errors.New("error")
				},
			},
			expectedError: "error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error find destination account",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               100,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: &mockAccountRepo{
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						5000,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.Account{}, errors.New("error")
				},
				invokedFind: &invoked{call: false},
			},
			expectedError: "error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error update origin account",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               250,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceOriginFake: func() error {
					return errors.New("error")
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						5999,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Test2",
						"13098565491",
						2999,
						time.Time{},
					), nil
				},
			},
			expectedError: "error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer error update destination account",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               100,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return errors.New("error")
				},
				invokedUpdate: &invoked{call: false},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						200,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Test2",
						"13098565491",
						100,
						time.Time{},
					), nil
				},
			},
			expectedError: "error",
			expected:      TransferOutput{},
		},
		{
			name: "Create transfer amount not have sufficient",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               200,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				AccountRepository: nil,
				updateBalanceOriginFake: func() error {
					return nil
				},
				updateBalanceDestinationFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"Test",
						"08098565895",
						0,
						time.Time{},
					), nil
				},
				findByIDDestinationFake: func() (domain.Account, error) {
					return domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Test2",
						"13098565491",
						0,
						time.Time{},
					), nil
				},
			},
			expectedError: "origin account does not have sufficient balance",
			expected:      TransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewTransfer(tt.transferRepo, tt.accountRepo, time.Second)

			got, err := uc.Store(
				context.Background(),
				tt.args.accountOriginID,
				tt.args.accountDestinationID,
				tt.args.amount,
			)
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

type mockTransferRepoFindAll struct {
	domain.TransferRepository

	result []domain.Transfer
	err    error
}

func (m mockTransferRepoFindAll) FindAll(_ context.Context) ([]domain.Transfer, error) {
	return m.result, m.err
}

func TestTransfer_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expected      []TransferOutput
		transferRepo  domain.TransferRepository
		accountRepo   domain.AccountRepository
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
			accountRepo: mockAccountRepo{},
			expected: []TransferOutput{
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
			accountRepo: mockAccountRepo{},
			expected:    []TransferOutput{},
		},
		{
			name: "Error when returning the transfer list",
			transferRepo: mockTransferRepoFindAll{
				result: []domain.Transfer{},
				err:    errors.New("error"),
			},
			accountRepo:   mockAccountRepo{},
			expected:      []TransferOutput{},
			expectedError: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewTransfer(tt.transferRepo, tt.accountRepo, time.Second)

			result, err := uc.FindAll(context.Background())
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
