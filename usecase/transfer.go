package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

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
func (t Transfer) Store(data domain.Transfer) (domain.Transfer, error) {
	if err := t.processTransfer(data); err != nil {
		return domain.Transfer{}, err
	}

	var transfer = domain.NewTransfer(data.AccountOriginID, data.AccountDestinationID, data.Amount)

	result, err := t.transferRepository.Store(transfer)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (t Transfer) processTransfer(transfer domain.Transfer) error {
	origin, err := t.accountRepository.FindOne(transfer.GetAccountOriginID())
	if err != nil {
		return err
	}

	if err := origin.Withdraw(transfer.GetAmount()); err != nil {
		return err
	}

	destination, err := t.accountRepository.FindOne(transfer.GetAccountDestinationID())
	if err != nil {
		return err
	}

	destination.Deposit(transfer.GetAmount())

	if err = t.accountRepository.UpdateBalance(transfer.GetAccountOriginID(), origin.GetBalance()); err != nil {
		return err
	}

	if err = t.accountRepository.UpdateBalance(transfer.GetAccountDestinationID(), destination.GetBalance()); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll() ([]domain.Transfer, error) {
	result, err := t.transferRepository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}
