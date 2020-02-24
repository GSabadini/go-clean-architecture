package domain

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Cpf       string        `json:"cpf" bson:"cpf"`
	Balance   float64       `json:"balance" bson:"balance"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) Deposit(amount float64) {
	a.Balance += amount
}

func (a *Account) Withdraw(amount float64) error {
	if a.Balance < amount {
		return errors.New("source account does not have sufficient balance")
	}

	a.Balance -= amount

	return nil
}
