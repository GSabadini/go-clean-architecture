package database

import "errors"

//MongoHandlerSuccessStub implementa a interface de NoSQLDbHandler com resultados de sucesso
type MongoHandlerSuccessStub struct{}

//Store retorna sucesso ao criar um recurso
func (m MongoHandlerSuccessStub) Store(_ string, _ interface{}) error {
	return nil
}

//Update retorna sucesso ao atualizar um recurso
func (m MongoHandlerSuccessStub) Update(_ string, _ interface{}, _ interface{}) error {
	return nil
}

//FindAll retorna sucesso ao listar recursos
func (m MongoHandlerSuccessStub) FindAll(_ string, _ interface{}, _ interface{}) error {
	return nil
}

//FindOne retorna sucesso ao obter um recurso
func (m MongoHandlerSuccessStub) FindOne(_ string, _ interface{}, _ interface{}, _ interface{}) error {
	return nil
}

//MongoHandlerErrorStub implementa a interface de NoSQLDbHandler com resultados de sucesso
type MongoHandlerErrorStub struct {
	TypeErr error
}

//Store retorna erro ao criar um recurso
func (m MongoHandlerErrorStub) Store(_ string, _ interface{}) error {
	return errors.New("error")
}

//Update retorna erro ao atualizar um recurso
func (m MongoHandlerErrorStub) Update(_ string, _ interface{}, _ interface{}) error {
	return errors.New("error")
}

//FindAll retorna erro ao listar recursos
func (m MongoHandlerErrorStub) FindAll(_ string, _ interface{}, _ interface{}) error {
	return errors.New("error")
}

//FindOne retorna erro ao obter um recurso
func (m MongoHandlerErrorStub) FindOne(_ string, _ interface{}, _ interface{}, _ interface{}) error {
	if m.TypeErr != nil {
		return m.TypeErr
	}
	return errors.New("error")
}
