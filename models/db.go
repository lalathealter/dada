package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", produceConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

func produceConnectionString() string {
	dbhost := os.Getenv("dbhost")
	dbname := os.Getenv("dbname")
	dbport := os.Getenv("dbport")
	dbuser := os.Getenv("dbuser")
	dbpass := os.Getenv("dbpass")
	return fmt.Sprintf(`host=%s port=%s user=%s
    password=%s dbname=%s sslmode=disable`,
		dbhost, dbport, dbuser, dbpass, dbname,
	)
}
