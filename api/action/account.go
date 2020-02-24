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

type Account struct {
	dbHandler database.NoSQLDBHandler
}

func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

//Store é um handler para criação de account
func (a Account) Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account domain.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	err := usecase.Store(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//Index é um handler para retornar a lista de accounts
func (a Account) Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var accountRepository = repository.NewAccount(a.dbHandler)
	result, err := usecase.FindAll(accountRepository)
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

type ReturnBallance struct {
	Ballance float64 `json:"ballance"`
}

//Show é um handler para buscar uma account específica
func (a Account) Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	accountId, ok := vars["account_id"]
	if !ok {
		fmt.Println("Not Parameter")
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	result, err := usecase.FindOne(accountRepository, accountId)
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
