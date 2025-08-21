package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContextLayer represents different layers of conversation context
type ContextLayer struct {
	Type      string         `json:"type" bson:"type"`
	Content   string         `json:"content" bson:"content"`
	Priority  int            `json:"priority" bson:"priority"`
	Metadata  map[string]any `json:"metadata" bson:"metadata"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// ConversationContext represents the full context for AI conversations
type ConversationContext struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	ConversationID primitive.ObjectID `json:"conversation_id" bson:"conversation_id"`
	UserID         string             `json:"user_id" bson:"user_id"`
	CompanionID    string             `json:"companion_id" bson:"companion_id"`

	// Context layers
	BaseIdentityLayer  *ContextLayer `json:"base_identity_layer" bson:"base_identity_layer"`
	RelationshipLayer  *ContextLayer `json:"relationship_layer" bson:"relationship_layer"`
	ConversationLayer  *ContextLayer `json:"conversation_layer" bson:"conversation_layer"`
	SituationalLayer   *ContextLayer `json:"situational_layer" bson:"situational_layer"`
	ResponseStyleLayer *ContextLayer `json:"response_style_layer" bson:"response_style_layer"`

	// Emotional context
	UserEmotionalState      *EmotionalState     `json:"user_emotional_state" bson:"user_emotional_state"`
	CompanionEmotionalState *EmotionalState     `json:"companion_emotional_state" bson:"companion_emotional_state"`
	EmotionalHistory        []EmotionalSnapshot `json:"emotional_history" bson:"emotional_history"`

	// Memory and relationship
	ActiveMemories    []AIEnhancedMemoryEntry `json:"active_memories" bson:"active_memories"`
	RelationshipStage string                  `json:"relationship_stage" bson:"relationship_stage"`
	TrustLevel        float64                 `json:"trust_level" bson:"trust_level"`
	IntimacyLevel     float64                 `json:"intimacy_level" bson:"intimacy_level"`

	// Conversation flow
	CurrentTopic       string   `json:"current_topic" bson:"current_topic"`
	TopicHistory       []string `json:"topic_history" bson:"topic_history"`
	ConversationPacing string   `json:"conversation_pacing" bson:"conversation_pacing"`

	// Performance tracking
	TokenUsage       int     `json:"token_usage" bson:"token_usage"`
	ResponseQuality  float64 `json:"response_quality" bson:"response_quality"`
	UserSatisfaction float64 `json:"user_satisfaction" bson:"user_satisfaction"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// EmotionalState represents the current emotional state
type EmotionalState struct {
	PrimaryEmotion   string         `json:"primary_emotion" bson:"primary_emotion"`
	SecondaryEmotion string         `json:"secondary_emotion" bson:"secondary_emotion"`
	Intensity        float64        `json:"intensity" bson:"intensity"`
	Confidence       float64        `json:"confidence" bson:"confidence"`
	MixedEmotions    []string       `json:"mixed_emotions" bson:"mixed_emotions"`
	Triggers         []string       `json:"triggers" bson:"triggers"`
	Metadata         map[string]any `json:"metadata" bson:"metadata"`
	DetectedAt       time.Time      `json:"detected_at" bson:"detected_at"`
}

// EmotionalSnapshot represents a point-in-time emotional state
type EmotionalSnapshot struct {
	EmotionalState *EmotionalState    `json:"emotional_state" bson:"emotional_state"`
	MessageID      primitive.ObjectID `json:"message_id" bson:"message_id"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
	Context        string             `json:"context" bson:"context"`
}

// AIEnhancedMemoryEntry represents an enhanced memory entry for AI context
type AIEnhancedMemoryEntry struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id"`
	ConversationID  primitive.ObjectID   `json:"conversation_id" bson:"conversation_id"`
	Type            string               `json:"type" bson:"type"` // factual, emotional, conversational, behavioral, shared
	Category        string               `json:"category" bson:"category"`
	Content         string               `json:"content" bson:"content"`
	Importance      float64              `json:"importance" bson:"importance"`
	EmotionalWeight float64              `json:"emotional_weight" bson:"emotional_weight"`
	Frequency       int                  `json:"frequency" bson:"frequency"`
	LastReferenced  time.Time            `json:"last_referenced" bson:"last_referenced"`
	RelatedMemories []primitive.ObjectID `json:"related_memories" bson:"related_memories"`
	Metadata        map[string]any       `json:"metadata" bson:"metadata"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at" bson:"updated_at"`
}

// PromptTemplate represents a reusable prompt template
type PromptTemplate struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Name             string             `json:"name" bson:"name"`
	Description      string             `json:"description" bson:"description"`
	Template         string             `json:"template" bson:"template"`
	Variables        []string           `json:"variables" bson:"variables"`
	Scenario         string             `json:"scenario" bson:"scenario"`
	Version          string             `json:"version" bson:"version"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	PerformanceScore float64            `json:"performance_score" bson:"performance_score"`
	UsageCount       int                `json:"usage_count" bson:"usage_count"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// ResponseQuality represents quality metrics for AI responses
type ResponseQuality struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	MessageID      primitive.ObjectID `json:"message_id" bson:"message_id"`
	ConversationID primitive.ObjectID `json:"conversation_id" bson:"conversation_id"`

	// Quality metrics
	PersonalityConsistency      float64 `json:"personality_consistency" bson:"personality_consistency"`
	EmotionalAppropriateness    float64 `json:"emotional_appropriateness" bson:"emotional_appropriateness"`
	FactualAccuracy             float64 `json:"factual_accuracy" bson:"factual_accuracy"`
	RelationshipAppropriateness float64 `json:"relationship_appropriateness" bson:"relationship_appropriateness"`
	SafetyScore                 float64 `json:"safety_score" bson:"safety_score"`
	OverallQuality              float64 `json:"overall_quality" bson:"overall_quality"`

	// User feedback
	UserRating   *int   `json:"user_rating,omitempty" bson:"user_rating,omitempty"`
	UserFeedback string `json:"user_feedback" bson:"user_feedback"`

	// Analysis metadata
	AnalysisModel      string   `json:"analysis_model" bson:"analysis_model"`
	AnalysisConfidence float64  `json:"analysis_confidence" bson:"analysis_confidence"`
	IssuesDetected     []string `json:"issues_detected" bson:"issues_detected"`
	Suggestions        []string `json:"suggestions" bson:"suggestions"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ConversationIntelligence represents conversation flow analysis
type ConversationIntelligence struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	ConversationID primitive.ObjectID `json:"conversation_id" bson:"conversation_id"`

	// Topic analysis
	CurrentTopics    []string           `json:"current_topics" bson:"current_topics"`
	TopicDepth       map[string]float64 `json:"topic_depth" bson:"topic_depth"`
	TopicTransitions []TopicTransition  `json:"topic_transitions" bson:"topic_transitions"`

	// Flow analysis
	ConversationPacing  string  `json:"conversation_pacing" bson:"conversation_pacing"`
	QuestionBalance     float64 `json:"question_balance" bson:"question_balance"`
	InformationExchange float64 `json:"information_exchange" bson:"information_exchange"`
	EngagementLevel     float64 `json:"engagement_level" bson:"engagement_level"`

	// Relationship dynamics
	IntimacyProgression float64 `json:"intimacy_progression" bson:"intimacy_progression"`
	TrustBuilding       float64 `json:"trust_building" bson:"trust_building"`
	VulnerabilityLevel  float64 `json:"vulnerability_level" bson:"vulnerability_level"`

	// Recommendations
	SuggestedTopics         []string `json:"suggested_topics" bson:"suggested_topics"`
	PacingRecommendations   []string `json:"pacing_recommendations" bson:"pacing_recommendations"`
	RelationshipSuggestions []string `json:"relationship_suggestions" bson:"relationship_suggestions"`

	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// TopicTransition represents a transition between conversation topics
type TopicTransition struct {
	FromTopic      string    `json:"from_topic" bson:"from_topic"`
	ToTopic        string    `json:"to_topic" bson:"to_topic"`
	TransitionType string    `json:"transition_type" bson:"transition_type"` // natural, forced, user-initiated
	Smoothness     float64   `json:"smoothness" bson:"smoothness"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
}

// AIUsageMetrics represents usage and performance metrics
type AIUsageMetrics struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         string             `json:"user_id" bson:"user_id"`
	ConversationID primitive.ObjectID `json:"conversation_id" bson:"conversation_id"`

	// Usage tracking
	TotalTokens      int            `json:"total_tokens" bson:"total_tokens"`
	PromptTokens     int            `json:"prompt_tokens" bson:"prompt_tokens"`
	CompletionTokens int            `json:"completion_tokens" bson:"completion_tokens"`
	ModelUsage       map[string]int `json:"model_usage" bson:"model_usage"`

	// Performance metrics
	AverageResponseTime float64 `json:"average_response_time" bson:"average_response_time"`
	ResponseQuality     float64 `json:"response_quality" bson:"response_quality"`
	UserSatisfaction    float64 `json:"user_satisfaction" bson:"user_satisfaction"`

	// Cost tracking
	EstimatedCost  float64 `json:"estimated_cost" bson:"estimated_cost"`
	CostPerMessage float64 `json:"cost_per_message" bson:"cost_per_message"`

	// Session data
	SessionStart      time.Time  `json:"session_start" bson:"session_start"`
	SessionEnd        *time.Time `json:"session_end,omitempty" bson:"session_end,omitempty"`
	MessagesProcessed int        `json:"messages_processed" bson:"messages_processed"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
