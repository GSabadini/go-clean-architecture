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

//TransferInput armazena a estruturas de dados de entrada da API
type TransferInput struct {
	AccountOriginID      string  `json:"account_origin_id" validate:"required"`
	AccountDestinationID string  `json:"account_destination_id" validate:"required"`
	Amount               float64 `json:"amount" validate:"gt=0,required"`
}

//Transfer armazena as dependências de uma transferência
type Transfer struct {
	validator validator.Validator
	log       logger.Logger
	usecase   usecase.TransferUseCase
}

//NewTransfer constrói uma transferência com suas dependências
func NewTransfer(usecase usecase.TransferUseCase, log logger.Logger, v validator.Validator) Transfer {
	return Transfer{usecase: usecase, log: log, validator: v}
}

//Store é um handler para criação de transferência
func (t Transfer) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_transfer"

	var input TransferInput
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

	if errMsg := t.validateInput(input); len(errMsg) > 0 {
		t.logError(
			logKey,
			"input invalid",
			http.StatusBadRequest,
			errors.New("validate"),
		)

		response.NewMessagesError(errMsg, http.StatusBadRequest).Send(w)
		return
	}

	result, err := t.usecase.Store(
		input.AccountOriginID,
		input.AccountDestinationID,
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

	response.NewSuccess(result, http.StatusCreated).Send(w)
}

//Index é um handler para retornar a lista de transferências
func (t Transfer) Index(w http.ResponseWriter, _ *http.Request) {
	const logKey = "index_transfer"

	result, err := t.usecase.FindAll()
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

	response.NewSuccess(result, http.StatusOK).Send(w)
}

func (t Transfer) validateInput(input TransferInput) []string {
	var (
		messages          []string
		errAccountsEquals = errors.New("account origin equals destination account")
	)

	if input.AccountOriginID == input.AccountDestinationID {
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
