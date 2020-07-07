package domain

import "time"

//AccountRepository expõe os métodos disponíveis para as abstrações do repositório de Account
type AccountRepository interface {
	Store(Account) (Account, error)
	UpdateBalance(string, float64) error
	FindAll() ([]Account, error)
	FindByID(string) (*Account, error)
	FindBalance(string) (Account, error)
}

//Account armazena a estrutura de uma conta
type Account struct {
	ID        string
	Name      string
	CPF       string
	Balance   float64
	CreatedAt time.Time
}

//NewAccount cria um Account
func NewAccount(ID, name, CPF string, balance float64, createdAt time.Time) Account {
	return Account{
		ID:        ID,
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}

//Deposit adiciona um valor no Balance
func (a *Account) Deposit(amount float64) {
	a.Balance += amount
}

//Withdraw remove um valor no Balance
func (a *Account) Withdraw(amount float64) error {
	if a.Balance < amount {
		return ErrInsufficientBalance
	}

	a.Balance -= amount

	return nil
}
