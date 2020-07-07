package usecase

import (
	"errors"
	"time"
)

//TransferUseCaseStubSuccess implementa a interface de TransferUseCase com resultados de sucesso
type TransferUseCaseStubSuccess struct{}

//Store
func (t TransferUseCaseStubSuccess) Store(accountOriginID, accountDestinationID string, amount float64) (TransferOutput, error) {
	return TransferOutput{
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll
func (t TransferUseCaseStubSuccess) FindAll() ([]TransferOutput, error) {
	return []TransferOutput{}, nil
}

//TransferUseCaseStubError implementa a interface de TransferUseCase com resultados de erro
type TransferUseCaseStubError struct {
	TypeErr error
}

//Store
func (t TransferUseCaseStubError) Store(_, _ string, _ float64) (TransferOutput, error) {
	var err = errors.New("Errors")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return TransferOutput{}, err
}

//FindAll
func (t TransferUseCaseStubError) FindAll() ([]TransferOutput, error) {
	return []TransferOutput{}, errors.New("Errors")
}
