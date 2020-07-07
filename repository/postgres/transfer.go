package postgres

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//TransferRepository armazena a estrutura de dados de um repositório de Transfer
type TransferRepository struct {
	handler repository.SQLHandler
}

//NewTransferRepository constrói um TransferRepository com suas dependências
func NewTransferRepository(handler repository.SQLHandler) TransferRepository {
	return TransferRepository{handler: handler}
}

//Store insere uma Transfer no database
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

//FindAll busca todas as Transfer no database
func (t TransferRepository) FindAll() ([]domain.Transfer, error) {
	var (
		transfers = make([]domain.Transfer, 0)
		query     = "SELECT * FROM transfers"
	)

	rows, err := t.handler.Query(query)
	if err != nil {
		return transfers, errors.Wrap(err, "error listing transfers")
	}

	defer rows.Close()
	for rows.Next() {
		var (
			ID                   string
			accountOriginID      string
			accountDestinationID string
			amount               float64
			createdAt            time.Time
		)

		if err = rows.Scan(&ID, &accountOriginID, &accountDestinationID, &amount, &createdAt); err != nil {
			return []domain.Transfer{}, errors.Wrap(err, "error listing transfers")
		}

		transfers = append(transfers, domain.Transfer{
			ID:                   ID,
			AccountOriginID:      accountOriginID,
			AccountDestinationID: accountDestinationID,
			Amount:               amount,
			CreatedAt:            createdAt,
		})
	}

	if err = rows.Err(); err != nil {
		return []domain.Transfer{}, err
	}

	return transfers, nil
}
