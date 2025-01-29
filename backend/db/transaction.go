package db

import (
	"context"
	"crypto-wallet/config"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gorm.io/gorm"
)

type Transaction struct {
	TxHash      string `gorm:"primaryKey"`
	FromAddress string
	ToAddress   string
	Amount      string
	Timestamp   *timestamp.Timestamp
}

func NewTransaction(txHash, fromAddress, toAddress, amount string) *Transaction {
	return &Transaction{
		TxHash:      txHash,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
		Timestamp: &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
}

func SynchronizeTransactions(db *gorm.DB) {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	client, err := ethclient.Dial(cfg.RPC_URL)
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("Subscription error:", err)
		return
	}

	for {
		select {
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("Error getting block: %v", err)
				continue
			}

			for _, tx := range block.Transactions() {
				from, _ := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)

				transaction := NewTransaction(
					tx.Hash().Hex(),
					from.Hex(),
					tx.To().Hex(),
					tx.Value().String(),
				)

				if err := db.Create(transaction).Error; err != nil {
					log.Printf("DB insert error: %v", err)
				}
			}

		case err := <-sub.Err():
			log.Fatal("Subscription terminated:", err)
			return
		}
	}
}
