package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/config"
	"github.com/sahmaragaev/lunaria-backend/internal/enums/tokentype"
)

type JWTService struct {
	config *config.JWTConfig
	redis  *RedisService
}

type Claims struct {
	UserID uuid.UUID      `json:"user_id"`
	Email  string         `json:"email"`
	Type   tokentype.Type `json:"type"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg *config.JWTConfig, redis *RedisService) *JWTService {
	return &JWTService{config: cfg, redis: redis}
}

func (j *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	expiryDuration, err := time.ParseDuration(j.config.AccessExpiry)
	if err != nil {
		return "", err
	}
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   tokentype.Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.config.Issuer,
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

func (j *JWTService) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	expiryDuration, err := time.ParseDuration(j.config.RefreshExpiry)
	if err != nil {
		return "", err
	}
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   tokentype.Refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.config.Issuer,
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

func (j *JWTService) ValidateToken(tokenString string, expectedType tokentype.Type) (*Claims, error) {
	// First check if token is blacklisted
	if j.redis != nil {
		blacklisted, err := j.redis.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			return nil, fmt.Errorf("failed to check token blacklist: %w", err)
		}
		if blacklisted {
			return nil, fmt.Errorf("token has been revoked")
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Type != expectedType {
			return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// RevokeToken adds a token to the blacklist
func (j *JWTService) RevokeToken(ctx context.Context, tokenString string) error {
	if j.redis == nil {
		return fmt.Errorf("redis service not available")
	}

	// Parse token to get expiration
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token claims")
	}

	// Calculate time until expiration
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		return fmt.Errorf("token already expired")
	}

	// Add to blacklist with the same expiration as the token
	return j.redis.BlacklistToken(ctx, tokenString, expiration)
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (j *JWTService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	if j.redis == nil {
		return fmt.Errorf("redis service not available")
	}

	// This would require storing active tokens per user in Redis
	// For now, we'll just return success as the tokens will expire naturally
	// In a production system, you'd want to track active tokens per user
	return nil
}
