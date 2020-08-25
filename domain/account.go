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
	ErrInsufficientBalance = errors.New("origin validator does not have sufficient balance")

	//ErrUpdateBalance é um erro ao atualizar o saldo de uma conta
	ErrUpdateBalance = errors.New("error update validator balance")
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
	id        AccountID
	name      string
	cpf       string
	balance   Money
	createdAt time.Time
}

//NewAccount cria um Account somento com o Balance
func NewAccountBalance(balance Money) Account {
	return Account{balance: balance}
}

//NewAccount cria um Account
func NewAccount(ID AccountID, name, CPF string, balance Money, createdAt time.Time) Account {
	return Account{
		id:        ID,
		name:      name,
		cpf:       CPF,
		balance:   balance,
		createdAt: createdAt,
	}
}

//Deposit adiciona um valor no Balance
func (a *Account) Deposit(amount Money) {
	a.balance += amount
}

//Withdraw remove um valor no Balance
func (a *Account) Withdraw(amount Money) error {
	if a.balance < amount {
		return ErrInsufficientBalance
	}

	a.balance -= amount

	return nil
}

//ID
func (a Account) ID() AccountID {
	return a.id
}

//Name
func (a Account) Name() string {
	return a.name
}

//CPF
func (a Account) CPF() string {
	return a.cpf
}

//Balance
func (a Account) Balance() Money {
	return a.balance
}

//CreatedAt
func (a Account) CreatedAt() time.Time {
	return a.createdAt
}
