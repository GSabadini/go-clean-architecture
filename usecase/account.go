package usecase

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

type AccountService struct {
	repository repository.AccountRepository
}

func NewAccountService(repository repository.AccountRepository) AccountService {
	return AccountService{repository: repository}
}

//StoreAccount cria uma nova conta
func (a AccountService) StoreAccount(account domain.Account) (domain.Account, error) {
	t := time.Now()
	account.CreatedAt = &t
	account.ID = bson.NewObjectId()
	account.CPF = a.cleanCPF(account.CPF)

	result, err := a.repository.Store(account)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (a AccountService) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}

//FindAllAccount retorna uma lista de contas
func (a AccountService) FindAllAccount() ([]domain.Account, error) {
	result, err := a.repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindBalanceAccount retorna o saldo de uma conta
func (a AccountService) FindBalanceAccount(ID string) (domain.Account, error) {
	var (
		query    = bson.M{"_id": bson.ObjectIdHex(ID)}
		selector = bson.M{"balance": 1, "_id": 0}
	)

	result, err := a.repository.FindOneWithSelector(query, selector)
	if err != nil {
		return result, err
	}

	return result, nil
}
