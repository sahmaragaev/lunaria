package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	userRepo    *repositories.UserRepository
	validator   *validator.Validate
}

func NewAuthHandler(authService *services.AuthService, userRepo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}
	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "validation error") ||
			strings.Contains(errMsg, "password must") ||
			strings.Contains(errMsg, "already exists") {
			response.BadRequest(c, err, nil)
			return
		}
		response.InternalServerError(c, err, nil)
		return
	}

	response.Created(c, resp, "User registered successfully")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			response.Error(c, 401, err, nil)
			return
		}
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, resp, "Login successful")
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid refresh token") {
			response.Error(c, 401, err, nil)
			return
		}
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, resp, "Token refreshed successfully")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, 401, fmt.Errorf("authorization header required"), gin.H{"error": "Authorization header required"})
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		response.Error(c, 401, fmt.Errorf("invalid authorization header format"), gin.H{"error": "Invalid authorization header format"})
		return
	}

	// Revoke the token
	if err := h.authService.Logout(c.Request.Context(), bearerToken[1]); err != nil {
		response.InternalServerError(c, err, gin.H{"error": "Failed to logout"})
		return
	}

	response.Success(c, nil, "Logged out successfully")
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, fmt.Errorf("unauthorized"), gin.H{"error": "Unauthorized"})
		return
	}

	response.Success(c, user, "Profile retrieved successfully")
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}

	user := userInterface.(*models.User)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Validation error"})
		return
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	updatedUser, err := h.userRepo.UpdateProfile(c.Request.Context(), user.ID, updates)
	if err != nil {
		response.InternalServerError(c, err, gin.H{"error": "Failed to update profile"})
		return
	}
	response.Success(c, updatedUser, "Profile updated successfully")
}
