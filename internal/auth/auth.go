package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type User struct {
	ID       int      `db:"ID"`
	Login    string   `db:"login"  json:"login"`
	Name     string   `db:"name"`
	Password Password `json:"password"`
}

type Password struct {
	hash   [32]byte
	cipher *Cipher
}

type Cipher struct {
	gsm cipher.AEAD
	key [32]byte
}

func generateRandom(size int) (string, error) {

	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func NewCipher(key []byte) (*Cipher, error) {

	var k [32]byte

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	copy(k[:], key)

	return &Cipher{aesgcm, k}, nil
}

func (c *Cipher) GetKey() string {
	return string(c.key[:])
}

func (c *Cipher) GenerateToken() (string, error) {

	salt, err := generateRandom(16)
	if err != nil {
		return "", err
	}

	nonce := c.key[7 : 7+c.gsm.NonceSize()]
	token := c.gsm.Seal(nil, nonce, []byte(salt), nil)

	return hex.EncodeToString(token), nil
}

func (c *Cipher) GetSalt(token string) (string, error) {

	ciphertext, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}

	nonce := c.key[7 : 7+c.gsm.NonceSize()]
	salt, err := c.gsm.Open(nil, nonce, ciphertext, nil)

	return string(salt), err
}

func CreatePassword(password string) (*Password, error) {

	random, err := generateRandom(32)
	if err != nil {
		return nil, err
	}

	key := sha256.Sum256([]byte(random))

	gsm, err := NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	token, err := gsm.GenerateToken()
	if err != nil {
		return nil, err
	}

	salt, _ := gsm.GetSalt(token)

	return &Password{cipher: gsm, hash: sha256.Sum256([]byte(password + salt))}, err
}

func (p *Password) GetHash() string {
	return hex.EncodeToString(p.hash[:])
}

func (p *Password) Verify(password, token string) bool {

	salt, err := p.cipher.GetSalt(token)
	if err != nil {
		return false
	}
	return sha256.Sum256([]byte(password+salt)) == p.hash
}

func NewPassword(hash string, key string) (*Password, error) {

	var h [32]byte
	copy(h[:], []byte(hash))

	cipher, err := NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &Password{cipher: cipher, hash: h}, nil
}
