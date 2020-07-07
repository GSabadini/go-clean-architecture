package mongodb

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//transferBSON armazena a estrutura de dados do MongoDB
type transferBSON struct {
	ID                   string    `bson:"id"`
	AccountOriginID      string    `bson:"account_origin_id"`
	AccountDestinationID string    `bson:"account_destination_id"`
	Amount               float64   `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

//TransferRepository armazena a estrutura de dados de um repositório de Transfer
type TransferRepository struct {
	handler        repository.NoSQLHandler
	collectionName string
}

//NewTransferRepository constrói um repository com suas dependências
func NewTransferRepository(handler repository.NoSQLHandler) TransferRepository {
	return TransferRepository{handler: handler, collectionName: "transfers"}
}

//Store insere uma Transfer no database
func (t TransferRepository) Store(transfer domain.Transfer) (domain.Transfer, error) {
	transferBson := &transferBSON{
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

//FindAll busca todas as Transfer no database
func (t TransferRepository) FindAll() ([]domain.Transfer, error) {
	var (
		transfersBson = make([]transferBSON, 0)
		transfers     = make([]domain.Transfer, 0)
	)

	if err := t.handler.FindAll(t.collectionName, nil, &transfersBson); err != nil {
		return transfers, errors.Wrap(err, "error listing transfers")
	}

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
