package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
