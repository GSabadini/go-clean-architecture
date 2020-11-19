package domain

import (
	"context"
	"time"
)

type TransferID string

func (t TransferID) String() string {
	return string(t)
}

type (
	TransferRepository interface {
		Create(context.Context, Transfer) (Transfer, error)
		FindAll(context.Context) ([]Transfer, error)
		WithTransaction(context.Context, func(context.Context) error) error
	}

	Transfer struct {
		id                   TransferID
		accountOriginID      AccountID
		accountDestinationID AccountID
		amount               Money
		createdAt            time.Time
	}
)

func NewTransfer(
	ID TransferID,
	accountOriginID AccountID,
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

func (t Transfer) ID() TransferID {
	return t.id
}

func (t Transfer) AccountOriginID() AccountID {
	return t.accountOriginID
}

func (t Transfer) AccountDestinationID() AccountID {
	return t.accountDestinationID
}

func (t Transfer) Amount() Money {
	return t.amount
}

func (t Transfer) CreatedAt() time.Time {
	return t.createdAt
}
