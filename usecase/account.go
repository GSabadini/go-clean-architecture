package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//StoreAccount cria uma nova account
func StoreAccount(repository repository.AccountRepository, account *domain.Account) (*domain.Account, error) {
	result, err := repository.Store(account)
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindAllAccount retorna uma lista de accounts
func FindAllAccount(repository repository.AccountRepository) ([]domain.Account, error) {
	result, err := repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindOneAccount retorna uma account com base em um ID
func FindOneAccount(repository repository.AccountRepository, id string) (*domain.Account, error) {
	var query = bson.M{"_id": bson.ObjectIdHex(id)}

	result, err := repository.FindOne(query)
	if err != nil {
		return result, err
	}

	return result, nil
}
