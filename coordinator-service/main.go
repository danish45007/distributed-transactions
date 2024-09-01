package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/2pc", handle2PC)
	http.HandleFunc("/saga", handleSaga)
	http.ListenAndServe(":8080", nil)
}

// Handle 2PC Transaction
func handle2PC(w http.ResponseWriter, r *http.Request) {
	// Assume transaction involves debit in Account Service and order creation in Order Service
	accountID := "acc123"
	orderID := "ord123"
	amount := 100.0

	// Prepare phase (both services are asked if they can commit)
	if !debitAccount(accountID, amount) || !createOrder(orderID) {
		// Rollback if either fails
		creditAccount(accountID, amount)
		cancelOrder(orderID)
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	// Commit phase (do nothing in this simple case since changes are already made)
	w.WriteHeader(http.StatusOK)
}

// Handle Saga Transaction
func handleSaga(w http.ResponseWriter, r *http.Request) {
	// Assume transaction involves debit in Account Service and order creation in Order Service
	accountID := "acc123"
	orderID := "ord123"
	amount := 100.0

	// Saga pattern: Perform each step, rolling back if any step fails
	if !debitAccount(accountID, amount) {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	if !createOrder(orderID) {
		// Compensate by crediting back the account
		creditAccount(accountID, amount)
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	// If everything succeeds, transaction is complete
	w.WriteHeader(http.StatusOK)
}

// Helper functions to communicate with the Account and Order services
func debitAccount(accountID string, amount float64) bool {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"account_id": accountID,
		"amount":     amount,
	})

	resp, err := http.Post("http://localhost:8081/debit", "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to debit account")
		return false
	}
	return true
}

func creditAccount(accountID string, amount float64) bool {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"account_id": accountID,
		"amount":     amount,
	})

	resp, err := http.Post("http://localhost:8081/credit", "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to credit account")
		return false
	}
	return true
}

func createOrder(orderID string) bool {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"order_id": orderID,
	})

	resp, err := http.Post("http://localhost:8082/create", "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to create order")
		return false
	}
	return true
}

func cancelOrder(orderID string) bool {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"order_id": orderID,
	})

	resp, err := http.Post("http://localhost:8082/cancel", "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to cancel order")
		return false
	}
	return true
}
