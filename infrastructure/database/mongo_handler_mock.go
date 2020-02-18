package database

type MongoHandlerMock struct{}

func (m MongoHandlerMock) Insert(_ string, _ interface{}) error {
	return nil
}
