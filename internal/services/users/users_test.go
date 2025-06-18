package users

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
	"github.com/Neeeooshka/gopher-club/internal/storage/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginUserHandler(t *testing.T) {

	ctx := context.Background()
	mockRepo := &mocks.MockRepository{
		Users: make(map[string]models.User),
	}
	service := NewUserService(ctx, mockRepo)

	testUser := models.User{
		Login:       "newuser",
		Password:    "bf9760c303b7fbb093352d0e892c054c7b7a1db4fa26d690511cdd9602cdec5f",
		Credentials: "f3dcb06e549ee9732d9f86579310ad297151343fff0089b2ad1ada4fd0aff8c6e4ef359306dda5a77d0028235cacf704",
	}
	mockRepo.Users[testUser.Login] = testUser

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

			creds := credentials{
				Login:    tt.login,
				Password: tt.password,
			}
			body, _ := json.Marshal(creds)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			service.LoginUserHandler(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK && rec.Header().Get("Authorization") == "" {
				t.Error("expected Authorization header")
			}
		})
	}
}

func TestRegisterUserHandler(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockRepository{
		Users: make(map[string]models.User),
	}
	service := NewUserService(ctx, mockRepo)

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

	mockRepo.Users["existinguser"] = models.User{Login: "existinguser"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			creds := credentials{
				Login:    tt.login,
				Password: tt.password,
			}
			body, _ := json.Marshal(creds)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			service.RegisterUserHandler(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				_, err := mockRepo.GetUserByLogin(tt.login)
				if err != nil {
					t.Errorf("user %s was not saved", tt.login)
				}
			}
		})
	}
}
