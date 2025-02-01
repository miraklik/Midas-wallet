package main

import (
	"crypto-wallet/config"
	"crypto-wallet/db"
	"crypto-wallet/handlers"
	"crypto-wallet/services"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
				▄              ▄
                ▌▒█           ▄▀▒▌
                ▌▒▒█        ▄▀▒▒▒▐
               ▐▄█▒▒▀▀▀▀▄▄▄▀▒▒▒▒▒▐
             ▄▄▀▒▒▒▒▒▒▒▒▒▒▒█▒▒▄█▒▐
           ▄▀▒▒▒░░░▒▒▒░░░▒▒▒▀██▀▒▌
          ▐▒▒▒▄▄▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▀▄▒▌
          ▌░░▌█▀▒▒▒▒▒▄▀█▄▒▒▒▒▒▒▒█▒▐
         ▐░░░▒▒▒▒▒▒▒▒▌██▀▒▒░░░▒▒▒▀▄▌
         ▌░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░▒▒▒▒▌
        ▌▒▒▒▄██▄▒▒▒▒▒▒▒▒░░░░░░░░▒▒▒▐
        ▐▒▒▐▄█▄█▌▒▒▒▒▒▒▒▒▒▒░▒░▒░▒▒▒▒▌
        ▐▒▒▐▀▐▀▒▒▒▒▒▒▒▒▒▒▒▒▒░▒░▒░▒▒▐
         ▌▒▒▀▄▄▄▄▄▄▀▒▒▒▒▒▒▒░▒░▒░▒▒▒▌
         ▐▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░▒░▒▒▄▒▒▐
          ▀▄▒▒▒▒▒▒▒▒▒▒▒▒▒░▒░▒▄▒▒▒▒▌
            ▀▄▒▒▒▒▒▒▒▒▒▒▄▄▄▀▒▒▒▒▄▀
              ▀▄▄▄▄▄▄▀▀▀▒▒▒▒▒▄▄▀
                 ▀▀▀▀▀▀▀▀▀▀▀▀

*/

func initDB() *gorm.DB {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}

func main() {
	db := initDB()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := ethclient.Dial(cfg.RPC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(cfg.SK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	ethService := &services.EthereumService{
		Client:          client,
		ContractAddress: common.HexToAddress(cfg.CONTRACT_ADDRESS),
		PrivateKey:      privateKey,
		Contract:        nil,
	}

	router := gin.Default()

	s := handlers.NewServers(db)

	router.POST("/withdraw", s.Withdraw(ethService))
	router.POST("/transfer", s.Transfer(ethService))
	router.GET("/balance", handlers.GetBalance)
	router.GET("/transaction-status", handlers.GetTransactionStatusHandler)

	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
