package validator

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
)

//Validator é uma abstração para os validator da aplicação
type Validator interface {
	Validate(interface{}) error
	Messages() []string
}

var (
	errInvalidValidatorInstance = errors.New("invalid validator instance")
)

const (
	InstanceGoPlayground int = iota
)

//NewValidatorFactory retorna a instância de um validator
func NewValidatorFactory(instance int, log logger.Logger) (Validator, error) {
	switch instance {
	case InstanceGoPlayground:
		return NewGoPlayground(log), nil
	default:
		return nil, errInvalidValidatorInstance
	}
}
