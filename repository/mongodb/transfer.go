package mongodb

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/pkg/errors"
	"time"
)

type transferBson struct {
	ID                   string    `bson:"id"`
	AccountOriginID      string    `bson:"account_origin_id"`
	AccountDestinationID string    `bson:"account_destination_id"`
	Amount               float64   `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

//TransferRepository representa um repositório para manipulação de dados de transferências utilizando MongoDB
type TransferRepository struct {
	handler        repository.NoSQLHandler
	collectionName string
}

//NewTransferRepository cria um repository com suas dependências
func NewTransferRepository(handler repository.NoSQLHandler) TransferRepository {
	return TransferRepository{handler: handler, collectionName: "transfers"}
}

//Store cria uma transferência através da implementação real do database
func (t TransferRepository) Store(transfer domain.Transfer) (domain.Transfer, error) {
	transferBson := &transferBson{
		ID:                   transfer.ID,
		AccountOriginID:      transfer.AccountOriginID,
		AccountDestinationID: transfer.AccountDestinationID,
		Amount:               transfer.Amount,
		CreatedAt:            transfer.CreatedAt,
	}

	if err := t.handler.Store(t.collectionName, transferBson); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll realiza uma busca através da implementação real do database
func (t TransferRepository) FindAll() ([]domain.Transfer, error) {
	var transfersBson = make([]transferBson, 0)

	if err := t.handler.FindAll(t.collectionName, nil, &transfersBson); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	var transfers []domain.Transfer
	for _, transferBson := range transfersBson {
		var transfer = domain.Transfer{
			ID:                   transferBson.ID,
			AccountOriginID:      transferBson.AccountOriginID,
			AccountDestinationID: transferBson.AccountDestinationID,
			Amount:               transferBson.Amount,
			CreatedAt:            transferBson.CreatedAt,
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
