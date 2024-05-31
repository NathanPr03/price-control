package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"price-control/api/generated"
	"price-control/pkg/db"
)

func SetProductPrice(w http.ResponseWriter, request *http.Request) {
	var product generated.PostProductPriceJSONBody

	// Decode the JSON request body into the struct
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Print the extracted values
	fmt.Printf("Product Name: %s, Price: %f\n", product.ProductName, product.Price)

	dbConnection, err := db.ConnectToDb()
	if err != nil {
		fmt.Fprintf(w, "<h1>Error connecting to database: %v(</h1>", err)
		return
	}

	defer dbConnection.Close()
	query := fmt.Sprintf("UPDATE products SET price = %s WHERE name = '%s'", product.Price, product.ProductName)
	_, err = dbConnection.Exec(query)
	if err != nil {
		fmt.Fprintf(w, "<h1>Error inserting product price: %v</h1>", err)
		return

	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Product price added successfully"}`))
}

func init() {
	http.HandleFunc("/productPrice", SetProductPrice)
}
