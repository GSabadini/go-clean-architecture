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
	accountOrigin, err := accountRepository.FindOne(bson.M{"_id": transfer.AccountOriginId})
	if err != nil {
		return errors.Wrap(err, "erro ao buscar conta de origem")
	}

	accountDestination, err := accountRepository.FindOne(bson.M{"_id": transfer.AccountDestinationId})
	if err != nil {
		return errors.Wrap(err, "erro ao buscar conta de destino")
	}

	if err = transferAccountBalance(accountOrigin, accountDestination, transfer.Amount); err != nil {
		return err
	}

	if err = accountRepository.Update(
		bson.M{"_id": transfer.AccountOriginId},
		bson.M{"$set": bson.M{"balance": accountOrigin.Balance}},
	); err != nil {
		return errors.Wrap(err, "erro ao atualizar conta de origem")
	}

	if err = accountRepository.Update(
		bson.M{"_id": transfer.AccountDestinationId},
		bson.M{"$set": bson.M{"balance": accountDestination.Balance}},
	); err != nil {
		return errors.Wrap(err, "erro ao atualizar conta de destino")
	}

	if err = transferRepository.Store(transfer); err != nil {
		return err
	}

	return nil
}

func transferAccountBalance(accountOrigin *domain.Account, accountDestination *domain.Account, amount float64) error {
	if accountOrigin.Balance < amount {
		return errors.New("A conta de origem não tem saldo suficiente")
	}

	accountOrigin.SubtractBalance(amount)

	accountDestination.SumBalance(amount)

	return nil
}
