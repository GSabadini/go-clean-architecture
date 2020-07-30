package domain

import (
	"context"
	"time"
)

//TransferRepository expõe os métodos disponíveis para as abstrações do repositório de Transfer
type TransferRepository interface {
	Store(context.Context, Transfer) (Transfer, error)
	FindAll(context.Context) ([]Transfer, error)
}

//TransferID define o tipo identificador de uma Transfer
type TransferID string

//Transfer armazena a estrutura de transferência
type Transfer struct {
	ID                   TransferID
	AccountOriginID      AccountID
	AccountDestinationID AccountID
	Amount               float64
	CreatedAt            time.Time
}

//NewTransfer cria um Transfer
func NewTransfer(
	ID TransferID,
	accountOriginID,
	accountDestinationID AccountID,
	amount float64,
	createdAt time.Time,
) Transfer {
	return Transfer{
		ID:                   ID,
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            createdAt,
	}
}
