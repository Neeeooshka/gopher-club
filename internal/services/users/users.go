package users

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserRepository interface {
	AddUser(context.Context, User, string) error
	GetUserByLogin(string) (User, error)
}

type UserService struct {
	Errors  []error
	Inited  bool
	User    User
	storage UserRepository
	ctx     context.Context
}

func NewUserService(ctx context.Context, ur interface{}) UserService {

	var us UserService

	userRepo, ok := ur.(UserRepository)

	if !ok {
		us.Errors = append(us.Errors, fmt.Errorf("2th argument expected UserRepository, got %T", ur))
	}

	if len(us.Errors) > 0 {
		return us
	}

	us.ctx = ctx
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

	user := User{
		Login:    cr.Login,
		Password: password,
	}

	err = u.storage.AddUser(u.ctx, user, salt)
	var ce *ConflictUserError
	if err != nil {
		if errors.As(err, &ce) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u.LoginUserHandler(w, r)
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

func (u *UserService) Authenticate(jwtToken string) error {

	login, err := VerifyJWTToken(jwtToken)
	if err != nil {
		return err
	}

	u.User, err = u.storage.GetUserByLogin(login)
	if err != nil {
		return err
	}

	return nil
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

	u.User = user

	return CreateJWTToken(u.User.Login)
}

type User struct {
	ID          int     `db:"ID"`
	Login       string  `db:"login"`
	Password    string  `db:"password"`
	Credentials string  `db:"credentials"`
	Balance     float64 `db:"balance"`
}

type ConflictUserError struct {
	ID    int
	login string
}

func NewConflictUserError(ID int, login string) *ConflictUserError {
	return &ConflictUserError{ID, login}
}

func (e *ConflictUserError) Error() string {
	return "User with login " + e.login + " already exsists"
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

	return string(hash[:]), token, nil
}

func (cr *credentials) verifyPassword(user User) error {

	var hash [32]byte
	copy(hash[:], user.Password)

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
