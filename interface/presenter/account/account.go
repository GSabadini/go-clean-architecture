package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	usecase "github.com/gsabadini/go-bank-transfer/usecase/account"
)

type accountPresenter struct{}

//NewAccountPresenter
func NewAccountPresenter() accountPresenter {
	return accountPresenter{}
}

//Output
func (a accountPresenter) Output(account domain.Account) usecase.AccountOutput {
	return usecase.AccountOutput{
		ID:        account.ID().String(),
		Name:      account.Name(),
		CPF:       account.CPF(),
		Balance:   account.Balance().Float64(),
		CreatedAt: account.CreatedAt(),
	}
}

//OutputList
func (a accountPresenter) OutputList(accounts []domain.Account) []usecase.AccountOutput {
	var output = make([]usecase.AccountOutput, 0)

	for _, account := range accounts {
		output = append(output, usecase.AccountOutput{
			ID:        account.ID().String(),
			Name:      account.Name(),
			CPF:       account.CPF(),
			Balance:   account.Balance().Float64(),
			CreatedAt: account.CreatedAt(),
		})
	}

	return output
}

//OutputBalance
func (a accountPresenter) OutputBalance(balance domain.Money) usecase.AccountBalanceOutput {
	return usecase.AccountBalanceOutput{Balance: balance.Float64()}
}