package domain

import (
	"errors"
	"time"
)

/* TODO rever errors */
var (
	//ErrNotFound é um erro de Account não encontrado
	ErrNotFound = errors.New("not found")

	//ErrInsufficientBalance é um erro de saldo insuficiente
	ErrInsufficientBalance = errors.New("origin account does not have sufficient balance")

	//ErrInsufficientBalance é um erro ao atualizar o saldo de uma conta
	ErrUpdateBalance = errors.New("error update account balance")
)

//AccountRepository expõe os métodos disponíveis para as abstrações do repositório de Account
type AccountRepository interface {
	Store(Account) (Account, error)
	UpdateBalance(AccountID, float64) error
	FindAll() ([]Account, error)
	FindByID(AccountID) (Account, error)
	FindBalance(AccountID) (Account, error)
}

//AccountID define o tipo identificador de uma Account
type AccountID string

//Account armazena a estrutura de uma conta
type Account struct {
	ID        AccountID
	Name      string
	CPF       string
	Balance   float64
	CreatedAt time.Time
}

//NewAccount cria um Account
func NewAccount(ID AccountID, name, CPF string, balance float64, createdAt time.Time) Account {
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
