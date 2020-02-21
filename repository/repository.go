package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

//DbRepository encapsula um repositório que contém um handler para determinado banco de dados
type DbRepository struct {
	dbHandler database.NoSQLDBHandler
}

type Repository interface {
	Store(domain.Account) error
	FindAll([]domain.Account) ([]domain.Account, error)
}
