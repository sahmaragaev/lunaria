package models

import (
	"time"

	"github.com/google/uuid"
)

type Companion struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Gender    string    `db:"gender" json:"gender"`
	Age       int       `db:"age" json:"age"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CompanionRelationship struct {
	ID                    uuid.UUID `db:"id" json:"id"`
	UserID                uuid.UUID `db:"user_id" json:"user_id"`
	CompanionID           uuid.UUID `db:"companion_id" json:"companion_id"`
	RelationshipStage     string    `db:"relationship_stage" json:"relationship_stage"`
	IntimacyLevel         int       `db:"intimacy_level" json:"intimacy_level"`
	MessageCount          int       `db:"message_count" json:"message_count"`
	LastInteractionAt     time.Time `db:"last_interaction_at" json:"last_interaction_at"`
	RelationshipStartedAt time.Time `db:"relationship_started_at" json:"relationship_started_at"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}
