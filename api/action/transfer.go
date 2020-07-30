package action

import (
	"encoding/json"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/api/response"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/pkg/errors"
)

//transferInput armazena a estrutura de dados de entrada da API
type transferInput struct {
	AccountOriginID      string  `json:"account_origin_id" validate:"required,uuid4"`
	AccountDestinationID string  `json:"account_destination_id" validate:"required,uuid4"`
	Amount               float64 `json:"amount" validate:"gt=0,required"`
}

//Transfer armazena as dependências para as ações de Transfer
type Transfer struct {
	validator validator.Validator
	log       logger.Logger
	uc        usecase.TransferUseCase
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransfer(uc usecase.TransferUseCase, l logger.Logger, v validator.Validator) Transfer {
	return Transfer{uc: uc, log: l, validator: v}
}

//Store é um handler para criação de Transfer
func (t Transfer) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_transfer"

	var input transferInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		t.logError(
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := t.validateInput(input); len(errs) > 0 {
		t.logError(
			logKey,
			"invalid input",
			http.StatusBadRequest,
			errors.New("invalid input"),
		)

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := t.uc.Store(
		r.Context(),
		domain.AccountID(input.AccountOriginID),
		domain.AccountID(input.AccountDestinationID),
		input.Amount,
	)
	if err != nil {
		switch err {
		case domain.ErrInsufficientBalance:
			t.logError(
				logKey,
				"insufficient balance",
				http.StatusUnprocessableEntity,
				err,
			)

			response.NewError(err, http.StatusUnprocessableEntity).Send(w)
			return
		default:
			t.logError(
				logKey,
				"error when creating a new transfer",
				http.StatusInternalServerError,
				err,
			)

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}
	t.logSuccess(logKey, "success create transfer", http.StatusCreated)

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

//Index é um handler para retornar todas as Transfer
func (t Transfer) Index(w http.ResponseWriter, r *http.Request) {
	const logKey = "index_transfer"

	output, err := t.uc.FindAll(r.Context())
	if err != nil {
		t.logError(
			logKey,
			"error when returning the transfer list",
			http.StatusInternalServerError,
			err)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	t.logSuccess(logKey, "success when returning transfer list", http.StatusOK)

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (t Transfer) validateInput(input transferInput) []string {
	var (
		messages          []string
		errAccountsEquals = errors.New("account origin equals destination account")
		accountIsEquals   = input.AccountOriginID == input.AccountDestinationID
		accountsIsEmpty   = input.AccountOriginID == "" && input.AccountDestinationID == ""
	)

	if !accountsIsEmpty && accountIsEquals {
		messages = append(messages, errAccountsEquals.Error())
	}

	err := t.validator.Validate(input)
	if err != nil {
		for _, msg := range t.validator.Messages() {
			messages = append(messages, msg)
		}
	}

	return messages
}

func (t Transfer) logSuccess(key string, message string, httpStatus int) {
	t.log.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
	}).Infof(message)
}

func (t Transfer) logError(key string, message string, httpStatus int, err error) {
	t.log.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
		"error":       err.Error(),
	}).Errorf(message)
}
