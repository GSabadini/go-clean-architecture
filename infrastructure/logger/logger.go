package logger

import "errors"

//Logger é uma abstração para os loggers da aplicação
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

//NewLoggerFactory retorna a instância de um logger
func NewLoggerFactory(instance int) (Logger, error) {
	switch instance {
	case InstanceZapLogger:
		logger, err := NewZapLogger()
		if err != nil {
			return nil, err
		}
		return logger, nil
	case InstanceLogrusLogger:
		var logger = NewLogrusLogger()
		return logger, nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
