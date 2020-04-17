package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type PostgresHandler struct {
	Database *sql.DB
}

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
