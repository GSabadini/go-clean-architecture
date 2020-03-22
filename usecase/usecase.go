package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
)

type AccountUseCase interface {
	StoreAccount(domain.Account) (domain.Account, error)
	FindAllAccount() ([]domain.Account, error)
	FindBalanceAccount(string) (domain.Account, error)
}

type TransferUseCase interface {
	FindAllTransfer() ([]domain.Transfer, error)
	StoreTransfer(domain.Transfer) (domain.Transfer, error)
}
