package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

type DbRepository struct {
	dbHandler database.NoSQLDBHandler
}

type Repository interface {
	Store(domain.Account) error
	FindAll() ([]domain.Account, error)
}
