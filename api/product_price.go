package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"price-control/api/generated"
	"price-control/pkg/db"
)

func SetProductPrice(w http.ResponseWriter, request *http.Request) {
	var product generated.PostProductPriceJSONBody

	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConnection, err := db.ConnectToDb()
	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error connecting to database: %v(</h1>", err)
		return
	}

	defer func(dbConnection *sql.DB) {
		_ = dbConnection.Close()
	}(dbConnection)

	query := "UPDATE products SET price = $1 WHERE name = $2"
	fmt.Println(query)
	_, err = dbConnection.Exec(query, product.Price, product.ProductName)
	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error inserting product price: </h1>")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Product price added successfully"}`))
}

func init() {
	http.HandleFunc("/productPrice", SetProductPrice)
}
