package usecase

import (
	"errors"
	"time"
)

//TransferUseCaseStubSuccess implementa a interface de TransferUseCase com resultados de sucesso
type TransferUseCaseStubSuccess struct{}

//Store
func (t TransferUseCaseStubSuccess) Store(accountOriginID, accountDestinationID string, amount float64) (transferOutput, error) {
	return transferOutput{
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll
func (t TransferUseCaseStubSuccess) FindAll() ([]transferOutput, error) {
	return []transferOutput{}, nil
}

//TransferUseCaseStubError implementa a interface de TransferUseCase com resultados de erro
type TransferUseCaseStubError struct {
	TypeErr error
}

//Store
func (t TransferUseCaseStubError) Store(_, _ string, _ float64) (transferOutput, error) {
	var err = errors.New("Error")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return transferOutput{}, err
}

//FindAll
func (t TransferUseCaseStubError) FindAll() ([]transferOutput, error) {
	return []transferOutput{}, errors.New("Error")
}
