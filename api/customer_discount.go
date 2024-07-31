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
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	discount := calculateDiscount(customer)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(DiscountResponse{Discount: discount})
}

func fetchCustomer(customerID int) (Customer, error) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		println("Error connecting to database: " + err.Error())
		return Customer{}, err
	}

	var customer Customer
	query := "SELECT id, email, has_loyalty_card, total_purchases as SUM(customer_purchases.amount_purchased) FROM customer JOIN customer_purchases on customer.id = customer_purchases.customer_id WHERE id = $1"
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
