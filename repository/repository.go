package repository

import (
	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/infrastructure/database"
)

type DbRepository struct {
	dbHandler database.NoSQLDBHandler
}

type Repository interface {
	Store(domain.Account) error
}
