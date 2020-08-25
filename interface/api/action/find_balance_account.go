package action

import (
	"net/http"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/interface/api/logging"
	"github.com/gsabadini/go-bank-transfer/interface/api/response"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

//FindAllAccountAction armazena as dependências para as ações de Account
type FindBalanceAccountAction struct {
	uc  usecase.FindBalanceAccount
	log logger.Logger
}

//NewFindBalanceAccountAction constrói um FindBalanceAccountAction com suas dependências
func NewFindBalanceAccountAction(uc usecase.FindBalanceAccount, l logger.Logger) FindBalanceAccountAction {
	return FindBalanceAccountAction{uc: uc, log: l}
}

//Execute é um handler para retornar o Balance de uma Account
func (a FindBalanceAccountAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "find_balance"

	var accountID = r.URL.Query().Get("account_id")
	if !domain.IsValidUUID(accountID) {
		var err = response.ErrParameterInvalid
		logging.NewError(
			a.log,
			logKey,
			"parameter invalid",
			http.StatusBadRequest,
			err,
		).Log()

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.Execute(r.Context(), domain.AccountID(accountID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			logging.NewError(
				a.log,
				logKey,
				"error fetching validator",
				http.StatusBadRequest,
				err,
			).Log()

			response.NewError(err, http.StatusBadRequest).Send(w)
			return
		default:
			logging.NewError(
				a.log,
				logKey,
				"error when returning validator balance",
				http.StatusInternalServerError,
				err,
			).Log()

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}
	logging.NewInfo(a.log, logKey, "success when returning validator balance", http.StatusOK).Log()

	response.NewSuccess(output, http.StatusOK).Send(w)
}
