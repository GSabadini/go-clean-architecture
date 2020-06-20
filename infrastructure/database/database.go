package database

import "errors"

//NoSQLHandler expõe os métodos disponíveis para as abstrações de banco NoSQL
type NoSQLHandler interface {
	Store(string, interface{}) error
	Update(string, interface{}, interface{}) error
	FindAll(string, interface{}, interface{}) error
	FindOne(string, interface{}, interface{}, interface{}) error
}

//SQLHandler expõe os métodos disponíveis para as abstrações de banco SQL
type SQLHandler interface {
	Execute(string, ...interface{}) error
	Query(string, ...interface{}) (Row, error)
}

//Row
type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
}

var (
	errInvalidDatabaseInstance = errors.New("invalid database instance")
)

const (
	InstanceMongoDB int = iota
)

//NewDatabaseNoSQL retorna o handler de um banco de dados NoSQL
func NewDatabaseNoSQL(instance int, host, dbName string) (NoSQLHandler, error) {
	switch instance {
	case InstanceMongoDB:
		db, err := NewMongoHandler(host, dbName)
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errInvalidDatabaseInstance
	}
}

const (
	InstancePostgres int = iota
)

//NewDatabaseSQL retorna o handler de um banco de dados SQL
func NewDatabaseSQL(instance int, dataSource string) (SQLHandler, error) {
	switch instance {
	case InstancePostgres:
		db, err := NewPostgresHandler(dataSource)
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errInvalidDatabaseInstance
	}
}
