package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/interface/logger"
	"github.com/gsabadini/go-bank-transfer/interface/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
)

type CreateAccountAction struct {
	uc        usecase.CreateAccount
	log       logger.Logger
	validator validator.Validator
}

func NewCreateAccountAction(uc usecase.CreateAccount, log logger.Logger, v validator.Validator) CreateAccountAction {
	return CreateAccountAction{uc: uc, log: log, validator: v}
}

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
			"invalid input",
			http.StatusBadRequest,
			errors.New("invalid input"),
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
			"error when creating a new account",
			http.StatusInternalServerError,
			err,
		).Log()

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, "success creating account", http.StatusCreated).Log()

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
