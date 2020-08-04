package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/api/response"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

//accountInput armazena a estrutura de dados de entrada da API
type accountInput struct {
	Name    string `json:"name" validate:"required"`
	CPF     string `json:"cpf" validate:"required"`
	Balance int64  `json:"balance" validate:"gt=0,required"`
}

//Account armazena as dependências para as ações de Account
type Account struct {
	uc        usecase.AccountUseCase
	log       logger.Logger
	validator validator.Validator
}

//NewAccount constrói um Account com suas dependências
func NewAccount(uc usecase.AccountUseCase, l logger.Logger, v validator.Validator) Account {
	return Account{uc: uc, log: l, validator: v}
}

//Store é um handler para criação de Account
func (a Account) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_account"

	var input accountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		a.logError(
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := a.validateInput(input); len(errs) > 0 {
		a.logError(
			logKey,
			"invalid input",
			http.StatusBadRequest,
			errors.New("invalid input"),
		)

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.Store(
		r.Context(),
		input.Name,
		a.cleanCPF(input.CPF),
		domain.Money(input.Balance),
	)
	if err != nil {
		a.logError(
			logKey,
			"error when creating a new account",
			http.StatusInternalServerError,
			err,
		)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	a.logSuccess(logKey, "success creating account", http.StatusCreated)

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

//Index é um handler para retornar todas as Account
func (a Account) Index(w http.ResponseWriter, r *http.Request) {
	const logKey = "index_account"

	output, err := a.uc.FindAll(r.Context())
	if err != nil {
		a.logError(
			logKey,
			"error when returning account list",
			http.StatusInternalServerError,
			err,
		)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	a.logSuccess(logKey, "success when returning account list", http.StatusOK)

	response.NewSuccess(output, http.StatusOK).Send(w)
}

//FindBalance é um handler para retornar o Balance de uma Account
func (a Account) FindBalance(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_balance"

	var accountID = r.URL.Query().Get("account_id")
	if !domain.IsValidUUID(accountID) {
		var err = response.ErrParameterInvalid
		a.logError(
			logKey,
			"parameter invalid",
			http.StatusBadRequest,
			err,
		)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.FindBalance(r.Context(), domain.AccountID(accountID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			a.logError(
				logKey,
				"error fetching account",
				http.StatusBadRequest,
				err,
			)

			response.NewError(err, http.StatusBadRequest).Send(w)
			return
		default:
			a.logError(
				logKey,
				"error when returning account balance",
				http.StatusInternalServerError,
				err,
			)

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}
	a.logSuccess(logKey, "success when returning account balance", http.StatusOK)

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (a Account) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}

func (a Account) validateInput(input accountInput) []string {
	var messages []string

	err := a.validator.Validate(input)
	if err != nil {
		for _, msg := range a.validator.Messages() {
			messages = append(messages, msg)
		}
	}

	return messages
}

func (a Account) logSuccess(key string, message string, httpStatus int) {
	a.log.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
	}).Infof(message)
}

func (a Account) logError(key string, message string, httpStatus int, err error) {
	a.log.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
		"error":       err.Error(),
	}).Errorf(message)
}
