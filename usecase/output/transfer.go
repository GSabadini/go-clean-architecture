package output

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type TransferPresenter interface {
	Output(domain.Transfer) Transfer
	OutputList([]domain.Transfer) []Transfer
}

type Transfer struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
