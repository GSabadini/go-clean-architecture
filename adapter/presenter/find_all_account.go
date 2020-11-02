package presenter

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type findAllAccountPresenter struct{}

func NewFindAllAccountPresenter() usecase.FindAllAccountPresenter {
	return findAllAccountPresenter{}
}

func (a findAllAccountPresenter) Output(accounts []domain.Account) []usecase.FindAllAccountOutput {
	var o = make([]usecase.FindAllAccountOutput, 0)

	for _, account := range accounts {
		o = append(o, usecase.FindAllAccountOutput{
			ID:        account.ID().String(),
			Name:      account.Name(),
			CPF:       account.CPF(),
			Balance:   account.Balance().Float64(),
			CreatedAt: account.CreatedAt().Format(time.RFC3339),
		})
	}

	return o
}
