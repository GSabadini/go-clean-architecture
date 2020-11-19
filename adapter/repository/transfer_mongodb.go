package repository

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type transferBSON struct {
	ID                   string    `bson:"id"`
	AccountOriginID      string    `bson:"account_origin_id"`
	AccountDestinationID string    `bson:"account_destination_id"`
	Amount               int64     `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

type TransferNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewTransferNoSQL(db NoSQL) TransferNoSQL {
	return TransferNoSQL{
		db:             db,
		collectionName: "transfers",
	}
}

func (t TransferNoSQL) Create(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	transferBSON := &transferBSON{
		ID:                   transfer.ID().String(),
		AccountOriginID:      transfer.AccountOriginID().String(),
		AccountDestinationID: transfer.AccountDestinationID().String(),
		Amount:               transfer.Amount().Int64(),
		CreatedAt:            transfer.CreatedAt(),
	}

	if err := t.db.Store(ctx, t.collectionName, transferBSON); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

func (t TransferNoSQL) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var transfersBSON = make([]transferBSON, 0)

	if err := t.db.FindAll(ctx, t.collectionName, bson.M{}, &transfersBSON); err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	var transfers = make([]domain.Transfer, 0)

	for _, transferBSON := range transfersBSON {
		var transfer = domain.NewTransfer(
			domain.TransferID(transferBSON.ID),
			domain.AccountID(transferBSON.AccountOriginID),
			domain.AccountID(transferBSON.AccountDestinationID),
			domain.Money(transferBSON.Amount),
			transferBSON.CreatedAt,
		)

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (t TransferNoSQL) WithTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	session, err := t.db.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = session.WithTransaction(ctx, fn)
	if err != nil {
		return err
	}

	return nil
}
