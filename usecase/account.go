package usecase

import (
	"strings"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferOutput armazena a estrutura de dados de retorno do caso de uso
type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

//TransferOutput armazena a estrutura de dados de retorno do caso de uso
type AccountBalanceOutput struct {
	Balance float64 `json:"balance"`
}

//Account armazena as dependências para os casos de uso de Account
type Account struct {
	repository domain.AccountRepository
}

//NewAccount constrói um Account com suas dependências
func NewAccount(repository domain.AccountRepository) Account {
	return Account{repository: repository}
}

//Store cria uma nova Account
func (a Account) Store(name, CPF string, balance float64) (AccountOutput, error) {
	var account = domain.NewAccount(
		domain.NewUUID(),
		name,
		a.cleanCPF(CPF),
		balance,
		time.Now(),
	)

	account, err := a.repository.Store(account)
	if err != nil {
		return AccountOutput{}, err
	}

	return AccountOutput{
		ID:        account.ID,
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
	}, nil
}

//FindAll retorna uma lista de Accounts
func (a Account) FindAll() ([]AccountOutput, error) {
	var output = make([]AccountOutput, 0)

	accounts, err := a.repository.FindAll()
	if err != nil {
		return output, err
	}

	for _, account := range accounts {
		var account = AccountOutput{
			ID:        account.ID,
			Name:      account.Name,
			CPF:       account.CPF,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
		}

		output = append(output, account)
	}

	return output, nil
}

//FindBalance retorna o saldo de uma Account
func (a Account) FindBalance(ID string) (AccountBalanceOutput, error) {
	account, err := a.repository.FindBalance(ID)
	if err != nil {
		return AccountBalanceOutput{}, err
	}

	return AccountBalanceOutput{
		Balance: account.Balance,
	}, nil
}

func (a Account) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}
