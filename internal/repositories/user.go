package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash, name, age, gender, avatar_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	user.ID = uuid.New()
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
		user.Age, user.Gender, user.AvatarURL, user.IsActive).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("user with email %s already exists", user.Email)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, gender, avatar_url, is_active, created_at, updated_at
		FROM users 
		WHERE email = $1 AND is_active = true`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.Age, &user.Gender, &user.AvatarURL, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, gender, avatar_url, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1 AND is_active = true`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.Age, &user.Gender, &user.AvatarURL, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, userID)
	}
	setParts := []string{}
	args := []any{userID}
	argIndex := 2
	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}
	query := fmt.Sprintf(`
		UPDATE users 
		SET %s, updated_at = NOW()
		WHERE id = $1 AND is_active = true
		RETURNING id, email, name, age, gender, avatar_url, is_active, created_at, updated_at`,
		strings.Join(setParts, ", "))
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Email, &user.Name,
		&user.Age, &user.Gender, &user.AvatarURL, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}
