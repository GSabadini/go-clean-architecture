package database

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/repository"
)

var (
	errInvalidDatabaseInstance = errors.New("invalid database instance")
)

const (
	InstanceMongoDB int = iota
)

//NewDatabaseNoSQL retorna o handler de um banco de dados NoSQL
func NewDatabaseNoSQL(instance int, host, dbName string) (repository.NoSQLHandler, error) {
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
func NewDatabaseSQL(instance int, dataSource string) (repository.SQLHandler, error) {
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
