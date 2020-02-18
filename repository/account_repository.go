package repository

import (
	"time"

	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/infrastructure/database"
)

const accountsCollectionName = "accounts"

type Account DbRepository

//NewAccount cria um repository com suas dependências
func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store realiza uma inserção no banco de dados através da implementação real do database
func (a Account) Store(account domain.Account) error {
	account.CreatedAt = time.Now()

	return a.dbHandler.Insert(accountsCollectionName, account)
}
