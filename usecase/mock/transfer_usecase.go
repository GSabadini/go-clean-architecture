package mock

import (
	"errors"
	"time"

	"github.com/gsabadini/go-bank-transfer/usecase"
)

//TransferUseCaseStubSuccess implementa a interface de TransferUseCase com resultados de sucesso
type TransferUseCaseStubSuccess struct{}

//Store
func (t TransferUseCaseStubSuccess) Store(input usecase.TransferInput) (usecase.TransferResult, error) {
	return usecase.TransferResult{
		ID:                   "",
		AccountOriginID:      input.AccountOriginID,
		AccountDestinationID: input.AccountDestinationID,
		Amount:               input.Amount,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll
func (t TransferUseCaseStubSuccess) FindAll() ([]usecase.TransferResult, error) {
	return []usecase.TransferResult{}, nil
}

//TransferUseCaseStubError implementa a interface de TransferUseCase com resultados de erro
type TransferUseCaseStubError struct {
	TypeErr error
}

//Store
func (t TransferUseCaseStubError) Store(_ usecase.TransferInput) (usecase.TransferResult, error) {
	var err = errors.New("Error")
	if t.TypeErr != nil {
		err = t.TypeErr
	}

	return usecase.TransferResult{}, err
}

//FindAll
func (t TransferUseCaseStubError) FindAll() ([]usecase.TransferResult, error) {
	return []usecase.TransferResult{}, errors.New("Error")
}
