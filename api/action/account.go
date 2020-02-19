package action

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type Account struct {
	dbHandler database.NoSQLDBHandler
}

func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

func (a Account) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var account domain.Account
	if err = json.Unmarshal(reqBody, &account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	err = usecase.Create(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))
}

func (a Account) Index(w http.ResponseWriter, _ *http.Request) {
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
