package database

import (
	"errors"
)

type MongoHandlerSuccessMock struct{}

func (m MongoHandlerSuccessMock) FindOne(string, interface{}, interface{}) error {
	return nil
}

func (m MongoHandlerSuccessMock) Store(_ string, _ interface{}) error {
	return nil
}

func (m MongoHandlerSuccessMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return nil
}

type MongoHandlerErrorMock struct{}

func (m MongoHandlerErrorMock) Store(_ string, _ interface{}) error {
	return errors.New("Error")
}

func (m MongoHandlerErrorMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return errors.New("Error")
}

func (m MongoHandlerErrorMock) FindOne(string, interface{}, interface{}) error {
	return errors.New("Error")
}
