package action

import (
	"encoding/json"
	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/repository"
	"github.com/gsabadini/go-stone/usecase"
	"io/ioutil"
	"net/http"
)

type AccountAction struct {
	dbHandler string
}

func NewAccountAction(dbHandler string) AccountAction {
	return AccountAction{dbHandler: dbHandler}
}

func (a AccountAction) CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	var accountRepository = repository.NewAccountRepository(a.dbHandler)

	err = usecase.CreateAccount(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))
}

