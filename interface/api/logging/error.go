package logging

import (
	"github.com/gsabadini/go-bank-transfer/interface/logger"
)

//Error
type Error struct {
	log        logger.Logger
	key        string
	msg        string
	httpStatus int
	err        error
}

//NewError
func NewError(log logger.Logger, key string, msg string, httpStatus int, err error) Error {
	return Error{log: log, key: key, msg: msg, httpStatus: httpStatus, err: err}
}

//Log
func (e Error) Log() {
	e.log.WithFields(logger.Fields{
		"key":         e.key,
		"http_status": e.httpStatus,
		"error":       e.err.Error(),
	}).Errorf(e.msg)
}
