package users

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const passKey = "supergophermarket"

const jwtKey = "aB3vC6dF9gJ2kM5nQ8rS1uV4xZ7yT0wE4hH7jK9lL0pO7iU"

var JWTLiveTime = time.Hour * 720

type Cipher struct {
	gsm cipher.AEAD
}

func NewCipher() (*Cipher, error) {

	key := sha256.Sum256([]byte(passKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	return &Cipher{aesgcm}, nil
}

func (c *Cipher) GenerateSaltToken() (string, error) {

	salt, err := generateRandom(16)
	if err != nil {
		return "", err
	}

	token := c.gsm.Seal(nil, c.getNonce(), []byte(salt), nil)

	return hex.EncodeToString(token), nil
}

func (c *Cipher) DecodeSalt(token string) (string, error) {

	ciphertext, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}

	salt, err := c.gsm.Open(nil, c.getNonce(), ciphertext, nil)

	return string(salt), err
}

func (c *Cipher) getNonce() []byte {
	key := sha256.Sum256([]byte(passKey))
	return key[7 : 7+c.gsm.NonceSize()]
}

func generateRandom(size int) (string, error) {

	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

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
