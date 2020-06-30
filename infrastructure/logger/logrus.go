package logger

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Logger
}

//NewLogrusLogger constrói um logrus logger
func NewLogrusLogger(isJSON bool) Logger {
	log := logrus.New()

	if isJSON {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return &logrusLogger{logger: log}
}

//Infof
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

//Warnf
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

//Errorf
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

//Fatalln
func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

//WithFields
func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

//WithError
func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithError(err),
	}
}

type logrusLogEntry struct {
	entry *logrus.Entry
}

//Infof
func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

//Warnf
func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

//Errorf
func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

//Fatalln
func (l *logrusLogEntry) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

//WithFields
func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

//WithError
func (l *logrusLogEntry) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithError(err),
	}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}

	return logrusFields
}
