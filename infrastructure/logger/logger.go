package logger

import "errors"

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

//NewLogger retorna a instancia de um logger
func NewLogger(loggerInstance int) (Logger, error) {
	switch loggerInstance {
	case InstanceZapLogger:
		//logger, err := newZapLogger(config)
		//if err != nil {
		//	return err
		//}
		//log = logger
		return nil, nil

	case InstanceLogrusLogger:
		var logger = NewLogrusLogger()
		return logger, nil

	default:
		return nil, errInvalidLoggerInstance
	}
}
