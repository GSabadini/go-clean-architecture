package action

import (
	"net/http"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

//FindAllTransferAction armazena as dependências para as ações de Transfer
type FindAllTransferAction struct {
	uc  usecase.FindAllTransfer
	log logger.Logger
}

//NewFindAllTransferAction constrói um Transfer com suas dependências
func NewFindAllTransferAction(uc usecase.FindAllTransfer, l logger.Logger) FindAllTransferAction {
	return FindAllTransferAction{uc: uc, log: l}
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
