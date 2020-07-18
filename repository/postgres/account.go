package postgres

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"

	"github.com/pkg/errors"
)

//AccountRepository armazena a estrutura de dados de um repositório de Account
type AccountRepository struct {
	handler repository.SQLHandler
}

//NewAccountRepository constrói um AccountRepository com suas dependências
func NewAccountRepository(handler repository.SQLHandler) AccountRepository {
	return AccountRepository{handler: handler}
}

//Store insere uma Account no database
func (a AccountRepository) Store(account domain.Account) (domain.Account, error) {
	query := `
		INSERT INTO 
			accounts (id, name, cpf, balance, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := a.handler.Execute(
		query,
		account.ID,
		account.Name,
		account.CPF,
		account.Balance,
		account.CreatedAt,
	); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

//UpdateBalance atualiza o Balance de uma Account no database
func (a AccountRepository) UpdateBalance(ID string, balance float64) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"

	if err := a.handler.Execute(query, balance, ID); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAlL busca todas as Account no database
func (a AccountRepository) FindAll() ([]domain.Account, error) {
	var (
		accounts = make([]domain.Account, 0)
		query    = "SELECT * FROM accounts"
	)

	rows, err := a.handler.Query(query)
	if err != nil {
		return accounts, errors.Wrap(err, "error listing accounts")
	}

	defer rows.Close()
	for rows.Next() {
		var (
			ID        string
			name      string
			CPF       string
			balance   float64
			createdAt time.Time
		)

		if err = rows.Scan(&ID, &name, &CPF, &balance, &createdAt); err != nil {
			return accounts, errors.Wrap(err, "error listing accounts")
		}

		accounts = append(accounts, domain.Account{
			ID:        ID,
			Name:      name,
			CPF:       CPF,
			Balance:   balance,
			CreatedAt: createdAt,
		})
	}

	if err = rows.Err(); err != nil {
		return []domain.Account{}, err
	}

	return accounts, nil
}

//FindByID busca uma Account por ID no database
func (a AccountRepository) FindByID(ID string) (domain.Account, error) {
	var (
		account   = domain.Account{}
		query     = "SELECT * FROM accounts WHERE id = $1"
		id        string
		name      string
		CPF       string
		balance   float64
		createdAt time.Time
	)

	row, err := a.handler.Query(query, ID)
	if err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	defer row.Close()
	row.Next()
	if err = row.Scan(&id, &name, &CPF, &balance, &createdAt); err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	if err = row.Err(); err != nil {
		return domain.Account{}, err
	}

	account.ID = id
	account.Name = name
	account.CPF = CPF
	account.Balance = balance
	account.CreatedAt = createdAt

	return account, nil
}

//FindBalance busca o Balance de uma Account no database
func (a AccountRepository) FindBalance(ID string) (domain.Account, error) {
	var (
		account = domain.Account{}
		query   = "SELECT balance FROM accounts WHERE id = $1"
		balance float64
	)

	row, err := a.handler.Query(query, ID)
	if err != nil {
		return account, errors.Wrap(err, "error fetching account balance")
	}

	defer row.Close()
	row.Next()
	if err := row.Scan(&balance); err != nil {
		return account, errors.Wrap(err, "error fetching account balance")
	}

	if err = row.Err(); err != nil {
		return domain.Account{}, err
	}

	account.Balance = balance

	return account, nil
}
