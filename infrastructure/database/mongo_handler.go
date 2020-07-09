package database

import mongo "gopkg.in/mgo.v2"

//mongoHandler armazena a estrutura para MongoDB
type mongoHandler struct {
	database *mongo.Database
	session  *mongo.Session
}

//NewMongoHandler constrói um novo handler de banco para MongoDB
func NewMongoHandler(host, databaseName string) (*mongoHandler, error) {
	session, err := mongo.Dial(host)
	if err != nil {
		return &mongoHandler{}, err
	}

	handler := new(mongoHandler)
	handler.session = session
	handler.database = handler.session.DB(databaseName)

	return handler, nil
}

//Store realiza uma inserção no banco de dados
func (mgo mongoHandler) Store(collection string, data interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Insert(data)
}

//Update realiza uma atualização no banco de dados
func (mgo mongoHandler) Update(collection string, query interface{}, update interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Update(query, update)
}

//FindAll realiza uma busca por todos os registros no banco de dados
func (mgo mongoHandler) FindAll(collection string, query interface{}, result interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).All(result)
}

//FindOne realiza a busca de um item específico no banco de dados
func (mgo mongoHandler) FindOne(collection string, query interface{}, selector interface{}, result interface{}) error {
	session := mgo.session.Clone()
	defer session.Close()

	return mgo.database.C(collection).With(session).Find(query).Select(selector).One(result)
}
