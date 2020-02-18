package repository

import "github.com/gsabadini/go-stone/domain"

type Repository interface {
	Store(domain.Account) error
}
