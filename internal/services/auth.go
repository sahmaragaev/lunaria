package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
)

type AuthService struct {
	userRepo        *repositories.UserRepository
	jwtService      *JWTService
	passwordService *PasswordService
	validator       *validator.Validate
}

func NewAuthService(userRepo *repositories.UserRepository, jwtService *JWTService, passwordService *PasswordService) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		validator:       validator.New(),
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Age:          req.Age,
		IsActive:     true,
	}
	if req.Gender != "" {
		user.Gender = &req.Gender
	}
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	accessToken, err := s.jwtService.GenerateAccessToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiryDuration, _ := time.ParseDuration(s.jwtService.config.AccessExpiry)
	return &dto.AuthResponse{
		User:         createdUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expiryDuration.Seconds()),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err := s.passwordService.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiryDuration, _ := time.ParseDuration(s.jwtService.config.AccessExpiry)
	return &dto.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expiryDuration.Seconds()),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := s.jwtService.ValidateToken(req.RefreshToken, "refresh")
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiryDuration, _ := time.ParseDuration(s.jwtService.config.AccessExpiry)
	return &dto.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(expiryDuration.Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Revoke the access token
	if err := s.jwtService.RevokeToken(ctx, accessToken); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}
