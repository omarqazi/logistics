package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20151005012132(txn *sql.Tx) {
	sql := `
	create table orders (
		id uuid not null,
		createdat timestamp with time zone not null,
		updatedat timestamp with time zone not null,
		orderedat timestamp with time zone not null,
		delivered_at timestamp with time zone,
		seller_id uuid not null,
		buyer_id uuid not null,
		quantity_purchased integer default 0,
		constraint order_pk primary key (id)
	)
	with (
		OIDS=FALSE
	);
	create index orders_placed on orders(orderedat);
	create index orders_delivered on orders(delivered_at);
	create index orders_sellers on orders(seller_id);
	create index orders_buyers on orders(buyer_id);
	`
	if _, err := txn.Exec(sql); err != nil {
		fmt.Println("Error creating inventory table:", err)
	}
}

// Down is executed when this migration is rolled back
func Down_20151005012132(txn *sql.Tx) {
	if _, err := txn.Exec("drop table orders;"); err != nil {
		fmt.Println("Error dropping orders table:", err)
	}
}
