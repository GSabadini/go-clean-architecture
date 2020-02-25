package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Transfer armazena a estrutura de transferência
type Transfer struct {
	Id                   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AccountOriginId      bson.ObjectId `json:"account_origin_id" bson:"account_origin_id"`
	AccountDestinationId bson.ObjectId `json:"account_destination_id" bson:"account_destination_id"`
	Amount               float64       `json:"amount" bson:"amount"`
	CreatedAt            time.Time     `json:"created_at" bson:"created_at"`
}

//GetAccountOriginId retorna o id da conta de origem
func (t *Transfer) GetAccountOriginId() bson.ObjectId {
	return t.AccountOriginId
}

//GetAccountDestinationId retorna o id da conta de destino
func (t *Transfer) GetAccountDestinationId() bson.ObjectId {
	return t.AccountDestinationId
}

//GetAmount retorna o valor
func (t *Transfer) GetAmount() float64 {
	return t.Amount
}
