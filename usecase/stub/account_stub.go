package stub

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type AccountUseCaseStubSuccess struct{}

func (a AccountUseCaseStubSuccess) StoreAccount(account domain.Account) (domain.Account, error) {
	return account, nil
}

func (a AccountUseCaseStubSuccess) FindAllAccount() ([]domain.Account, error) {
	return []domain.Account{}, nil
}

func (a AccountUseCaseStubSuccess) FindBalanceAccount(_ string) (domain.Account, error) {
	return domain.Account{}, nil
}

type AccountUseCaseStubError struct {
	TypeErr error
}

func (a AccountUseCaseStubError) StoreAccount(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}

func (a AccountUseCaseStubError) FindAllAccount() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Error")
}

func (a AccountUseCaseStubError) FindBalanceAccount(_ string) (domain.Account, error) {
	var err = errors.New("Error")
	if a.TypeErr != nil {
		err = a.TypeErr
	}

	return domain.Account{}, err
}
