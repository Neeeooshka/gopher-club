package dto

// AuthData accepts authorization data in requests
type AuthData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
