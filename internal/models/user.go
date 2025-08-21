package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	Age          *int      `db:"age" json:"age,omitempty"`
	Gender       *string   `db:"gender" json:"gender,omitempty"`
	AvatarURL    *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserPreferences struct {
	ID                    uuid.UUID `db:"id" json:"id"`
	UserID                uuid.UUID `db:"user_id" json:"user_id"`
	PreferredCompanionAge *int      `db:"preferred_companion_age" json:"preferred_companion_age,omitempty"`
	PreferredGender       *string   `db:"preferred_gender" json:"preferred_gender,omitempty"`
	NotificationSettings  any       `db:"notification_settings" json:"notification_settings"`
	PrivacySettings       any       `db:"privacy_settings" json:"privacy_settings"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}
