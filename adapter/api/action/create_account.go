package action

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/adapter/api/logging"
	"github.com/gsabadini/go-bank-transfer/adapter/api/response"
	"github.com/gsabadini/go-bank-transfer/adapter/logger"
	"github.com/gsabadini/go-bank-transfer/adapter/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type CreateAccountAction struct {
	uc        usecase.CreateAccountUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateAccountAction(uc usecase.CreateAccountUseCase, log logger.Logger, v validator.Validator) CreateAccountAction {
	return CreateAccountAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a CreateAccountAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_account"

	var input usecase.CreateAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error when decoding json")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := a.validateInput(input); len(errs) > 0 {
		logging.NewError(
			a.log,
			response.ErrInvalidInput,
			logKey,
			http.StatusBadRequest,
		).Log("invalid input")

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	a.cleanCPF(input.CPF)

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when creating a new account")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success creating account")

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (a CreateAccountAction) validateInput(input usecase.CreateAccountInput) []string {
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
