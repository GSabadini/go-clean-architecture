package database

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/interface/repository"
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

const (
	InstancePostgres int = iota
)

func NewDatabaseSQLFactory(instance int) (repository.SQLHandler, error) {
	switch instance {
	case InstancePostgres:
		return NewPostgresHandler(newConfigPostgres())
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
}
