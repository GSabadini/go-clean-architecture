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

	return &PostgresHandler{Database: db}, nil
}

func (p PostgresHandler) Store(statement string, args ...interface{}) error {
	//sqlStatement := `
	//	INSERT INTO accounts (id, name, cpf, balance, created_at)
	//	VALUES ($1, $2, $3, $4, $5)
	//	RETURNING id`

	//time := time.Now()
	//id := 0
	_, err := p.Database.Exec(statement, args...)
	if err != nil {
		return err
	}

	//p.Database.

	//fmt.Println("New record ID is:", id)
	return nil
}

func (p PostgresHandler) Update(statement string, args ...interface{}) error {
	//sqlStatementUpdate := `
	//	UPDATE accounts
	//	SET name = $2
	//	WHERE id = $1;`
	_, err := p.Database.Exec(statement, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresHandler) FindAll(statement string, args ...interface{}) error {
	_, err := p.Database.Exec(statement, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresHandler) FindOne(statement string, args ...interface{}) error {
	_, err := p.Database.Exec(statement, args...)
	if err != nil {
		return err
	}

	return nil
}
