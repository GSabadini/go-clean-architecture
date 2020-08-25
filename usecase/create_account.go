package usecase

import (
	"context"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

//CreateAccount é uma abstração de caso de uso de Account
type CreateAccount interface {
	Execute(context.Context, input.Account) (output.AccountOutput, error)
}

//CreateAccountInteractor armazena as dependências para os casos de uso de Account
type CreateAccountInteractor struct {
	repo       domain.AccountRepository
	presenter  output.AccountPresenter
	ctxTimeout time.Duration
}

//NewCreateAccountInteractor constrói um Account com suas dependências
func NewCreateAccountInteractor(
	repo domain.AccountRepository,
	presenter output.AccountPresenter,
	t time.Duration,
) CreateAccountInteractor {
	return CreateAccountInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

//Store cria uma nova Account
func (a CreateAccountInteractor) Execute(ctx context.Context, input input.Account) (output.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var account = domain.NewAccount(
		domain.AccountID(domain.NewUUID()),
		input.Name,
		input.CPF,
		domain.Money(input.Balance),
		time.Now(),
	)

	account, err := a.repo.Store(ctx, account)
	if err != nil {
		return a.presenter.Output(domain.Account{}), err
	}

	return a.presenter.Output(account), nil
}
