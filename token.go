package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

var KEY = os.Getenv("KEY")

type UserClaims struct {
	Id       float64 `json:"id"`
	Username string  `json:"username"`
	jwt.StandardClaims
}

func NewUserClaims(user User) *UserClaims {
	return &UserClaims{
		Id:             user.Id,
		Username:       user.Username,
		StandardClaims: jwt.StandardClaims{},
	}
}

func generateToken(user User) (string, error) {
	signingKey := []byte(KEY)

	claims := NewUserClaims(user)
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
