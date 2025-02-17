package services

import (
	"context"
	"crypto-wallet/utils"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
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

func (es *EthereumService) loadContractr(PathABI string) (*bind.BoundContract, error) {
	contractABI, err := os.ReadFile(PathABI)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	parseABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(es.ContractAddress, parseABI, es.Client, es.Client, es.Client)
	return contract, nil
}

func (es *EthereumService) createAuth(gasLimit uint64, gasPrice *big.Int) (*bind.TransactOpts, error) {
	chainID, err := es.Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get Chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(es.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %v", err)
	}

	if gasLimit > 0 {
		auth.GasLimit = gasLimit
	}

	if gasPrice != nil {
		auth.GasPrice = gasPrice
	} else {
		suggestedGasPrice, err := es.Client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %v", err)
		}
		auth.GasPrice = suggestedGasPrice
	}

	return auth, nil
}

func (es *EthereumService) Withdraw(to, amount string) error {
	if err := utils.CheckAddress(to); err != nil {
		log.Printf("Failed to check address: %v", err)
		return fmt.Errorf("failed to check address: %v", err)
	}

	auth, err := es.createAuth(21000, nil)
	if err != nil {
		return err
	}

	contract, err := es.loadContractr("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
	}

	tx, err := contract.Transact(auth, "withdraw", to, amount)
	if err != nil {
		log.Printf("Failed to withdraw: %v", err)
		return fmt.Errorf("failed to withdraw: %v", err)
	}

	fmt.Printf("Withdrawal transaction sent: %s\n", tx.Hash().Hex())
	return nil
}

func (es *EthereumService) Deposit(amount string) error {
	auth, err := es.createAuth(21000, nil)
	if err != nil {
		return err
	}

	contract, err := es.loadContractr("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
	}

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

	auth, err := es.createAuth(21000, nil)
	if err != nil {
		return err
	}

	contract, err := es.loadContractr("./build/walletABI.json")
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
	}

	tx, err := contract.Transact(auth, "transfer", to, amount)
	if err != nil {
		log.Printf("Failed to transfer: %v", err)
		return fmt.Errorf("failed to transfer: %v", err)
	}

	fmt.Printf("Transfer transaction sent: %s\n", tx.Hash().Hex())
	return nil
}
