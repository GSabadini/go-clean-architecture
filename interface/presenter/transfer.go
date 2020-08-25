package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type transferPresenter struct{}

//NewTransferPresenter
func NewTransferPresenter() transferPresenter {
	return transferPresenter{}
}

//Output
func (tp transferPresenter) Output(transfer domain.Transfer) output.TransferOutput {
	return output.TransferOutput{
		ID:                   transfer.ID().String(),
		AccountOriginID:      transfer.AccountOriginID().String(),
		AccountDestinationID: transfer.AccountDestinationID().String(),
		Amount:               transfer.Amount().Float64(),
		CreatedAt:            transfer.CreatedAt(),
	}
}

//OutputList
func (tp transferPresenter) OutputList(transfers []domain.Transfer) []output.TransferOutput {
	var o = make([]output.TransferOutput, 0)

	for _, transfer := range transfers {
		o = append(o, output.TransferOutput{
			ID:                   transfer.ID().String(),
			AccountOriginID:      transfer.AccountOriginID().String(),
			AccountDestinationID: transfer.AccountDestinationID().String(),
			Amount:               transfer.Amount().Float64(),
			CreatedAt:            transfer.CreatedAt(),
		})
	}

	return o
}
