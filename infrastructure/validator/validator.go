package validator

import (
	"errors"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
)

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

//NewValidator
func NewValidator(instance int, log logger.Logger) (Validator, error) {
	switch instance {
	case InstanceGoPlayground:
		return NewGoPlayground(log), nil
	default:
		return nil, errInvalidValidatorInstance
	}
}
