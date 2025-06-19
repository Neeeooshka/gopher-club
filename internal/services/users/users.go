package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"net/http"
	"time"
)

type UserRepository interface {
	AddUser(context.Context, models.User, string) error
	GetUserByLogin(string) (models.User, error)
}

type UserService struct {
	Errors  []error
	Inited  bool
	storage UserRepository
}

func NewUserService(ur interface{}) UserService {

	var us UserService

	userRepo, ok := ur.(UserRepository)

	if !ok {
		us.Errors = append(us.Errors, fmt.Errorf("2th argument expected UserRepository, got %T", ur))
	}

	if len(us.Errors) > 0 {
		return us
	}

	us.storage = userRepo
	us.Inited = true

	return us
}

func (u *UserService) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var cr credentials

	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !cr.validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, salt, err := cr.createPassword()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	user := models.User{
		Login:    cr.Login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err = u.storage.AddUser(ctx, user, salt)
	var ce *storage.ConflictUserError
	if err != nil {
		if errors.As(err, &ce) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := u.Authorize(cr)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (u *UserService) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var cr credentials

	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !cr.validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := u.Authorize(cr)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (u *UserService) Authenticate(jwtToken string) (models.User, error) {

	var user models.User

	login, err := VerifyJWTToken(jwtToken)
	if err != nil {
		return user, err
	}

	user, err = u.storage.GetUserByLogin(login)
	if err != nil {
		return user, fmt.Errorf("error authentication: %w", err)
	}

	return user, nil
}

func (u *UserService) Authorize(cr credentials) (string, error) {

	user, err := u.storage.GetUserByLogin(cr.Login)
	if err != nil {
		return "", fmt.Errorf("error authorization: %w", err)
	}

	err = cr.verifyPassword(user)
	if err != nil {
		return "", fmt.Errorf("error authorization: %w", err)
	}

	return CreateJWTToken(user.Login)
}

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
