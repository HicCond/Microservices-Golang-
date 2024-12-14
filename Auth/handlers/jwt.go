package handlers

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

var JwtSecret = []byte("zxckey")

// GenerateJWT generates a JWT token for the given username
func GenerateJWT(username string) (string, error) {
    // Create the Claims
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret key
    return token.SignedString(JwtSecret)
}
