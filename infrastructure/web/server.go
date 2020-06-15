package web

import (
	"errors"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
)

type Server interface {
	Listen()
}

var (
	errInvalidWebServerInstance = errors.New("invalid web server instance")
)

const (
	InstanceGorillaWithNegroni int = iota
	InstanceGin
)

//NewWebServer
func NewWebServer(
	webServerInstance int,
	log logger.Logger,
	dbConnSQL database.SQLHandler,
	dbConnNoSQL database.NoSQLHandler,
	port int64,
) (Server, error) {
	switch webServerInstance {
	case InstanceGorillaWithNegroni:
		return NewGorillaMux(log, dbConnSQL, dbConnNoSQL, port), nil
	case InstanceGin:
		return NewGin(log, dbConnNoSQL, port), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
