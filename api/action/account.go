package action

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	err := usecase.Store(accountRepository, account)
	if err != nil {
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//Index é um handler para retornar a lista de accounts
func (a Account) Index(w http.ResponseWriter, _ *http.Request) {
	var accountRepository = repository.NewAccount(a.dbHandler)

	result, err := usecase.FindAll(accountRepository)
	if err != nil {
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}

	if err := Success(result, http.StatusOK).Send(w); err != nil {
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}
}

type ReturnBallance struct {
	Ballance float64 `json:"ballance"`
}

//Show é um handler para buscar uma account específica
func (a Account) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountId, ok := vars["account_id"]
	if !ok {
		fmt.Println("Not Parameter")
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	result, err := usecase.FindOne(accountRepository, accountId)
	if err != nil {
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}

	resultBallance := ReturnBallance{Ballance: result.Ballance}

	if err := Success(resultBallance, http.StatusOK).Send(w); err != nil {
		log.Println([]byte(err.Error()))
		ErrInternalServer.Send(w)
		return
	}
}
