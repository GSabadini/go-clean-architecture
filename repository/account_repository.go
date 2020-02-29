package repository

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

const accountsCollectionName = "accounts"

//Account representa um repositório para dados de uma conta
type Account DbRepository

//NewAccount constrói um repository com suas dependências
func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store realiza uma inserção no banco de dados através da implementação real do database
func (a Account) Store(account domain.Account) (domain.Account, error) {
	t := time.Now()
	account.CreatedAt = &t
	account.ID = bson.NewObjectId()

	if err := a.dbHandler.Store(accountsCollectionName, &account); err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

//Update realiza uma atualização no banco de dados através da implementação real do database
func (a Account) Update(query bson.M, update bson.M) error {
	return a.dbHandler.Update(accountsCollectionName, query, update)
}

//FindAll realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindAll() ([]domain.Account, error) {
	var account = make([]domain.Account, 0)

	if err := a.dbHandler.FindAll(accountsCollectionName, nil, &account); err != nil {
		return []domain.Account{}, err
	}

	return account, nil
}

//FindOne realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindOne(query bson.M) (*domain.Account, error) {
	var account = &domain.Account{}

	if err := a.dbHandler.FindOne(accountsCollectionName, query, nil, &account); err != nil {
		return account, err
	}

	return account, nil
}

//FindOneWithSelector realiza uma busca no banco de dados através da implementação real do database
func (a Account) FindOneWithSelector(query bson.M, selector interface{}) (domain.Account, error) {
	var account = domain.Account{}

	if err := a.dbHandler.FindOne(accountsCollectionName, query, selector, &account); err != nil {
		return account, err
	}

	return account, nil
}

//AccountRepositoryMockSuccess
type AccountRepositoryMockSuccess struct{}

//Store
func (a AccountRepositoryMockSuccess) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{
		ID:        "5e570851adcef50116aa7a5c",
		Name:      "Test",
		CPF:       "028.155.170-78",
		Balance:   100,
		CreatedAt: nil,
	}, nil
}

//Update
func (a AccountRepositoryMockSuccess) Update(_ bson.M, _ bson.M) error {
	return nil
}

//FindAll
func (a AccountRepositoryMockSuccess) FindAll() ([]domain.Account, error) {
	var account = []domain.Account{
		{
			ID:      "5e570851adcef50116aa7a5c",
			Name:    "Test-0",
			CPF:     "028.155.170-78",
			Balance: 0,
		},
		{
			ID:      "5e570854adcef50116aa7a5d",
			Name:    "Test-1",
			CPF:     "028.155.170-78",
			Balance: 50.25,
		},
	}

	return account, nil
}

//FindOne
func (a AccountRepositoryMockSuccess) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{
		ID:      "5e570854adcef50116aa7a5d",
		Name:    "Test-1",
		CPF:     "028.155.170-78",
		Balance: 50.25,
	}, nil
}

//FindOneWithSelector
func (a AccountRepositoryMockSuccess) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{
		Balance: 100.00,
	}, nil
}

//AccountRepositoryMockError
type AccountRepositoryMockError struct{}

//Store
func (a AccountRepositoryMockError) Store(_ domain.Account) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}

//Update
func (a AccountRepositoryMockError) Update(_ bson.M, _ bson.M) error {
	return errors.New("Error")
}

//FindAll
func (a AccountRepositoryMockError) FindAll() ([]domain.Account, error) {
	return []domain.Account{}, errors.New("Error")
}

//FindOne
func (a AccountRepositoryMockError) FindOne(_ bson.M) (*domain.Account, error) {
	return &domain.Account{}, errors.New("Error")
}

//FindOneWithSelector
func (a AccountRepositoryMockError) FindOneWithSelector(_ bson.M, _ interface{}) (domain.Account, error) {
	return domain.Account{}, errors.New("Error")
}
