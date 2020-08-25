package database

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/interface/repository"
)

var (
	errInvalidDatabaseInstance = errors.New("invalid db instance")
)

const (
	InstanceMongoDB int = iota
)

//NewDatabaseNoSQLFactory retorna a instância de um banco de dados NoSQL
func NewDatabaseNoSQLFactory(instance int) (repository.NoSQLHandler, error) {
	switch instance {
	case InstanceMongoDB:
		db, err := NewMongoHandler(newConfigMongoDB())
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

//NewDatabaseSQLFactory retorna a instância de um banco de dados SQL
func NewDatabaseSQLFactory(instance int) (repository.SQLHandler, error) {
	switch instance {
	case InstancePostgres:
		db, err := NewPostgresHandler(newConfigPostgres())
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errInvalidDatabaseInstance
	}
}
