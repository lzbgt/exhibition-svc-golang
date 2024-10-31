// Author: Bruce Lu
// Email: lzbgt_AT_icloud.com

package main

import (
	"time"

	"go-http-svc/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("0xfabefaffeaab")

// Generate JWT Token
func GenerateToken(username string, uid, eid int) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(4 * 7 * 24 * time.Hour)
	claims := &models.Claims{
		Username: username,
		UserId:   uid,
		Eid:      eid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate JWT Token
func ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
