package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferPresenter é uma abstração para a apresentação de Account
type TransferPresenter interface {
	Output(domain.Transfer) TransferOutput
	OutputList([]domain.Transfer) []TransferOutput
}

//TransferOutput armazena a estrutura de dados de retorno do caso de uso
type TransferOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
