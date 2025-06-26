package users

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"net/http"
	"time"
)

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
