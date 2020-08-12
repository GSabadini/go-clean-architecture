package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gsabadini/go-bank-transfer/repository"
	_ "github.com/lib/pq"
)

//postgresHandler armazena a estrutura para o Postgres
type postgresHandler struct {
	db *sql.DB
}

//NewPostgresHandler constrói um novo handler de banco para Postgres
func NewPostgresHandler(c *config) (*postgresHandler, error) {
	var ds = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.host,
		c.port,
		c.user,
		c.database,
		c.password,
	)

	db, err := sql.Open(os.Getenv("POSTGRES_DRIVER"), ds)
	if err != nil {
		return &postgresHandler{}, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &postgresHandler{db: db}, nil
}

func (p postgresHandler) BeginTx(ctx context.Context) (repository.Tx, error) {
	return p.db.BeginTx(ctx, nil)
}

//postgresRow
type postgresTx struct {
	tx *sql.Tx
}

func (p postgresTx) Commit() error {
	return p.tx.Commit()
}
func (p postgresTx) Rollback() error {
	return p.tx.Rollback()
}

//newPostgresTx
func newPostgresTx(tx *sql.Tx) postgresTx {
	return postgresTx{tx: tx}
}

//ExecuteContext
func (p postgresHandler) ExecuteContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

//Query
func (p postgresHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (repository.Row, error) {
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	row := newPostgresRow(rows)

	return row, nil
}

//postgresRow
type postgresRow struct {
	rows *sql.Rows
}

//newPostgresRow
func newPostgresRow(rows *sql.Rows) postgresRow {
	return postgresRow{rows: rows}
}

//Scan
func (pr postgresRow) Scan(dest ...interface{}) error {
	if err := pr.rows.Scan(dest...); err != nil {
		return err
	}

	return nil
}

//Next retorna o método next
func (pr postgresRow) Next() bool {
	return pr.rows.Next()
}

//Err retorna o método err
func (pr postgresRow) Err() error {
	return pr.rows.Err()
}

//Close retorna o método close
func (pr postgresRow) Close() error {
	return pr.rows.Close()
}
