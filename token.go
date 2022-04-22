package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

var KEY = os.Getenv("KEY")

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewUserClaims(username string) *UserClaims {
	return &UserClaims{
		Username:       username,
		StandardClaims: jwt.StandardClaims{},
	}
}

func generateToken(username string) (string, error) {
	signingKey := []byte(KEY)

	claims := NewUserClaims(username)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(KEY), nil
	})

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return nil
	} else {
		return err
	}
}
