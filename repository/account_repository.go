package repository

import (
	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/infrastructure/database"
)

const stoneCollectionName = "stone"

type Account struct {
	dbHandler database.NoSQLDBHandler
}

func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

func (a Account) Store(account domain.Account) error {
	return a.dbHandler.Insert(stoneCollectionName, account)
}
