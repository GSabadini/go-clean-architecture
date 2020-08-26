package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type FindBalanceAccount interface {
	Execute(context.Context, domain.AccountID) (output.AccountBalance, error)
}

type FindBalanceAccountInteractor struct {
	repo       domain.AccountRepository
	presenter  output.AccountPresenter
	ctxTimeout time.Duration
}

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

func (a FindBalanceAccountInteractor) Execute(ctx context.Context, ID domain.AccountID) (output.AccountBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	account, err := a.repo.FindBalance(ctx, ID)
	if err != nil {
		return a.presenter.OutputBalance(domain.Money(0)), err
	}

	return a.presenter.OutputBalance(account.Balance()), nil
}
