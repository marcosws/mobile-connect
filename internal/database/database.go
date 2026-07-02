package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func New() *sql.DB {

	err := os.MkdirAll(
		"storage/database",
		0755,
	)

	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(
		"sqlite",
		"storage/database/mobile-connect.db",
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
