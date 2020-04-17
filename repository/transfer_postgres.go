package repository

import (
	"database/sql"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/pkg/errors"
)

type TransferPostgres struct {
	handler *sql.DB
}

func NewTransferPostgres(handler *sql.DB) TransferPostgres {
	return TransferPostgres{handler: handler}
}

func (t TransferPostgres) Store(transfer domain.Transfer) (domain.Transfer, error) {
	query := `
		INSERT INTO 
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	_, err := t.handler.Exec(
		query,
		transfer.ID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
		transfer.CreatedAt,
	)
	if err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

func (t TransferPostgres) FindAll() ([]domain.Transfer, error) {
	var (
		transfers = make([]domain.Transfer, 0)
		query     = "SELECT * FROM transfers"
	)

	rows, err := t.handler.Query(query)
	if err != nil {
		return transfers, errors.Wrap(err, "error listing transfers")
	}

	for rows.Next() {
		var (
			ID                   string
			accountOriginID      string
			accountDestinationID string
			amount               float64
			createdAt            time.Time
		)

		if err = rows.Scan(&ID, &accountOriginID, &accountDestinationID, &amount, &createdAt); err != nil {
			return transfers, errors.Wrap(err, "error listing transfers")
		}

		transfer := domain.Transfer{
			ID:                   ID,
			AccountOriginID:      accountOriginID,
			AccountDestinationID: accountDestinationID,
			Amount:               amount,
			CreatedAt:            createdAt,
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
