package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferOutput armazena a estrutura de dados de retorno do caso de uso
type TransferOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

//Transfer armazena as dependências para os casos de uso de Transfer
type Transfer struct {
	transferRepository domain.TransferRepository
	accountRepository  domain.AccountRepository
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransfer(
	transferRepository domain.TransferRepository,
	accountRepository domain.AccountRepository,
) Transfer {
	return Transfer{
		transferRepository: transferRepository,
		accountRepository:  accountRepository,
	}
}

//Store cria uma nova Transfer
func (t Transfer) Store(accountOriginID, accountDestinationID string, amount float64) (TransferOutput, error) {
	if err := t.process(accountOriginID, accountDestinationID, amount); err != nil {
		return TransferOutput{}, err
	}

	var transfer = domain.NewTransfer(domain.NewUUID(), accountOriginID, accountDestinationID, amount, time.Now())

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

/* TODO melhorar processsamento de transação */
func (t Transfer) process(accountOriginID, accountDestinationID string, amount float64) error {
	origin, err := t.accountRepository.FindByID(accountOriginID)
	if err != nil {
		return err
	}

	destination, err := t.accountRepository.FindByID(accountDestinationID)
	if err != nil {
		return err
	}

	if err := origin.Withdraw(amount); err != nil {
		return err
	}

	destination.Deposit(amount)

	if err = t.accountRepository.UpdateBalance(origin.ID, origin.Balance); err != nil {
		return err
	}

	if err = t.accountRepository.UpdateBalance(destination.ID, destination.Balance); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll() ([]TransferOutput, error) {
	var output = make([]TransferOutput, 0)

	transfers, err := t.transferRepository.FindAll()
	if err != nil {
		return output, err
	}

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
