package handler

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

type Product struct {
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Discount       string  `json:"discount"`
	RemainingStock int     `json:"remaining_stock"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer dbConnection.Close()

	query := `INSERT INTO products (name, price, discount, remaining_stock) VALUES ($1, $2, $3, $4)`
	_, err = dbConnection.Exec(query, product.Name, product.Price, product.Discount, product.RemainingStock)
	if err != nil {
		http.Error(w, "Error inserting new product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product added successfully"})
}

func init() {
	http.HandleFunc("/product", AddProduct)
}
