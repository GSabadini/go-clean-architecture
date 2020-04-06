package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestNewPostgresHandler(t *testing.T) {
	const (
		host     = "127.0.0.1"
		port     = "5432"
		user     = "dev"
		password = "dev"
		dbname   = "bank"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	sqlStatement := `
		INSERT INTO accounts (id, name, cpf, balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	time := time.Now()
	id := 0
	err = db.QueryRow(sqlStatement, "2", "Gabriel", "070905114", 10.42, &time).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)


	sqlStatementUpdate := `
		UPDATE accounts
		SET name = $2
		WHERE id = $1;`
	res, err := db.Exec(sqlStatementUpdate, "2", "Sabadini")
	if err != nil {
		panic(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)

	sqlStatementDelete := `
		DELETE FROM accounts
		WHERE id = $1;`
	_, err = db.Exec(sqlStatementDelete, "1")
	if err != nil {
		panic(err)
	}
}