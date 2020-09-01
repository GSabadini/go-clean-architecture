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

	logKey, logMsg string
}

func NewCreateTransferAction(uc usecase.CreateTransfer, log logger.Logger, v validator.Validator) CreateTransferAction {
	return CreateTransferAction{
		uc:        uc,
		log:       log,
		validator: v,
		logKey:    "create_transfer",
		logMsg:    "creating a new transfer",
	}
}

func (t CreateTransferAction) Execute(w http.ResponseWriter, r *http.Request) {
	var transferInput input.Transfer
	if err := json.NewDecoder(r.Body).Decode(&transferInput); err != nil {
		logging.NewError(
			t.log,
			err,
			t.logKey,
			http.StatusBadRequest,
		).Log(t.logMsg)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := t.validateInput(transferInput); len(errs) > 0 {
		logging.NewError(
			t.log,
			response.ErrInvalidInput,
			t.logKey,
			http.StatusBadRequest,
		).Log(t.logMsg)

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := t.uc.Execute(r.Context(), transferInput)
	if err != nil {
		t.handleErr(w, err)
		return
	}

	logging.NewInfo(t.log, t.logKey, http.StatusCreated).Log(t.logMsg)

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (t CreateTransferAction) handleErr(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInsufficientBalance:
		logging.NewError(
			t.log,
			err,
			t.logKey,
			http.StatusUnprocessableEntity,
		).Log(t.logMsg)

		response.NewError(err, http.StatusUnprocessableEntity).Send(w)
		return
	case domain.ErrAccountOriginNotFound:
		logging.NewError(
			t.log,
			err,
			t.logKey,
			http.StatusUnprocessableEntity,
		).Log(t.logMsg)

		response.NewError(err, http.StatusUnprocessableEntity).Send(w)
		return
	case domain.ErrAccountDestinationNotFound:
		logging.NewError(
			t.log,
			err,
			t.logKey,
			http.StatusUnprocessableEntity,
		).Log(t.logMsg)

		response.NewError(err, http.StatusUnprocessableEntity).Send(w)
		return
	default:
		logging.NewError(
			t.log,
			err,
			t.logKey,
			http.StatusInternalServerError,
		).Log(t.logMsg)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
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
