package mongodb

import (
	"context"
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
	Amount               int64     `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

//TransferRepository armazena a estrutura de dados de um repositório de Transfer
type TransferRepository struct {
	collectionName string
	handler        repository.NoSQLHandler
}

//NewTransferRepository constrói um repository com suas dependências
func NewTransferRepository(h repository.NoSQLHandler) TransferRepository {
	return TransferRepository{handler: h, collectionName: "transfers"}
}

//Store insere uma Transfer no database
func (t TransferRepository) Store(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	transferBson := &transferBSON{
		ID:                   transfer.ID.String(),
		AccountOriginID:      transfer.AccountOriginID.String(),
		AccountDestinationID: transfer.AccountDestinationID.String(),
		Amount:               transfer.Amount.Int64(),
		CreatedAt:            transfer.CreatedAt,
	}

	if err := t.handler.Store(ctx, t.collectionName, transferBson); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll busca todas as Transfer no database
func (t TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var (
		transfersBson = make([]transferBSON, 0)
		transfers     = make([]domain.Transfer, 0)
	)

	if err := t.handler.FindAll(ctx, t.collectionName, nil, &transfersBson); err != nil {
		return transfers, errors.Wrap(err, "error listing transfers")
	}

	for _, transferBson := range transfersBson {
		var transfer = domain.Transfer{
			ID:                   domain.TransferID(transferBson.ID),
			AccountOriginID:      domain.AccountID(transferBson.AccountOriginID),
			AccountDestinationID: domain.AccountID(transferBson.AccountDestinationID),
			Amount:               domain.Money(transferBson.Amount),
			CreatedAt:            transferBson.CreatedAt,
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
