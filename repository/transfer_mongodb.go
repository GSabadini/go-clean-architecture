package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/pkg/errors"
)

//const transfersCollectionName = "transfers"

//TransferMongoDB representa um repositório para manipulação de dados de transferências utilizando MongoDB
type TransferMongoDB struct {
	handler        database.NoSQLHandler
	collectionName string
}

//NewTransferMongoDB cria um repository com suas dependências
func NewTransferMongoDB(handler database.NoSQLHandler) TransferMongoDB {
	return TransferMongoDB{handler: handler, collectionName: "transfers"}
}

//Store cria uma transferência através da implementação real do database
func (t TransferMongoDB) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.handler.Store(t.collectionName, &transfer); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll realiza uma busca através da implementação real do database
func (t TransferMongoDB) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.handler.FindAll(t.collectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	return transfer, nil
}
