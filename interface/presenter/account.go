package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type accountPresenter struct{}

//NewAccountPresenter
func NewAccountPresenter() accountPresenter {
	return accountPresenter{}
}

//Output
func (a accountPresenter) Output(account domain.Account) output.AccountOutput {
	return output.AccountOutput{
		ID:        account.ID().String(),
		Name:      account.Name(),
		CPF:       account.CPF(),
		Balance:   account.Balance().Float64(),
		CreatedAt: account.CreatedAt(),
	}
}

//OutputList
func (a accountPresenter) OutputList(accounts []domain.Account) []output.AccountOutput {
	var o = make([]output.AccountOutput, 0)

	for _, account := range accounts {
		o = append(o, output.AccountOutput{
			ID:        account.ID().String(),
			Name:      account.Name(),
			CPF:       account.CPF(),
			Balance:   account.Balance().Float64(),
			CreatedAt: account.CreatedAt(),
		})
	}

	return o
}

//OutputBalance
func (a accountPresenter) OutputBalance(balance domain.Money) output.AccountBalanceOutput {
	return output.AccountBalanceOutput{Balance: balance.Float64()}
}
