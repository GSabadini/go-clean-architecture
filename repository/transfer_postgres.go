package repository

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

//TransferPostgres representa um repositório de dados para transferências utilizando Postgres
type TransferPostgres struct {
	handler database.SQLDbHandler
}

//NewTransferPostgres cria um repositório utilizando Postgres
func NewTransferPostgres(handler database.SQLDbHandler) TransferPostgres {
	return TransferPostgres{handler: handler}
}

//Store cria uma transferência
func (t TransferPostgres) Store(transfer domain.Transfer) (domain.Transfer, error) {
	query := `
		INSERT INTO 
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := t.handler.Execute(
		query,
		transfer.ID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
		transfer.CreatedAt,
	); err != nil {
		return domain.Transfer{}, errors.Wrap(err, "error creating transfer")
	}

	return transfer, nil
}

//FindAll retorna uma lista de transferências
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
