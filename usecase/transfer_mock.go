package usecase

import (
	"errors"
	"time"
)

//TransferUseCaseStubSuccess implementa a interface de TransferUseCase com resultados de sucesso
type TransferUseCaseStubSuccess struct{}

//Store
func (t TransferUseCaseStubSuccess) Store(input TransferInput) (TransferOutput, error) {
	return TransferOutput{
		ID:                   "",
		AccountOriginID:      input.AccountOriginID,
		AccountDestinationID: input.AccountDestinationID,
		Amount:               input.Amount,
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
func (t TransferUseCaseStubError) Store(_ TransferInput) (TransferOutput, error) {
	var err = errors.New("Error")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return TransferOutput{}, err
}

//FindAll
func (t TransferUseCaseStubError) FindAll() ([]TransferOutput, error) {
	return []TransferOutput{}, errors.New("Error")
}
