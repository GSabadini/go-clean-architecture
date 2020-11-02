package log

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/adapter/logger"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid log instance")
)

func NewLoggerFactory(instance int) (logger.Logger, error) {
	switch instance {
	case InstanceZapLogger:
		return NewZapLogger()
	case InstanceLogrusLogger:
		return NewLogrusLogger(), nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
