package usecase

import (
	"errors"
	"github.com/gsabadini/go-bank-transfer/domain"
	"reflect"
	"testing"
	"time"
)

type mockTransferRepoStore struct {
	domain.TransferRepository

	result domain.Transfer
	err    error
}

func (m mockTransferRepoStore) Store(_ domain.Transfer) (domain.Transfer, error) {
	return m.result, m.err
}

type mockAccountRepo struct {
	domain.AccountRepository

	updateBalanceFake func() error

	findByIDOriginFake      func() (*domain.Account, error)
	findByIDDestinationFake func() (*domain.Account, error)
	findByIDOriginInvoked   bool
}

func (m mockAccountRepo) UpdateBalance(_ string, _ float64) error {
	return m.updateBalanceFake()
}

func (m mockAccountRepo) FindByID(_ string) (*domain.Account, error) {
	if m.findByIDOriginInvoked {
		return m.findByIDDestinationFake()
	}

	m.findByIDOriginInvoked = true
	return m.findByIDOriginFake()
}

func TestTransfer_Store(t *testing.T) {
	t.Parallel()

	type args struct {
		accountOriginID      string
		accountDestinationID string
		amount               float64
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
				amount:               20,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               20,
					CreatedAt:            time.Time{},
				},
				err: nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "08098565895",
						Balance:   50,
						CreatedAt: time.Time{},
					}, nil
				},
				findByIDDestinationFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04682",
						Name:      "Test2",
						CPF:       "13098565491",
						Balance:   30,
						CreatedAt: time.Time{},
					}, nil
				},
			},
			expected: TransferOutput{
				ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               20,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "Create transfer generic error transfer repository",
			args: args{
				accountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
				accountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04681",
				amount:               20,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    errors.New("error"),
			},
			accountRepo: mockAccountRepo{
				AccountRepository: nil,
				updateBalanceFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "08098565895",
						Balance:   50,
						CreatedAt: time.Time{},
					}, nil
				},
				findByIDDestinationFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04682",
						Name:      "Test2",
						CPF:       "13098565491",
						Balance:   30,
						CreatedAt: time.Time{},
					}, nil
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
				amount:               20,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				AccountRepository: nil,
				updateBalanceFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (*domain.Account, error) {
					return &domain.Account{}, errors.New("error")
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
				amount:               20,
			},
			transferRepo: mockTransferRepoStore{
				result: domain.Transfer{},
				err:    nil,
			},
			accountRepo: mockAccountRepo{
				updateBalanceFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "08098565895",
						Balance:   50,
						CreatedAt: time.Time{},
					}, nil
				},
				findByIDDestinationFake: func() (*domain.Account, error) {
					return &domain.Account{}, errors.New("error")
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
				updateBalanceFake: func() error {
					return nil
				},
				findByIDOriginFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04681",
						Name:      "Test",
						CPF:       "08098565895",
						Balance:   0,
						CreatedAt: time.Time{},
					}, nil
				},
				findByIDDestinationFake: func() (*domain.Account, error) {
					return &domain.Account{
						ID:        "3c096a40-ccba-4b58-93ed-57379ab04682",
						Name:      "Test2",
						CPF:       "13098565491",
						Balance:   0,
						CreatedAt: time.Time{},
					}, nil
				},
			},
			expectedError: "origin account does not have sufficient balance",
			expected:      TransferOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewTransfer(tt.transferRepo, tt.accountRepo)
			got, err := uc.Store(tt.args.accountOriginID, tt.args.accountDestinationID, tt.args.amount)

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

func (m mockTransferRepoFindAll) FindAll() ([]domain.Transfer, error) {
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
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
						Amount:               100,
						CreatedAt:            time.Time{},
					},
					{
						ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
						AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
						AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
						Amount:               500,
						CreatedAt:            time.Time{},
					},
				},
				err: nil,
			},
			accountRepo: mockAccountRepo{},
			expected: []TransferOutput{
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               100,
					CreatedAt:            time.Time{},
				},
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               500,
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
			var uc = NewTransfer(tt.transferRepo, tt.accountRepo)
			result, err := uc.FindAll()

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
