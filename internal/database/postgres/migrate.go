package postgres

import (
	"context"
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Create tables first
	createTables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			age INTEGER,
			gender VARCHAR(50),
			avatar_url TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// User preferences table
		`CREATE TABLE IF NOT EXISTS user_preferences (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			preferred_companion_age INTEGER,
			preferred_gender VARCHAR(50),
			notification_settings JSONB,
			privacy_settings JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Companions table
		`CREATE TABLE IF NOT EXISTS companions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			gender VARCHAR(50) NOT NULL,
			age INTEGER NOT NULL,
			avatar_url TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Companion relationships table
		`CREATE TABLE IF NOT EXISTS companion_relationships (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			companion_id UUID NOT NULL REFERENCES companions(id) ON DELETE CASCADE,
			relationship_stage VARCHAR(100) DEFAULT 'acquaintance',
			intimacy_level INTEGER DEFAULT 0,
			message_count INTEGER DEFAULT 0,
			last_interaction_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			relationship_started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Conversations table (PostgreSQL version for analytics/summary data)
		`CREATE TABLE IF NOT EXISTS conversations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			companion_id VARCHAR(255) NOT NULL,
			message_count INTEGER DEFAULT 0,
			last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			intimacy_level INTEGER DEFAULT 0,
			relationship_stage VARCHAR(100),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Messages table (PostgreSQL version for analytics)
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
			sender_id UUID NOT NULL,
			type VARCHAR(50) NOT NULL,
			sentiment VARCHAR(50),
			tokens INTEGER,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Media files table
		`CREATE TABLE IF NOT EXISTS media_files (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			s3_url TEXT NOT NULL,
			format VARCHAR(50),
			size BIGINT,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// User engagement analytics table
		`CREATE TABLE IF NOT EXISTS user_engagement_analytics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR(255) NOT NULL,
			companion_id VARCHAR(255) NOT NULL,
			conversation_id VARCHAR(255) NOT NULL,
			session_duration INTERVAL,
			messages_per_session INTEGER DEFAULT 0,
			response_time INTERVAL,
			engagement_score DECIMAL(5,2) DEFAULT 0.0,
			conversation_depth DECIMAL(5,2) DEFAULT 0.0,
			emotional_intensity DECIMAL(5,2) DEFAULT 0.0,
			topic_diversity DECIMAL(5,2) DEFAULT 0.0,
			vulnerability_level DECIMAL(5,2) DEFAULT 0.0,
			peak_activity_time TIMESTAMP WITH TIME ZONE,
			session_frequency INTEGER DEFAULT 0,
			preferred_topics JSONB,
			interaction_style VARCHAR(100),
			intimacy_growth DECIMAL(5,2) DEFAULT 0.0,
			trust_building DECIMAL(5,2) DEFAULT 0.0,
			relationship_stage VARCHAR(100),
			milestone_progress JSONB,
			sentiment_trend JSONB,
			emotional_regulation DECIMAL(5,2) DEFAULT 0.0,
			empathy_response DECIMAL(5,2) DEFAULT 0.0,
			mood_impact DECIMAL(5,2) DEFAULT 0.0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	// Create tables
	for _, stmt := range createTables {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			log.Printf("Failed to create table: %v", err)
			return err
		}
	}

	// Create indexes after tables exist
	createIndexes := []string{
		// Conversations table indexes
		`CREATE INDEX IF NOT EXISTS idx_conversations_user_companion ON conversations(user_id, companion_id, last_activity DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at DESC);`,

		// Messages table indexes
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_sender_type ON messages(sender_id, type);`,

		// User engagement analytics indexes
		`CREATE INDEX IF NOT EXISTS idx_analytics_user_companion_conversation_created ON user_engagement_analytics(user_id, companion_id, conversation_id, created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_engagement_score ON user_engagement_analytics(engagement_score DESC);`,

		// Users table indexes
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);`,

		// User preferences indexes
		`CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);`,

		// Companions table indexes
		`CREATE INDEX IF NOT EXISTS idx_companions_user_id ON companions(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_companions_created_at ON companions(created_at DESC);`,

		// Companion relationships indexes
		`CREATE INDEX IF NOT EXISTS idx_companion_relationships_user_companion ON companion_relationships(user_id, companion_id);`,
		`CREATE INDEX IF NOT EXISTS idx_companion_relationships_last_interaction ON companion_relationships(last_interaction_at DESC);`,

		// Media files indexes
		`CREATE INDEX IF NOT EXISTS idx_media_files_user_id ON media_files(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_media_files_type_status ON media_files(type, status);`,
	}

	// Create indexes
	for _, stmt := range createIndexes {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			log.Printf("Failed to create index: %v", err)
			return err
		}
	}

	log.Println("Postgres migrations applied successfully.")
	return nil
}
