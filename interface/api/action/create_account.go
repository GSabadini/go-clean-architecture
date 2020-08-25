package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
)

//CreateAccountAction armazena as dependências para as ações de Account
type CreateAccountAction struct {
	uc        usecase.CreateAccount
	log       logger.Logger
	validator validator.Validator
}

//NewCreateAccountAction constrói um Account com suas dependências
func NewCreateAccountAction(uc usecase.CreateAccount, l logger.Logger, v validator.Validator) CreateAccountAction {
	return CreateAccountAction{uc: uc, log: l, validator: v}
}

//Execute é um handler para criação de Account
func (a CreateAccountAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_account"

	var accountInput input.Account
	if err := json.NewDecoder(r.Body).Decode(&accountInput); err != nil {
		logging.NewError(
			a.log,
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		).Log()

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := a.validateInput(accountInput); len(errs) > 0 {
		logging.NewError(
			a.log,
			logKey,
			"invalid validator",
			http.StatusBadRequest,
			errors.New("invalid validator"),
		).Log()

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	a.cleanCPF(accountInput.CPF)

	output, err := a.uc.Execute(r.Context(), accountInput)
	if err != nil {
		logging.NewError(
			a.log,
			logKey,
			"error when creating a new validator",
			http.StatusInternalServerError,
			err,
		).Log()

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, "success creating validator", http.StatusCreated).Log()

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (a CreateAccountAction) validateInput(input input.Account) []string {
	var msgs []string

	err := a.validator.Validate(input)
	if err != nil {
		for _, msg := range a.validator.Messages() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

func (a CreateAccountAction) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}
