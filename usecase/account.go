package usecase

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

func Create(repository repository.AccountRepository, account domain.Account) error {
	if err := repository.Store(account); err != nil {
		return err
	}

	return nil
}

func FindAll(repository repository.AccountRepository, account []domain.Account) ([]domain.Account, error) {
	result, err := repository.FindAll(account)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOne(repository repository.AccountRepository, account *domain.Account, id string) (*domain.Account, error) {
	var query = bson.M{"_id": bson.ObjectIdHex(id)}

	err := repository.FindOne(query, account)
	if err != nil {
		return account, err
	}

	return account, nil
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
