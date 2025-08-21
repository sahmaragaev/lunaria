package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
)

type RelationshipRepository struct {
	db *sql.DB
}

func NewRelationshipRepository(db *sql.DB) *RelationshipRepository {
	return &RelationshipRepository{db: db}
}

func (r *RelationshipRepository) Create(ctx context.Context, rel *models.CompanionRelationship) (*models.CompanionRelationship, error) {
	query := `
		INSERT INTO companion_relationships (id, user_id, companion_id, relationship_stage, intimacy_level, message_count, last_interaction_at, relationship_started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW(), NOW())
		RETURNING id, created_at, updated_at`
	rel.ID = uuid.New()
	err := r.db.QueryRowContext(ctx, query,
		rel.ID, rel.UserID, rel.CompanionID, rel.RelationshipStage, rel.IntimacyLevel, rel.MessageCount).
		Scan(&rel.ID, &rel.CreatedAt, &rel.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}
	return rel, nil
}

func (r *RelationshipRepository) GetByUserAndCompanion(ctx context.Context, userID, companionID uuid.UUID) (*models.CompanionRelationship, error) {
	rel := &models.CompanionRelationship{}
	query := `
		SELECT id, user_id, companion_id, relationship_stage, intimacy_level, message_count, last_interaction_at, relationship_started_at, created_at, updated_at
		FROM companion_relationships
		WHERE user_id = $1 AND companion_id = $2`
	err := r.db.QueryRowContext(ctx, query, userID, companionID).Scan(
		&rel.ID, &rel.UserID, &rel.CompanionID, &rel.RelationshipStage, &rel.IntimacyLevel, &rel.MessageCount,
		&rel.LastInteractionAt, &rel.RelationshipStartedAt, &rel.CreatedAt, &rel.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("relationship not found")
		}
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}
	return rel, nil
}
