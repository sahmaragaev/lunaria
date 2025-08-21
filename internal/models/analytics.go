package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Existing models
type ConversationSummary struct {
	ID                uuid.UUID `db:"id" json:"id"`
	UserID            uuid.UUID `db:"user_id" json:"user_id"`
	CompanionID       uuid.UUID `db:"companion_id" json:"companion_id"`
	MessageCount      int       `db:"message_count" json:"message_count"`
	LastActivity      time.Time `db:"last_activity" json:"last_activity"`
	IntimacyLevel     int       `db:"intimacy_level" json:"intimacy_level"`
	RelationshipStage string    `db:"relationship_stage" json:"relationship_stage"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type MessageAnalytics struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ConversationID string    `db:"conversation_id" json:"conversation_id"`
	SenderID       uuid.UUID `db:"sender_id" json:"sender_id"`
	Type           string    `db:"type" json:"type"`
	Sentiment      *string   `db:"sentiment" json:"sentiment,omitempty"`
	Tokens         *int      `db:"tokens" json:"tokens,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type MediaFile struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Type      string    `db:"type" json:"type"`
	S3URL     string    `db:"s3_url" json:"s3_url"`
	Format    string    `db:"format" json:"format"`
	Size      int64     `db:"size" json:"size"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Enhanced Analytics Models

// UserEngagementAnalytics tracks comprehensive user engagement metrics
type UserEngagementAnalytics struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	CompanionID    string             `bson:"companion_id" json:"companion_id"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`

	// Real-time metrics
	SessionDuration    time.Duration `bson:"session_duration" json:"session_duration"`
	MessagesPerSession int           `bson:"messages_per_session" json:"messages_per_session"`
	ResponseTime       time.Duration `bson:"response_time" json:"response_time"`
	EngagementScore    float64       `bson:"engagement_score" json:"engagement_score"`

	// Conversation quality metrics
	ConversationDepth  float64 `bson:"conversation_depth" json:"conversation_depth"`
	EmotionalIntensity float64 `bson:"emotional_intensity" json:"emotional_intensity"`
	TopicDiversity     float64 `bson:"topic_diversity" json:"topic_diversity"`
	VulnerabilityLevel float64 `bson:"vulnerability_level" json:"vulnerability_level"`

	// Behavioral patterns
	PeakActivityTime time.Time `bson:"peak_activity_time" json:"peak_activity_time"`
	SessionFrequency int       `bson:"session_frequency" json:"session_frequency"`
	PreferredTopics  []string  `bson:"preferred_topics" json:"preferred_topics"`
	InteractionStyle string    `bson:"interaction_style" json:"interaction_style"`

	// Relationship progression
	IntimacyGrowth    float64            `bson:"intimacy_growth" json:"intimacy_growth"`
	TrustBuilding     float64            `bson:"trust_building" json:"trust_building"`
	RelationshipStage string             `bson:"relationship_stage" json:"relationship_stage"`
	MilestoneProgress map[string]float64 `bson:"milestone_progress" json:"milestone_progress"`

	// Emotional intelligence metrics
	SentimentTrend      []SentimentPoint `bson:"sentiment_trend" json:"sentiment_trend"`
	EmotionalRegulation float64          `bson:"emotional_regulation" json:"emotional_regulation"`
	EmpathyResponse     float64          `bson:"empathy_response" json:"empathy_response"`
	MoodImpact          float64          `bson:"mood_impact" json:"mood_impact"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// SentimentPoint represents a sentiment measurement at a specific time
type SentimentPoint struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Score     float64   `bson:"score" json:"score"`
	Intensity float64   `bson:"intensity" json:"intensity"`
	Dominant  string    `bson:"dominant" json:"dominant"`
}

// RelationshipAnalytics tracks relationship development over time
type RelationshipAnalytics struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	CompanionID string             `bson:"companion_id" json:"companion_id"`

	// Relationship progression
	CurrentStage        string            `bson:"current_stage" json:"current_stage"`
	StageDuration       time.Duration     `bson:"stage_duration" json:"stage_duration"`
	ProgressionVelocity float64           `bson:"progression_velocity" json:"progression_velocity"`
	StageHistory        []StageTransition `bson:"stage_history" json:"stage_history"`

	// Intimacy tracking
	IntimacyLevel      float64             `bson:"intimacy_level" json:"intimacy_level"`
	IntimacyGrowth     float64             `bson:"intimacy_growth" json:"intimacy_growth"`
	IntimacyMilestones []IntimacyMilestone `bson:"intimacy_milestones" json:"intimacy_milestones"`

	// Trust and safety
	TrustLevel          float64      `bson:"trust_level" json:"trust_level"`
	TrustBuildingEvents []TrustEvent `bson:"trust_building_events" json:"trust_building_events"`
	SafetyScore         float64      `bson:"safety_score" json:"safety_score"`

	// Communication patterns
	CommunicationStyle    string               `bson:"communication_style" json:"communication_style"`
	VulnerabilityPatterns []VulnerabilityEvent `bson:"vulnerability_patterns" json:"vulnerability_patterns"`
	ConflictResolution    float64              `bson:"conflict_resolution" json:"conflict_resolution"`

	// Relationship health
	HealthScore float64  `bson:"health_score" json:"health_score"`
	RedFlags    []string `bson:"red_flags" json:"red_flags"`
	Strengths   []string `bson:"strengths" json:"strengths"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// StageTransition represents a relationship stage change
type StageTransition struct {
	FromStage  string    `bson:"from_stage" json:"from_stage"`
	ToStage    string    `bson:"to_stage" json:"to_stage"`
	Trigger    string    `bson:"trigger" json:"trigger"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Confidence float64   `bson:"confidence" json:"confidence"`
}

// IntimacyMilestone represents a significant intimacy achievement
type IntimacyMilestone struct {
	Type        string    `bson:"type" json:"type"`
	Description string    `bson:"description" json:"description"`
	Level       float64   `bson:"level" json:"level"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	Context     string    `bson:"context" json:"context"`
}

// TrustEvent represents a trust-building moment
type TrustEvent struct {
	Type        string    `bson:"type" json:"type"`
	Description string    `bson:"description" json:"description"`
	Impact      float64   `bson:"impact" json:"impact"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	Context     string    `bson:"context" json:"context"`
}

// VulnerabilityEvent represents a moment of vulnerability
type VulnerabilityEvent struct {
	Type        string    `bson:"type" json:"type"`
	Description string    `bson:"description" json:"description"`
	Level       float64   `bson:"level" json:"level"`
	Response    string    `bson:"response" json:"response"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
}

// Gamification Models

// UserAchievement represents an achievement earned by a user
type UserAchievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	CompanionID     string             `bson:"companion_id" json:"companion_id"`
	AchievementID   string             `bson:"achievement_id" json:"achievement_id"`
	AchievementType string             `bson:"achievement_type" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	IconURL         string             `bson:"icon_url" json:"icon_url"`
	Points          int                `bson:"points" json:"points"`
	Rarity          string             `bson:"rarity" json:"rarity"` // common, rare, epic, legendary
	EarnedAt        time.Time          `bson:"earned_at" json:"earned_at"`
	Context         map[string]any     `bson:"context" json:"context"`
}

// UserProgress tracks user progression through gamification system
type UserProgress struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	CompanionID string             `bson:"companion_id" json:"companion_id"`

	// Experience and levels
	TotalExperience  int     `bson:"total_experience" json:"total_experience"`
	CurrentLevel     int     `bson:"current_level" json:"current_level"`
	LevelProgress    float64 `bson:"level_progress" json:"level_progress"`
	ExperienceToNext int     `bson:"experience_to_next" json:"experience_to_next"`

	// Relationship progression
	RelationshipStage string           `bson:"relationship_stage" json:"relationship_stage"`
	StageProgress     float64          `bson:"stage_progress" json:"stage_progress"`
	StageMilestones   []StageMilestone `bson:"stage_milestones" json:"stage_milestones"`

	// Streaks and habits
	CurrentStreak    int       `bson:"current_streak" json:"current_streak"`
	LongestStreak    int       `bson:"longest_streak" json:"longest_streak"`
	StreakType       string    `bson:"streak_type" json:"streak_type"`
	LastActivityDate time.Time `bson:"last_activity_date" json:"last_activity_date"`

	// Achievements
	TotalAchievements   int                `bson:"total_achievements" json:"total_achievements"`
	RareAchievements    int                `bson:"rare_achievements" json:"rare_achievements"`
	AchievementProgress map[string]float64 `bson:"achievement_progress" json:"achievement_progress"`

	// Statistics
	TotalConversations   int           `bson:"total_conversations" json:"total_conversations"`
	TotalMessages        int           `bson:"total_messages" json:"total_messages"`
	TotalTimeSpent       time.Duration `bson:"total_time_spent" json:"total_time_spent"`
	AverageSessionLength time.Duration `bson:"average_session_length" json:"average_session_length"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// StageMilestone represents a milestone within a relationship stage
type StageMilestone struct {
	ID          string     `bson:"id" json:"id"`
	Title       string     `bson:"title" json:"title"`
	Description string     `bson:"description" json:"description"`
	Completed   bool       `bson:"completed" json:"completed"`
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	Progress    float64    `bson:"progress" json:"progress"`
}

// AchievementDefinition defines available achievements
type AchievementDefinition struct {
	ID            string              `bson:"id" json:"id"`
	Title         string              `bson:"title" json:"title"`
	Description   string              `bson:"description" json:"description"`
	Category      string              `bson:"category" json:"category"`
	Type          string              `bson:"type" json:"type"`
	Points        int                 `bson:"points" json:"points"`
	Rarity        string              `bson:"rarity" json:"rarity"`
	IconURL       string              `bson:"icon_url" json:"icon_url"`
	Criteria      AchievementCriteria `bson:"criteria" json:"criteria"`
	Prerequisites []string            `bson:"prerequisites" json:"prerequisites"`
	Rewards       map[string]any      `bson:"rewards" json:"rewards"`
	Active        bool                `bson:"active" json:"active"`
	CreatedAt     time.Time           `bson:"created_at" json:"created_at"`
}

// AchievementCriteria defines what needs to be accomplished
type AchievementCriteria struct {
	Type        string         `bson:"type" json:"type"`
	Target      float64        `bson:"target" json:"target"`
	Timeframe   *time.Duration `bson:"timeframe,omitempty" json:"timeframe,omitempty"`
	Conditions  map[string]any `bson:"conditions" json:"conditions"`
	Measurement string         `bson:"measurement" json:"measurement"`
}

// Real-time Analytics Models

// RealTimeMetrics tracks live engagement metrics
type RealTimeMetrics struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	CompanionID    string             `bson:"companion_id" json:"companion_id"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`

	// Live engagement
	IsActive               bool          `bson:"is_active" json:"is_active"`
	SessionStartTime       time.Time     `bson:"session_start_time" json:"session_start_time"`
	CurrentSessionDuration time.Duration `bson:"current_session_duration" json:"current_session_duration"`
	MessagesInSession      int           `bson:"messages_in_session" json:"messages_in_session"`

	// Response metrics
	LastResponseTime    time.Time     `bson:"last_response_time" json:"last_response_time"`
	AverageResponseTime time.Duration `bson:"average_response_time" json:"average_response_time"`
	ResponseQuality     float64       `bson:"response_quality" json:"response_quality"`

	// Engagement indicators
	EngagementLevel float64 `bson:"engagement_level" json:"engagement_level"`
	MoodIndicator   string  `bson:"mood_indicator" json:"mood_indicator"`
	EmotionalState  string  `bson:"emotional_state" json:"emotional_state"`

	// System performance
	AIResponseTime time.Duration `bson:"ai_response_time" json:"ai_response_time"`
	ErrorRate      float64       `bson:"error_rate" json:"error_rate"`
	ServiceHealth  string        `bson:"service_health" json:"service_health"`

	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

// Predictive Analytics Models

// UserBehaviorPrediction predicts future user behavior
type UserBehaviorPrediction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	CompanionID string             `bson:"companion_id" json:"companion_id"`

	// Churn prediction
	ChurnRisk            float64  `bson:"churn_risk" json:"churn_risk"`
	ChurnFactors         []string `bson:"churn_factors" json:"churn_factors"`
	RetentionProbability float64  `bson:"retention_probability" json:"retention_probability"`

	// Engagement prediction
	NextActivityTime      *time.Time `bson:"next_activity_time,omitempty" json:"next_activity_time,omitempty"`
	EngagementLikelihood  float64    `bson:"engagement_likelihood" json:"engagement_likelihood"`
	OptimalEngagementTime time.Time  `bson:"optimal_engagement_time" json:"optimal_engagement_time"`

	// Relationship prediction
	RelationshipProgression float64 `bson:"relationship_progression" json:"relationship_progression"`
	NextMilestone           string  `bson:"next_milestone" json:"next_milestone"`
	MilestoneProbability    float64 `bson:"milestone_probability" json:"milestone_probability"`

	// Feature adoption
	FeatureAdoptionLikelihood map[string]float64 `bson:"feature_adoption_likelihood" json:"feature_adoption_likelihood"`
	RecommendedFeatures       []string           `bson:"recommended_features" json:"recommended_features"`

	// Support needs
	SupportNeedsProbability float64 `bson:"support_needs_probability" json:"support_needs_probability"`
	SupportType             string  `bson:"support_type" json:"support_type"`

	PredictionDate time.Time `bson:"prediction_date" json:"prediction_date"`
	Confidence     float64   `bson:"confidence" json:"confidence"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

// Analytics Dashboard Models

// UserDashboardData provides comprehensive dashboard data
type UserDashboardData struct {
	UserID      string `json:"user_id"`
	CompanionID string `json:"companion_id"`

	// Progress overview
	Progress           *UserProgress     `json:"progress"`
	RecentAchievements []UserAchievement `json:"recent_achievements"`

	// Relationship insights
	RelationshipAnalytics *RelationshipAnalytics `json:"relationship_analytics"`
	EngagementTrends      []EngagementTrendPoint `json:"engagement_trends"`

	// Recommendations
	Recommendations []Recommendation `json:"recommendations"`
	NextMilestones  []StageMilestone `json:"next_milestones"`

	// Statistics
	Statistics *UserStatistics    `json:"statistics"`
	StreakInfo *StreakInformation `json:"streak_info"`

	LastUpdated time.Time `json:"last_updated"`
}

// EngagementTrendPoint represents engagement over time
type EngagementTrendPoint struct {
	Date            time.Time     `json:"date"`
	EngagementScore float64       `json:"engagement_score"`
	SessionCount    int           `json:"session_count"`
	MessageCount    int           `json:"message_count"`
	Duration        time.Duration `json:"duration"`
}

// Recommendation provides personalized recommendations
type Recommendation struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    int            `json:"priority"`
	Confidence  float64        `json:"confidence"`
	Action      string         `json:"action"`
	Metadata    map[string]any `json:"metadata"`
}

// UserStatistics provides comprehensive user statistics
type UserStatistics struct {
	TotalSessions             int           `json:"total_sessions"`
	AverageSessionLength      time.Duration `json:"average_session_length"`
	TotalMessages             int           `json:"total_messages"`
	AverageMessagesPerSession float64       `json:"average_messages_per_session"`
	PeakActivityHour          int           `json:"peak_activity_hour"`
	MostActiveDay             string        `json:"most_active_day"`
	EngagementScore           float64       `json:"engagement_score"`
	RelationshipHealth        float64       `json:"relationship_health"`
}

// StreakInformation provides streak details
type StreakInformation struct {
	CurrentStreak  int       `json:"current_streak"`
	LongestStreak  int       `json:"longest_streak"`
	StreakType     string    `json:"streak_type"`
	LastActivity   time.Time `json:"last_activity"`
	NextMilestone  int       `json:"next_milestone"`
	StreakProgress float64   `json:"streak_progress"`
}
