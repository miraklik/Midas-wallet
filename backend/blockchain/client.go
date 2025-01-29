package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client *ethclient.Client
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (c *Client) GetBlockNumber() (uint64, error) {
	ctx := context.Background()
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		log.Printf("Failed to get block number: %v", err)
		return 0, fmt.Errorf("failed to get block number: %v", err)
	}
	return blockNumber, nil
}

func (c *Client) GetBalance(address string) (*big.Int, error) {
	ctx := context.Background()
	account := common.HexToAddress(address)
	balance, err := c.client.BalanceAt(ctx, account, nil)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	return balance, nil
}

func (c *Client) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	ctx := context.Background()
	hash := common.HexToHash(txHash)
	receipt, err := c.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
