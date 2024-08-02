package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

type WholeProduct struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Discount       *string `json:"discount"`
	RemainingStock int     `json:"remaining_stock"`
}

func AllProducts(w http.ResponseWriter, r *http.Request) {
	// Implement all products endpoint based on the defined specification
	dbConnection, _ := db.ConnectToDb()
	defer dbConnection.Close()
	rows, err := dbConnection.Query("SELECT * FROM products")
	if err != nil {

	}

	var products []WholeProduct

	for rows.Next() {
		var product WholeProduct
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Discount, &product.RemainingStock); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	response := map[string][]WholeProduct{"products": products}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func init() {
	http.HandleFunc("/products", AllProducts)
}
