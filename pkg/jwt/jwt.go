package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Create(secret string, expirationInSec int, userId int) (string, error) {
	expiration := time.Second * time.Duration(expirationInSec)
	claims := jwt.MapClaims{
		"userId":    strconv.Itoa(int(userId)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Validate(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func GetTokenUserId(tokenString, secret string) (int, error) {
	token, err := Validate(tokenString, secret)
	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, err := strconv.Atoi(claims["userId"].(string))
	if err != nil {
		return -1, err
	}

	return userId, nil
}
