package usecase

import (
	"errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//StoreTransfer cria uma nova transação
func StoreTransfer(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
	transfer *domain.Transfer,
) error {
	var queryAccountOriginId = bson.M{"_id": transfer.AccountOriginId}
	accountOrigin, err := accountRepository.FindOne(queryAccountOriginId)
	if err != nil {
		return err
	}

	//@TODO CRIAR ERRO ESPECIFICO
	if accountOrigin.Ballance < transfer.Amount {
		return errors.New("A conta de origin não tem saldo suficiente")
	}

	accountOrigin.Ballance = accountOrigin.Ballance - transfer.Amount
	var updateBallanceAccountOriginId = bson.M{"$set": bson.M{"ballance": accountOrigin.Ballance}}

	err = accountRepository.Update(queryAccountOriginId, updateBallanceAccountOriginId)
	if err != nil {
		return err
	}

	//-------------------
	var queryAccountDestinationId = bson.M{"_id": transfer.AccountDestinationId}
	accountDestinationId, err := accountRepository.FindOne(queryAccountDestinationId)
	if err != nil {
		return err
	}

	accountDestinationId.Ballance = accountDestinationId.Ballance + transfer.Amount
	var updateBallanceaccountDestinationId = bson.M{"$set": bson.M{"ballance": accountDestinationId.Ballance}}
	err = accountRepository.Update(queryAccountDestinationId, updateBallanceaccountDestinationId)
	if err != nil {
		return err
	}

	//------------------------

	err = transferRepository.Store(transfer)
	if err != nil {
		return err
	}

	return nil
}

func TransferAccount(accountOrigin domain.Account, accountDestination domain.Account) {

}
