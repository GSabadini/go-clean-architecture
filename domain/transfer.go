package domain

import "time"

//TransferRepository expõe os métodos disponíveis para as abstrações do repositório de Transfer
type TransferRepository interface {
	Store(Transfer) (Transfer, error)
	FindAll() ([]Transfer, error)
}

//Transfer armazena a estrutura de transferência
type Transfer struct {
	ID                   string
	AccountOriginID      string
	AccountDestinationID string
	Amount               float64
	CreatedAt            time.Time
}

//NewTransfer cria um Transfer
func NewTransfer(ID, accountOriginID, accountDestinationID string, amount float64, createdAt time.Time) Transfer {
	return Transfer{
		ID:                   ID,
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            createdAt,
	}
}
