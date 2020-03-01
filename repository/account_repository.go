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
func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

var ErrNotFound = errors.New("Not found")

//Store realiza a inserção de uma conta no banco de dados
func (a Account) Store(account domain.Account) (domain.Account, error) {
	if err := a.dbHandler.Store(accountsCollectionName, &account); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//Update realiza a atualização de uma conta no banco de dados
func (a Account) Update(query bson.M, update bson.M) error {
	return a.dbHandler.Update(accountsCollectionName, query, update)
}

//FindAll realiza a busca de todas as contas no banco de dados
func (a Account) FindAll() ([]domain.Account, error) {
	var account = make([]domain.Account, 0)

	if err := a.dbHandler.FindAll(accountsCollectionName, nil, &account); err != nil {
		return account, errors.Wrap(err, "error listing accounts")
	}

	return account, nil
}

//FindOne realiza a busca de uma conta no banco de dados
func (a Account) FindOne(query bson.M) (*domain.Account, error) {
	var account = &domain.Account{}

	if err := a.dbHandler.FindOne(accountsCollectionName, query, nil, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, ErrNotFound
		default:
			return account, errors.Wrap(err, "error fetching account")
		}
	}

	return account, nil
}

//FindOneWithSelector realiza a busca de uma conta com campos específicos no banco de dados
func (a Account) FindOneWithSelector(query bson.M, selector interface{}) (domain.Account, error) {
	var account = domain.Account{}

	if err := a.dbHandler.FindOne(accountsCollectionName, query, selector, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return account, ErrNotFound
		default:
			return account, errors.Wrap(err, "error fetching account")
		}
	}

	return account, nil
}
