package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testUser = models.User{
	Login:    "newuser",
	Password: "bf9760c303b7fbb093352d0e892c054c7b7a1db4fa26d690511cdd9602cdec5f",
}

func TestLoginUserHandler(t *testing.T) {

	repo := mocks.NewUserRepository(t)

	repo.On("GetUserByLogin", testUser.Login).Return(testUser, nil)
	repo.On("GetUserByLogin", "nonexistent").Return(models.User{}, errors.New("user not found"))

	service := NewUserServer(repo)

	tests := []struct {
		name           string
		login          string
		password       string
		expectedStatus int
	}{
		{
			name:           "Successful login",
			login:          "newuser",
			password:       "validpassword123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid password",
			login:          "newuser",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Non-existent user",
			login:          "nonexistent",
			password:       "anypassword",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty credentials",
			login:          "",
			password:       "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			creds := dto.AuthData{
				Login:    tt.login,
				Password: tt.password,
			}
			body, _ := json.Marshal(creds)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			service.LoginUserHandler(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				auth := rec.Header().Get("Authorization")
				assert.NotEmpty(t, auth)
			}
		})
	}
}

func TestRegisterUserHandler(t *testing.T) {

	ce := &storage.ConflictUserError{}

	repo := mocks.NewUserRepository(t)
	repo.On("AddUser", mock.Anything, mock.MatchedBy(func(u models.User) bool { return u.Login == testUser.Login }), mock.Anything).Return(nil)
	repo.On("AddUser", mock.Anything, mock.MatchedBy(func(u models.User) bool { return u.Login == "existinguser" }), mock.Anything).Return(ce)
	repo.On("GetUserByLogin", testUser.Login).Return(testUser, nil)

	svs := NewUserServer(repo)

	tests := []struct {
		name           string
		login          string
		password       string
		expectedStatus int
	}{
		{
			name:           "Successful registration",
			login:          "newuser",
			password:       "validpassword123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty credentials",
			login:          "",
			password:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "User already exists",
			login:          "existinguser",
			password:       "anypassword",
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			creds := dto.AuthData{
				Login:    tt.login,
				Password: tt.password,
			}
			body, _ := json.Marshal(creds)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			svs.RegisterUserHandler(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.NotEmpty(t, rec.Header().Get("Authorization"))
			}
		})
	}

	repo.AssertExpectations(t)
}
