package repository

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

//DbRepository encapsula um repositório que contém um handler para determinado banco de dados
type DbRepository struct {
	dbHandler database.NoSQLDBHandler
}

//AccountRepository expõe os métodos disponíveis para as abstrações de banco
type AccountRepository interface {
	Store(*domain.Account) (*domain.Account, error)
	Update(bson.M, bson.M) error
	FindAll() ([]domain.Account, error)
	FindOne(bson.M) (*domain.Account, error)
}

//TransferRepository expõe os métodos disponíveis para as abstrações de banco
type TransferRepository interface {
	Store(*domain.Transfer) error
}
