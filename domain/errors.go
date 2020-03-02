package domain

import "github.com/pkg/errors"

var (
	ErrNotFound            = errors.New("Not found")
	ErrInsufficientBalance = errors.New("Origin account does not have sufficient balance")
)
