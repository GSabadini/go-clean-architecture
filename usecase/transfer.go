package usecase

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

type TransferService struct {
	transferRepository repository.TransferRepository
	accountRepository  repository.AccountRepository
}

func NewTransferService(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
) TransferService {
	return TransferService{
		transferRepository: transferRepository,
		accountRepository:  accountRepository,
	}
}

//StoreTransfer cria uma nova transferência
func (t TransferService) StoreTransfer(transfer domain.Transfer) (domain.Transfer, error) {
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

func (t TransferService) processTransfer(transfer domain.Transfer) error {
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

func (t TransferService) updateAccount(query bson.M, update bson.M) error {
	if err := t.accountRepository.Update(query, update); err != nil {
		return errors.Wrap(err, "error updating account")
	}

	return nil
}

//FindAllTransfer retorna uma lista de transferências
func (t TransferService) FindAllTransfer() ([]domain.Transfer, error) {
	result, err := t.transferRepository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}
