package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/api/action"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//ValidateTransfer armazena a estrutura de validação de entrada de dados
type ValidateTransfer struct {
	logger *logrus.Logger
}

//NewValidateTransfer constrói um ValidateTransfer com suas dependências
func NewValidateTransfer(log *logrus.Logger) ValidateTransfer {
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
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error("error read body")

		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	var transfer transferRequest
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error("error when decoding json")

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := transfer.validateAmount(); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error(messageInvalidField)

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	if err := transfer.validateOriginEqualsDestination(); err != nil {
		v.logger.WithFields(logrus.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Error(messageInvalidField)

		action.ErrorMessage(err, http.StatusBadRequest).Send(w)
		return
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	next.ServeHTTP(w, r)
}

var (
	errAmountInvalid                  = errors.New("Amount invalid")
	errAccountOriginEqualsDestination = errors.New("Account origin equals destination account")
)

type transferRequest struct {
	AccountOriginID      bson.ObjectId `json:"account_origin_id"`
	AccountDestinationID bson.ObjectId `json:"account_destination_id"`
	Amount               float64       `json:"amount"`
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
