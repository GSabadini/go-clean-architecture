package action

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/sirupsen/logrus"
)

//Account armazena as dependências de uma conta
type Account struct {
	usecase usecase.AccountUseCase
	logger  *logrus.Logger
}

//NewAccount constrói uma conta com suas dependências
func NewAccount(usecase usecase.AccountUseCase, log *logrus.Logger) Account {
	return Account{usecase: usecase, logger: log}
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

	result, err := a.usecase.Store(account)
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

	result, err := a.usecase.FindAll()
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
	if !ok || !domain.IsValidUUID(accountID) {
		var err = errParameterInvalid

		a.logError(
			logKey,
			"parameter invalid",
			http.StatusBadRequest,
			err,
		)

		ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	result, err := a.usecase.FindBalance(accountID)
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
