package logger

type LoggerMock struct{}

func (l LoggerMock) Infof(_ string, _ ...interface{}) {}

func (l LoggerMock) Warnf(_ string, _ ...interface{}) {}

func (l LoggerMock) Errorf(_ string, _ ...interface{}) {}

func (l LoggerMock) Fatalln(_ ...interface{}) {}

func (l LoggerMock) WithFields(_ Fields) Logger {
	return LoggerEntryMock{}
}

func (l LoggerMock) WithError(_ error) Logger {
	return LoggerEntryMock{}
}

type LoggerEntryMock struct{}

func (l LoggerEntryMock) Infof(_ string, _ ...interface{}) {}

func (l LoggerEntryMock) Warnf(_ string, _ ...interface{}) {}

func (l LoggerEntryMock) Errorf(_ string, _ ...interface{}) {}

func (l LoggerEntryMock) Fatalln(_ ...interface{}) {}

func (l LoggerEntryMock) WithFields(_ Fields) Logger {
	return LoggerEntryMock{}
}

func (l LoggerEntryMock) WithError(_ error) Logger {
	return LoggerEntryMock{}
}
