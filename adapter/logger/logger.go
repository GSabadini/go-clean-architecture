package logger

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalln(args ...interface{})
	WithFields(keyValues Fields) Logger
	WithError(err error) Logger
}

type Fields map[string]interface{}
