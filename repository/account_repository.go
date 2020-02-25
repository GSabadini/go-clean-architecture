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
func (a Account) Store(account *domain.Account) (*domain.Account, error) {
	account.CreatedAt = time.Now()
	account.Id = bson.NewObjectId()

	err := a.dbHandler.Store(accountsCollectionName, &account)
	if err != nil {
		return nil, err
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

	err := a.dbHandler.FindAll(accountsCollectionName, nil, &account)
	if err != nil {
		return []domain.Account{}, err
	}

	return account, nil
}

//FindOne realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindOne(query bson.M) (*domain.Account, error) {
	var account *domain.Account

	err := a.dbHandler.FindOne(accountsCollectionName, query, &account)
	if err != nil {
		return &domain.Account{}, err
	}

	return account, nil
}
