package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnalyticsRepository struct {
	db    *sql.DB
	mongo *mongo.Database
}

func NewAnalyticsRepository(db *sql.DB, mongo *mongo.Database) *AnalyticsRepository {
	return &AnalyticsRepository{
		db:    db,
		mongo: mongo,
	}
}

// Existing PostgreSQL methods
func (r *AnalyticsRepository) UpsertConversationSummary(ctx context.Context, summary *models.ConversationSummary) error {
	query := `INSERT INTO conversation_summaries (id, user_id, companion_id, message_count, last_activity, intimacy_level, relationship_stage, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		ON CONFLICT (id) DO UPDATE SET message_count=$4, last_activity=$5, intimacy_level=$6, relationship_stage=$7, updated_at=NOW()`
	_, err := r.db.ExecContext(ctx, query, summary.ID, summary.UserID, summary.CompanionID, summary.MessageCount, summary.LastActivity, summary.IntimacyLevel, summary.RelationshipStage)
	return err
}

func (r *AnalyticsRepository) InsertMessageAnalytics(ctx context.Context, analytics *models.MessageAnalytics) error {
	query := `INSERT INTO message_analytics (id, conversation_id, sender_id, type, sentiment, tokens, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,NOW())`
	_, err := r.db.ExecContext(ctx, query, analytics.ID, analytics.ConversationID, analytics.SenderID, analytics.Type, analytics.Sentiment, analytics.Tokens)
	return err
}

func (r *AnalyticsRepository) InsertMediaFile(ctx context.Context, file *models.MediaFile) error {
	query := `INSERT INTO media_files (id, user_id, type, s3_url, format, size, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`
	_, err := r.db.ExecContext(ctx, query, file.ID, file.UserID, file.Type, file.S3URL, file.Format, file.Size, file.Status)
	return err
}

func (r *AnalyticsRepository) UpdateMediaFileStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE media_files SET status=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

// Enhanced Analytics Methods (MongoDB)

// User Engagement Analytics
func (r *AnalyticsRepository) UpsertUserEngagementAnalytics(ctx context.Context, analytics *models.UserEngagementAnalytics) error {
	collection := r.mongo.Collection("user_engagement_analytics")

	filter := bson.M{
		"user_id":         analytics.UserID,
		"companion_id":    analytics.CompanionID,
		"conversation_id": analytics.ConversationID,
	}

	update := bson.M{
		"$set": bson.M{
			"session_duration":     analytics.SessionDuration,
			"messages_per_session": analytics.MessagesPerSession,
			"response_time":        analytics.ResponseTime,
			"engagement_score":     analytics.EngagementScore,
			"conversation_depth":   analytics.ConversationDepth,
			"emotional_intensity":  analytics.EmotionalIntensity,
			"topic_diversity":      analytics.TopicDiversity,
			"vulnerability_level":  analytics.VulnerabilityLevel,
			"peak_activity_time":   analytics.PeakActivityTime,
			"session_frequency":    analytics.SessionFrequency,
			"preferred_topics":     analytics.PreferredTopics,
			"interaction_style":    analytics.InteractionStyle,
			"intimacy_growth":      analytics.IntimacyGrowth,
			"trust_building":       analytics.TrustBuilding,
			"relationship_stage":   analytics.RelationshipStage,
			"milestone_progress":   analytics.MilestoneProgress,
			"sentiment_trend":      analytics.SentimentTrend,
			"emotional_regulation": analytics.EmotionalRegulation,
			"empathy_response":     analytics.EmpathyResponse,
			"mood_impact":          analytics.MoodImpact,
			"updated_at":           time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":             primitive.NewObjectID(),
			"user_id":         analytics.UserID,
			"companion_id":    analytics.CompanionID,
			"conversation_id": analytics.ConversationID,
			"created_at":      time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *AnalyticsRepository) GetUserEngagementAnalytics(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID) (*models.UserEngagementAnalytics, error) {
	collection := r.mongo.Collection("user_engagement_analytics")

	filter := bson.M{
		"user_id":         userID,
		"companion_id":    companionID,
		"conversation_id": conversationID,
	}

	var analytics models.UserEngagementAnalytics
	err := collection.FindOne(ctx, filter).Decode(&analytics)
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

// Relationship Analytics
func (r *AnalyticsRepository) UpsertRelationshipAnalytics(ctx context.Context, analytics *models.RelationshipAnalytics) error {
	collection := r.mongo.Collection("relationship_analytics")

	filter := bson.M{
		"user_id":      analytics.UserID,
		"companion_id": analytics.CompanionID,
	}

	update := bson.M{
		"$set": bson.M{
			"current_stage":          analytics.CurrentStage,
			"stage_duration":         analytics.StageDuration,
			"progression_velocity":   analytics.ProgressionVelocity,
			"stage_history":          analytics.StageHistory,
			"intimacy_level":         analytics.IntimacyLevel,
			"intimacy_growth":        analytics.IntimacyGrowth,
			"intimacy_milestones":    analytics.IntimacyMilestones,
			"trust_level":            analytics.TrustLevel,
			"trust_building_events":  analytics.TrustBuildingEvents,
			"safety_score":           analytics.SafetyScore,
			"communication_style":    analytics.CommunicationStyle,
			"vulnerability_patterns": analytics.VulnerabilityPatterns,
			"conflict_resolution":    analytics.ConflictResolution,
			"health_score":           analytics.HealthScore,
			"red_flags":              analytics.RedFlags,
			"strengths":              analytics.Strengths,
			"updated_at":             time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      analytics.UserID,
			"companion_id": analytics.CompanionID,
			"created_at":   time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *AnalyticsRepository) GetRelationshipAnalytics(ctx context.Context, userID, companionID string) (*models.RelationshipAnalytics, error) {
	collection := r.mongo.Collection("relationship_analytics")

	filter := bson.M{
		"user_id":      userID,
		"companion_id": companionID,
	}

	var analytics models.RelationshipAnalytics
	err := collection.FindOne(ctx, filter).Decode(&analytics)
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

// Real-time Analytics
func (r *AnalyticsRepository) UpsertRealTimeMetrics(ctx context.Context, metrics *models.RealTimeMetrics) error {
	collection := r.mongo.Collection("real_time_metrics")

	filter := bson.M{
		"user_id":         metrics.UserID,
		"companion_id":    metrics.CompanionID,
		"conversation_id": metrics.ConversationID,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":                metrics.IsActive,
			"session_start_time":       metrics.SessionStartTime,
			"current_session_duration": metrics.CurrentSessionDuration,
			"messages_in_session":      metrics.MessagesInSession,
			"last_response_time":       metrics.LastResponseTime,
			"average_response_time":    metrics.AverageResponseTime,
			"response_quality":         metrics.ResponseQuality,
			"engagement_level":         metrics.EngagementLevel,
			"mood_indicator":           metrics.MoodIndicator,
			"emotional_state":          metrics.EmotionalState,
			"ai_response_time":         metrics.AIResponseTime,
			"error_rate":               metrics.ErrorRate,
			"service_health":           metrics.ServiceHealth,
			"timestamp":                time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":             primitive.NewObjectID(),
			"user_id":         metrics.UserID,
			"companion_id":    metrics.CompanionID,
			"conversation_id": metrics.ConversationID,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Gamification Methods

// User Progress
func (r *AnalyticsRepository) UpsertUserProgress(ctx context.Context, progress *models.UserProgress) error {
	collection := r.mongo.Collection("user_progress")

	filter := bson.M{
		"user_id":      progress.UserID,
		"companion_id": progress.CompanionID,
	}

	update := bson.M{
		"$set": bson.M{
			"total_experience":       progress.TotalExperience,
			"current_level":          progress.CurrentLevel,
			"level_progress":         progress.LevelProgress,
			"experience_to_next":     progress.ExperienceToNext,
			"relationship_stage":     progress.RelationshipStage,
			"stage_progress":         progress.StageProgress,
			"stage_milestones":       progress.StageMilestones,
			"current_streak":         progress.CurrentStreak,
			"longest_streak":         progress.LongestStreak,
			"streak_type":            progress.StreakType,
			"last_activity_date":     progress.LastActivityDate,
			"total_achievements":     progress.TotalAchievements,
			"rare_achievements":      progress.RareAchievements,
			"achievement_progress":   progress.AchievementProgress,
			"total_conversations":    progress.TotalConversations,
			"total_messages":         progress.TotalMessages,
			"total_time_spent":       progress.TotalTimeSpent,
			"average_session_length": progress.AverageSessionLength,
			"updated_at":             time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      progress.UserID,
			"companion_id": progress.CompanionID,
			"created_at":   time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *AnalyticsRepository) GetUserProgress(ctx context.Context, userID, companionID string) (*models.UserProgress, error) {
	collection := r.mongo.Collection("user_progress")

	filter := bson.M{
		"user_id":      userID,
		"companion_id": companionID,
	}

	var progress models.UserProgress
	err := collection.FindOne(ctx, filter).Decode(&progress)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

// User Achievements
func (r *AnalyticsRepository) InsertUserAchievement(ctx context.Context, achievement *models.UserAchievement) error {
	collection := r.mongo.Collection("user_achievements")

	achievement.ID = primitive.NewObjectID()
	achievement.EarnedAt = time.Now()

	_, err := collection.InsertOne(ctx, achievement)
	return err
}

func (r *AnalyticsRepository) GetUserAchievements(ctx context.Context, userID, companionID string, limit int) ([]models.UserAchievement, error) {
	collection := r.mongo.Collection("user_achievements")

	filter := bson.M{
		"user_id":      userID,
		"companion_id": companionID,
	}

	opts := options.Find().
		SetSort(bson.M{"earned_at": -1}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.UserAchievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AnalyticsRepository) CheckAchievementEarned(ctx context.Context, userID, companionID, achievementID string) (bool, error) {
	collection := r.mongo.Collection("user_achievements")

	filter := bson.M{
		"user_id":        userID,
		"companion_id":   companionID,
		"achievement_id": achievementID,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Achievement Definitions
func (r *AnalyticsRepository) GetAchievementDefinitions(ctx context.Context, category string) ([]models.AchievementDefinition, error) {
	collection := r.mongo.Collection("achievement_definitions")

	filter := bson.M{"active": true}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var definitions []models.AchievementDefinition
	if err = cursor.All(ctx, &definitions); err != nil {
		return nil, err
	}

	return definitions, nil
}

func (r *AnalyticsRepository) GetAchievementDefinition(ctx context.Context, achievementID string) (*models.AchievementDefinition, error) {
	collection := r.mongo.Collection("achievement_definitions")

	filter := bson.M{
		"id":     achievementID,
		"active": true,
	}

	var definition models.AchievementDefinition
	err := collection.FindOne(ctx, filter).Decode(&definition)
	if err != nil {
		return nil, err
	}

	return &definition, nil
}

// InsertAchievementDefinition inserts a new achievement definition
func (r *AnalyticsRepository) InsertAchievementDefinition(ctx context.Context, definition *models.AchievementDefinition) error {
	collection := r.mongo.Collection("achievement_definitions")

	_, err := collection.InsertOne(ctx, definition)
	return err
}

// GetMongoCollection returns a MongoDB collection by name
func (r *AnalyticsRepository) GetMongoCollection(collectionName string) *mongo.Collection {
	return r.mongo.Collection(collectionName)
}

// Predictive Analytics
func (r *AnalyticsRepository) UpsertUserBehaviorPrediction(ctx context.Context, prediction *models.UserBehaviorPrediction) error {
	collection := r.mongo.Collection("user_behavior_predictions")

	filter := bson.M{
		"user_id":      prediction.UserID,
		"companion_id": prediction.CompanionID,
	}

	update := bson.M{
		"$set": bson.M{
			"churn_risk":                  prediction.ChurnRisk,
			"churn_factors":               prediction.ChurnFactors,
			"retention_probability":       prediction.RetentionProbability,
			"next_activity_time":          prediction.NextActivityTime,
			"engagement_likelihood":       prediction.EngagementLikelihood,
			"optimal_engagement_time":     prediction.OptimalEngagementTime,
			"relationship_progression":    prediction.RelationshipProgression,
			"next_milestone":              prediction.NextMilestone,
			"milestone_probability":       prediction.MilestoneProbability,
			"feature_adoption_likelihood": prediction.FeatureAdoptionLikelihood,
			"recommended_features":        prediction.RecommendedFeatures,
			"support_needs_probability":   prediction.SupportNeedsProbability,
			"support_type":                prediction.SupportType,
			"prediction_date":             time.Now(),
			"confidence":                  prediction.Confidence,
		},
		"$setOnInsert": bson.M{
			"_id":          primitive.NewObjectID(),
			"user_id":      prediction.UserID,
			"companion_id": prediction.CompanionID,
			"created_at":   time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Analytics Queries and Aggregations

// Get engagement trends for a user
func (r *AnalyticsRepository) GetEngagementTrends(ctx context.Context, userID, companionID string, days int) ([]models.EngagementTrendPoint, error) {
	collection := r.mongo.Collection("user_engagement_analytics")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id":      userID,
				"companion_id": companionID,
				"created_at": bson.M{
					"$gte": time.Now().AddDate(0, 0, -days),
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},
				"engagement_score": bson.M{"$avg": "$engagement_score"},
				"session_count":    bson.M{"$sum": 1},
				"message_count":    bson.M{"$sum": "$messages_per_session"},
				"duration":         bson.M{"$avg": "$session_duration"},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var trends []models.EngagementTrendPoint
	for _, result := range results {
		date, _ := time.Parse("2006-01-02", result["_id"].(string))
		trend := models.EngagementTrendPoint{
			Date:            date,
			EngagementScore: result["engagement_score"].(float64),
			SessionCount:    int(result["session_count"].(int32)),
			MessageCount:    int(result["message_count"].(int32)),
			Duration:        time.Duration(result["duration"].(int64)),
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// Get user statistics
func (r *AnalyticsRepository) GetUserStatistics(ctx context.Context, userID, companionID string) (*models.UserStatistics, error) {
	collection := r.mongo.Collection("user_engagement_analytics")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id":      userID,
				"companion_id": companionID,
			},
		},
		{
			"$group": bson.M{
				"_id":                      nil,
				"total_sessions":           bson.M{"$sum": 1},
				"total_messages":           bson.M{"$sum": "$messages_per_session"},
				"avg_session_length":       bson.M{"$avg": "$session_duration"},
				"avg_messages_per_session": bson.M{"$avg": "$messages_per_session"},
				"engagement_score":         bson.M{"$avg": "$engagement_score"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.UserStatistics{}, nil
	}

	result := results[0]
	stats := &models.UserStatistics{
		TotalSessions:             int(result["total_sessions"].(int32)),
		TotalMessages:             int(result["total_messages"].(int32)),
		AverageSessionLength:      time.Duration(result["avg_session_length"].(int64)),
		AverageMessagesPerSession: result["avg_messages_per_session"].(float64),
		EngagementScore:           result["engagement_score"].(float64),
	}

	return stats, nil
}

// Get streak information
func (r *AnalyticsRepository) GetStreakInformation(ctx context.Context, userID, companionID string) (*models.StreakInformation, error) {
	progress, err := r.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, err
	}

	// Calculate next milestone (next multiple of 7 for weekly streaks)
	nextMilestone := ((progress.CurrentStreak / 7) + 1) * 7

	streakInfo := &models.StreakInformation{
		CurrentStreak:  progress.CurrentStreak,
		LongestStreak:  progress.LongestStreak,
		StreakType:     progress.StreakType,
		LastActivity:   progress.LastActivityDate,
		NextMilestone:  nextMilestone,
		StreakProgress: float64(progress.CurrentStreak%7) / 7.0,
	}

	return streakInfo, nil
}

// Get users at risk of churn
func (r *AnalyticsRepository) GetUsersAtChurnRisk(ctx context.Context, threshold float64) ([]models.UserBehaviorPrediction, error) {
	collection := r.mongo.Collection("user_behavior_predictions")

	filter := bson.M{
		"churn_risk": bson.M{"$gte": threshold},
		"prediction_date": bson.M{
			"$gte": time.Now().AddDate(0, 0, -7), // Predictions from last 7 days
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var predictions []models.UserBehaviorPrediction
	if err = cursor.All(ctx, &predictions); err != nil {
		return nil, err
	}

	return predictions, nil
}

// Get platform-wide analytics
func (r *AnalyticsRepository) GetPlatformAnalytics(ctx context.Context, days int) (map[string]any, error) {
	collection := r.mongo.Collection("user_engagement_analytics")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": time.Now().AddDate(0, 0, -days),
				},
			},
		},
		{
			"$group": bson.M{
				"_id":                  nil,
				"total_users":          bson.M{"$addToSet": "$user_id"},
				"total_sessions":       bson.M{"$sum": 1},
				"total_messages":       bson.M{"$sum": "$messages_per_session"},
				"avg_engagement":       bson.M{"$avg": "$engagement_score"},
				"avg_session_duration": bson.M{"$avg": "$session_duration"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]any{}, nil
	}

	result := results[0]
	analytics := map[string]any{
		"total_users":          len(result["total_users"].(primitive.A)),
		"total_sessions":       result["total_sessions"].(int32),
		"total_messages":       result["total_messages"].(int32),
		"avg_engagement":       result["avg_engagement"].(float64),
		"avg_session_duration": time.Duration(result["avg_session_duration"].(int64)),
	}

	return analytics, nil
}
