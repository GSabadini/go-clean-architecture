package validator

import (
	"errors"
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
func NewValidatorFactory(instance int) (Validator, error) {
	switch instance {
	case InstanceGoPlayground:
		return NewGoPlayground()
	default:
		return nil, errInvalidValidatorInstance
	}
}
