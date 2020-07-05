package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//transferOutput armazena a estrutura de dados de retorno da API
type transferOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

//Transfer armazena as dependências para ações de uma transferência
type Transfer struct {
	transferRepository domain.TransferRepository
	accountRepository  domain.AccountRepository
}

//NewTransfer cria uma transferência com suas dependências
func NewTransfer(
	transferRepository domain.TransferRepository,
	accountRepository domain.AccountRepository,
) Transfer {
	return Transfer{
		transferRepository: transferRepository,
		accountRepository:  accountRepository,
	}
}

//Store cria uma nova transferência
func (t Transfer) Store(accountOriginID, accountDestinationID string, amount float64) (transferOutput, error) {
	if err := t.process(accountOriginID, accountDestinationID, amount); err != nil {
		return transferOutput{}, err
	}

	var transfer = domain.NewTransfer(accountOriginID, accountDestinationID, amount)

	transfer, err := t.transferRepository.Store(transfer)
	if err != nil {
		return transferOutput{}, err
	}

	return transferOutput{
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
func (t Transfer) FindAll() ([]transferOutput, error) {
	var output = make([]transferOutput, 0)

	transfers, err := t.transferRepository.FindAll()
	if err != nil {
		return output, err
	}

	for _, transfer := range transfers {
		var transfer = transferOutput{
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
