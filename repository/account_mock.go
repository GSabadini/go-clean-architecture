package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"time"

	"github.com/pkg/errors"
)

//AccountRepositoryStubSuccess implementa a interface de AccountRepository com resultados de sucesso
type AccountRepositoryStubSuccess struct{}

//Store cria uma conta pré-definida
func (a AccountRepositoryStubSuccess) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{
		ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
		Name:      "Test",
		CPF:       "02815517078",
		Balance:   100,
		CreatedAt: time.Time{},
	}, nil
}

//UpdateBalance retorna sucesso ao atualizar uma conta
func (a AccountRepositoryStubSuccess) UpdateBalance(_ string, _ float64) error {
	return nil
}

//FindAll retorna uma lista de contas pré-definidas
func (a AccountRepositoryStubSuccess) FindAll() ([]domain.Account, error) {
	var account = []domain.Account{
		{
			ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
			Name:    "Test-0",
			CPF:     "02815517078",
			Balance: 0,
		},
		{
			ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
			Name:    "Test-1",
			CPF:     "02815517078",
			Balance: 50.25,
		},
	}

	return account, nil
}

//FindByID retorna uma conta pré-definida
func (a AccountRepositoryStubSuccess) FindByID(_ string) (*domain.Account, error) {
	return &domain.Account{
		ID:      "3c096a40-ccba-4b58-93ed-57379ab04680",
		Name:    "Test",
		CPF:     "02815517078",
		Balance: 50.25,
	}, nil
}

//FindBalance retorna apenas o saldo da conta
func (a AccountRepositoryStubSuccess) FindBalance(_ string) (domain.Account, error) {
	return domain.Account{
		Balance: 100.00,
	}, nil
}

//AccountRepositoryStubError implementa a interface de AccountRepository com resultados de erro
type AccountRepositoryStubError struct{}

//Store retorna um error ao criar uma conta
func (a AccountRepositoryStubError) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Errors")
}

//UpdateBalance retorna um error ao atualizar uma conta
func (a AccountRepositoryStubError) UpdateBalance(_ string, _ float64) error {
	return errors.New("Errors")
}

//FindAll retorna um error ao atualizar uma conta
func (a AccountRepositoryStubError) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Errors")
}

//FindByID retorna um error ao buscar uma conta
func (a AccountRepositoryStubError) FindByID(_ string) (*domain.Account, error) {
	return &domain.Account{}, errors.New("Errors")
}

//FindBalance retorna um error ao buscar uma conta com campo específico
func (a AccountRepositoryStubError) FindBalance(_ string) (domain.Account, error) {
	return domain.Account{}, errors.New("Errors")
}
