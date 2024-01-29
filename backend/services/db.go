package services

import (
	"database/sql"
	"log"

	env "palhost/env"

	_ "github.com/lib/pq"
)

func ConnectToDB(config *env.EnvConfig) *sql.DB {
	db, err := sql.Open("postgres", config.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
