package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type transferPresenter struct{}

func NewTransferPresenter() transferPresenter {
	return transferPresenter{}
}

func (tp transferPresenter) Output(transfer domain.Transfer) output.Transfer {
	return output.Transfer{
		ID:                   transfer.ID().String(),
		AccountOriginID:      transfer.AccountOriginID().String(),
		AccountDestinationID: transfer.AccountDestinationID().String(),
		Amount:               transfer.Amount().Float64(),
		CreatedAt:            transfer.CreatedAt(),
	}
}

func (tp transferPresenter) OutputList(transfers []domain.Transfer) []output.Transfer {
	var o = make([]output.Transfer, 0)

	for _, transfer := range transfers {
		o = append(o, output.Transfer{
			ID:                   transfer.ID().String(),
			AccountOriginID:      transfer.AccountOriginID().String(),
			AccountDestinationID: transfer.AccountDestinationID().String(),
			Amount:               transfer.Amount().Float64(),
			CreatedAt:            transfer.CreatedAt(),
		})
	}

	return o
}
