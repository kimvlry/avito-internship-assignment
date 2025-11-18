package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var secret string

func main() {
	secret = "there-definitely-should-not-be-default-value-but-for-demonstration-simplicity-its-there"

	userID := flag.String("user", "admin1", "User ID")
	isAdmin := flag.Bool("admin", false, "Is admin user")
	hours := flag.Int("hours", 24, "Token validity in hours")
	secret := flag.String("secret", secret, "JWT secret")

	flag.Parse()

	token, expiresAt, err := generateJWT(*userID, *isAdmin, *hours, *secret)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("JWT Generated")
	fmt.Printf("User ID:    %s\n", *userID)
	fmt.Printf("Is Admin:   %v\n", *isAdmin)
	fmt.Printf("Expires:    %s\n", expiresAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(token)
	fmt.Println()
	fmt.Println("Use in requests:")
	fmt.Printf("Authorization: Bearer %s\n", token)
	fmt.Println()
}

func generateJWT(userID string, isAdmin bool, hours int, secret string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}
