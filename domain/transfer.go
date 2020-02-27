package domain

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Transfer armazena a estrutura de transferência
type Transfer struct {
	ID                   bson.ObjectId `json:"id" bson:"_id"`
	AccountOriginID      bson.ObjectId `json:"account_origin_id" bson:"account_origin_id"`
	AccountDestinationID bson.ObjectId `json:"account_destination_id" bson:"account_destination_id"`
	Amount               float64       `json:"amount" bson:"amount"`
	CreatedAt            time.Time     `json:"created_at" bson:"created_at"`
}

//ValidateAmount verifica se o saldo é valido
func (t *Transfer) ValidateAmount() error {
	if t.Amount < 0 {
		return errors.New("amount invalid")
	}

	return nil
}

//GetAccountOriginID retorna o id da conta de origem
func (t *Transfer) GetAccountOriginID() bson.ObjectId {
	return t.AccountOriginID
}

//GetAccountDestinationID retorna o id da conta de destino
func (t *Transfer) GetAccountDestinationID() bson.ObjectId {
	return t.AccountDestinationID
}

//GetAmount retorna o valor
func (t *Transfer) GetAmount() float64 {
	return t.Amount
}
