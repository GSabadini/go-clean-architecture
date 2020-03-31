package usecase

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
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
func (t Transfer) Store(transfer domain.Transfer) (domain.Transfer, error) {
	if err := t.processTransfer(transfer); err != nil {
		return domain.Transfer{}, err
	}

	transfer.CreatedAt = time.Now()
	transfer.ID = bson.NewObjectId()

	result, err := t.transferRepository.Store(transfer)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (t Transfer) processTransfer(transfer domain.Transfer) error {
	origin, err := t.accountRepository.FindOne(bson.M{"_id": transfer.GetAccountOriginID()})
	if err != nil {
		return err
	}

	if err := origin.Withdraw(transfer.GetAmount()); err != nil {
		return err
	}

	destination, err := t.accountRepository.FindOne(bson.M{"_id": transfer.GetAccountDestinationID()})
	if err != nil {
		return err
	}

	destination.Deposit(transfer.GetAmount())

	if err = t.updateAccount(
		bson.M{"_id": transfer.GetAccountOriginID()},
		bson.M{"$set": bson.M{"balance": origin.GetBalance()}},
	); err != nil {
		return err
	}

	if err = t.updateAccount(
		bson.M{"_id": transfer.GetAccountDestinationID()},
		bson.M{"$set": bson.M{"balance": destination.GetBalance()}},
	); err != nil {
		return err
	}

	return nil
}

func (t Transfer) updateAccount(query bson.M, update bson.M) error {
	if err := t.accountRepository.Update(query, update); err != nil {
		return errors.Wrap(err, "error updating account")
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
