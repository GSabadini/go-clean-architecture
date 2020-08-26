package postgres

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/interface/repository"

	"github.com/pkg/errors"
)

type TransferRepository struct {
	handler repository.SQLHandler
}

func NewTransferRepository(h repository.SQLHandler) TransferRepository {
	return TransferRepository{handler: h}
}

func (t TransferRepository) Create(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	var query = `
		INSERT INTO 
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := t.handler.ExecuteContext(
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

func (t TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var query = "SELECT * FROM transfers"

	rows, err := t.handler.QueryContext(ctx, query)
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
