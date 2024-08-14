package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NathanPr03/price-control/api/generated"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

var dbConnection, _ = db.ConnectToDb()

func ProductDiscountHandler(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		SetProductDiscount(w, request)
	case http.MethodGet:
		GetDiscountedProducts(w, request)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetDiscountedProducts(w http.ResponseWriter, r *http.Request) {
	discountType := r.URL.Query().Get("discountType")
	if discountType == "" {
		http.Error(w, "Missing discountType parameter", http.StatusBadRequest)
		return
	}

	query := "SELECT name FROM products WHERE discount = $1"
	rows, err := dbConnection.Query(query, discountType)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var products []string
	for rows.Next() {
		var productName string
		if err := rows.Scan(&productName); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, productName)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]string{"products": products}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func SetProductDiscount(w http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

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

	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error connecting to database: %v(</h1>", err)
		return
	}

	defer func(dbConnection *sql.DB) {
		_ = dbConnection.Close()
	}(dbConnection)

	query := "UPDATE products SET discount = $1 WHERE name = $2"
	_, err = dbConnection.Exec(query, product.DiscountType, product.ProductName)
	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error inserting product discount: </h1>")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Product discount added successfully"}`))
}

func init() {
	http.HandleFunc("/productDiscount", ProductDiscountHandler)
}
