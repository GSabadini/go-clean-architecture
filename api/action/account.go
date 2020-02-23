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
)

type AccountAction struct {
	dbHandler database.NoSQLDBHandler
}

func NewAccountAction(dbHandler database.NoSQLDBHandler) AccountAction {
	return AccountAction{dbHandler: dbHandler}
}

func (a AccountAction) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account domain.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	err := usecase.Create(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a AccountAction) Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account []domain.Account
	var accountRepository = repository.NewAccount(a.dbHandler)
	result, err := usecase.FindAll(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

func (a AccountAction) Show(w http.ResponseWriter, r *http.Request) {
	type ReturnBallance struct {
		Ballance float64 `json:"ballance"`
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	accountId, ok := vars["account_id"]
	if !ok {
		fmt.Println("Not Parameter")
		return
	}

	var account domain.Account

	var accountRepository = repository.NewAccount(a.dbHandler)
	result, err := usecase.FindOne(accountRepository, account, accountId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	resultBallance := ReturnBallance{Ballance: result.Ballance}

	if err := json.NewEncoder(w).Encode(resultBallance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}
