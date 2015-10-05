package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20151005010829(txn *sql.Tx) {
	sql := `
	create table inventory (
		id uuid not null,
		createdat timestamp with time zone not null,
		updatedat timestamp with time zone not null,
		quantity integer default 0,
		title text not null,
		description text not null,
		unit_price numeric default 0.00,
		photo_url text,
		user_id uuid not null,
		constraint inventory_pk primary key (id)
	)
	with (
		OIDS=FALSE
	);
	create index inventory_title on inventory(title);
	create index inventory_user on inventory(user_id);
	`
	if _, err := txn.Exec(sql); err != nil {
		fmt.Println("Error creating inventory table:", err)
	}
}

// Down is executed when this migration is rolled back
func Down_20151005010829(txn *sql.Tx) {
	if _, err := txn.Exec("drop table inventory"); err != nil {
		fmt.Println("Error dropping inventory table:", err)
	}
}
