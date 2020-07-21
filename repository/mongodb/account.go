package mongodb

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//accountBSON armazena a estrutura de dados do MongoDB
type accountBSON struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	CPF       string    `bson:"cpf"`
	Balance   float64   `bson:"balance"`
	CreatedAt time.Time `bson:"created_at"`
}

//AccountRepository armazena a estrutura de dados de um repositório de Account
type AccountRepository struct {
	handler        repository.NoSQLHandler
	collectionName string
}

//NewAccountRepository constrói um repository com suas dependências
func NewAccountRepository(dbHandler repository.NoSQLHandler) AccountRepository {
	return AccountRepository{handler: dbHandler, collectionName: "accounts"}
}

//Store insere uma Account no database
func (a AccountRepository) Store(account domain.Account) (domain.Account, error) {
	var accountBSON = &accountBSON{
		ID:        string(account.ID),
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
	}

	if err := a.handler.Store(a.collectionName, accountBSON); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance atualiza o Balance de uma Account no database
func (a AccountRepository) UpdateBalance(ID domain.AccountID, balance float64) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.handler.Update(a.collectionName, query, update); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAll busca todas as Account no database
func (a AccountRepository) FindAll() ([]domain.Account, error) {
	var (
		accountsBson = make([]accountBSON, 0)
		accounts     = make([]domain.Account, 0)
	)

	if err := a.handler.FindAll(a.collectionName, nil, &accountsBson); err != nil {
		return accounts, errors.Wrap(err, "error listing accounts")
	}

	for _, accountBSON := range accountsBson {
		var account = domain.Account{
			ID:        domain.AccountID(accountBSON.ID),
			Name:      accountBSON.Name,
			CPF:       accountBSON.CPF,
			Balance:   accountBSON.Balance,
			CreatedAt: accountBSON.CreatedAt,
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

//FindByID busca uma Account por ID no database
func (a AccountRepository) FindByID(ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
	)

	if err := a.handler.FindOne(a.collectionName, query, nil, accountBSON); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account")
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account")
		}
	}

	return domain.Account{
		ID:        domain.AccountID(accountBSON.ID),
		Name:      accountBSON.Name,
		CPF:       accountBSON.CPF,
		Balance:   accountBSON.Balance,
		CreatedAt: accountBSON.CreatedAt,
	}, nil
}

//FindBalance busca o Balance de uma Account no database
func (a AccountRepository) FindBalance(ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
		selector    = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.handler.FindOne(a.collectionName, query, selector, accountBSON); err != nil {
		switch err {
		case mgo.ErrNotFound:
			return domain.Account{}, errors.Wrap(domain.ErrNotFound, "error fetching account balance")
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account balance")
		}
	}

	return domain.Account{
		ID:        domain.AccountID(accountBSON.ID),
		Name:      accountBSON.Name,
		CPF:       accountBSON.CPF,
		Balance:   accountBSON.Balance,
		CreatedAt: accountBSON.CreatedAt,
	}, nil
}
