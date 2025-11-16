package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/joho/godotenv"
)

// token-generator for demo use

func main() {
    userID := flag.String("user", "admin1", "User ID")
    isAdmin := flag.Bool("admin", false, "Is admin user")
    hours := flag.Int("hours", 24, "Token validity in hours")

    flag.Parse()

    _ = godotenv.Load()
    secret := os.Getenv("JWT_SECRET")

    token, expiresAt := generateJWT(*userID, *isAdmin, *hours, secret)

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

func generateJWT(userID string, isAdmin bool, hours int, secret string) (string, time.Time) {
    expiresAt := time.Now().Add(time.Duration(hours) * time.Hour)

    header := map[string]string{
        "alg": "HS256",
        "typ": "JWT",
    }
    headerJSON, _ := json.Marshal(header)
    headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

    payload := map[string]interface{}{
        "user_id":  userID,
        "is_admin": isAdmin,
        "exp":      expiresAt.Unix(),
    }
    payloadJSON, _ := json.Marshal(payload)
    payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

    message := headerB64 + "." + payloadB64
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(message))
    signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

    token := message + "." + signature
    return token, expiresAt
}
