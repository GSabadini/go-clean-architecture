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

//String converte o tipo TranferID para uma string
func (t TransferID) String() string {
	return string(t)
}

//Transfer armazena a estrutura de transferência
type Transfer struct {
	id                   TransferID
	accountOriginID      AccountID
	accountDestinationID AccountID
	amount               Money
	createdAt            time.Time
}

//NewTransfer cria um Transfer
func NewTransfer(
	ID TransferID,
	accountOriginID,
	accountDestinationID AccountID,
	amount Money,
	createdAt time.Time,
) Transfer {
	return Transfer{
		id:                   ID,
		accountOriginID:      accountOriginID,
		accountDestinationID: accountDestinationID,
		amount:               amount,
		createdAt:            createdAt,
	}
}

//CreatedAt
func (t Transfer) ID() TransferID {
	return t.id
}

//AccountOriginID
func (t Transfer) AccountOriginID() AccountID {
	return t.accountOriginID
}

//AccountDestinationID
func (t Transfer) AccountDestinationID() AccountID {
	return t.accountDestinationID
}

//Amount
func (t Transfer) Amount() Money {
	return t.amount
}

//CreatedAt
func (t Transfer) CreatedAt() time.Time {
	return t.createdAt
}
