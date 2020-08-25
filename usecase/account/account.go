package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//Account armazena as dependências para os casos de uso de Account
type Account struct {
	repo       domain.AccountRepository
	presenter  AccountPresenter
	ctxTimeout time.Duration
}

//NewAccount constrói um Account com suas dependências
func NewAccount(repo domain.AccountRepository, presenter AccountPresenter, t time.Duration) Account {
	return Account{repo: repo, presenter: presenter, ctxTimeout: t}
}

//Store cria uma nova Account
func (a Account) Store(ctx context.Context, name, CPF string, balance domain.Money) (AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var account = domain.NewAccount(
		domain.AccountID(domain.NewUUID()),
		name,
		CPF,
		balance,
		time.Now(),
	)

	account, err := a.repo.Store(ctx, account)
	if err != nil {
		return a.presenter.Output(domain.Account{}), err
	}

	return a.presenter.Output(account), nil
}

//FindAll retorna uma lista de Accounts
func (a Account) FindAll(ctx context.Context) ([]AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	accounts, err := a.repo.FindAll(ctx)
	if err != nil {
		return a.presenter.OutputList([]domain.Account{}), err
	}

	return a.presenter.OutputList(accounts), nil
}

//FindBalance retorna o saldo de uma Account
func (a Account) FindBalance(ctx context.Context, ID domain.AccountID) (AccountBalanceOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	account, err := a.repo.FindBalance(ctx, ID)
	if err != nil {
		return a.presenter.OutputBalance(domain.Money(0)), err
	}

	return a.presenter.OutputBalance(account.Balance()), nil
}
