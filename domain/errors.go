package domain

import "github.com/pkg/errors"

var (
	//ErrNotFound é um erro de não encontrado
	ErrNotFound = errors.New("Not found")

	//ErrInsufficientBalance é um erro de saldo insuficiente
	ErrInsufficientBalance = errors.New("Origin account does not have sufficient balance")
)
