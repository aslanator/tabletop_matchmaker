package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func GetDB() *sql.DB {
	connStr := "postgres://root:root@fullstack-postgres:5432/test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
