package database

import (
	"database/sql"
	"os"

	"github.com/gsabadini/go-bank-transfer/repository"

	_ "github.com/lib/pq"
)

//postgresHandler armazena a estrutura para o Postgres
type postgresHandler struct {
	database *sql.DB
}

//NewPostgresHandler constrói um novo handler de banco para Postgres
func NewPostgresHandler(dataSource string) (*postgresHandler, error) {
	db, err := sql.Open(os.Getenv("POSTGRES_DRIVER"), dataSource)
	if err != nil {
		return &postgresHandler{}, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &postgresHandler{database: db}, nil
}

//Execute
func (p postgresHandler) Execute(query string, args ...interface{}) error {
	_, err := p.database.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

//Query
func (p postgresHandler) Query(query string, args ...interface{}) (repository.Row, error) {
	rows, err := p.database.Query(query, args...)
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

//Next
func (pr postgresRow) Next() bool {
	return pr.rows.Next()
}

//Err
func (pr postgresRow) Err() error {
	return pr.rows.Err()
}

//CLose
func (pr postgresRow) Close() error {
	return pr.rows.Close()
}
