package action

import (
	"encoding/json"
	"errors"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
)

//CreateTransferAction armazena as dependências para as ações de Transfer
type CreateTransferAction struct {
	log       logger.Logger
	uc        usecase.CreateTransfer
	validator validator.Validator
}

//NewCreateTransferAction constrói um Transfer com suas dependências
func NewCreateTransferAction(uc usecase.CreateTransfer, l logger.Logger, v validator.Validator) CreateTransferAction {
	return CreateTransferAction{uc: uc, log: l, validator: v}
}

//Execute é um handler para criação de Transfer
func (t CreateTransferAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_transfer"

	var transferInput input.Transfer
	if err := json.NewDecoder(r.Body).Decode(&transferInput); err != nil {
		logging.NewError(
			t.log,
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		).Log()

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := t.validateInput(transferInput); len(errs) > 0 {
		logging.NewError(
			t.log,
			logKey,
			"invalid validator",
			http.StatusBadRequest,
			errors.New("invalid validator"),
		).Log()

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := t.uc.Execute(r.Context(), transferInput)
	if err != nil {
		switch err {
		case domain.ErrInsufficientBalance:
			logging.NewError(
				t.log,
				logKey,
				"insufficient balance",
				http.StatusUnprocessableEntity,
				err,
			).Log()

			response.NewError(err, http.StatusUnprocessableEntity).Send(w)
			return
		default:
			logging.NewError(
				t.log,
				logKey,
				"error when creating a new transfer",
				http.StatusInternalServerError,
				err,
			).Log()

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}

	logging.NewInfo(t.log, logKey, "success create transfer", http.StatusCreated).Log()

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (t CreateTransferAction) validateInput(input input.Transfer) []string {
	var (
		msgs              []string
		errAccountsEquals = errors.New("account origin equals destination account")
		accountIsEquals   = input.AccountOriginID == input.AccountDestinationID
		accountsIsEmpty   = input.AccountOriginID == "" && input.AccountDestinationID == ""
	)

	if !accountsIsEmpty && accountIsEquals {
		msgs = append(msgs, errAccountsEquals.Error())
	}

	err := t.validator.Validate(input)
	if err != nil {
		for _, msg := range t.validator.Messages() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}
