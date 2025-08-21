package postgres

import (
	"testing"

	"github.com/sahmaragaev/lunaria-backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	cfg := config.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test_user",
		Password: "test_pass",
		DBName:   "test_db",
		SSLMode:  "disable",
	}
	db, err := NewPostgresConnection(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	if db != nil {
		_ = db.Close()
	}
}
