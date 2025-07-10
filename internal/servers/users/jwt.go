package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtKey = "aB3vC6dF9gJ2kM5nQ8rS1uV4xZ7yT0wE4hH7jK9lL0pO7iU"

var JWTLiveTime = time.Hour * 720

func CreateJWTToken(login string) (string, error) {

	expirationTime := time.Now().Add(JWTLiveTime)

	claims := jwt.RegisteredClaims{
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	JWTToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return JWTToken, nil
}

// VerifyJWTToken return login, error
func VerifyJWTToken(JWTToken string) (string, error) {

	if JWTToken != "" {
		token, err := jwt.ParseWithClaims(JWTToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil {
			return "", err
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			return claims.Subject, nil
		}
	}

	return "", fmt.Errorf("invalid token or token expired")
}
