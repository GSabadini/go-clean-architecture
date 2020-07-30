package postgres

import (
	"context"
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
func NewTransferRepository(h repository.SQLHandler) TransferRepository {
	return TransferRepository{handler: h}
}

//Store insere uma Transfer no database
func (t TransferRepository) Store(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	query := `
		INSERT INTO 
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := t.handler.ExecuteContext(
		ctx,
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
func (t TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	var (
		transfers = make([]domain.Transfer, 0)
		query     = "SELECT * FROM transfers"
	)

	rows, err := t.handler.QueryContext(ctx, query)
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
			ID:                   domain.TransferID(ID),
			AccountOriginID:      domain.AccountID(accountOriginID),
			AccountDestinationID: domain.AccountID(accountDestinationID),
			Amount:               amount,
			CreatedAt:            createdAt,
		})
	}

	if err = rows.Err(); err != nil {
		return []domain.Transfer{}, err
	}

	return transfers, nil
}
