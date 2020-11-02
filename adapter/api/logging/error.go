package logging

import (
	"github.com/gsabadini/go-bank-transfer/adapter/logger"
)

type Error struct {
	log        logger.Logger
	err        error
	key        string
	httpStatus int
}

func NewError(log logger.Logger, err error, key string, httpStatus int) Error {
	return Error{
		log:        log,
		err:        err,
		key:        key,
		httpStatus: httpStatus,
	}
}

func (e Error) Log(msg string) {
	e.log.WithFields(logger.Fields{
		"key":         e.key,
		"error":       e.err.Error(),
		"http_status": e.httpStatus,
	}).Errorf(msg)
}
