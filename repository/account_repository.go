package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

const accountsCollectionName = "accounts"

//Account representa um repositório para dados de uma conta
type Account DbRepository

//NewAccount constrói um repository com suas dependências
func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store realiza uma inserção no banco de dados através da implementação real do database
func (a Account) Store(account domain.Account) (domain.Account, error) {
	t := time.Now()
	account.CreatedAt = &t
	account.ID = bson.NewObjectId()

	if err := a.dbHandler.Store(accountsCollectionName, &account); err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

//Update realiza uma atualização no banco de dados através da implementação real do database
func (a Account) Update(query bson.M, update bson.M) error {
	return a.dbHandler.Update(accountsCollectionName, query, update)
}

//FindAll realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindAll() ([]domain.Account, error) {
	var account = make([]domain.Account, 0)

	if err := a.dbHandler.FindAll(accountsCollectionName, nil, &account); err != nil {
		return []domain.Account{}, err
	}

	return account, nil
}

//FindOne realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindOne(query bson.M, selector interface{}) (*domain.Account, error) {
	var account = &domain.Account{}

	if err := a.dbHandler.FindOne(accountsCollectionName, query, selector, &account); err != nil {
		return &domain.Account{}, err
	}

	return account, nil
}

type AccountRepositoryMock struct{}

func (a AccountRepositoryMock) Store(account domain.Account) (domain.Account, error) {
	return domain.Account{
		ID:        "1",
		Name:      "1",
		CPF:       "1",
		Balance:   100,
		CreatedAt: nil,
	}, nil
}

func (a AccountRepositoryMock) Update(query bson.M, update bson.M) error {
	return nil
}

func (a AccountRepositoryMock) FindAll() ([]domain.Account, error) {
	var account = make([]domain.Account, 0)

	return account, nil
}

//FindOne realiza uma busca no banco de dados através da implementação real do database
func (a AccountRepositoryMock) FindOne(query bson.M, selector interface{}) (*domain.Account, error) {
	var account = &domain.Account{}

	return account, nil
}
