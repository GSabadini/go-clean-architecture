package domain

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

//Account armazena a estrutura de uma conta
type Account struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name      string        `json:"name,omitempty" bson:"name"`
	CPF       string        `json:"cpf,omitempty" bson:"cpf"`
	Balance   float64       `json:"balance" bson:"balance"`
	CreatedAt *time.Time    `json:"created_at,omitempty" bson:"created_at"`
}

//ValidateBalance verifica se o saldo é valido
func (a *Account) ValidateBalance() error {
	if a.Balance < 0 {
		return errors.New("balance invalid")
	}

	return nil
}

//GetBalance retorna o saldo
func (a *Account) GetBalance() float64 {
	return a.Balance
}

//Deposit adiciona um valor no saldo
func (a *Account) Deposit(amount float64) {
	a.Balance += amount
}

//Withdraw remove um valor do saldo
func (a *Account) Withdraw(amount float64) error {
	if a.Balance < amount {
		return errors.New("source account does not have sufficient balance")
	}

	a.Balance -= amount

	return nil
}
