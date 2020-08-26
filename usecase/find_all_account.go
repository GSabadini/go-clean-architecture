package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type FindAllAccount interface {
	Execute(context.Context) ([]output.Account, error)
}

type FindAllAccountInteractor struct {
	repo       domain.AccountRepository
	presenter  output.AccountPresenter
	ctxTimeout time.Duration
}

func NewFindAllAccountInteractor(
	repo domain.AccountRepository,
	presenter output.AccountPresenter,
	t time.Duration,
) FindAllAccountInteractor {
	return FindAllAccountInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

func (a FindAllAccountInteractor) Execute(ctx context.Context) ([]output.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	accounts, err := a.repo.FindAll(ctx)
	if err != nil {
		return a.presenter.OutputList([]domain.Account{}), err
	}

	return a.presenter.OutputList(accounts), nil
}
