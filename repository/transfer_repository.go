package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

const transfersCollectionName = "transfers"

//Transfer representa um repositório para dados de transferência
type Transfer DbRepository

//NewTransfer cria um repository com suas dependências
func NewTransfer(dbHandler database.NoSQLDBHandler) Transfer {
	return Transfer{dbHandler: dbHandler}
}

//Store cria uma transferência através da implementação real do database
func (t Transfer) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.dbHandler.Store(transfersCollectionName, &transfer); err != nil {
		return domain.Transfer{}, err
	}

	return transfer, nil
}

//FindAll realiza uma busca no banco de dados através da implementação real do database
func (t Transfer) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.dbHandler.FindAll(transfersCollectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, err
	}

	return transfer, nil
}
