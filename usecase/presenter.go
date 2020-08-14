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

//AccountPresenter é uma abstração para os apresentação de Account
type AccountPresenter interface {
	Output(domain.Account) AccountOutput
	OutputList([]domain.Account) []AccountOutput
	OutputBalance(domain.Money) AccountBalanceOutput
}

//AccountOutput armazena a estrutura de dados de retorno do caso de uso
type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

//AccountBalanceOutput armazena a estrutura de dados de retorno do caso de uso
type AccountBalanceOutput struct {
	Balance float64 `json:"balance"`
}
