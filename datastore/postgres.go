package datastore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const postgresEnvironmentVariable = "POSTGRES_ADDRESS"
const postgresDefaultAddress = "user=postgres password=postgres dbname=logistics sslmode=disabled"

var Postgres *sql.DB

func init() {
	if err := postgresConnect(); err != nil {
		log.Fatalln(err)
	}
}

func postgresAddress() string {
	if pg := os.Getenv(postgresEnvironmentVariable); pg != "" {
		return pg
	}
	return postgresDefaultAddress
}

func postgresConnect() (err error) {
	Postgres, err = sql.Open("postgres", postgresAddress())
	return
}
