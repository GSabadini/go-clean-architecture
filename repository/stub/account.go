package stub

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"

	"github.com/pkg/errors"
)

//AccountRepositoryStubSuccess implementa a interface de AccountRepository com resultados de sucesso
type AccountRepositoryStubSuccess struct{}

//Store cria uma conta pré-definida
func (a AccountRepositoryStubSuccess) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{
		ID:        "5e570851adcef50116aa7a5c",
		Name:      "Test",
		CPF:       "02815517078",
		Balance:   100,
		CreatedAt: nil,
	}, nil
}

//Update retorna sucesso ao atualizar uma conta
func (a AccountRepositoryStubSuccess) Update(_ bson.M, _ bson.M) error {
	return nil
}

//FindAll retorna uma lista de contas pré-definidas
func (a AccountRepositoryStubSuccess) FindAll() ([]domain.Account, error) {
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
func (a AccountRepositoryStubSuccess) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{
		ID:      "5e570854adcef50116aa7a5d",
		Name:    "Test",
		CPF:     "02815517078",
		Balance: 50.25,
	}, nil
}

//FindOneWithSelector retorna apenas o saldo da conta
func (a AccountRepositoryStubSuccess) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{
		Balance: 100.00,
	}, nil
}

//AccountRepositoryStubError implementa a interface de AccountRepository com resultados de erro
type AccountRepositoryStubError struct{}

//Store retorna um error ao criar uma conta
func (a AccountRepositoryStubError) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}

//Update retorna um error ao atualizar uma conta
func (a AccountRepositoryStubError) Update(_ bson.M, _ bson.M) error {
	return errors.New("Error")
}

//FindAll retorna um error ao atualizar uma conta
func (a AccountRepositoryStubError) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Error")
}

//FindOne retorna um error ao buscar uma conta
func (a AccountRepositoryStubError) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{}, errors.New("Error")
}

//FindOneWithSelector retorna um error ao buscar uma conta com campo específico
func (a AccountRepositoryStubError) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}
