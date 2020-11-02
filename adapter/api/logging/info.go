package logging

import (
	"github.com/gsabadini/go-bank-transfer/adapter/logger"
)

type Info struct {
	log        logger.Logger
	key        string
	httpStatus int
}

func NewInfo(log logger.Logger, key string, httpStatus int) Info {
	return Info{
		log:        log,
		key:        key,
		httpStatus: httpStatus,
	}
}

func (i Info) Log(msg string) {
	i.log.WithFields(logger.Fields{
		"key":         i.key,
		"http_status": i.httpStatus,
	}).Infof(msg)
}
