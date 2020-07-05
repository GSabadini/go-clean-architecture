package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//AccountRepository representa um repositório para manipulação de dados de contas utilizando MongoDB
type AccountRepository struct {
	handler        repository.NoSQLHandler
	collectionName string
}

//NewAccountRepository constrói um repository com suas dependências
func NewAccountRepository(dbHandler repository.NoSQLHandler) AccountRepository {
	return AccountRepository{handler: dbHandler, collectionName: "accounts"}
}

//Store realiza a inserção de uma conta no banco de dados
func (a AccountRepository) Store(account domain.Account) (domain.Account, error) {
	if err := a.handler.Store(a.collectionName, &account); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance realiza a atualização do saldo de uma conta no banco de dados
func (a AccountRepository) UpdateBalance(ID string, balance float64) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.handler.Update(a.collectionName, query, update); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAll realiza a busca de todas as contas no banco de dados
func (a AccountRepository) FindAll() ([]domain.Account, error) {
	var accounts = make([]domain.Account, 0)

	if err := a.handler.FindAll(a.collectionName, nil, &accounts); err != nil {
		return []domain.Account{}, errors.Wrap(err, "error listing accounts")
	}

	return accounts, nil
}

//FindByID realiza a busca de uma conta no banco de dados
func (a AccountRepository) FindByID(ID string) (*domain.Account, error) {
	var (
		account = &domain.Account{}
		query   = bson.M{"id": ID}
	)

	if err := a.handler.FindOne(a.collectionName, query, nil, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return &domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account")
		default:
			return &domain.Account{}, errors.Wrap(err, "error fetching account")
		}
	}

	return account, nil
}

//FindBalance realiza a busca do saldo de uma conta no banco de dados
func (a AccountRepository) FindBalance(ID string) (domain.Account, error) {
	var (
		account  = domain.Account{}
		query    = bson.M{"id": ID}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.handler.FindOne(a.collectionName, query, selector, &account); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account balance")
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account balance")
		}
	}

	return account, nil
}
