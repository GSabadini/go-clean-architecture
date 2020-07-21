package usecase

import "github.com/gsabadini/go-bank-transfer/domain"

//AccountUseCase é uma abstração para os casos de uso de Account
type AccountUseCase interface {
	Store(string, string, float64) (AccountOutput, error)
	FindAll() ([]AccountOutput, error)
	FindBalance(domain.AccountID) (AccountBalanceOutput, error)
}

//TransferUseCase é uma abstração para os casos de uso de Transfer
type TransferUseCase interface {
	Store(domain.AccountID, domain.AccountID, float64) (TransferOutput, error)
	FindAll() ([]TransferOutput, error)
}
