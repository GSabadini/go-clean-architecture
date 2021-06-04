package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gsabadini/go-bank-transfer/adapter/repository"

	_ "github.com/lib/pq"
)

type postgresHandler struct {
	db *sql.DB
}

func NewPostgresHandler(c *config) (*postgresHandler, error) {
	var ds = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.host,
		c.port,
		c.user,
		c.database,
		c.password,
	)

	fmt.Println(ds)
	db, err := sql.Open(c.driver, ds)
	if err != nil {
		return &postgresHandler{}, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return &postgresHandler{db: db}, nil
}

func (p postgresHandler) BeginTx(ctx context.Context) (repository.Tx, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return postgresTx{}, err
	}

	return newPostgresTx(tx), nil
}

func (p postgresHandler) ExecuteContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p postgresHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (repository.Rows, error) {
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	row := newPostgresRows(rows)

	return row, nil
}

func (p postgresHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) repository.Row {
	row := p.db.QueryRowContext(ctx, query, args...)

	return newPostgresRow(row)
}

type postgresRow struct {
	row *sql.Row
}

func newPostgresRow(row *sql.Row) postgresRow {
	return postgresRow{row: row}
}

func (pr postgresRow) Scan(dest ...interface{}) error {
	if err := pr.row.Scan(dest...); err != nil {
		return err
	}

	return nil
}

type postgresRows struct {
	rows *sql.Rows
}

func newPostgresRows(rows *sql.Rows) postgresRows {
	return postgresRows{rows: rows}
}

func (pr postgresRows) Scan(dest ...interface{}) error {
	if err := pr.rows.Scan(dest...); err != nil {
		return err
	}

	return nil
}

func (pr postgresRows) Next() bool {
	return pr.rows.Next()
}

func (pr postgresRows) Err() error {
	return pr.rows.Err()
}

func (pr postgresRows) Close() error {
	return pr.rows.Close()
}

type postgresTx struct {
	tx *sql.Tx
}

func newPostgresTx(tx *sql.Tx) postgresTx {
	return postgresTx{tx: tx}
}

func (p postgresTx) ExecuteContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p postgresTx) QueryContext(ctx context.Context, query string, args ...interface{}) (repository.Rows, error) {
	rows, err := p.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	row := newPostgresRows(rows)

	return row, nil
}

func (p postgresTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) repository.Row {
	row := p.tx.QueryRowContext(ctx, query, args...)

	return newPostgresRow(row)
}

func (p postgresTx) Commit() error {
	return p.tx.Commit()
}

func (p postgresTx) Rollback() error {
	return p.tx.Rollback()
}
