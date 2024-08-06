package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/NathanPr03/price-control/pkg/db"
	"log"
	"net/http"
)

type BigCustomer struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HasLoyaltyCard bool   `json:"has_loyalty_card"`
	FavoriteItem   string `json:"favorite_item"`
}

func AllCustomers(w http.ResponseWriter, r *http.Request) {
	var dbConnection, err = db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	customersQuery := "SELECT id, email, has_loyalty_card FROM customer"
	rows, err := dbConnection.Query(customersQuery)
	if err != nil {
		http.Error(w, "Failed to retrieve customers", http.StatusInternalServerError)
		log.Printf("Failed to execute query: %v", err)
		return
	}
	defer rows.Close()

	// Prepare a slice to hold all customer data
	var customers []BigCustomer

	// Iterate over each customer
	for rows.Next() {
		var customer BigCustomer
		if err := rows.Scan(&customer.ID, &customer.Email, &customer.HasLoyaltyCard); err != nil {
			http.Error(w, "Failed to parse customer data", http.StatusInternalServerError)
			log.Printf("Failed to scan row: %v", err)
			return
		}

		// Query to find the customer's favorite item
		favoriteItemQuery := `
			SELECT product_purchased
			FROM customer_purchases
			WHERE customer_id = ?
			GROUP BY product_purchased
			ORDER BY SUM(amount_purchased) DESC
			LIMIT 1
		`
		log.Printf("Customer is: %v", customer)
		err := dbConnection.QueryRow(favoriteItemQuery, customer.ID).Scan(&customer.FavoriteItem)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				customer.FavoriteItem = "None" // Default value if no purchases are found
			} else {
				http.Error(w, "Failed to retrieve favorite item", http.StatusInternalServerError)
				log.Printf("Failed to execute favorite item query: %v", err)
				return
			}
		}

		// Add the customer to the slice
		customers = append(customers, customer)
	}

	// Check for any row iteration error
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating over customer rows", http.StatusInternalServerError)
		log.Printf("Row iteration error: %v", err)
		return
	}

	// Convert the slice of customers to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(customers); err != nil {
		http.Error(w, "Failed to encode customers to JSON", http.StatusInternalServerError)
		log.Printf("JSON encoding error: %v", err)
	}
}

func init() {
	http.HandleFunc("/customers", AllCustomers)
}
