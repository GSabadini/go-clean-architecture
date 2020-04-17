package repository

import (
	"database/sql"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/pkg/errors"
)

type AccountPostgres struct {
	handler *sql.DB
}

func NewAccountPostgres(handler *sql.DB) AccountPostgres {
	return AccountPostgres{handler: handler}
}

func (a AccountPostgres) Store(account domain.Account) (domain.Account, error) {
	query := `
		INSERT INTO 
			accounts (id, name, cpf, balance, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	_, err := a.handler.Exec(
		query,
		account.ID,
		account.Name,
		account.CPF,
		account.Balance,
		account.CreatedAt,
	)
	if err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

func (a AccountPostgres) UpdateBalance(ID string, balance float64) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"

	if _, err := a.handler.Exec(query, balance, ID); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

func (a AccountPostgres) FindAll() ([]domain.Account, error) {
	var (
		accounts = make([]domain.Account, 0)
		query    = "SELECT * FROM accounts"
	)

	rows, err := a.handler.Query(query)
	if err != nil {
		return accounts, errors.Wrap(err, "error listing accounts")
	}

	for rows.Next() {
		var (
			ID        string
			name      string
			CPF       string
			balance   float64
			createdAt *time.Time
		)

		if err = rows.Scan(&ID, &name, &CPF, &balance, &createdAt); err != nil {
			return accounts, errors.Wrap(err, "error listing accounts")
		}

		account := domain.Account{
			ID:        ID,
			Name:      name,
			CPF:       CPF,
			Balance:   balance,
			CreatedAt: createdAt,
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (a AccountPostgres) FindByID(ID string) (*domain.Account, error) {
	var (
		account   = &domain.Account{}
		query     = "SELECT * FROM accounts WHERE id = $1"
		id        string
		name      string
		CPF       string
		balance   float64
		createdAt *time.Time
	)

	if err := a.handler.QueryRow(query, ID).Scan(&id, &name, &CPF, &balance, &createdAt); err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	account.ID = id
	account.Name = name
	account.CPF = CPF
	account.Balance = balance
	account.CreatedAt = createdAt

	return account, nil
}

func (a AccountPostgres) FindBalance(ID string) (domain.Account, error) {
	var (
		account = domain.Account{}
		query   = "SELECT balance FROM accounts WHERE id = $1"
		balance float64
	)

	if err := a.handler.QueryRow(query, ID).Scan(&balance); err != nil {
		return account, errors.Wrap(err, "error fetching account balance")
	}

	account.Balance = balance

	return account, nil
}
