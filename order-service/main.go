package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var orders = map[string]*Order{}
var mu sync.Mutex

func main() {
	http.HandleFunc("/create", handleCreateOrder)
	http.HandleFunc("/cancel", handleCancelOrder)
	http.ListenAndServe(":8082", nil)
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID string `json:"order_id"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	orders[req.OrderID] = &Order{ID: req.OrderID, Status: "created"}
	w.WriteHeader(http.StatusOK)
}

func handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID string `json:"order_id"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	if order, ok := orders[req.OrderID]; ok {
		order.Status = "canceled"
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Order not found", http.StatusBadRequest)
	}
}
