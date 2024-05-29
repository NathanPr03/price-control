package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"price-control/pkg/db"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		fmt.Fprintf(w, "<h1>Error connecting to database: %v(</h1>", err)
		return
	}
	defer dbConnection.Close()
	dbConnection.Exec("INSERT INTO products (name) VALUES ('Washing Machine')")
	var id int
	var name string

	query := `SELECT id, name FROM products LIMIT 1`
	row := dbConnection.QueryRow(query)
	switch err := row.Scan(&id, &name); err {
	case sql.ErrNoRows:
		fmt.Fprintln(w, "<h1>No rows were returned!</h1>")
	case nil:
		fmt.Fprintf(w, "<h1>Hello from Go!</h1><p>ID: %d, Name: %s</p>", id, name)
	default:
		fmt.Fprintf(w, "<h1>Error :(: %v</h1>", err)
	}
}
