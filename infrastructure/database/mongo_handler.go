package database

import (
	mongo "gopkg.in/mgo.v2"
)

//MongoHandler implementação para banco de dados MongoDb
type MongoHandler struct {
	Database *mongo.Database
	Session  *mongo.Session
}

//NewMongoHandler constrói um novo handler de banco para MongoDb
func NewMongoHandler(host, databaseName string) (*MongoHandler, error) {
	session, err := mongo.Dial(host)
	if err != nil {
		return &MongoHandler{}, err
	}

	handler := new(MongoHandler)
	handler.Session = session
	handler.Database = handler.Session.DB(databaseName)

	return handler, nil
}

//Store realiza uma inserção no banco de dados
func (mgo MongoHandler) Store(collection string, data interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Insert(data)
}

//FindAll realiza uma inserção no banco de dados
func (mgo MongoHandler) FindAll(collection string, query interface{}, data interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Find(query).All(data)
}
