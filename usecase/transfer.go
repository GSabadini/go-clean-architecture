package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//StoreTransfer cria uma nova transação
func StoreTransfer(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
	transfer *domain.Transfer,
) error {
	if err := transferAccountBalance(accountRepository, transfer); err != nil {
		return err
	}

	if err := transferRepository.Store(transfer); err != nil {
		return err
	}

	return nil
}

func transferAccountBalance(accountRepository repository.AccountRepository, transfer *domain.Transfer) error {
	accountOrigin, err := findAccount(accountRepository, bson.M{"_id": transfer.GetAccountOrigin()})
	if err != nil {
		return err
	}

	accountDestination, err := findAccount(accountRepository, bson.M{"_id": transfer.GetAccountOrigin()})
	if err != nil {
		return err
	}

	if err := accountOrigin.Withdraw(transfer.GetAmount()); err != nil {
		return err
	}

	accountDestination.Deposit(transfer.GetAmount())

	if err = updateAccount(
		accountRepository,
		bson.M{"_id": transfer.GetAccountOrigin()},
		bson.M{"$set": bson.M{"balance": accountOrigin.GetBalance()}},
	); err != nil {
		return err
	}

	if err = updateAccount(
		accountRepository,
		bson.M{"_id": transfer.GetAccountDestination()},
		bson.M{"$set": bson.M{"balance": accountDestination.GetBalance()}},
	); err != nil {
		return err
	}

	return nil
}

func findAccount(accountRepository repository.AccountRepository, query bson.M) (*domain.Account, error) {
	account, err := accountRepository.FindOne(query)
	if err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	return account, nil
}

func updateAccount(accountRepository repository.AccountRepository, query bson.M, update bson.M) error {
	if err := accountRepository.Update(query, update); err != nil {
		return errors.Wrap(err, "error updating account")
	}

	return nil
}
