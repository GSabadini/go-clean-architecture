package action

import (
	"encoding/json"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/interface/api/input"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	usecase "github.com/gsabadini/go-bank-transfer/usecase/transfer"


	"github.com/pkg/errors"
)

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

	var inputTransfer input.Transfer
	if err := json.NewDecoder(r.Body).Decode(&inputTransfer); err != nil {
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

	if errs := inputTransfer.Validate(t.validator); len(errs) > 0 {
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

	output, err := t.uc.Store(
		r.Context(),
		domain.AccountID(inputTransfer.AccountOriginID),
		domain.AccountID(inputTransfer.AccountDestinationID),
		domain.Money(inputTransfer.Amount),
	)
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

//FindAll é um handler para retornar todas as Transfer
func (t Transfer) FindAll(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_transfer"

	output, err := t.uc.FindAll(r.Context())
	if err != nil {
		logging.NewError(
			t.log,
			logKey,
			"error when returning the transfer list",
			http.StatusInternalServerError,
			err,
		).Log()

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(t.log, logKey, "success when returning transfer list", http.StatusOK).Log()

	response.NewSuccess(output, http.StatusOK).Send(w)
}
