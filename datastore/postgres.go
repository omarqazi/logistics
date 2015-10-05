package datastore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"log"
	"os"
)

const postgresEnvironmentVariable = "POSTGRES_ADDRESS"
const postgresDefaultAddress = "user=postgres password=postgres dbname=logistics sslmode=disable"

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

func NewUUID() string {
	return uuid.NewV4().String()
}
