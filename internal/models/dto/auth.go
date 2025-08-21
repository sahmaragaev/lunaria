package dto

import (
	"github.com/sahmaragaev/lunaria-backend/internal/models"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Age      *int   `json:"age,omitempty" validate:"omitempty,min=18,max=120"`
	Gender   string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Age       *int    `json:"age,omitempty" validate:"omitempty,min=18,max=120"`
	Gender    *string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}
