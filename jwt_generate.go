package main

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("supersecretkey")

func main() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "test_user",
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	fmt.Println("Generated JWT Token:", tokenString)
}
