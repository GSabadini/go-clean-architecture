package logging

import "github.com/gsabadini/go-bank-transfer/infrastructure/logger"

//Info
type Info struct {
	log        logger.Logger
	key        string
	msg        string
	httpStatus int
}

//NewInfo
func NewInfo(log logger.Logger, key string, msg string, httpStatus int) Info {
	return Info{
		log:        log,
		key:        key,
		msg:        msg,
		httpStatus: httpStatus,
	}
}

//Log
func (i Info) Log() {
	i.log.WithFields(logger.Fields{
		"key":         i.key,
		"http_status": i.httpStatus,
	}).Infof(i.msg)
}
