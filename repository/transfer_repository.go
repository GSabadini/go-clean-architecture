package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
)

const transfersCollectionName = "transfers"

//Transfer representa um repositório para dados de transferência
type Transfer DbRepository

//NewTransfer cria um repository com suas dependências
func NewTransfer(dbHandler database.NoSQLDBHandler) Transfer {
	return Transfer{dbHandler: dbHandler}
}

//Store cria uma transferência
func (t Transfer) Store(transfer *domain.Transfer) error {
	transfer.CreatedAt = time.Now()
	transfer.ID = bson.NewObjectId()

	return t.dbHandler.Store(transfersCollectionName, &transfer)
}
