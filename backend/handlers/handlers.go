package handlers

import (
	"crypto-wallet/blockchain"
	"crypto-wallet/config"
	"crypto-wallet/db"
	"crypto-wallet/services"
	"crypto-wallet/utils"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Balance struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}

type TransactionStatusResponse struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

type DB_Server struct {
	db *gorm.DB
}

var (
	database, _ = db.ConnectDB()

	cfg, _ = config.Load()
)

func NewServers(db *gorm.DB) *DB_Server {
	return &DB_Server{db: db}
}

func GetBalance(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	client, err := blockchain.NewClient(cfg.RPC_URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := client.GetBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responce := Balance{
		Address: address,
		Balance: balance,
	}

	c.JSON(http.StatusOK, responce)
}

func GetTransactionStatusHandler(c *gin.Context) {
	txHash := c.Query("txHash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing txHash parameter"})
		return
	}

	client, err := blockchain.NewClient(cfg.RPC_URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Ethereum node"})
		return
	}

	receipt, err := client.GetTransactionReceipt(txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction receipt"})
		return
	}

	var status string
	if receipt == nil {
		status = "Pending"
	} else if receipt.Status == uint64(1) {
		status = "Confirmed"
	} else {
		status = "Failed"
	}

	response := TransactionStatusResponse{
		TxHash: txHash,
		Status: status,
	}

	c.JSON(http.StatusOK, response)
}

func (s *DB_Server) Withdraw(ethService *services.EthereumService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			To     string `json:"to"`
			Amount string `json:"amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.To == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to address is required"})
			return
		}

		if utils.CheckAddress(req.To) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to address"})
			return
		}

		err := ethService.Withdraw(req.To, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := database.Create(&db.Transaction{Amount: req.Amount, ToAddress: req.To}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "withdrawal successful"})
	}
}

func (s *DB_Server) Transfer(ethService *services.EthereumService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			To     string `json:"to"`
			Amount string `json:"amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.To == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to address is required"})
			return
		}

		if utils.CheckAddress(req.To) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to address"})
			return
		}

		if err := ethService.Transfer(req.To, req.Amount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := database.Create(&db.Transaction{Amount: req.Amount, ToAddress: req.To}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
	}
}
