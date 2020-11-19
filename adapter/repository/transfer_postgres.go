package repository

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/pkg/errors"
)

type TransferSQL struct {
	db SQL
}

func NewTransferSQL(db SQL) TransferSQL {
	return TransferSQL{
		db: db,
	}
}

func (t TransferSQL) Create(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	tx, ok := ctx.Value("TransactionContextKey").(Tx)
	if !ok {
		var err error
		tx, err = t.db.BeginTx(ctx)
		if err != nil {
			return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
		}
	}

	var query = `
		INSERT INTO 
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := tx.ExecuteContext(
		ctx,
		query,
		transfer.ID(),
		transfer.AccountOriginID(),
		transfer.AccountDestinationID(),
		transfer.Amount(),
		transfer.CreatedAt(),
	); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

func (t TransferSQL) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var query = "SELECT * FROM transfers"

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
	}

	var transfers = make([]domain.Transfer, 0)
	for rows.Next() {
		var (
			ID                   string
			accountOriginID      string
			accountDestinationID string
			amount               int64
			createdAt            time.Time
		)

		if err = rows.Scan(&ID, &accountOriginID, &accountDestinationID, &amount, &createdAt); err != nil {
			return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
		}

		transfers = append(transfers, domain.NewTransfer(
			domain.TransferID(ID),
			domain.AccountID(accountOriginID),
			domain.AccountID(accountDestinationID),
			domain.Money(amount),
			createdAt,
		))
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return []domain.Transfer{}, err
	}

	return transfers, nil
}

func (t TransferSQL) WithTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "error begin tx")
	}

	ctxTx := context.WithValue(ctx, "TransactionContextKey", tx)
	err = fn(ctxTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(err, "rollback error")
		}
		return err
	}

	return tx.Commit()
}
