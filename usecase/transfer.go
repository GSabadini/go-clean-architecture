package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

type TransferInput struct {
	AccountOriginID      string  `json:"account_origin_id" validate:"required"`
	AccountDestinationID string  `json:"account_destination_id" validate:"required"`
	Amount               float64 `json:"amount" validate:"gt=0,required"`
}

type TransferOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

//Transfer armazena as depedências para ações de uma transferência
type Transfer struct {
	transferRepository repository.TransferRepository
	accountRepository  repository.AccountRepository
}

//NewTransfer cria uma transferência com suas dependências
func NewTransfer(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
) Transfer {
	return Transfer{
		transferRepository: transferRepository,
		accountRepository:  accountRepository,
	}
}

//Store cria uma nova transferência
func (t Transfer) Store(input TransferInput) (TransferOutput, error) {
	if err := t.processTransfer(input); err != nil {
		return TransferOutput{}, err
	}

	var transfer = domain.NewTransfer(input.AccountOriginID, input.AccountDestinationID, input.Amount)

	transfer, err := t.transferRepository.Store(transfer)
	if err != nil {
		return TransferOutput{}, err
	}

	return TransferOutput{
		ID:                   transfer.ID,
		AccountOriginID:      transfer.AccountOriginID,
		AccountDestinationID: transfer.AccountDestinationID,
		Amount:               transfer.Amount,
		CreatedAt:            transfer.CreatedAt,
	}, nil
}

func (t Transfer) processTransfer(transfer TransferInput) error {
	origin, err := t.accountRepository.FindByID(transfer.AccountOriginID)
	if err != nil {
		return err
	}

	destination, err := t.accountRepository.FindByID(transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	if err := origin.Withdraw(transfer.Amount); err != nil {
		return err
	}

	destination.Deposit(transfer.Amount)

	if err = t.accountRepository.UpdateBalance(origin.ID, origin.GetBalance()); err != nil {
		return err
	}

	if err = t.accountRepository.UpdateBalance(destination.ID, destination.GetBalance()); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll() ([]TransferOutput, error) {
	transfers, err := t.transferRepository.FindAll()
	if err != nil {
		return []TransferOutput{}, err
	}

	var output []TransferOutput
	for _, transfer := range transfers {
		var transfer = TransferOutput{
			ID:                   transfer.ID,
			AccountOriginID:      transfer.AccountOriginID,
			AccountDestinationID: transfer.AccountDestinationID,
			Amount:               transfer.Amount,
			CreatedAt:            transfer.CreatedAt,
		}

		output = append(output, transfer)
	}

	return output, nil
}
