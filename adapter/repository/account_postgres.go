package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/pkg/errors"
)

type AccountSQL struct {
	db SQL
}

func NewAccountSQL(db SQL) AccountSQL {
	return AccountSQL{
		db: db,
	}
}

func (a AccountSQL) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	var query = `
		INSERT INTO 
			accounts (id, name, cpf, balance, created_at)
		VALUES 
			($1, $2, $3, $4, $5)
	`

	if err := a.db.ExecuteContext(
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

func (a AccountSQL) UpdateBalance(ctx context.Context, ID domain.AccountID, balance domain.Money) error {
	tx, ok := ctx.Value("TransactionContextKey").(Tx)
	if !ok {
		var err error
		tx, err = a.db.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "error updating account balance")
		}
	}

	query := "UPDATE accounts SET balance = $1 WHERE id = $2"

	if err := tx.ExecuteContext(ctx, query, balance, ID); err != nil {
		return errors.Wrap(err, "error updating account balance")
	}

	return nil
}

func (a AccountSQL) FindAll(ctx context.Context) ([]domain.Account, error) {
	var query = "SELECT * FROM accounts"

	rows, err := a.db.QueryContext(ctx, query)
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

func (a AccountSQL) FindByID(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	tx, ok := ctx.Value("TransactionContextKey").(Tx)
	if !ok {
		var err error
		tx, err = a.db.BeginTx(ctx)
		if err != nil {
			return domain.Account{}, errors.Wrap(err, "error find account by id")
		}
	}

	var (
		query     = "SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE"
		id        string
		name      string
		CPF       string
		balance   int64
		createdAt time.Time
	)

	err := tx.QueryRowContext(ctx, query, ID).Scan(&id, &name, &CPF, &balance, &createdAt)
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

func (a AccountSQL) FindBalance(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		query   = "SELECT balance FROM accounts WHERE id = $1"
		balance int64
	)

	err := a.db.QueryRowContext(ctx, query, ID).Scan(&balance)
	switch {
	case err == sql.ErrNoRows:
		return domain.Account{}, domain.ErrAccountNotFound
	default:
		return domain.NewAccountBalance(domain.Money(balance)), err
	}
}
