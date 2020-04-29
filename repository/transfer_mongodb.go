package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/pkg/errors"
)

const transfersCollectionName = "transfers"

//TransferMongoDB representa um repositório de dados para transferências utilizando MongoDB
type TransferMongoDB struct {
	handler database.NoSQLDbHandler
}

//NewTransfer cria um repository com suas dependências
func NewTransferMongoDB(handler database.NoSQLDbHandler) TransferMongoDB {
	return TransferMongoDB{handler: handler}
}

//Store cria uma transferência através da implementação real do database
func (t TransferMongoDB) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.handler.Store(transfersCollectionName, &transfer); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll realiza uma busca através da implementação real do database
func (t TransferMongoDB) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.handler.FindAll(transfersCollectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	return transfer, nil
}
