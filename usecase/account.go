package usecase

import (
	"strings"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//Account armazena as depedências para ações de uma conta
type Account struct {
	repository domain.AccountRepository
}

//NewAccount cria uma conta com suas dependências
func NewAccount(repository domain.AccountRepository) Account {
	return Account{repository: repository}
}

//Store cria uma nova conta
func (a Account) Store(data domain.Account) (domain.Account, error) {
	var account = domain.NewAccount(data.Name, a.cleanCPF(data.CPF), data.Balance)

	result, err := a.repository.Store(account)
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindAll retorna uma lista de contas
func (a Account) FindAll() ([]domain.Account, error) {
	result, err := a.repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindBalance retorna o saldo de uma conta
func (a Account) FindBalance(ID string) (domain.Account, error) {
	result, err := a.repository.FindBalance(ID)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (a Account) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}
