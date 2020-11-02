package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type findAccountBalancePresenter struct{}

func NewFindAccountBalancePresenter() usecase.FindAccountBalancePresenter {
	return findAccountBalancePresenter{}
}

func (a findAccountBalancePresenter) Output(balance domain.Money) usecase.FindAccountBalanceOutput {
	return usecase.FindAccountBalanceOutput{Balance: balance.Float64()}
}
