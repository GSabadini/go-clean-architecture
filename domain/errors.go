package domain

import "github.com/pkg/errors"

var (
	//ErrNotFound é um erro de não encontrado
	ErrNotFound = errors.New("not found")

	//ErrInsufficientBalance é um erro de saldo insuficiente
	ErrInsufficientBalance = errors.New("origin account does not have sufficient balance")

	//ErrInsufficientBalance é um erro ao atualizar o saldo de uma conta
	ErrUpdateBalance = errors.New("error update account balance")
)
