package action

import (
	"encoding/json"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Account armazena as dependências de uma conta
type Account struct {
	dbHandler database.NoSQLDBHandler
	logger    *logrus.Logger
}

//NewAccount constrói uma conta com suas dependências
func NewAccount(dbHandler database.NoSQLDBHandler, log *logrus.Logger) Account {
	return Account{dbHandler: dbHandler, logger: log}
}

//Store é um handler para criação de conta
func (a Account) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_account"
	var account *domain.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		a.logError(
			logKey,
			"error when decoding json",
			http.StatusInternalServerError,
			err,
		)
		ErrInternalServer.Send(w)
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.StoreAccount(accountRepository, account)
	if err != nil {
		a.logError(
			logKey,
			"error when creating a new account",
			http.StatusInternalServerError,
			err,
		)
		ErrInternalServer.Send(w)
		return
	}

	a.logInfoSuccess(logKey, "success create account", http.StatusCreated)

	Success(result, http.StatusCreated).Send(w)
}

//Index é um handler para retornar a lista de contas
func (a Account) Index(w http.ResponseWriter, _ *http.Request) {
	const logKey = "index_account"
	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindAllAccount(accountRepository)
	if err != nil {
		a.logError(
			logKey,
			"error when returning account list",
			http.StatusInternalServerError,
			err,
		)
		ErrInternalServer.Send(w)
		return
	}

	a.logInfoSuccess(logKey, "success return list accounts", http.StatusOK)

	Success(result, http.StatusOK).Send(w)
}

type ReturnBalance struct {
	Balance float64 `json:"balance"`
}

//ShowBalance é um handler para retornar o saldo de uma conta
func (a Account) ShowBalance(w http.ResponseWriter, r *http.Request) {
	const logKey = "show_balance"
	var vars = mux.Vars(r)

	accountId, ok := vars["account_id"]
	if !ok {
		a.logError(
			logKey,
			"not parameter",
			http.StatusNotFound,
			nil,
		)

		ErrNotFound.Send(w)
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindOneAccount(accountRepository, accountId)
	if err != nil {
		a.logError(
			logKey,
			"error when returning account balance",
			http.StatusInternalServerError,
			err,
		)

		ErrInternalServer.Send(w)
		return
	}

	resultBalance := ReturnBalance{Balance: result.Balance}

	a.logInfoSuccess(logKey, "success when returning account balance", http.StatusOK)

	Success(resultBalance, http.StatusOK).Send(w)
}

func (a Account) logInfoSuccess(key string, description string, httpStatus int) {
	a.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
		"description": description,
	}).Info()
}

func (a Account) logError(key string, description string, httpStatus int, err error) {
	a.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
		"description": description,
		"error":       err.Error(),
	}).Error()
}
