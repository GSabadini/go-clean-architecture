package usecase

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//StoreAccount cria uma nova conta
func StoreAccount(repository repository.AccountRepository, account domain.Account) (domain.Account, error) {
	t := time.Now()
	account.CreatedAt = &t
	account.ID = bson.NewObjectId()
	account.CPF = cleanCPF(account.CPF)

	result, err := repository.Store(account)
	if err != nil {
		return result, err
	}

	return result, nil
}

func cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}

//FindAllAccount retorna uma lista de contas
func FindAllAccount(repository repository.AccountRepository) ([]domain.Account, error) {
	result, err := repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindBalanceAccount retorna o saldo de uma conta
func FindBalanceAccount(repository repository.AccountRepository, ID string) (domain.Account, error) {
	var (
		query    = bson.M{"_id": bson.ObjectIdHex(ID)}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	result, err := repository.FindOneWithSelector(query, selector)
	if err != nil {
		return result, err
	}

	return result, nil
}
