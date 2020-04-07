package database

import mongo "gopkg.in/mgo.v2"

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

func (mgo MongoHandler) Store(collection string, args ...interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	return mgo.Database.C(collection).With(session).Insert(args...)
}

func (mgo MongoHandler) Update(collection string, args ...interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	var (
		query  = args[0]
		update = args[1]
	)

	return mgo.Database.C(collection).With(session).Update(query, update)
}

func (mgo MongoHandler) FindAll(collection string, args ...interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	var (
		query  = args[0]
		result = args[1]
	)

	return mgo.Database.C(collection).With(session).Find(query).All(result)
}

func (mgo MongoHandler) FindOne(collection string, args ...interface{}) error {
	session := mgo.Session.Clone()
	defer session.Close()

	var (
		query    = args[0]
		selector = args[1]
		result   = args[2]
	)

	return mgo.Database.C(collection).With(session).Find(query).Select(selector).One(result)
}
