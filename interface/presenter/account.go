package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type accountPresenter struct{}

func NewAccountPresenter() accountPresenter {
	return accountPresenter{}
}

func (a accountPresenter) Output(account domain.Account) output.Account {
	return output.Account{
		ID:        account.ID().String(),
		Name:      account.Name(),
		CPF:       account.CPF(),
		Balance:   account.Balance().Float64(),
		CreatedAt: account.CreatedAt(),
	}
}

func (a accountPresenter) OutputList(accounts []domain.Account) []output.Account {
	var o = make([]output.Account, 0)

	for _, account := range accounts {
		o = append(o, output.Account{
			ID:        account.ID().String(),
			Name:      account.Name(),
			CPF:       account.CPF(),
			Balance:   account.Balance().Float64(),
			CreatedAt: account.CreatedAt(),
		})
	}

	return o
}

func (a accountPresenter) OutputBalance(balance domain.Money) output.AccountBalance {
	return output.AccountBalance{Balance: balance.Float64()}
}
