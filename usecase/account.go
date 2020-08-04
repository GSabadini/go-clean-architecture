package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//AccountOutput armazena a estrutura de dados de retorno do caso de uso
type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

//AccountBalanceOutput armazena a estrutura de dados de retorno do caso de uso
type AccountBalanceOutput struct {
	Balance float64 `json:"balance"`
}

//Account armazena as dependências para os casos de uso de Account
type Account struct {
	repo       domain.AccountRepository
	ctxTimeout time.Duration
}

//NewAccount constrói um Account com suas dependências
func NewAccount(repo domain.AccountRepository, t time.Duration) Account {
	return Account{repo: repo, ctxTimeout: t}
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
		return AccountOutput{}, err
	}

	return AccountOutput{
		ID:        account.ID.String(),
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance.Float64(),
		CreatedAt: account.CreatedAt,
	}, nil
}

//FindAll retorna uma lista de Accounts
func (a Account) FindAll(ctx context.Context) ([]AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var output = make([]AccountOutput, 0)

	accounts, err := a.repo.FindAll(ctx)
	if err != nil {
		return output, err
	}

	for _, account := range accounts {
		output = append(output, AccountOutput{
			ID:        account.ID.String(),
			Name:      account.Name,
			CPF:       account.CPF,
			Balance:   account.Balance.Float64(),
			CreatedAt: account.CreatedAt,
		})
	}

	return output, nil
}

//FindBalance retorna o saldo de uma Account
func (a Account) FindBalance(ctx context.Context, ID domain.AccountID) (AccountBalanceOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	account, err := a.repo.FindBalance(ctx, ID)
	if err != nil {
		return AccountBalanceOutput{}, err
	}

	return AccountBalanceOutput{
		Balance: account.Balance.Float64(),
	}, nil
}
