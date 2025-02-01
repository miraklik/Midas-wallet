package main

import (
	"crypto-wallet/config"
	"crypto-wallet/db"
	"crypto-wallet/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	r := gin.Default()

	router := r.Group("/auth")

	server := handlers.NewServer(db)

	router.POST("/register", server.Register)
	router.POST("/login", server.Login)

	if err := r.Run(cfg.PORT); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
