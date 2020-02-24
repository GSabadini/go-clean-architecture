package action

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Account struct {
	dbHandler database.NoSQLDBHandler
	logger    *logrus.Logger
}

func NewAccount(dbHandler database.NoSQLDBHandler, log *logrus.Logger) Account {
	return Account{dbHandler: dbHandler, logger: log}
}

//Store é um handler para criação de account
func (a Account) Store(w http.ResponseWriter, r *http.Request) {
	var account domain.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		a.logger.WithField("error", err).Error("error when creating a new account")
		ErrInternalServer.Send(w)
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.Store(accountRepository, account)
	if err != nil {
		a.logger.WithField("error", err).Error("error when creating a new account")
		ErrInternalServer.Send(w)
		return
	}

	Success(result, http.StatusCreated).Send(w)
}

//Index é um handler para retornar a lista de accounts
func (a Account) Index(w http.ResponseWriter, _ *http.Request) {
	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindAll(accountRepository)
	if err != nil {
		a.logger.WithField("error", err).Error("error fetching account list")
		ErrInternalServer.Send(w)
		return
	}

	Success(result, http.StatusOK).Send(w)
}

type ReturnBallance struct {
	Ballance float64 `json:"ballance"`
}

//ShowBallance é um handler para buscar o ballance de uma account
func (a Account) ShowBallance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountId, ok := vars["account_id"]
	if !ok {
		fmt.Println("Not Parameter")
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindOne(accountRepository, accountId)
	if err != nil {
		a.logger.WithField("error", err).Error("error fetching account balance")
		ErrInternalServer.Send(w)
		return
	}

	resultBallance := ReturnBallance{Ballance: result.Ballance}

	Success(resultBallance, http.StatusOK).Send(w)
}
