package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	usecase "github.com/gsabadini/go-bank-transfer/usecase/transfer"
)

type transferPresenter struct{}

//NewTransferPresenter
func NewTransferPresenter() transferPresenter {
	return transferPresenter{}
}

//Output
func (tp transferPresenter) Output(transfer domain.Transfer) usecase.TransferOutput {
	return usecase.TransferOutput{
		ID:                   transfer.ID().String(),
		AccountOriginID:      transfer.AccountOriginID().String(),
		AccountDestinationID: transfer.AccountDestinationID().String(),
		Amount:               transfer.Amount().Float64(),
		CreatedAt:            transfer.CreatedAt(),
	}
}

//OutputList
func (tp transferPresenter) OutputList(transfers []domain.Transfer) []usecase.TransferOutput {
	var output = make([]usecase.TransferOutput, 0)

	for _, transfer := range transfers {
		output = append(output, usecase.TransferOutput{
			ID:                   transfer.ID().String(),
			AccountOriginID:      transfer.AccountOriginID().String(),
			AccountDestinationID: transfer.AccountDestinationID().String(),
			Amount:               transfer.Amount().Float64(),
			CreatedAt:            transfer.CreatedAt(),
		})
	}

	return output
}