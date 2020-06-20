package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gsabadini/go-bank-transfer/api/response"
	"io/ioutil"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"

	"github.com/pkg/errors"
)

//ValidateTransfer armazena a estrutura de validação de entrada de dados
type ValidateTransfer struct {
	logger logger.Logger
}

//NewValidateTransfer constrói um ValidateTransfer com suas dependências
func NewValidateTransfer(log logger.Logger) ValidateTransfer {
	return ValidateTransfer{logger: log}
}

//Execute válida os dados de criação de transferência
func (v ValidateTransfer) Execute(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	const (
		logKey              = "validate_transfer_middleware"
		messageInvalidField = "Invalid field"
	)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		v.logger.WithFields(logger.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Errorf("error read body")

		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	var transfer transferRequest
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		v.logger.WithFields(logger.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Errorf("error when decoding json")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := transfer.validateAmount(); err != nil {
		v.logger.WithFields(logger.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Errorf(messageInvalidField)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := transfer.validateOriginEqualsDestination(); err != nil {
		v.logger.WithFields(logger.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Errorf(messageInvalidField)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	next.ServeHTTP(w, r)
}

var (
	errAmountInvalid                  = errors.New("amount invalid")
	errAccountOriginEqualsDestination = errors.New("account origin equals destination account")
)

type transferRequest struct {
	AccountOriginID      string  `json:"account_origin_id"`
	AccountDestinationID string  `json:"account_destination_id"`
	Amount               float64 `json:"amount"`
}

func (t *transferRequest) validateAmount() error {
	if t.Amount < 0 {
		return errAmountInvalid
	}

	return nil
}

func (t *transferRequest) validateOriginEqualsDestination() error {
	if t.AccountOriginID == t.AccountDestinationID {
		return errAccountOriginEqualsDestination
	}

	return nil
}
