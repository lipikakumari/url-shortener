package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("your-secret-key")

func HashPassword(password string) ([]byte, error) {

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Failed to hash password ", err)
		return nil, err

	}
	return hashPwd, err

}

func CheckPassword(password string, hashPassword []byte) error {

	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))

	return err

}

func GenerateToken(userID int32) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(30 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateToken(tokenString string) (int, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return secretKey, nil })

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userId := int(claims["userID"].(float64))
	return userId, nil

}
