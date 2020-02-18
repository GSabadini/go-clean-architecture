package database

import "errors"

type MongoHandlerSuccessMock struct{}

func (m MongoHandlerSuccessMock) Insert(_ string, _ interface{}) error {
	return nil
}

type MongoHandlerErrorMock struct{}

func (m MongoHandlerErrorMock) Insert(_ string, _ interface{}) error {
	return errors.New("Error")
}
