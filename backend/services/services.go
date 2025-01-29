package services

import (
	"context"
	"crypto-wallet/utils"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	Client          *ethclient.Client
	ContractAddress common.Address
	PrivateKey      *ecdsa.PrivateKey
	Contract        *bind.BoundContract
}

func NewEthreumService(rpcURL, contractAddress, privateKeyHex, abiJSON string) (*EthereumService, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Printf("Failed to parse private key: %v", err)
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddressHex := common.HexToAddress(contractAddress)

	parseABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(contractAddressHex, parseABI, client, client, client)
	if err != nil {
		log.Printf("Failed to bind contract: %v", err)
	}

	return &EthereumService{
		Client:          client,
		ContractAddress: contractAddressHex,
		PrivateKey:      privateKey,
		Contract:        contract,
	}, nil
}

func (es *EthereumService) createAuth() (*bind.TransactOpts, error) {
	chainID, err := es.Client.ChainID(context.Background())
	if err != nil {
		log.Printf("Failed to get Chain id: %v", err)
		return nil, fmt.Errorf("failed to get Chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(es.PrivateKey, chainID)
	if err != nil {
		log.Printf("Failed to create auth: %v", err)
		return nil, fmt.Errorf("failed to create auth: %v", err)
	}

	auth.GasLimit = uint64(21000)
	gasPrice, err := es.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("Failed to suggest gas price: %v", err)
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	return auth, nil
}

func (es *EthereumService) Withdraw(to, amount string) error {
	if err := utils.CheckAddress(to); err != nil {
		log.Printf("Failed to check address: %v", err)
		return fmt.Errorf("failed to check address: %v", err)
	}

	auth, err := es.createAuth()
	if err != nil {
		return err
	}

	contractABI, err := os.ReadFile("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
		return fmt.Errorf("failed to read ABI file: %v", err)
	}
	parseABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(es.ContractAddress, parseABI, es.Client, es.Client, es.Client)
	tx, err := contract.Transact(auth, "withdraw", to, amount)
	if err != nil {
		log.Printf("Failed to withdraw: %v", err)
		return fmt.Errorf("failed to withdraw: %v", err)
	}

	fmt.Printf("Withdrawal transaction sent: %s\n", tx.Hash().Hex())
	return nil
}

func (es *EthereumService) Deposit(amount string) error {
	auth, err := es.createAuth()
	if err != nil {
		return err
	}

	contractABI, err := os.ReadFile("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
		return fmt.Errorf("failed to read ABI file: %v", err)
	}
	parseABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(es.ContractAddress, parseABI, es.Client, es.Client, es.Client)
	tx, err := contract.Transact(auth, "deposit", amount)
	if err != nil {
		log.Printf("Failed to deposit: %v", err)
		return fmt.Errorf("failed to deposit: %v", err)
	}

	fmt.Printf("Deposit transaction sent: %s\n", tx.Hash().Hex())
	return nil
}

func (es *EthereumService) Transfer(to, amount string) error {
	if err := utils.CheckAddress(to); err != nil {
		log.Printf("Failed to check address: %v", err)
		return fmt.Errorf("failed to check address: %v", err)
	}

	auth, err := es.createAuth()
	if err != nil {
		return err
	}

	contractABI, err := os.ReadFile("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
		return fmt.Errorf("failed to read ABI file: %v", err)
	}
	parseABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(es.ContractAddress, parseABI, es.Client, es.Client, es.Client)
	tx, err := contract.Transact(auth, "transfer", to, amount)
	if err != nil {
		log.Printf("Failed to transfer: %v", err)
		return fmt.Errorf("failed to transfer: %v", err)
	}

	fmt.Printf("Transfer transaction sent: %s\n", tx.Hash().Hex())
	return nil
}
