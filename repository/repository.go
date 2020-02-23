package repository

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"gopkg.in/mgo.v2/bson"
)

//DbRepository encapsula um repositório que contém um handler para determinado banco de dados
type DbRepository struct {
	dbHandler database.NoSQLDBHandler
}

//AccountRepository expõe os métodos disponíveis para as abstrações de banco
type AccountRepository interface {
	Store(domain.Account) error
	FindAll([]domain.Account) ([]domain.Account, error)
	FindOne(bson.M, domain.Account) (domain.Account, error)
}
