package repository

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"

	"github.com/pkg/errors"
)

//AccountRepositoryMockSuccess implementa a interface de AccountRepository com resultados de sucesso
type AccountRepositoryMockSuccess struct{}

//Store cria uma conta pré-definida
func (a AccountRepositoryMockSuccess) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{
		ID:        "5e570851adcef50116aa7a5c",
		Name:      "Test",
		CPF:       "02815517078",
		Balance:   100,
		CreatedAt: nil,
	}, nil
}

//Update retorna sucesso ao atualizar uma conta
func (a AccountRepositoryMockSuccess) Update(_ bson.M, _ bson.M) error {
	return nil
}

//FindAll retorna uma lista de contas pré-definidas
func (a AccountRepositoryMockSuccess) FindAll() ([]domain.Account, error) {
	var account = []domain.Account{
		{
			ID:      "5e570851adcef50116aa7a5c",
			Name:    "Test-0",
			CPF:     "02815517078",
			Balance: 0,
		},
		{
			ID:      "5e570854adcef50116aa7a5d",
			Name:    "Test-1",
			CPF:     "02815517078",
			Balance: 50.25,
		},
	}

	return account, nil
}

//FindOne retorna uma conta pré-definida
func (a AccountRepositoryMockSuccess) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{
		ID:      "5e570854adcef50116aa7a5d",
		Name:    "Test",
		CPF:     "02815517078",
		Balance: 50.25,
	}, nil
}

//FindOneWithSelector retorna apenas o saldo da conta
func (a AccountRepositoryMockSuccess) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{
		Balance: 100.00,
	}, nil
}

//AccountRepositoryMockError implementa a interface de AccountRepository com resultados de erro
type AccountRepositoryMockError struct{}

//Store retorna um error ao criar uma conta
func (a AccountRepositoryMockError) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}

//Update retorna um error ao atualizar uma conta
func (a AccountRepositoryMockError) Update(_ bson.M, _ bson.M) error {
	return errors.New("Error")
}

//FindAll retorna um error ao atualizar uma conta
func (a AccountRepositoryMockError) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Error")
}

//FindOne retorna um error ao buscar uma conta
func (a AccountRepositoryMockError) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{}, errors.New("Error")
}

//FindOneWithSelector retorna um error ao buscar uma conta com campo específico
func (a AccountRepositoryMockError) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}
