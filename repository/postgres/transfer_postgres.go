package postgres

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

//TransferRepository representa um repositório para manipulação de dados de transferências utilizando Postgres
type TransferRepository struct {
	handler database.SQLHandler
}

//NewTransferRepository cria um repositório utilizando Postgres
func NewTransferRepository(handler database.SQLHandler) TransferRepository {
	return TransferRepository{handler: handler}
}

//Store cria uma transferência
func (t TransferRepository) Store(transfer domain.Transfer) (domain.Transfer, error) {
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
func (t TransferRepository) FindAll() ([]domain.Transfer, error) {
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
