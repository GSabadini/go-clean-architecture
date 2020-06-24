package domain

import (
	"time"
)

//Transfer armazena a estrutura de transferência
type Transfer struct {
	ID                   string    `json:"id" bson:"id"`
	AccountOriginID      string    `json:"account_origin_id" bson:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id" bson:"account_destination_id"`
	Amount               float64   `json:"amount" bson:"amount"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at"`
}

//NewTransfer cria uma transferência
func NewTransfer(accountOriginID string, accountDestinationID string, amount float64) Transfer {
	return Transfer{
		ID:                   uuid(),
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}
}

//GetAccountOriginID retorna o id da conta de origem
func (t *Transfer) GetAccountOriginID() string {
	return t.AccountOriginID
}

//GetAccountDestinationID retorna o id da conta de destino
func (t *Transfer) GetAccountDestinationID() string {
	return t.AccountDestinationID
}

//GetAmount retorna o valor
func (t *Transfer) GetAmount() float64 {
	return t.Amount
}
