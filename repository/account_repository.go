package repository

import (
	"time"

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
func (a Account) FindAll(account []domain.Account) ([]domain.Account, error) {
	err := a.dbHandler.FindAll(accountsCollectionName, nil, &account)
	if err != nil {
		return []domain.Account{}, err
	}

	return account, nil
}
