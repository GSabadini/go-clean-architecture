package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type (
	// CreateAccountUseCase input port
	CreateAccountUseCase interface {
		Execute(context.Context, CreateAccountInput) (CreateAccountOutput, error)
	}

	// CreateAccountInput input data
	CreateAccountInput struct {
		Name    string `json:"name" validate:"required"`
		CPF     string `json:"cpf" validate:"required"`
		Balance int64  `json:"balance" validate:"gt=0,required"`
	}

	// CreateAccountPresenter output port
	CreateAccountPresenter interface {
		Output(domain.Account) CreateAccountOutput
	}

	// CreateAccountOutput output data
	CreateAccountOutput struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		CPF       string  `json:"cpf"`
		Balance   float64 `json:"balance"`
		CreatedAt string  `json:"created_at"`
	}

	createAccountInteractor struct {
		repo       domain.AccountRepository
		presenter  CreateAccountPresenter
		ctxTimeout time.Duration
	}
)

// NewCreateAccountInteractor creates new createAccountInteractor with its dependencies
func NewCreateAccountInteractor(
	repo domain.AccountRepository,
	presenter CreateAccountPresenter,
	t time.Duration,
) CreateAccountUseCase {
	return createAccountInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a createAccountInteractor) Execute(ctx context.Context, input CreateAccountInput) (CreateAccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var account = domain.NewAccount(
		domain.AccountID(domain.NewUUID()),
		input.Name,
		input.CPF,
		domain.Money(input.Balance),
		time.Now(),
	)

	account, err := a.repo.Create(ctx, account)
	if err != nil {
		return a.presenter.Output(domain.Account{}), err
	}

	return a.presenter.Output(account), nil
}
