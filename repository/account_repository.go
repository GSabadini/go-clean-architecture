package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

const accountsCollectionName = "accounts"

//Account representa um repositório para dados da account
type Account DbRepository

//NewAccount cria um repository com suas dependências
func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store realiza uma inserção no banco de dados através da implementação real do database
func (a Account) Store(account domain.Account) error {
	account.CreatedAt = time.Now()

	return a.dbHandler.Store(accountsCollectionName, account)
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
func (a Account) FindOne(query bson.M) (domain.Account, error) {
	var account domain.Account

	err := a.dbHandler.FindOne(accountsCollectionName, query, &account)
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}
