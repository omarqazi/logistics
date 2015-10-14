package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20151013175608(txn *sql.Tx) {
	sql := `
	create table locations (
		id uuid not null,
		createdat timestamp with time zone not null,
		updatedat timestamp with time zone not null,
		recordedat timestamp with time zone not null,
		user_id uuid not null,
		constraint location_pk primary key (id)
	)
	with (
		OIDS=FALSE
	);
	create index locations_recorded on locations(recordedat);
	alter table locations add column location geometry(Point,4326);
	create index locations_gist on locations using GIST(location);
	alter table users add column location_recorded timestamp with time zone;
	alter table users add column latest_location_id uuid;
	`
	if _, err := txn.Exec(sql); err != nil {
		fmt.Println("Error creating locations table:", err)
	}
}

// Down is executed when this migration is rolled back
func Down_20151013175608(txn *sql.Tx) {
	if _, err := txn.Exec("drop table locations"); err != nil {
		fmt.Println("Error dropping locations table:", err)
	}
}
