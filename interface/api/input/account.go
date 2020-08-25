package input

import "github.com/gsabadini/go-bank-transfer/infrastructure/validator"

//Account armazena a estrutura de dados de entrada da API
type Account struct {
	Name    string `json:"name" validate:"required"`
	CPF     string `json:"cpf" validate:"required"`
	Balance int64  `json:"balance" validate:"gt=0,required"`
}

type AccountID struct {
	ID string `json:"id" validate:"required"`
}

func (a Account) Validate(validator validator.Validator) []string {
	var msgs []string

	err := validator.Validate(a)
	if err != nil {
		for _, msg := range validator.Messages() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}
