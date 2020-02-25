package action

import (
	"encoding/json"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/repository"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/sirupsen/logrus"
)

//Transfer armazena as dependências de uma transferência
type Transfer struct {
	dbHandler database.NoSQLDBHandler
	logger    *logrus.Logger
}

//NewTransfer constrói uma transferência com suas dependências
func NewTransfer(dbHandler database.NoSQLDBHandler, log *logrus.Logger) Transfer {
	return Transfer{dbHandler: dbHandler, logger: log}
}

//Store é um handler para criação de transferência
func (t Transfer) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_transfer"
	var transfer *domain.Transfer

	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		t.logError(
			logKey,
			"error when decoding json",
			http.StatusInternalServerError,
			err,
		)

		ErrorMessage(err, http.StatusInternalServerError).Send(w)
		return
	}

	var transferRepository = repository.NewTransfer(t.dbHandler)
	var accountRepository = repository.NewAccount(t.dbHandler)

	err := usecase.StoreTransfer(transferRepository, accountRepository, transfer)
	if err != nil {
		t.logError(
			logKey,
			"error when creating a new transfer",
			http.StatusInternalServerError,
			err,
		)

		ErrorMessage(err, http.StatusInternalServerError).Send(w)
		return
	}

	t.logInfoSuccess(logKey, "success create transfer", http.StatusCreated)

	w.WriteHeader(http.StatusNoContent)
}

func (t Transfer) logInfoSuccess(key string, description string, httpStatus int) {
	t.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
		"description": description,
	}).Info()
}

func (t Transfer) logError(key string, description string, httpStatus int, err error) {
	t.logger.WithFields(logrus.Fields{
		"key":         key,
		"http_status": httpStatus,
		"description": description,
		"error":       err.Error(),
	}).Error()
}
