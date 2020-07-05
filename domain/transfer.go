package domain

import "time"

//TransferRepository expõe os métodos disponíveis para as abstrações do repositório de transferências
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
