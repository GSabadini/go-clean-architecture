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

func (m mockTransferRepoStore) Create(_ context.Context, _ domain.Transfer) (domain.Transfer, error) {
	return m.result, m.err
}

func (m mockTransferRepoStore) WithTransaction(_ context.Context, fn func(context.Context) error) error {
	if err := fn(context.Background()); err != nil {
		return err
	}

	return nil
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

type mockCreateTransferPresenter struct {
	result CreateTransferOutput
}

func (m mockCreateTransferPresenter) Output(_ domain.Transfer) CreateTransferOutput {
	return m.result
}

func TestTransferCreateInteractor_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		input CreateTransferInput
	}

	tests := []struct {
		name          string
		args          args
		transferRepo  domain.TransferRepository
		accountRepo   domain.AccountRepository
		presenter     CreateTransferPresenter
		expected      CreateTransferOutput
		expectedError string
	}{
		{
			name: "Create transfer successful",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               2999,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               29.99,
					CreatedAt:            time.Time{}.String(),
				},
			},
			expected: CreateTransferOutput{
				ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               29.99,
				CreatedAt:            time.Time{}.String(),
			},
		},
		{
			name: "Create transfer generic error transfer gateway",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               200,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "error",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error find origin account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               1999,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "error",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error not found find origin account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               1999,
			}},
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
					return domain.Account{}, domain.ErrAccountOriginNotFound
				},
			},
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "account origin not found",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error find destination account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               100,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "error",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error not found find destination account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               100,
			}},
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
					return domain.Account{}, domain.ErrAccountDestinationNotFound
				},
				invokedFind: &invoked{call: false},
			},
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "account destination not found",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error update origin account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               250,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "error",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer error update destination account",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               100,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "error",
			expected:      CreateTransferOutput{},
		},
		{
			name: "Create transfer amount not have sufficient",
			args: args{input: CreateTransferInput{
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				Amount:               200,
			}},
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
			presenter: mockCreateTransferPresenter{
				result: CreateTransferOutput{},
			},
			expectedError: "origin account does not have sufficient balance",
			expected:      CreateTransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateTransferInteractor(tt.transferRepo, tt.accountRepo, tt.presenter, time.Second)

			got, err := uc.Execute(context.Background(), tt.args.input)
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
