package domain

import "time"

//Account armazena a estrutura de uma conta
type Account struct {
	ID        string     `json:"id,omitempty" bson:"id"`
	Name      string     `json:"name,omitempty" bson:"name"`
	CPF       string     `json:"cpf,omitempty" bson:"cpf"`
	Balance   float64    `json:"balance" bson:"balance"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
}

//NewAccount cria uma conta
func NewAccount(name string, CPF string, balance float64) Account {
	timeNow := time.Now()

	return Account{
		ID:        uuid(),
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: &timeNow,
	}
}

//GetBalance retorna o saldo
func (a *Account) GetBalance() float64 {
	return a.Balance
}

//Deposit adiciona um valor no saldo
func (a *Account) Deposit(amount float64) {
	a.Balance += amount
}

//Withdraw remove um valor do saldo
func (a *Account) Withdraw(amount float64) error {
	if a.Balance < amount {
		return ErrInsufficientBalance
	}

	a.Balance -= amount

	return nil
}
