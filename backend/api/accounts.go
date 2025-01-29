package api

import (
	"crypto-wallet/blockchain"
	"encoding/json"
	"net/http"
)

func CreateAccounts(w http.ResponseWriter, r *http.Request) {
	address, _, err := blockchain.GenerateNewAccount()
	if err != nil {
		http.Error(w, "Failed to generate new account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	seed, err := blockchain.GenerateNewSeedPhrase()
	if err != nil {
		http.Error(w, "Failed to generate seed phrase: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"address": address,
		"seed":    seed,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
