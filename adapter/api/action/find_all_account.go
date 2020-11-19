package action

import (
	"net/http"

	"github.com/gsabadini/go-bank-transfer/adapter/api/logging"
	"github.com/gsabadini/go-bank-transfer/adapter/api/response"
	"github.com/gsabadini/go-bank-transfer/adapter/logger"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type FindAllAccountAction struct {
	uc  usecase.FindAllAccountUseCase
	log logger.Logger
}

func NewFindAllAccountAction(uc usecase.FindAllAccountUseCase, log logger.Logger) FindAllAccountAction {
	return FindAllAccountAction{
		uc:  uc,
		log: log,
	}
}

func (a FindAllAccountAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_account"

	output, err := a.uc.Execute(r.Context())
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when returning account list")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning account list")

	response.NewSuccess(output, http.StatusOK).Send(w)
}
