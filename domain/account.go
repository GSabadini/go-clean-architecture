package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Account struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Cpf       string        `json:"cpf" bson:"cpf"`
	Ballance  int           `json:"ballance" bson:"ballance"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}
