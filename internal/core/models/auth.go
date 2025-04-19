package models

type RegistrationCredentials struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8"`
	Name         string `json:"name" validate:"required"`
	Role         string `json:"role" validate:"required"`
	Description  string `json:"description" validate:"required"`
	ActivityArea string `json:"activity_area" validate:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
