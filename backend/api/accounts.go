package api

import (
	"crypto-wallet/blockchain"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func GenerateAddress() (walletAddress string, err error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	pubKey := privKey.PublicKey
	publickKeyBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)

	hash := sha256.Sum256(publickKeyBytes)
	walletAddress = hex.EncodeToString(hash[:])

	return walletAddress, nil
}
