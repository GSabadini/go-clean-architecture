package stub

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//AccountUseCaseStubSuccess implementa a interface de AccountUseCase com resultados de sucesso
type AccountUseCaseStubSuccess struct{}

//Store
func (a AccountUseCaseStubSuccess) Store(account domain.Account) (domain.Account, error) {
	return account, nil
}

//FindAll
func (a AccountUseCaseStubSuccess) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, nil
}

//FindBalance
func (a AccountUseCaseStubSuccess) FindBalance(_ string) (domain.Account, error) {
	return domain.Account{}, nil
}

//AccountUseCaseStubSuccess implementa a interface de AccountUseCase com resultados de erro
type AccountUseCaseStubError struct {
	TypeErr error
}

//Store
func (a AccountUseCaseStubError) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}

//FindAll
func (a AccountUseCaseStubError) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Error")
}

//FindBalance
func (a AccountUseCaseStubError) FindBalance(_ string) (domain.Account, error) {
	var err = errors.New("Error")
	if a.TypeErr != nil {
		err = a.TypeErr
	}

	return domain.Account{}, err
}
