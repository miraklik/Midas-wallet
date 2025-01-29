package blockchain

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func GenerateNewAccount() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("Failed to generate private key: %v", err)
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	privateKeyHex := crypto.FromECDSA(privateKey)
	return address, string(privateKeyHex), nil
}

func GenerateNewSeedPhrase() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Printf("Failed to generate entropy: %v", err)
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Printf("Failed to generate mnemonic: %v", err)
		return "", err
	}

	return mnemonic, nil
}
