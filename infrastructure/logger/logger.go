package logger

import (
	"errors"
)

type Logger interface {
	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalln(args ...interface{})

	WithFields(keyValues Fields) Logger

	WithError(err error) Logger
}

type Fields map[string]interface{}

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

//NewLogger retorna a instância de um logger
func NewLogger(instance int, isJSON bool) (Logger, error) {
	switch instance {
	case InstanceZapLogger:
		logger, err := NewZapLogger(isJSON)
		if err != nil {
			return nil, err
		}
		return logger, nil
	case InstanceLogrusLogger:
		var logger = NewLogrusLogger(isJSON)
		return logger, nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
