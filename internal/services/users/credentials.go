package users

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

// credentials accepts authorization data in requests
type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (cr *credentials) validate() bool {
	return cr.Login != "" && cr.Password != ""
}

// createPassword return hashed salted password, hashed crypted salt and error
// salt generate random
func (cr *credentials) createPassword() (string, string, error) {

	gsm, err := NewCipher()
	if err != nil {
		return "", "", err
	}

	token, err := gsm.GenerateSaltToken()
	if err != nil {
		return "", "", err
	}

	salt, _ := gsm.DecodeSalt(token)
	hash := sha256.Sum256([]byte(cr.Password + salt))

	return hex.EncodeToString(hash[:]), token, nil
}

func (cr *credentials) verifyPassword(user models.User) error {

	pass, err := hex.DecodeString(user.Password)
	if err != nil {
		return fmt.Errorf("error verifying password: %w", err)
	}

	var hash [32]byte
	copy(hash[:], pass)

	gsm, err := NewCipher()
	if err != nil {
		return fmt.Errorf("error verifying password: %w", err)
	}

	salt, err := gsm.DecodeSalt(user.Credentials)
	if err != nil {
		return fmt.Errorf("error verifying password: %w", err)
	}

	if sha256.Sum256([]byte(cr.Password+salt)) != hash {
		return fmt.Errorf("error verifying password: %w", errors.New("password incorrect"))
	}

	return nil
}
