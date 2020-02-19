package database

import "github.com/gsabadini/go-bank-transfer/domain"

//NoSQLDbHandler expõe os métodos disponíveis para as abstrações de banco
type NoSQLDBHandler interface {
	Store(string, interface{}) error
	FindAll(string, []domain.Account) ([]domain.Account, error)
}
