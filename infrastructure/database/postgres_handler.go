package database

import (
	"database/sql"
	"os"

	"github.com/gsabadini/go-bank-transfer/repository"

	_ "github.com/lib/pq"
)

//PostgresHandler
type PostgresHandler struct {
	Database *sql.DB
}

//NewPostgresHandler
func NewPostgresHandler(dataSource string) (*PostgresHandler, error) {
	db, err := sql.Open(os.Getenv("POSTGRES_DRIVER"), dataSource)
	if err != nil {
		return &PostgresHandler{}, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &PostgresHandler{Database: db}, nil
}

//Execute
func (p PostgresHandler) Execute(query string, args ...interface{}) error {
	_, err := p.Database.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

//Query
func (p PostgresHandler) Query(query string, args ...interface{}) (repository.Row, error) {
	rows, err := p.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}

	row := NewPostgresRow(rows)

	return row, nil
}

//PostgresRow
type PostgresRow struct {
	Rows *sql.Rows
}

//NewPostgresRow
func NewPostgresRow(rows *sql.Rows) PostgresRow {
	return PostgresRow{Rows: rows}
}

//Scan
func (pr PostgresRow) Scan(dest ...interface{}) error {
	if err := pr.Rows.Scan(dest...); err != nil {
		return err
	}

	return nil
}

//Next
func (pr PostgresRow) Next() bool {
	return pr.Rows.Next()
}

//Next
func (pr PostgresRow) Err() error {
	return pr.Rows.Err()
}

func (pr PostgresRow) Close() error {
	return pr.Rows.Close()
}
