package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Account struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

var accounts = map[string]*Account{}

// add some accounts
func init() {
	accounts["1001"] = &Account{ID: "1001", Balance: 100}
	accounts["1002"] = &Account{ID: "1002", Balance: 200}
	accounts["1003"] = &Account{ID: "1003", Balance: 300}
	accounts["1004"] = &Account{ID: "1004", Balance: 400}
	accounts["1005"] = &Account{ID: "1005", Balance: 500}
	accounts["1006"] = &Account{ID: "1006", Balance: 600}
	accounts["1007"] = &Account{ID: "1007", Balance: 700}
	accounts["1008"] = &Account{ID: "1008", Balance: 800}
	accounts["1009"] = &Account{ID: "1009", Balance: 900}
	accounts["1010"] = &Account{ID: "1010", Balance: 1000}
}

var mu sync.Mutex

func main() {
	http.HandleFunc("/debit", handleDebit)
	http.HandleFunc("/credit", handleCredit)
	http.ListenAndServe(":8081", nil)
}

func handleDebit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	if account, ok := accounts[req.AccountID]; ok && account.Balance >= req.Amount {
		account.Balance -= req.Amount
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Insufficient balance or account not found", http.StatusBadRequest)
	}
}

func handleCredit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	if account, ok := accounts[req.AccountID]; ok {
		account.Balance += req.Amount
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Account not found", http.StatusBadRequest)
	}
}
