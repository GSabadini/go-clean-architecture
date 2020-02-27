package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//StoreTransfer cria uma nova transferência
func StoreTransfer(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
	transfer *domain.Transfer,
) (*domain.Transfer, error) {
	if err := transferAccountBalance(accountRepository, transfer); err != nil {
		return nil, err
	}

	result, err := transferRepository.Store(transfer)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func transferAccountBalance(accountRepository repository.AccountRepository, transfer *domain.Transfer) error {
	accountOrigin, err := findAccount(accountRepository, bson.M{"_id": transfer.GetAccountOriginID()})
	if err != nil {
		return err
	}

	if err := accountOrigin.Withdraw(transfer.GetAmount()); err != nil {
		return err
	}

	accountDestination, err := findAccount(accountRepository, bson.M{"_id": transfer.GetAccountDestinationID()})
	if err != nil {
		return err
	}

	accountDestination.Deposit(transfer.GetAmount())

	if err = updateAccountBalance(
		accountRepository,
		bson.M{"_id": transfer.GetAccountOriginID()},
		bson.M{"$set": bson.M{"balance": accountOrigin.GetBalance()}},
	); err != nil {
		return err
	}

	if err = updateAccountBalance(
		accountRepository,
		bson.M{"_id": transfer.GetAccountDestinationID()},
		bson.M{"$set": bson.M{"balance": accountDestination.GetBalance()}},
	); err != nil {
		return err
	}

	return nil
}

func findAccount(accountRepository repository.AccountRepository, query bson.M) (*domain.Account, error) {
	account, err := accountRepository.FindOne(query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching account")
	}

	return account, nil
}

func updateAccountBalance(accountRepository repository.AccountRepository, query bson.M, update bson.M) error {
	if err := accountRepository.Update(query, update); err != nil {
		return errors.Wrap(err, "error updating account")
	}

	return nil
}

//FindAllTransfer retorna uma lista de transferências
func FindAllTransfer(repository repository.TransferRepository) ([]domain.Transfer, error) {
	result, err := repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}
