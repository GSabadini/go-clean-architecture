package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

//FindAllAccount é uma abstração de caso de uso de Account
type FindAllAccount interface {
	Execute(context.Context) ([]output.AccountOutput, error)
}

//FindBalanceAccountInteractor armazena as dependências para os casos de uso de Account
type FindAllAccountInteractor struct {
	repo       domain.AccountRepository
	presenter  output.AccountPresenter
	ctxTimeout time.Duration
}

//NewFindAllAccountInteractor constrói um FindBalanceAccountInteractor com suas dependências
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

//Execute retorna uma lista de Accounts
func (a FindAllAccountInteractor) Execute(ctx context.Context) ([]output.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	accounts, err := a.repo.FindAll(ctx)
	if err != nil {
		return a.presenter.OutputList([]domain.Account{}), err
	}

	return a.presenter.OutputList(accounts), nil
}
