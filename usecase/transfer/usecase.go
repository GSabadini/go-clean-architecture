package usecase

import (
	"context"
	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferUseCase é uma abstração para os casos de uso de Transfer
type TransferUseCase interface {
	Store(context.Context, domain.AccountID, domain.AccountID, domain.Money) (TransferOutput, error)
	FindAll(context.Context) ([]TransferOutput, error)
}

type TransferCreate interface {
	Execute(context.Context, TransferInput) (TransferOutput, error)
}

type TransferFindAll interface {
	Execute(context.Context) ([]TransferOutput, error)
}
