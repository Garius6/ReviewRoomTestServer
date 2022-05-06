package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const TOKEN_TIME = time.Minute * 15

var KEY = []byte(os.Getenv("KEY"))

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	Id       float64 `json:"id"`
	Username string  `json:"username"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Id       float64 `json:"id"`
	Username string  `json:"username"`
	jwt.StandardClaims
}

func NewRefreshClaims(user User) *RefreshClaims {
	return &RefreshClaims{
		Id:       user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
}

func NewUserClaims(user User) *UserClaims {
	return &UserClaims{
		Id:       user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKEN_TIME).Unix(),
		},
	}
}

func GenerateUserToken(user User) (string, error) {
	mySigningKey := KEY

	claims := NewUserClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateUserToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return nil
	} else {
		return errors.New("invalid claims")
	}
}

func GenerateRefreshToken(user User) (string, error) {
	mySigningKey := KEY

	claims := NewRefreshClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if rc, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return rc, nil
	} else {
		return nil, errors.New("invalid claims")
	}
}

func GenerateTokenPair(user User) (*TokenPair, error) {
	var tokenPair TokenPair
	accessToken, err := GenerateUserToken(user)
	if err != nil {
		return nil, err
	}
	tokenPair.AccessToken = accessToken

	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	tokenPair.RefreshToken = refreshToken
	return &tokenPair, nil
}
