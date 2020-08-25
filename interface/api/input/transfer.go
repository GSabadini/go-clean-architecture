package input

import (
	"errors"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
)

//Transfer armazena a estrutura de dados de entrada da API
type Transfer struct {
	AccountOriginID      string `json:"account_origin_id" validate:"required,uuid4"`
	AccountDestinationID string `json:"account_destination_id" validate:"required,uuid4"`
	Amount               int64  `json:"amount" validate:"gt=0,required"`
}

func (t Transfer) Validate(validator validator.Validator) []string {
	var (
		msgs              []string
		errAccountsEquals = errors.New("account origin equals destination account")
		accountIsEquals   = t.AccountOriginID == t.AccountDestinationID
		accountsIsEmpty   = t.AccountOriginID == "" && t.AccountDestinationID == ""
	)

	if !accountsIsEmpty && accountIsEquals {
		msgs = append(msgs, errAccountsEquals.Error())
	}

	err := validator.Validate(t)
	if err != nil {
		for _, msg := range validator.Messages() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}
