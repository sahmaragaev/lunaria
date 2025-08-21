package dto

// ConversationIntelligenceRequest represents a request for conversation intelligence
type ConversationIntelligenceRequest struct {
	ConversationID string `json:"conversation_id" validate:"required"`
}

// ConversationIntelligenceResponse represents conversation intelligence insights
type ConversationIntelligenceResponse struct {
	ConversationID          string             `json:"conversation_id"`
	CurrentTopics           []string           `json:"current_topics"`
	TopicDepth              map[string]float64 `json:"topic_depth"`
	ConversationPacing      string             `json:"conversation_pacing"`
	QuestionBalance         float64            `json:"question_balance"`
	InformationExchange     float64            `json:"information_exchange"`
	EngagementLevel         float64            `json:"engagement_level"`
	IntimacyProgression     float64            `json:"intimacy_progression"`
	TrustBuilding           float64            `json:"trust_building"`
	VulnerabilityLevel      float64            `json:"vulnerability_level"`
	SuggestedTopics         []string           `json:"suggested_topics"`
	PacingRecommendations   []string           `json:"pacing_recommendations"`
	RelationshipSuggestions []string           `json:"relationship_suggestions"`
}

// TopicSuggestionRequest represents a request for topic suggestions
type TopicSuggestionRequest struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Context        string `json:"context"`
}

// TopicSuggestionResponse represents a topic suggestion response
type TopicSuggestionResponse struct {
	SuggestedTopic    string   `json:"suggested_topic"`
	AlternativeTopics []string `json:"alternative_topics"`
	Reasoning         string   `json:"reasoning"`
}

// EngagementAnalysisRequest represents a request for engagement analysis
type EngagementAnalysisRequest struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	TimeRange      string `json:"time_range"` // e.g., "last_hour", "last_day", "last_week"
}

// EngagementAnalysisResponse represents engagement analysis results
type EngagementAnalysisResponse struct {
	ConversationID    string   `json:"conversation_id"`
	EngagementLevel   float64  `json:"engagement_level"`
	EngagementTrend   string   `json:"engagement_trend"` // increasing, decreasing, stable
	EngagementFactors []string `json:"engagement_factors"`
	Recommendations   []string `json:"recommendations"`
}

// ResponseQualityRequest represents a request for response quality analysis
type ResponseQualityRequest struct {
	MessageID      string `json:"message_id" validate:"required"`
	ConversationID string `json:"conversation_id" validate:"required"`
}

// ResponseQualityResponse represents response quality analysis results
type ResponseQualityResponse struct {
	MessageID                   string   `json:"message_id"`
	ConversationID              string   `json:"conversation_id"`
	PersonalityConsistency      float64  `json:"personality_consistency"`
	EmotionalAppropriateness    float64  `json:"emotional_appropriateness"`
	FactualAccuracy             float64  `json:"factual_accuracy"`
	RelationshipAppropriateness float64  `json:"relationship_appropriateness"`
	SafetyScore                 float64  `json:"safety_score"`
	OverallQuality              float64  `json:"overall_quality"`
	IssuesDetected              []string `json:"issues_detected"`
	Suggestions                 []string `json:"suggestions"`
	AnalysisConfidence          float64  `json:"analysis_confidence"`
}

// EmotionalAnalysisRequest represents a request for emotional analysis
type EmotionalAnalysisRequest struct {
	MessageText    string `json:"message_text" validate:"required"`
	ConversationID string `json:"conversation_id"`
	Context        string `json:"context"`
}

// EmotionalAnalysisResponse represents emotional analysis results
type EmotionalAnalysisResponse struct {
	PrimaryEmotion   string         `json:"primary_emotion"`
	SecondaryEmotion string         `json:"secondary_emotion"`
	Intensity        float64        `json:"intensity"`
	Confidence       float64        `json:"confidence"`
	MixedEmotions    []string       `json:"mixed_emotions"`
	Triggers         []string       `json:"triggers"`
	EmotionalContext map[string]any `json:"emotional_context"`
}

// MemoryExtractionRequest represents a request for memory extraction
type MemoryExtractionRequest struct {
	ConversationID string   `json:"conversation_id" validate:"required"`
	MessageIDs     []string `json:"message_ids"`
	ExtractionType string   `json:"extraction_type"` // factual, emotional, conversational, behavioral, shared
}

// MemoryExtractionResponse represents memory extraction results
type MemoryExtractionResponse struct {
	ConversationID string            `json:"conversation_id"`
	Memories       []ExtractedMemory `json:"memories"`
	Summary        string            `json:"summary"`
	Metadata       map[string]any    `json:"metadata"`
}

// ExtractedMemory represents a single extracted memory
type ExtractedMemory struct {
	Type            string  `json:"type"`
	Category        string  `json:"category"`
	Content         string  `json:"content"`
	Importance      float64 `json:"importance"`
	EmotionalWeight float64 `json:"emotional_weight"`
	Confidence      float64 `json:"confidence"`
}

// PromptTemplateRequest represents a request for prompt template management
type PromptTemplateRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Template    string         `json:"template" validate:"required"`
	Variables   []string       `json:"variables"`
	Scenario    string         `json:"scenario"`
	Version     string         `json:"version"`
	IsActive    bool           `json:"is_active"`
	Metadata    map[string]any `json:"metadata"`
}

// PromptTemplateResponse represents a prompt template response
type PromptTemplateResponse struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Template         string         `json:"template"`
	Variables        []string       `json:"variables"`
	Scenario         string         `json:"scenario"`
	Version          string         `json:"version"`
	IsActive         bool           `json:"is_active"`
	PerformanceScore float64        `json:"performance_score"`
	UsageCount       int            `json:"usage_count"`
	Metadata         map[string]any `json:"metadata"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

// AIUsageMetricsRequest represents a request for AI usage metrics
type AIUsageMetricsRequest struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	TimeRange      string `json:"time_range"` // e.g., "last_hour", "last_day", "last_week", "last_month"
}

// AIUsageMetricsResponse represents AI usage metrics
type AIUsageMetricsResponse struct {
	UserID              string             `json:"user_id"`
	ConversationID      string             `json:"conversation_id"`
	TotalTokens         int                `json:"total_tokens"`
	PromptTokens        int                `json:"prompt_tokens"`
	CompletionTokens    int                `json:"completion_tokens"`
	ModelUsage          map[string]int     `json:"model_usage"`
	AverageResponseTime float64            `json:"average_response_time"`
	ResponseQuality     float64            `json:"response_quality"`
	UserSatisfaction    float64            `json:"user_satisfaction"`
	EstimatedCost       float64            `json:"estimated_cost"`
	CostPerMessage      float64            `json:"cost_per_message"`
	MessagesProcessed   int                `json:"messages_processed"`
	SessionStart        string             `json:"session_start"`
	SessionEnd          string             `json:"session_end"`
	UsageTrends         map[string]float64 `json:"usage_trends"`
}

// ConversationContextRequest represents a request for conversation context
type ConversationContextRequest struct {
	ConversationID string   `json:"conversation_id" validate:"required"`
	IncludeLayers  []string `json:"include_layers"` // base_identity, relationship, conversation, situational, response_style
}

// ConversationContextResponse represents conversation context
type ConversationContextResponse struct {
	ConversationID          string                  `json:"conversation_id"`
	UserID                  string                  `json:"user_id"`
	CompanionID             string                  `json:"companion_id"`
	BaseIdentityLayer       *ContextLayerResponse   `json:"base_identity_layer,omitempty"`
	RelationshipLayer       *ContextLayerResponse   `json:"relationship_layer,omitempty"`
	ConversationLayer       *ContextLayerResponse   `json:"conversation_layer,omitempty"`
	SituationalLayer        *ContextLayerResponse   `json:"situational_layer,omitempty"`
	ResponseStyleLayer      *ContextLayerResponse   `json:"response_style_layer,omitempty"`
	UserEmotionalState      *EmotionalStateResponse `json:"user_emotional_state,omitempty"`
	CompanionEmotionalState *EmotionalStateResponse `json:"companion_emotional_state,omitempty"`
	RelationshipStage       string                  `json:"relationship_stage"`
	TrustLevel              float64                 `json:"trust_level"`
	IntimacyLevel           float64                 `json:"intimacy_level"`
	CurrentTopic            string                  `json:"current_topic"`
	ConversationPacing      string                  `json:"conversation_pacing"`
	ActiveMemories          []MemoryEntryResponse   `json:"active_memories"`
}

// ContextLayerResponse represents a context layer
type ContextLayerResponse struct {
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Priority  int            `json:"priority"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt string         `json:"created_at"`
	ExpiresAt string         `json:"expires_at,omitempty"`
}

// EmotionalStateResponse represents an emotional state
type EmotionalStateResponse struct {
	PrimaryEmotion   string         `json:"primary_emotion"`
	SecondaryEmotion string         `json:"secondary_emotion"`
	Intensity        float64        `json:"intensity"`
	Confidence       float64        `json:"confidence"`
	MixedEmotions    []string       `json:"mixed_emotions"`
	Triggers         []string       `json:"triggers"`
	Metadata         map[string]any `json:"metadata"`
	DetectedAt       string         `json:"detected_at"`
}

// MemoryEntryResponse represents a memory entry
type MemoryEntryResponse struct {
	ID              string         `json:"id"`
	Type            string         `json:"type"`
	Category        string         `json:"category"`
	Content         string         `json:"content"`
	Importance      float64        `json:"importance"`
	EmotionalWeight float64        `json:"emotional_weight"`
	Frequency       int            `json:"frequency"`
	LastReferenced  string         `json:"last_referenced"`
	Metadata        map[string]any `json:"metadata"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}
