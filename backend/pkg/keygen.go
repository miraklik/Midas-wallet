package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKeys() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}
