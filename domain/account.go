package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Cpf       string        `json:"cpf" bson:"cpf"`
	Balance   float64       `json:"balance" bson:"balance"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

func (a *Account) SumBalance(amount float64) {
	a.Balance += amount
}

func (a *Account) SubtractBalance(amount float64) {
	a.Balance -= amount
}
