package postgres

import (
	"fmt"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"

	"github.com/pkg/errors"
)

//AccountRepository representa um repositório para manipulação de dados de uma conta
type AccountRepository struct {
	handler database.SQLHandler
}

//NewAccountRepository constrói um repository com suas dependências
func NewAccountRepository(handler database.SQLHandler) AccountRepository {
	return AccountRepository{handler: handler}
}

//Store realiza a inserção de uma conta no banco de dados
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

//UpdateBalance realiza a atualização do saldo de uma conta no banco de dados
func (a AccountRepository) UpdateBalance(ID string, balance float64) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"

	if err := a.handler.Execute(query, balance, ID); err != nil {
		return errors.Wrap(domain.ErrNotFound, "error updating account balance")
	}

	return nil
}

//FindAll realiza a busca de todas as contas no banco de dados
func (a AccountRepository) FindAll() ([]domain.Account, error) {
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

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	return accounts, nil
}

//FindByID realiza a busca de uma conta no banco de dados
func (a AccountRepository) FindByID(ID string) (*domain.Account, error) {
	var (
		account   = &domain.Account{}
		query     = "SELECT * FROM accounts WHERE id = $1"
		id        string
		name      string
		CPF       string
		balance   float64
		createdAt *time.Time
	)

	row, err := a.handler.Query(query, ID)
	if err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	row.Next()
	if err = row.Scan(&id, &name, &CPF, &balance, &createdAt); err != nil {
		return account, errors.Wrap(err, "error fetching account")
	}

	account.ID = id
	account.Name = name
	account.CPF = CPF
	account.Balance = balance
	account.CreatedAt = createdAt

	return account, nil
}

//FindBalance realiza a busca do saldo de uma conta no banco de dados
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

	row.Next()
	if err := row.Scan(&balance); err != nil {
		return account, errors.Wrap(err, "error fetching account balance")
	}

	account.Balance = balance

	return account, nil
}
