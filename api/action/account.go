package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/api/input"
	"github.com/gsabadini/go-bank-transfer/api/logging"
	"github.com/gsabadini/go-bank-transfer/api/response"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

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

	var inputAccount input.Account
	if err := json.NewDecoder(r.Body).Decode(&inputAccount); err != nil {
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

	if errs := inputAccount.Validate(a.validator); len(errs) > 0 {
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

	output, err := a.uc.Store(
		r.Context(),
		inputAccount.Name,
		a.cleanCPF(inputAccount.CPF),
		domain.Money(inputAccount.Balance),
	)
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

//FindAll é um handler para retornar todas as Account
func (a Account) FindAll(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_account"

	output, err := a.uc.FindAll(r.Context())
	if err != nil {
		logging.NewError(
			a.log,
			logKey,
			"error when returning account list",
			http.StatusInternalServerError,
			err,
		).Log()

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, "success when returning account list", http.StatusOK).Log()

	response.NewSuccess(output, http.StatusOK).Send(w)
}

//FindBalance é um handler para retornar o Balance de uma Account
func (a Account) FindBalance(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_balance"

	var accountID = r.URL.Query().Get("account_id")
	if !domain.IsValidUUID(accountID) {
		var err = response.ErrParameterInvalid
		logging.NewError(
			a.log,
			logKey,
			"parameter invalid",
			http.StatusBadRequest,
			err,
		).Log()

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.FindBalance(r.Context(), domain.AccountID(accountID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			logging.NewError(
				a.log,
				logKey,
				"error fetching account",
				http.StatusBadRequest,
				err,
			).Log()

			response.NewError(err, http.StatusBadRequest).Send(w)
			return
		default:
			logging.NewError(
				a.log,
				logKey,
				"error when returning account balance",
				http.StatusInternalServerError,
				err,
			).Log()

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}
	logging.NewInfo(a.log, logKey, "success when returning account balance", http.StatusOK).Log()

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (a Account) cleanCPF(cpf string) string {
	return strings.Replace(strings.Replace(cpf, ".", "", -1), "-", "", -1)
}
