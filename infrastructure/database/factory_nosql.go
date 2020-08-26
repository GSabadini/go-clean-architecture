package database

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/interface/repository"
)

var (
	errInvalidNoSQLDatabaseInstance = errors.New("invalid nosql db instance")
)

const (
	InstanceMongoDB int = iota
)

func NewDatabaseNoSQLFactory(instance int) (repository.NoSQLHandler, error) {
	switch instance {
	case InstanceMongoDB:
		return NewMongoHandler(newConfigMongoDB())
	default:
		return nil, errInvalidNoSQLDatabaseInstance
	}
}
