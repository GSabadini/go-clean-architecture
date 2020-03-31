package stub

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferUseCaseStubSuccess implementa a interface de TransferUseCase com resultados de sucesso
type TransferUseCaseStubSuccess struct{}

//Store
func (t TransferUseCaseStubSuccess) Store(transfer domain.Transfer) (domain.Transfer, error) {
	return transfer, nil
}

//FindAll
func (t TransferUseCaseStubSuccess) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{}, nil
}

//TransferUseCaseStubError implementa a interface de TransferUseCase com resultados de erro
type TransferUseCaseStubError struct {
	TypeErr error
}

//Store
func (t TransferUseCaseStubError) Store(domain.Transfer) (domain.Transfer, error) {
	var err = errors.New("Error")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return domain.Transfer{}, err
}

//FindAll
func (t TransferUseCaseStubError) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{}, errors.New("Error")
}
