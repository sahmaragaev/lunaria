package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
)

type AuthMiddleware struct {
	jwtService *services.JWTService
	userRepo   *repositories.UserRepository
}

func NewAuthMiddleware(jwtService *services.JWTService, userRepo *repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, fmt.Errorf("authorization header required"), gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			response.Error(c, 401, fmt.Errorf("invalid authorization header format"), gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateToken(bearerToken[1], "access")
		if err != nil {
			response.Error(c, 401, err, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		user, err := m.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Error(c, 401, err, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.Next()
			return
		}
		claims, err := m.jwtService.ValidateToken(bearerToken[1], "access")
		if err != nil {
			c.Next()
			return
		}
		user, err := m.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.Next()
			return
		}
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
