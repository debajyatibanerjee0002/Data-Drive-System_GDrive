package utility

import (
	"go-drive-project/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	logger.Log.Println("Inside CheckPasswordHash")
	logger.Log.Println("Password: ", password)
	logger.Log.Println("Hash: ", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Returns true if the password matches
}
