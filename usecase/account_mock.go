package usecase

import (
	"errors"
	"time"
)

//AccountUseCaseStubSuccess implementa a interface de AccountUseCase com resultados de sucesso
type AccountUseCaseStubSuccess struct{}

//Store
func (a AccountUseCaseStubSuccess) Store(name, CPF string, balance float64) (AccountOutput, error) {
	return AccountOutput{
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: time.Time{},
	}, nil
}

//FindAll
func (a AccountUseCaseStubSuccess) FindAll() ([]AccountOutput, error) {
	return []AccountOutput{}, nil
}

//FindBalance
func (a AccountUseCaseStubSuccess) FindBalance(_ string) (AccountBalanceOutput, error) {
	return AccountBalanceOutput{}, nil
}

//AccountUseCaseStubSuccess implementa a interface de AccountUseCase com resultados de erro
type AccountUseCaseStubError struct {
	TypeErr error
}

//Store
func (a AccountUseCaseStubError) Store(_, _ string, _ float64) (AccountOutput, error) {
	return AccountOutput{}, errors.New("Error")
}

//FindAll
func (a AccountUseCaseStubError) FindAll() ([]AccountOutput, error) {
	return []AccountOutput{}, errors.New("Error")
}

//FindBalance
func (a AccountUseCaseStubError) FindBalance(_ string) (AccountBalanceOutput, error) {
	var err = errors.New("Error")
	if a.TypeErr != nil {
		err = a.TypeErr
	}

	return AccountBalanceOutput{}, err
}
