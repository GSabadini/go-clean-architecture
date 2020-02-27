package database

import (
	mongo "gopkg.in/mgo.v2"
)

//MongoHandler armazena a estrutura para MongoDb
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

//Update realiza uma atualização no banco de dados
func (mgo MongoHandler) Update(collection string, query interface{}, update interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Update(query, update)
}

//FindAll realiza uma busca no banco de dados por todos os registros
func (mgo MongoHandler) FindAll(collection string, query interface{}, result interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Find(query).All(result)
}

//FindOne realiza uma busca específica no banco de dados
func (mgo MongoHandler) FindOne(collection string, query interface{}, selector interface{}, result interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Find(query).Select(selector).One(result)
}
