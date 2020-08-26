package logging

import (
	"github.com/gsabadini/go-bank-transfer/interface/logger"
)

type Info struct {
	log        logger.Logger
	key        string
	msg        string
	httpStatus int
}

func NewInfo(log logger.Logger, key string, msg string, httpStatus int) Info {
	return Info{
		log:        log,
		key:        key,
		msg:        msg,
		httpStatus: httpStatus,
	}
}

func (i Info) Log() {
	i.log.WithFields(logger.Fields{
		"key":         i.key,
		"http_status": i.httpStatus,
	}).Infof(i.msg)
}
