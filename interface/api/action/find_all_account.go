package action

import (
	"net/http"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

//FindAllAccountAction armazena as dependências para as ações de Account
type FindAllAccountAction struct {
	uc  usecase.FindAllAccount
	log logger.Logger
}

//NewFindAllAccountAction constrói um FindAllAccountAction com suas dependências
func NewFindAllAccountAction(uc usecase.FindAllAccount, l logger.Logger) FindAllAccountAction {
	return FindAllAccountAction{uc: uc, log: l}
}

//Execute é um handler para retornar todas as Account
func (a FindAllAccountAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_all_account"

	output, err := a.uc.Execute(r.Context())
	if err != nil {
		logging.NewError(
			a.log,
			logKey,
			"error when returning validator list",
			http.StatusInternalServerError,
			err,
		).Log()

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	logging.NewInfo(a.log, logKey, "success when returning validator list", http.StatusOK).Log()

	response.NewSuccess(output, http.StatusOK).Send(w)
}
