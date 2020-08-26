package action

import (
	"github.com/gsabadini/go-bank-transfer/interface/logger"
	"net/http"

	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

type FindAllTransferAction struct {
	uc  usecase.FindAllTransfer
	log logger.Logger
}

func NewFindAllTransferAction(uc usecase.FindAllTransfer, log logger.Logger) FindAllTransferAction {
	return FindAllTransferAction{uc: uc, log: log}
}

//Execute é um handler para retornar todas as Transfer
func (t FindAllTransferAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_transfer"

	output, err := t.uc.Execute(r.Context())
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
