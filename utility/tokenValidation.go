package utility

import (
	"go-drive-project/entities"
	"go-drive-project/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func ValidateToken(token string) (*entities.Claims, bool, error, int) {
	var claims = &entities.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_KEY")), nil
		})
	if err != nil {
		logger.Log.Print("\n\n err:", err)
		if err == jwt.ErrSignatureInvalid {
			return claims, false, err, 1
		}
		return claims, false, err, 2
	}
	if !tkn.Valid {
		return claims, false, nil, 3
	}
	return claims, true, nil, 0
}

func GenerateJWT(username, email string, userId int) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"userid":   userId,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
		"iat":      time.Now().Unix(),                    // Issued at
	}

	// Create the token with HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
