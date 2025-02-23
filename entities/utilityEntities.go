package entities

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId int64 `json:"userid"`
	jwt.StandardClaims
}
