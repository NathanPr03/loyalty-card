package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
	"strconv"
)

type Customer struct {
	ID             int     `json:"id"`
	Email          string  `json:"email"`
	HasLoyaltyCard bool    `json:"has_loyalty_card"`
	TotalPurchases float64 `json:"total_purchases"`
}

type DiscountResponse struct {
	Discount float64 `json:"discount"`
}

func Discount(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		http.Error(w, "customer_id query parameter is required", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		http.Error(w, "customer_id must be an integer", http.StatusBadRequest)
		return
	}

	customer, err := fetchCustomer(customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "customer not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	discount := calculateDiscount(customer)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(DiscountResponse{Discount: discount})
}

func fetchCustomer(customerID int) (Customer, error) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		println("Error connecting to database: " + err.Error())
		return Customer{}, err
	}

	var customer Customer

	query := "SELECT customer.id, email, has_loyalty_card, SUM(customer_purchases.amount_purchased) as total_purchases FROM customer JOIN customer_purchases on customer.id = customer_purchases.customer_id WHERE customer.id = $1 GROUP BY customer.id"
	row := dbConnection.QueryRow(query, customerID)
	err = row.Scan(&customer.ID, &customer.Email, &customer.HasLoyaltyCard, &customer.TotalPurchases)

	return customer, err
}

func calculateDiscount(customer Customer) float64 {
	discount := 0.0
	if customer.HasLoyaltyCard {
		discount = discount + 0.2
	} else if customer.TotalPurchases > 1000 {
		discount = discount + 0.1
	}

	return discount
}
