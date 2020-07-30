package web

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//Server é uma abstração para o server da aplicação
type Server interface {
	Listen()
}

//Port define uma porta para o servidor
type Port int64

var (
	errInvalidWebServerInstance = errors.New("invalid web server instance")
)

const (
	InstanceGorillaMux int = iota
	InstanceGin
)

//NewWebServerFactory retorna a instância de um web server
func NewWebServerFactory(
	instance int,
	log logger.Logger,
	dbSQL repository.SQLHandler,
	dbNoSQL repository.NoSQLHandler,
	validator validator.Validator,
	port Port,
) (Server, error) {
	switch instance {
	case InstanceGorillaMux:
		return newGorillaMux(log, dbSQL, validator, port), nil
	case InstanceGin:
		return newGinServer(log, dbNoSQL, validator, port), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
