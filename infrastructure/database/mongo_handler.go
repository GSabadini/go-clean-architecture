package database

import (
	"net/url"
	"strings"

	mongo "gopkg.in/mgo.v2"
)

//MongoHandler implementação para banco de dados MongoDb
type MongoHandler struct {
	Database *mongo.Database
	Session  *mongo.Session
}

//Insert realiza uma inserção no banco de dados
func (mgo MongoHandler) Insert(collection string, data interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Insert(data)
}

//NewMongoHandler constrói um novo handler de banco para MongoDb
func NewMongoHandler(uri string) (*MongoHandler, error) {
	session, err := mongo.Dial(uri)
	if err != nil {
		return &MongoHandler{}, err
	}

	handler := new(MongoHandler)
	handler.Session = session
	handler.Database = handler.Session.DB(getSchema(uri))

	return handler, nil
}

func getSchema(uri string) string {
	mongouri, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	return strings.Replace(mongouri.EscapedPath(), "/", "", -1)
}
