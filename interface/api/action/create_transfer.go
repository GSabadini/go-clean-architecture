package action

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/interface/logger"
	"github.com/gsabadini/go-bank-transfer/interface/validator"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
)

type CreateTransferAction struct {
	log       logger.Logger
	uc        usecase.CreateTransfer
	validator validator.Validator
}

func NewCreateTransferAction(uc usecase.CreateTransfer, log logger.Logger, v validator.Validator) CreateTransferAction {
	return CreateTransferAction{uc: uc, log: log, validator: v}
}

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
			"invalid input",
			http.StatusBadRequest,
			errors.New("invalid input"),
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
