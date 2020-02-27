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
func (t Transfer) Store(transfer domain.Transfer) (domain.Transfer, error) {
	transfer.CreatedAt = time.Now()
	transfer.ID = bson.NewObjectId()

	if err := t.dbHandler.Store(transfersCollectionName, &transfer); err != nil {
		return domain.Transfer{}, err
	}

	return transfer, nil
}

//FindAll realiza uma busca no banco de dados através da implementação real do database
func (t Transfer) FindAll() ([]domain.Transfer, error) {
	var transfer = make([]domain.Transfer, 0)

	if err := t.dbHandler.FindAll(transfersCollectionName, nil, &transfer); err != nil {
		return []domain.Transfer{}, err
	}

	return transfer, nil
}

type TransferRepositoryMockSuccess struct{}

//Store cria uma transferência
func (t TransferRepositoryMockSuccess) Store(transfer domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{
		ID:                   "5e570851adcef50116aa7a5a",
		AccountOriginID:      "5e570851adcef50116aa7a5d",
		AccountDestinationID: "5e570851adcef50116aa7a5c",
		Amount:               100,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll realiza uma busca no banco de dados através da implementação real do database
func (t TransferRepositoryMockSuccess) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{
		{
			ID:                   "5e570851adcef50116aa7a5a",
			AccountOriginID:      "5e570851adcef50116aa7a5d",
			AccountDestinationID: "5e570851adcef50116aa7a5c",
			Amount:               100,
			CreatedAt:            time.Time{},
		},
	}, nil
}
