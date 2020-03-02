package database

import "errors"

//MongoHandlerSuccessMock implementa a interface de NoSQLDbHandler com resultados de sucesso
type MongoHandlerSuccessMock struct{}

//Store retorna sucesso ao criar um recurso
func (m MongoHandlerSuccessMock) Store(_ string, _ interface{}) error {
	return nil
}

//Update retorna sucesso ao atualizar um recurso
func (m MongoHandlerSuccessMock) Update(_ string, _ interface{}, _ interface{}) error {
	return nil
}

//FindAll retorna sucesso ao listar recursos
func (m MongoHandlerSuccessMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return nil
}

//FindOne retorna sucesso ao obter um recurso
func (m MongoHandlerSuccessMock) FindOne(_ string, _ interface{}, _ interface{}, _ interface{}) error {
	return nil
}

//MongoHandlerErrorMock implementa a interface de NoSQLDbHandler com resultados de sucesso
type MongoHandlerErrorMock struct{}

//Store retorna erro ao criar um recurso
func (m MongoHandlerErrorMock) Store(_ string, _ interface{}) error {
	return errors.New("Error")
}

//Update retorna erro ao atualizar um recurso
func (m MongoHandlerErrorMock) Update(_ string, _ interface{}, _ interface{}) error {
	return errors.New("Error")
}

//FindAll retorna erro ao listar recursos
func (m MongoHandlerErrorMock) FindAll(_ string, _ interface{}, _ interface{}) error {
	return errors.New("Error")
}

//FindOne retorna erro ao obter um recurso
func (m MongoHandlerErrorMock) FindOne(_ string, _ interface{}, _ interface{}, _ interface{}) error {
	return errors.New("Error")
}
