package action

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Account armazena as dependências de uma conta
type Account struct {
	dbHandler database.NoSQLDbHandler
	logger    *logrus.Logger
}

//NewAccount constrói uma conta com suas dependências
func NewAccount(dbHandler database.NoSQLDbHandler, log *logrus.Logger) Account {
	return Account{dbHandler: dbHandler, logger: log}
}

//Store é um handler para criação de conta
func (a Account) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_account"

	var account domain.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		a.logError(
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		)

		ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.StoreAccount(accountRepository, account)
	if err != nil {
		a.logError(
			logKey,
			"error when creating a new account",
			http.StatusInternalServerError,
			err,
		)

		ErrorMessage(err, http.StatusInternalServerError).Send(w)
		return
	}

	a.logSuccess(logKey, "success creating account", http.StatusCreated)

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

		ErrorMessage(err, http.StatusInternalServerError).Send(w)
		return
	}

	a.logSuccess(logKey, "success when returning account list", http.StatusOK)

	Success(result, http.StatusOK).Send(w)
}

//FindBalance é um handler para retornar o saldo de uma conta
func (a Account) FindBalance(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_balance"

	var vars = mux.Vars(r)
	accountID, ok := vars["account_id"]
	if !ok || !bson.IsObjectIdHex(accountID) {
		var err = ErrParameterInvalid

		a.logError(
			logKey,
			"parameter invalid",
			http.StatusNotFound,
			err,
		)

		ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindBalanceAccount(accountRepository, accountID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			a.logError(
				logKey,
				"error fetching account",
				http.StatusBadRequest,
				err,
			)

			ErrorMessage(err, http.StatusBadRequest).Send(w)
			return
		default:
			a.logError(
				logKey,
				"error when returning account balance",
				http.StatusInternalServerError,
				err,
			)

			ErrorMessage(err, http.StatusInternalServerError).Send(w)
			return
		}
	}

	a.logSuccess(logKey, "success when returning account balance", http.StatusOK)

	Success(result, http.StatusOK).Send(w)
}

func (a Account) logSuccess(key string, message string, httpStatus int) {
	a.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
	}).Info(message)
}

func (a Account) logError(key string, message string, httpStatus int, err error) {
	a.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
		"error":       err.Error(),
	}).Error(message)
}
