package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20151005004856(txn *sql.Tx) {
	sql := `
	create extension postgis;
	create table users (
		id uuid not null,
		createdat timestamp with time zone not null,
		updatedat timestamp with time zone not null,
		public_key text not null,
		device_id text,
		constraint users_pk primary key (id)
	)
	with (
		OIDS=FALSE
	);
	alter table users add column location geometry(Point,4326);
	create index users_gist on users using GIST(location);
	create index users_updated on users(updatedat);
	`
	if _, err := txn.Exec(sql); err != nil {
		fmt.Println("Error creating users table:", err)
	}
}

// Down is executed when this migration is rolled back
func Down_20151005004856(txn *sql.Tx) {
	if _, err := txn.Exec("drop table users"); err != nil {
		fmt.Println("Error dropping users table:", err)
	}
}
