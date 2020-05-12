package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
)

//AccountRepository expõe os métodos disponíveis para as abstrações do repositório de contas
type AccountRepository interface {
	Store(domain.Account) (domain.Account, error)
	UpdateBalance(string, float64) error
	FindAll() ([]domain.Account, error)
	FindByID(string) (*domain.Account, error)
	FindBalance(string) (domain.Account, error)
}

//TransferRepository expõe os métodos disponíveis para as abstrações do repositório de transferências
type TransferRepository interface {
	Store(domain.Transfer) (domain.Transfer, error)
	FindAll() ([]domain.Transfer, error)
}
