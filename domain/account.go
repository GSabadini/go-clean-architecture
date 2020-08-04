package domain

import (
	"context"
	"errors"
	"time"
)

/* TODO rever errors */
var (
	//ErrNotFound é um erro de Account não encontrado
	ErrNotFound = errors.New("not found")

	//ErrInsufficientBalance é um erro de saldo insuficiente
	ErrInsufficientBalance = errors.New("origin account does not have sufficient balance")

	//ErrUpdateBalance é um erro ao atualizar o saldo de uma conta
	ErrUpdateBalance = errors.New("error update account balance")
)

//AccountRepository expõe os métodos disponíveis para as abstrações do repositório de Account
type AccountRepository interface {
	Store(context.Context, Account) (Account, error)
	UpdateBalance(context.Context, AccountID, Money) error
	FindAll(context.Context) ([]Account, error)
	FindByID(context.Context, AccountID) (Account, error)
	FindBalance(context.Context, AccountID) (Account, error)
}

//AccountID define o tipo identificador de uma Account
type AccountID string

//String converte o tipo AccountID para uma string
func (a AccountID) String() string {
	return string(a)
}

//Account armazena a estrutura de uma conta
type Account struct {
	ID        AccountID
	Name      string
	CPF       string
	Balance   Money
	CreatedAt time.Time
}

//NewAccount cria um Account
func NewAccount(ID AccountID, name, CPF string, balance Money, createdAt time.Time) Account {
	return Account{
		ID:        ID,
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}

//Deposit adiciona um valor no Balance
func (a *Account) Deposit(amount Money) {
	a.Balance += amount
}

//Withdraw remove um valor no Balance
func (a *Account) Withdraw(amount Money) error {
	if a.Balance < amount {
		return ErrInsufficientBalance
	}

	a.Balance -= amount

	return nil
}
