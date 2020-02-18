package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
)

func Create(repository repository.Account, account domain.Account) error {
	if err := repository.Store(account); err != nil {
		return err
	}

	return nil
}
