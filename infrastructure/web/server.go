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
	dbConnSQL repository.SQLHandler,
	dbConnNoSQL repository.NoSQLHandler,
	validator validator.Validator,
	port Port,
) (Server, error) {
	switch instance {
	case InstanceGorillaMux:
		return NewGorillaMux(log, dbConnSQL, validator, port), nil
	case InstanceGin:
		return NewGinServer(log, dbConnNoSQL, validator, port), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
