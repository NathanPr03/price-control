package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"price-control/api/generated"
	"price-control/pkg/db"
)

func SetProductDiscount(w http.ResponseWriter, request *http.Request) {
	var product generated.PostProductDiscountJSONBody

	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if product.DiscountType == "" || product.ProductName == "" {
		http.Error(w, "Discount type and product name cannot be empty", http.StatusBadRequest)
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

	query := "UPDATE products SET discount = $1 WHERE name = $2"
	fmt.Println(query)
	_, err = dbConnection.Exec(query, product.DiscountType, product.ProductName)
	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error inserting product discount: </h1>")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Product discount added successfully"}`))
}

func init() {
	http.HandleFunc("/productDiscount", SetProductDiscount)
}
