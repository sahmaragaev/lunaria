package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanionRepository struct {
	postgresDB *sql.DB
	mongoDB    *mongo.Database
}

func NewCompanionRepository(postgresDB *sql.DB, mongoDB *mongo.Database) *CompanionRepository {
	return &CompanionRepository{
		postgresDB: postgresDB,
		mongoDB:    mongoDB,
	}
}

func (r *CompanionRepository) Create(ctx context.Context, companion *models.Companion) (*models.Companion, error) {
	query := `
		INSERT INTO companions (id, user_id, name, gender, age, avatar_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	companion.ID = uuid.New()
	err := r.postgresDB.QueryRowContext(ctx, query,
		companion.ID, companion.UserID, companion.Name, companion.Gender,
		companion.Age, companion.AvatarURL, companion.IsActive).
		Scan(&companion.ID, &companion.CreatedAt, &companion.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create companion: %w", err)
	}
	return companion, nil
}

func (r *CompanionRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Companion, error) {
	companion := &models.Companion{}
	query := `
		SELECT id, user_id, name, gender, age, avatar_url, is_active, created_at, updated_at
		FROM companions 
		WHERE id = $1 AND user_id = $2 AND is_active = true`
	err := r.postgresDB.QueryRowContext(ctx, query, id, userID).Scan(
		&companion.ID, &companion.UserID, &companion.Name, &companion.Gender,
		&companion.Age, &companion.AvatarURL, &companion.IsActive,
		&companion.CreatedAt, &companion.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("companion not found")
		}
		return nil, fmt.Errorf("failed to get companion: %w", err)
	}
	return companion, nil
}

func (r *CompanionRepository) GetUserCompanions(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Companion, int, error) {
	offset := (page - 1) * pageSize
	countQuery := `SELECT COUNT(*) FROM companions WHERE user_id = $1 AND is_active = true`
	var total int
	err := r.postgresDB.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count companions: %w", err)
	}
	query := `
		SELECT id, user_id, name, gender, age, avatar_url, is_active, created_at, updated_at
		FROM companions 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.postgresDB.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get companions: %w", err)
	}
	defer rows.Close()
	var companions []models.Companion
	for rows.Next() {
		var companion models.Companion
		err := rows.Scan(
			&companion.ID, &companion.UserID, &companion.Name, &companion.Gender,
			&companion.Age, &companion.AvatarURL, &companion.IsActive,
			&companion.CreatedAt, &companion.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan companion: %w", err)
		}
		companions = append(companions, companion)
	}
	return companions, total, nil
}

func (r *CompanionRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]any) (*models.Companion, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id, userID)
	}
	setParts := []string{}
	args := []any{id, userID}
	argIndex := 3
	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}
	query := fmt.Sprintf(`
		UPDATE companions 
		SET %s, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_active = true
		RETURNING id, user_id, name, gender, age, avatar_url, is_active, created_at, updated_at`,
		strings.Join(setParts, ", "))
	companion := &models.Companion{}
	err := r.postgresDB.QueryRowContext(ctx, query, args...).Scan(
		&companion.ID, &companion.UserID, &companion.Name, &companion.Gender,
		&companion.Age, &companion.AvatarURL, &companion.IsActive,
		&companion.CreatedAt, &companion.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("companion not found")
		}
		return nil, fmt.Errorf("failed to update companion: %w", err)
	}
	return companion, nil
}

func (r *CompanionRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE companions SET is_active = false, updated_at = NOW() WHERE id = $1 AND user_id = $2`
	result, err := r.postgresDB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete companion: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("companion not found")
	}
	return nil
}

func (r *CompanionRepository) CreateProfile(ctx context.Context, profile *models.CompanionProfile) (*models.CompanionProfile, error) {
	collection := r.mongoDB.Collection("companion_profiles")
	profile.ID = primitive.NewObjectID()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create companion profile: %w", err)
	}
	return profile, nil
}

func (r *CompanionRepository) GetProfile(ctx context.Context, companionID string) (*models.CompanionProfile, error) {
	collection := r.mongoDB.Collection("companion_profiles")
	var profile models.CompanionProfile
	err := collection.FindOne(ctx, bson.M{"companion_id": companionID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("companion profile not found")
		}
		return nil, fmt.Errorf("failed to get companion profile: %w", err)
	}
	return &profile, nil
}

func (r *CompanionRepository) UpdateProfile(ctx context.Context, companionID string, updates bson.M) (*models.CompanionProfile, error) {
	collection := r.mongoDB.Collection("companion_profiles")
	updates["updated_at"] = time.Now()
	filter := bson.M{"companion_id": companionID}
	update := bson.M{"$set": updates}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update companion profile: %w", err)
	}
	return r.GetProfile(ctx, companionID)
}
