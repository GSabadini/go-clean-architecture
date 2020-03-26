package usecase

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//Account armazena as depedências para ações de uma conta
type Account struct {
	repository repository.AccountRepository
}

//NewAccount cria uma conta com suas dependências
func NewAccount(repository repository.AccountRepository) Account {
	return Account{repository: repository}
}

//Store cria uma nova conta
func (a Account) Store(account domain.Account) (domain.Account, error) {
	t := time.Now()
	account.CreatedAt = &t
	account.ID = bson.NewObjectId()
	account.CPF = a.cleanCPF(account.CPF)

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
	var (
		query    = bson.M{"_id": bson.ObjectIdHex(ID)}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	result, err := a.repository.FindOneWithSelector(query, selector)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (a Account) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}
