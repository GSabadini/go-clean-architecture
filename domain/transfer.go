package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Transfer struct {
	Id                   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AccountOriginId      bson.ObjectId `json:"account_origin_id" bson:"account_origin_id"`
	AccountDestinationId bson.ObjectId `json:"account_destination_id" bson:"account_destination_id"`
	Amount               float64       `json:"amount" bson:"amount"`
	CreatedAt            time.Time     `json:"created_at" bson:"created_at"`
}
