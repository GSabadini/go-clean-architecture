package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
)

type AccountUseCase interface {
	Store(domain.Account) (domain.Account, error)
	FindAll() ([]domain.Account, error)
	FindBalance(string) (domain.Account, error)
}

type TransferUseCase interface {
	Store(TransferInput) (TransferResult, error)
	FindAll() ([]TransferResult, error)
}
