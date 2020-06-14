package mongodb

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/pkg/errors"
)

//TransferRepository representa um repositório para manipulação de dados de transferências utilizando MongoDB
type TransferRepository struct {
	handler        database.NoSQLHandler
	collectionName string
}

//NewTransferRepository cria um repository com suas dependências
func NewTransferRepository(handler database.NoSQLHandler) TransferRepository {
	return TransferRepository{handler: handler, collectionName: "transfers"}
}

//Store cria uma transferência através da implementação real do database
func (t TransferRepository) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.handler.Store(t.collectionName, &transfer); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll realiza uma busca através da implementação real do database
func (t TransferRepository) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.handler.FindAll(t.collectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	return transfer, nil
}
