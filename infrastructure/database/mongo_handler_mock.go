package database

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type MongoHandlerSuccessMock struct{}

func (m MongoHandlerSuccessMock) Store(_ string, _ interface{}) error {
	return nil
}

func (m MongoHandlerSuccessMock) FindAll(_ string, account []domain.Account) ([]domain.Account, error) {
	return account, nil
}

type MongoHandlerErrorMock struct{}

func (m MongoHandlerErrorMock) Store(_ string, _ interface{}) error {
	return errors.New("Error")
}

func (m MongoHandlerErrorMock) FindAll(_ string, _ []domain.Account) ([]domain.Account, error) {
	return nil, errors.New("Error")
}
