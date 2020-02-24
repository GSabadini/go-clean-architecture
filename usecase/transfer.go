package usecase

import (
	"errors"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
	"gopkg.in/mgo.v2/bson"
)

//StoreTransfer cria uma nova transação
func StoreTransfer(
	transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository,
	transfer *domain.Transfer,
) error {

	var query = bson.M{"_id": transfer.AccountOriginId}
	accountOrigin, err := accountRepository.FindOne(query)
	if err != nil {
		return err
	}

	//@TODO CRIAR ERRO ESPECIFICO
	if accountOrigin.Ballance < transfer.Amount {
		return errors.New("A conta de origin não tem saldo suficiente")
	}

	err = transferRepository.Store(transfer)
	if err != nil {
		return err
	}

	return nil
}
