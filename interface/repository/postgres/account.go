package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/interface/repository"

	"github.com/pkg/errors"
)

type AccountRepository struct {
	handler repository.SQLHandler
}

func NewAccountRepository(h repository.SQLHandler) AccountRepository {
	return AccountRepository{
		handler: h,
	}
}

func (a AccountRepository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	var query = `
		INSERT INTO 
			accounts (id, name, cpf, balance, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := a.handler.ExecuteContext(
		ctx,
		query,
		account.ID(),
		account.Name(),
		account.CPF(),
		account.Balance(),
		account.CreatedAt(),
	); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

func (a AccountRepository) UpdateBalance(ctx context.Context, ID domain.AccountID, balance domain.Money) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"

	if err := a.handler.ExecuteContext(ctx, query, balance, ID); err != nil {
		return errors.Wrap(err, "error updating account balance")
	}

	return nil
}

func (a AccountRepository) FindAll(ctx context.Context) ([]domain.Account, error) {
	var query = "SELECT * FROM accounts"

	rows, err := a.handler.QueryContext(ctx, query)
	if err != nil {
		return []domain.Account{}, errors.Wrap(err, "error listing accounts")
	}

	var accounts = make([]domain.Account, 0)
	for rows.Next() {
		var (
			ID        string
			name      string
			CPF       string
			balance   int64
			createdAt time.Time
		)

		if err = rows.Scan(&ID, &name, &CPF, &balance, &createdAt); err != nil {
			return []domain.Account{}, errors.Wrap(err, "error listing accounts")
		}

		accounts = append(accounts, domain.NewAccount(
			domain.AccountID(ID),
			name,
			CPF,
			domain.Money(balance),
			createdAt,
		))
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return []domain.Account{}, err
	}

	return accounts, nil
}

func (a AccountRepository) FindByID(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		query     = "SELECT * FROM accounts WHERE id = $1"
		id        string
		name      string
		CPF       string
		balance   int64
		createdAt time.Time
	)

	err := a.handler.QueryRowContext(ctx, query, ID).Scan(&id, &name, &CPF, &balance, &createdAt)
	switch {
	case err == sql.ErrNoRows:
		return domain.Account{}, domain.ErrAccountNotFound
	default:
		return domain.NewAccount(
			domain.AccountID(id),
			name,
			CPF,
			domain.Money(balance),
			createdAt,
		), err
	}
}

func (a AccountRepository) FindBalance(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		query   = "SELECT balance FROM accounts WHERE id = $1"
		balance int64
	)

	err := a.handler.QueryRowContext(ctx, query, ID).Scan(&balance)
	switch {
	case err == sql.ErrNoRows:
		return domain.Account{}, domain.ErrAccountNotFound
	default:
		return domain.NewAccountBalance(domain.Money(balance)), err
	}
}
