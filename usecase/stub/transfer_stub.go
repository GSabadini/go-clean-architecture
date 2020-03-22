package stub

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type TransferUseCaseStubSuccess struct{}

func (t TransferUseCaseStubSuccess) FindAllTransfer() ([]domain.Transfer, error) {
	return []domain.Transfer{}, nil
}

func (t TransferUseCaseStubSuccess) StoreTransfer(transfer domain.Transfer) (domain.Transfer, error) {
	return transfer, nil
}

type TransferUseCaseStubError struct {
	TypeErr error
}

func (t TransferUseCaseStubError) FindAllTransfer() ([]domain.Transfer, error) {
	return []domain.Transfer{}, errors.New("Error")
}

func (t TransferUseCaseStubError) StoreTransfer(domain.Transfer) (domain.Transfer, error) {
	var err = errors.New("Error")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return domain.Transfer{}, err
}
