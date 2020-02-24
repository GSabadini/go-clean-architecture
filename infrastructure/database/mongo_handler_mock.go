package database

import (
	"errors"
)

//MongoHandlerSuccessMock
type MongoHandlerSuccessMock struct{}

//Store
func (m MongoHandlerSuccessMock) Store(_ string, _ interface{}) error {
	return nil
}

//FindAll
func (m MongoHandlerSuccessMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return nil
}

//FindOne
func (m MongoHandlerSuccessMock) FindOne(string, interface{}, interface{}) error {
	return nil
}

//MongoHandlerSuccessMock
type MongoHandlerErrorMock struct{}

//Store
func (m MongoHandlerErrorMock) Store(_ string, _ interface{}) error {
	return errors.New("Error")
}

//FindAll
func (m MongoHandlerErrorMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return errors.New("Error")
}

//FindOne
func (m MongoHandlerErrorMock) FindOne(string, interface{}, interface{}) error {
	return errors.New("Error")
}
