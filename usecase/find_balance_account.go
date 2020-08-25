package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

//FindBalanceAccount é uma abstração de caso de uso de Account
type FindBalanceAccount interface {
	Execute(context.Context, domain.AccountID) (output.AccountBalanceOutput, error)
}

//FindBalanceAccountInteractor armazena as dependências para os casos de uso de Account
type FindBalanceAccountInteractor struct {
	repo       domain.AccountRepository
	presenter  output.AccountPresenter
	ctxTimeout time.Duration
}

//NewFindBalanceAccountInteractor constrói um FindBalanceAccountInteractor com suas dependências
func NewFindBalanceAccountInteractor(
	repo domain.AccountRepository,
	presenter output.AccountPresenter,
	t time.Duration,
) FindBalanceAccountInteractor {
	return FindBalanceAccountInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

//Execute retorna o saldo de uma Account
func (a FindBalanceAccountInteractor) Execute(ctx context.Context, ID domain.AccountID) (output.AccountBalanceOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	account, err := a.repo.FindBalance(ctx, ID)
	if err != nil {
		return a.presenter.OutputBalance(domain.Money(0)), err
	}

	return a.presenter.OutputBalance(account.Balance()), nil
}
