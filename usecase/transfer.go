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
	transfer domain.Transfer,
) (domain.Transfer, error) {
	if err := processTransfer(accountRepository, transfer); err != nil {
		return domain.Transfer{}, err
	}

	result, err := transferRepository.Store(transfer)
	if err != nil {
		return domain.Transfer{}, err
	}

	return result, nil
}

func processTransfer(repository repository.AccountRepository, transfer domain.Transfer) error {
	origin, err := findAccount(repository, bson.M{"_id": transfer.GetAccountOriginID()})
	if err != nil {
		return err
	}

	if err := origin.Withdraw(transfer.GetAmount()); err != nil {
		return err
	}

	destination, err := findAccount(repository, bson.M{"_id": transfer.GetAccountDestinationID()})
	if err != nil {
		return err
	}

	destination.Deposit(transfer.GetAmount())

	if err = updateAccount(
		repository,
		bson.M{"_id": transfer.GetAccountOriginID()},
		bson.M{"$set": bson.M{"balance": origin.GetBalance()}},
	); err != nil {
		return err
	}

	if err = updateAccount(
		repository,
		bson.M{"_id": transfer.GetAccountDestinationID()},
		bson.M{"$set": bson.M{"balance": destination.GetBalance()}},
	); err != nil {
		return err
	}

	return nil
}

func findAccount(repository repository.AccountRepository, query bson.M) (*domain.Account, error) {
	account, err := repository.FindOne(query)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching account")
	}

	return account, nil
}

func updateAccount(repository repository.AccountRepository, query bson.M, update bson.M) error {
	if err := repository.Update(query, update); err != nil {
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
