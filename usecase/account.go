package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//Store cria uma nova account
func Store(repository repository.AccountRepository, account domain.Account) error {
	if err := repository.Store(account); err != nil {
		return err
	}

	return nil
}

//FindAll recupera uma lista de accounts
func FindAll(repository repository.AccountRepository) ([]domain.Account, error) {
	result, err := repository.FindAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

//FindOne recupera uma account com base em um ID
func FindOne(repository repository.AccountRepository, id string) (domain.Account, error) {
	var query = bson.M{"_id": bson.ObjectIdHex(id)}

	result, err := repository.FindOne(query)
	if err != nil {
		return result, err
	}

	return result, nil
}

//func RecoverUserTrackingData(hash assign.SenderHash, userTrackingRepository domain.UserTrackingRepository) (domain.UserTracking, error) {
//	var (
//		hashQuery        = bson.M{"senderHash": hash}
//		userTrackingData domain.UserTracking
//	)
//
//	err := userTrackingRepository.RecoverOne(hashQuery, &userTrackingData)
//	if err != nil && err.Error() == "not found" {
//		return userTrackingData, NewNotFoundError(err.Error())
//	}
//
//	return userTrackingData, err
//}
