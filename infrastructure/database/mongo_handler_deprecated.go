package database

import (
	"context"

	mongo "gopkg.in/mgo.v2"
)

type mongoHandlerDeprecated struct {
	database *mongo.Database
	session  *mongo.Session
}

func NewMongoHandlerDeprecated(c *config) (*mongoHandlerDeprecated, error) {
	session, err := mongo.DialWithTimeout(c.host, c.ctxTimeout)
	if err != nil {
		return &mongoHandlerDeprecated{}, err
	}

	handler := new(mongoHandlerDeprecated)
	handler.session = session
	handler.database = handler.session.DB(c.database)

	return handler, nil
}

func (mgo mongoHandlerDeprecated) Store(_ context.Context, collection string, data interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Insert(data)
}

func (mgo mongoHandlerDeprecated) Update(_ context.Context, collection string, query interface{}, update interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Update(query, update)
}

func (mgo mongoHandlerDeprecated) FindAll(_ context.Context, collection string, query interface{}, result interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).All(result)
}

func (mgo mongoHandlerDeprecated) FindOne(
	_ context.Context,
	collection string,
	query interface{},
	selector interface{},
	result interface{},
) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).Select(selector).One(result)
}
