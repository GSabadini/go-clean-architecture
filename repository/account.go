package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

const accountsCollectionName = "accounts"

//Account representa um repositório para dados de uma conta
type Account DbRepository

//NewAccount constrói um repository com suas dependências
func NewAccount(dbHandler database.DbHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store realiza a inserção de uma conta no banco de dados
func (a Account) Store(account domain.Account) (domain.Account, error) {
	if err := a.dbHandler.Store(accountsCollectionName, &account); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance realiza a atualização do saldo de uma conta no banco de dados
func (a Account) UpdateBalance(ID string, balance float64) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.dbHandler.Update(accountsCollectionName, query, update); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAll realiza a busca de todas as contas no banco de dados
func (a Account) FindAll() ([]domain.Account, error) {
	var accounts = make([]domain.Account, 0)

	if err := a.dbHandler.FindAll(accountsCollectionName, nil, &accounts); err != nil {
		return accounts, errors.Wrap(err, "error listing accounts")
	}

	return accounts, nil
}

//FindByID realiza a busca de uma conta no banco de dados
func (a Account) FindByID(ID string) (*domain.Account, error) {
	var (
		account = &domain.Account{}
		query   = bson.M{"id": ID}
	)

	if err := a.dbHandler.FindOne(accountsCollectionName, query, nil, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, errors.Wrap(domain.ErrNotFound, "error fetching account")
		default:
			return account, errors.Wrap(err, "error fetching account")
		}
	}

	return account, nil
}

//FindBalance realiza a busca do saldo de uma conta no banco de dados
func (a Account) FindBalance(ID string) (domain.Account, error) {
	var (
		account  = domain.Account{}
		query    = bson.M{"id": ID}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.dbHandler.FindOne(accountsCollectionName, query, selector, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, errors.Wrap(domain.ErrNotFound, "error fetching account balance")
		default:
			return account, errors.Wrap(err, "error fetching account balance")
		}
	}

	return account, nil
}
