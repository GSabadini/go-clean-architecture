package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(level string, isJSON bool) (Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := log.Sugar()
	defer log.Sync()

	var zapCoreLevel zapcore.Level
	if level != "" {
		zapCoreLevel = getZapLevel(level)
	}

	var encoderConfig zapcore.EncoderConfig
	if isJSON {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapcore.NewJSONEncoder(encoderConfig)
	}

	var cores []zapcore.Core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zapCoreLevel,
	)
	cores = append(cores, core)

	return &zapLogger{logger: sugar}, nil
}

//Infof
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

//Warnf
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

//Errorf
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

//Fatalln
func (l *zapLogger) Fatalln(args ...interface{}) {
	l.logger.Fatal(args)
}

//WithFields
func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}

	log := l.logger.With(f...)
	return &zapLogger{logger: log}
}

//WithError
func (l *zapLogger) WithError(err error) Logger {
	var log = l.logger.With(err.Error())
	return &zapLogger{logger: log}
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
