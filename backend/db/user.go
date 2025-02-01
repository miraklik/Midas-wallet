package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Address  string `json:"address"`
	Password string `json:"password"`
}

func CreateUser(db *gorm.DB, user *User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) HashedPassword() error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
		return err
	}
	u.Password = string(hashPass)

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
