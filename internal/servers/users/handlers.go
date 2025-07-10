package users

import (
	"encoding/json"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/dto"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
)

func (u *UserServer) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var auth dto.AuthData

	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{
		Login:    auth.Login,
		Password: auth.Password,
	}

	err := u.storage.AddUser(r.Context(), user)
	var ce *storage.ConflictUserError
	if err != nil {
		if errors.As(err, &ce) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := u.Authorize(auth)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (u *UserServer) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var auth dto.AuthData

	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := u.Authorize(auth)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
