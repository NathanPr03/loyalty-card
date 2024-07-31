package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NathanPr03/price-control/pkg/db"
	"loyalty-card/api/generated"
	"net/http"
)

func SetLoyaltyCard(w http.ResponseWriter, request *http.Request) {
	var product generated.PostLoyaltyCardJSONBody

	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if product.CustomerId == "" || product.LoyaltyCard {
		http.Error(w, "Customer ID and product discount cannot be empty", http.StatusBadRequest)
		return
	}

	var dbConnection, err2 = db.ConnectToDb()

	if err2 != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error connecting to database: %v(</h1>", err)
		return
	}

	defer func(dbConnection *sql.DB) {
		_ = dbConnection.Close()
	}(dbConnection)

	query := "UPDATE customer SET has_loyalty_card = $1 WHERE id = $2"
	_, err = dbConnection.Exec(query, product.LoyaltyCard, product.CustomerId)
	if err != nil {
		_, _ = fmt.Fprintf(w, "<h1>Error inserting product discount: </h1>")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Loyalty card set successfully"}`))
}
func init() {
	http.HandleFunc("/loyaltyCard", SetLoyaltyCard)
}
