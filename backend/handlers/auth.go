package handlers

import (
	"crypto-wallet/api"
	"crypto-wallet/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server { return &Server{db} }

func (s *Server) Register(c *gin.Context) {
	var Input RegisterRequest

	if err := c.ShouldBindJSON(&Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := api.GenerateAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := db.User{Address: address, Password: Input.Password}
	user.HashedPassword()

	if user.Password == "" || Input.Password != Input.PasswordAgain {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (s *Server) Login(c *gin.Context) {
	var Input LoginRequest

	if err := c.ShouldBindJSON(&Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.LoginCheck(Input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

func (s *Server) LoginCheck(password string) error {
	var err error

	user := db.User{}

	if err = s.db.Model(db.User{}).Where("password=?", password).Take(&user).Error; err != nil {
		return err
	}

	err = db.VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	return nil
}
