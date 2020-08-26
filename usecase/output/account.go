package output

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type AccountPresenter interface {
	Output(domain.Account) Account
	OutputList([]domain.Account) []Account
	OutputBalance(domain.Money) AccountBalance
}

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountBalance struct {
	Balance float64 `json:"balance"`
}
