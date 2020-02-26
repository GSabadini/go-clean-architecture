package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//StoreAccount cria uma nova conta
func StoreAccount(repository repository.AccountRepository, account *domain.Account) (*domain.Account, error) {
	result, err := repository.Store(account)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//FindAllAccount retorna uma lista de contas
func FindAllAccount(repository repository.AccountRepository) ([]domain.Account, error) {
	result, err := repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

type AccountBalance struct {
	Balance float64 `json:"balance"`
}

//FindBalanceAccount retorna o saldo de uma conta
func FindBalanceAccount(repository repository.AccountRepository, id string) (*AccountBalance, error) {
	var (
		query          = bson.M{"_id": bson.ObjectIdHex(id)}
		accountBalance = &AccountBalance{}
	)

	result, err := repository.FindOne(query)
	if err != nil {
		return accountBalance, err
	}

	accountBalance.Balance = result.Balance

	return accountBalance, nil
}
