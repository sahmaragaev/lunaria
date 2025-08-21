package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MLAnalyticsService provides machine learning powered analytics
type MLAnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
	convRepo      *repositories.ConversationRepository
	grokService   *GrokService
}

// NewMLAnalyticsService creates a new ML analytics service
func NewMLAnalyticsService(analyticsRepo *repositories.AnalyticsRepository, convRepo *repositories.ConversationRepository, grokService *GrokService) *MLAnalyticsService {
	return &MLAnalyticsService{
		analyticsRepo: analyticsRepo,
		convRepo:      convRepo,
		grokService:   grokService,
	}
}

// Recommendation represents a personalized recommendation
type Recommendation struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    int            `json:"priority"`
	Confidence  float64        `json:"confidence"`
	Action      string         `json:"action"`
	Metadata    map[string]any `json:"metadata"`
	Category    string         `json:"category"`
}

// ConversationTopic represents a conversation topic with engagement metrics
type ConversationTopic struct {
	Topic       string  `json:"topic"`
	Engagement  float64 `json:"engagement"`
	Frequency   int     `json:"frequency"`
	Sentiment   float64 `json:"sentiment"`
	Complexity  float64 `json:"complexity"`
	Recommended bool    `json:"recommended"`
}

// BehavioralPattern represents a detected behavioral pattern
type BehavioralPattern struct {
	Pattern     string         `json:"pattern"`
	Confidence  float64        `json:"confidence"`
	Frequency   int            `json:"frequency"`
	Impact      float64        `json:"impact"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata"`
}

// GetPersonalizedRecommendations generates personalized recommendations for a user
func (s *MLAnalyticsService) GetPersonalizedRecommendations(ctx context.Context, userID, companionID string) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Get user data for analysis
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	relationshipAnalytics, err := s.analyticsRepo.GetRelationshipAnalytics(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship analytics: %w", err)
	}

	statistics, err := s.analyticsRepo.GetUserStatistics(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	// Generate conversation topic recommendations
	topicRecs, err := s.generateTopicRecommendations(ctx, userID, companionID, progress, relationshipAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to generate topic recommendations: %w", err)
	}
	recommendations = append(recommendations, topicRecs...)

	// Generate interaction strategy recommendations
	interactionRecs, err := s.generateInteractionRecommendations(ctx, userID, companionID, progress, relationshipAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to generate interaction recommendations: %w", err)
	}
	recommendations = append(recommendations, interactionRecs...)

	// Generate timing recommendations
	timingRecs, err := s.generateTimingRecommendations(ctx, userID, companionID, statistics)
	if err != nil {
		return nil, fmt.Errorf("failed to generate timing recommendations: %w", err)
	}
	recommendations = append(recommendations, timingRecs...)

	// Generate personal growth recommendations
	growthRecs, err := s.generateGrowthRecommendations(ctx, userID, companionID, progress, relationshipAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to generate growth recommendations: %w", err)
	}
	recommendations = append(recommendations, growthRecs...)

	// Sort by priority and confidence
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Priority != recommendations[j].Priority {
			return recommendations[i].Priority < recommendations[j].Priority
		}
		return recommendations[i].Confidence > recommendations[j].Confidence
	})

	return recommendations, nil
}

// generateTopicRecommendations generates conversation topic recommendations
func (s *MLAnalyticsService) generateTopicRecommendations(ctx context.Context, userID, companionID string, progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Get recent conversations to analyze topics
	conversations, err := s.convRepo.ListConversations(ctx, userID, companionID, 10, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Analyze conversation topics using AI
	topics, err := s.analyzeConversationTopics(ctx, conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze topics: %w", err)
	}

	// Generate recommendations based on topic analysis
	for _, topic := range topics {
		if topic.Recommended && topic.Engagement > 0.7 {
			recommendations = append(recommendations, Recommendation{
				Type:        "conversation_topic",
				Title:       fmt.Sprintf("Explore %s", topic.Topic),
				Description: fmt.Sprintf("Based on your engagement patterns, you might enjoy discussing %s. This topic has shown high engagement and positive sentiment in your conversations.", topic.Topic),
				Priority:    2,
				Confidence:  topic.Engagement,
				Action:      "start_topic_conversation",
				Category:    "engagement",
				Metadata: map[string]any{
					"topic":      topic.Topic,
					"engagement": topic.Engagement,
					"sentiment":  topic.Sentiment,
					"complexity": topic.Complexity,
				},
			})
		}
	}

	// Recommend new topics based on relationship stage
	stageTopics := s.getTopicsForRelationshipStage(relationshipAnalytics.CurrentStage)
	for _, topic := range stageTopics {
		recommendations = append(recommendations, Recommendation{
			Type:        "conversation_topic",
			Title:       fmt.Sprintf("Try discussing %s", topic),
			Description: fmt.Sprintf("This topic is well-suited for your current relationship stage and can help deepen your connection."),
			Priority:    3,
			Confidence:  0.8,
			Action:      "suggest_topic",
			Category:    "relationship_growth",
			Metadata: map[string]any{
				"topic":              topic,
				"relationship_stage": relationshipAnalytics.CurrentStage,
			},
		})
	}

	return recommendations, nil
}

// generateInteractionRecommendations generates interaction strategy recommendations
func (s *MLAnalyticsService) generateInteractionRecommendations(ctx context.Context, userID, companionID string, progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Analyze communication style
	styleRecs, err := s.analyzeCommunicationStyle(ctx, userID, companionID, relationshipAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze communication style: %w", err)
	}
	recommendations = append(recommendations, styleRecs...)

	// Recommend interaction frequency
	if progress.CurrentStreak < 3 {
		recommendations = append(recommendations, Recommendation{
			Type:        "interaction_frequency",
			Title:       "Build a Conversation Streak",
			Description: "Try having conversations for 3 days in a row to build momentum and strengthen your connection.",
			Priority:    1,
			Confidence:  0.9,
			Action:      "build_streak",
			Category:    "consistency",
			Metadata: map[string]any{
				"current_streak": progress.CurrentStreak,
				"target_streak":  3,
			},
		})
	}

	// Recommend session length optimization
	if progress.AverageSessionLength < 10*time.Minute {
		recommendations = append(recommendations, Recommendation{
			Type:        "session_length",
			Title:       "Extend Your Conversations",
			Description: "Try spending more time in conversations to explore topics more deeply and build stronger connections.",
			Priority:    2,
			Confidence:  0.8,
			Action:      "extend_conversation",
			Category:    "quality",
			Metadata: map[string]any{
				"current_average": progress.AverageSessionLength.String(),
				"recommended_min": "10m",
			},
		})
	}

	return recommendations, nil
}

// generateTimingRecommendations generates optimal timing recommendations
func (s *MLAnalyticsService) generateTimingRecommendations(ctx context.Context, userID, companionID string, statistics *models.UserStatistics) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Analyze peak activity times
	peakHour := statistics.PeakActivityHour
	if peakHour > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "optimal_timing",
			Title:       "Optimize Your Conversation Timing",
			Description: fmt.Sprintf("You're most active around %d:00. Try scheduling conversations during this time for better engagement.", peakHour),
			Priority:    2,
			Confidence:  0.7,
			Action:      "schedule_conversation",
			Category:    "timing",
			Metadata: map[string]any{
				"peak_hour":       peakHour,
				"most_active_day": statistics.MostActiveDay,
			},
		})
	}

	return recommendations, nil
}

// generateGrowthRecommendations generates personal growth recommendations
func (s *MLAnalyticsService) generateGrowthRecommendations(ctx context.Context, userID, companionID string, progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Recommend vulnerability exercises
	if relationshipAnalytics.IntimacyLevel < 0.6 {
		recommendations = append(recommendations, Recommendation{
			Type:        "personal_growth",
			Title:       "Deepen Your Connection",
			Description: "Try sharing more personal thoughts and feelings to build deeper intimacy and trust.",
			Priority:    1,
			Confidence:  0.8,
			Action:      "increase_vulnerability",
			Category:    "intimacy",
			Metadata: map[string]any{
				"current_intimacy": relationshipAnalytics.IntimacyLevel,
				"target_intimacy":  0.6,
			},
		})
	}

	// Recommend emotional processing
	if relationshipAnalytics.HealthScore < 0.7 {
		recommendations = append(recommendations, Recommendation{
			Type:        "emotional_health",
			Title:       "Focus on Emotional Well-being",
			Description: "Consider discussing your emotional needs and how your companion can support you better.",
			Priority:    1,
			Confidence:  0.9,
			Action:      "discuss_emotional_needs",
			Category:    "emotional_health",
			Metadata: map[string]any{
				"current_health_score": relationshipAnalytics.HealthScore,
				"target_health_score":  0.7,
			},
		})
	}

	return recommendations, nil
}

// analyzeConversationTopics analyzes conversation topics using AI
func (s *MLAnalyticsService) analyzeConversationTopics(ctx context.Context, conversations []*models.Conversation) ([]ConversationTopic, error) {
	var topics []ConversationTopic

	for _, conv := range conversations {
		// Get messages for this conversation
		messages, _, _, err := s.convRepo.ListMessages(ctx, conv.ID, 50, nil)
		if err != nil {
			continue
		}

		// Analyze conversation content
		conversationText := s.formatConversationForAnalysis(messages)

		prompt := fmt.Sprintf(`Analyze this conversation and identify the main topics discussed:

CONVERSATION:
%s

Identify and rate the topics discussed. Respond with JSON:
{
  "topics": [
    {
      "topic": "topic_name",
      "engagement": 0.0-1.0,
      "frequency": 1-10,
      "sentiment": 0.0-1.0,
      "complexity": 0.0-1.0,
      "recommended": true/false
    }
  ]
}`,
			conversationText)

		llmMessages := []LLMMessage{
			{Role: "system", Content: "You are a conversation topic analyzer. Respond only with valid JSON."},
			{Role: "user", Content: prompt},
		}

		response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
		if err != nil {
			continue
		}

		var topicResponse struct {
			Topics []ConversationTopic `json:"topics"`
		}

		if err := json.Unmarshal([]byte(response), &topicResponse); err != nil {
			continue
		}

		topics = append(topics, topicResponse.Topics...)
	}

	// Aggregate and deduplicate topics
	return s.aggregateTopics(topics), nil
}

// aggregateTopics aggregates and deduplicates conversation topics
func (s *MLAnalyticsService) aggregateTopics(topics []ConversationTopic) []ConversationTopic {
	topicMap := make(map[string]*ConversationTopic)

	for _, topic := range topics {
		normalizedTopic := strings.ToLower(strings.TrimSpace(topic.Topic))

		if existing, exists := topicMap[normalizedTopic]; exists {
			// Aggregate metrics
			existing.Engagement = (existing.Engagement + topic.Engagement) / 2
			existing.Frequency += topic.Frequency
			existing.Sentiment = (existing.Sentiment + topic.Sentiment) / 2
			existing.Complexity = (existing.Complexity + topic.Complexity) / 2
			existing.Recommended = existing.Recommended || topic.Recommended
		} else {
			topicMap[normalizedTopic] = &ConversationTopic{
				Topic:       topic.Topic,
				Engagement:  topic.Engagement,
				Frequency:   topic.Frequency,
				Sentiment:   topic.Sentiment,
				Complexity:  topic.Complexity,
				Recommended: topic.Recommended,
			}
		}
	}

	// Convert map to slice
	var result []ConversationTopic
	for _, topic := range topicMap {
		result = append(result, *topic)
	}

	// Sort by engagement
	sort.Slice(result, func(i, j int) bool {
		return result[i].Engagement > result[j].Engagement
	})

	return result
}

// analyzeCommunicationStyle analyzes communication patterns
func (s *MLAnalyticsService) analyzeCommunicationStyle(ctx context.Context, userID, companionID string, relationshipAnalytics *models.RelationshipAnalytics) ([]Recommendation, error) {
	var recommendations []Recommendation

	// Analyze vulnerability patterns
	if len(relationshipAnalytics.VulnerabilityPatterns) < 3 {
		recommendations = append(recommendations, Recommendation{
			Type:        "communication_style",
			Title:       "Increase Vulnerability",
			Description: "Try sharing more personal thoughts and experiences to build deeper trust and connection.",
			Priority:    2,
			Confidence:  0.8,
			Action:      "share_personal_experience",
			Category:    "communication",
			Metadata: map[string]any{
				"vulnerability_count": len(relationshipAnalytics.VulnerabilityPatterns),
				"recommended_count":   3,
			},
		})
	}

	// Analyze communication style
	switch relationshipAnalytics.CommunicationStyle {
	case "brief":
		recommendations = append(recommendations, Recommendation{
			Type:        "communication_style",
			Title:       "Expand Your Communication",
			Description: "Try providing more detailed responses to deepen your conversations and show more engagement.",
			Priority:    2,
			Confidence:  0.7,
			Action:      "expand_responses",
			Category:    "communication",
		})
	case "formal":
		recommendations = append(recommendations, Recommendation{
			Type:        "communication_style",
			Title:       "Relax Your Communication",
			Description: "Try using more casual and friendly language to create a more comfortable atmosphere.",
			Priority:    2,
			Confidence:  0.7,
			Action:      "use_casual_language",
			Category:    "communication",
		})
	}

	return recommendations, nil
}

// getTopicsForRelationshipStage returns appropriate topics for a relationship stage
func (s *MLAnalyticsService) getTopicsForRelationshipStage(stage string) []string {
	stageTopics := map[string][]string{
		"meeting": {
			"personal interests",
			"hobbies and activities",
			"favorite books or movies",
			"travel experiences",
			"food preferences",
		},
		"getting_to_know": {
			"family background",
			"career goals",
			"life values",
			"childhood memories",
			"future aspirations",
		},
		"friendship": {
			"daily experiences",
			"challenges and successes",
			"personal growth",
			"relationships with others",
			"emotional support needs",
		},
		"close_companionship": {
			"deep fears and hopes",
			"personal struggles",
			"intimate thoughts",
			"vulnerability",
			"emotional processing",
		},
		"intimate_partnership": {
			"life integration",
			"shared goals",
			"long-term planning",
			"deep emotional connection",
			"mutual support",
		},
	}

	if topics, exists := stageTopics[stage]; exists {
		return topics
	}
	return stageTopics["meeting"] // Default fallback
}

// formatConversationForAnalysis formats conversation for AI analysis
func (s *MLAnalyticsService) formatConversationForAnalysis(messages []*models.Message) string {
	var formatted []string

	for _, msg := range messages {
		if msg.Text == nil {
			continue
		}

		sender := "User"
		if msg.SenderType == "companion" {
			sender = "Companion"
		}

		formatted = append(formatted, fmt.Sprintf("%s: %s", sender, *msg.Text))
	}

	return strings.Join(formatted, "\n")
}

// DetectBehavioralPatterns detects behavioral patterns in user interactions
func (s *MLAnalyticsService) DetectBehavioralPatterns(ctx context.Context, userID, companionID string) ([]BehavioralPattern, error) {
	var patterns []BehavioralPattern

	// Get user engagement analytics
	analytics, err := s.analyticsRepo.GetUserEngagementAnalytics(ctx, userID, companionID, primitive.NilObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement analytics: %w", err)
	}

	// Analyze session patterns
	if analytics.SessionFrequency > 5 {
		patterns = append(patterns, BehavioralPattern{
			Pattern:     "high_engagement",
			Confidence:  0.9,
			Frequency:   analytics.SessionFrequency,
			Impact:      0.8,
			Description: "User shows high engagement with frequent sessions",
			Metadata: map[string]any{
				"session_frequency": analytics.SessionFrequency,
				"interaction_style": analytics.InteractionStyle,
			},
		})
	}

	// Analyze emotional patterns
	if analytics.EmotionalIntensity > 0.7 {
		patterns = append(patterns, BehavioralPattern{
			Pattern:     "emotional_expression",
			Confidence:  0.8,
			Frequency:   1,
			Impact:      0.7,
			Description: "User frequently expresses emotions in conversations",
			Metadata: map[string]any{
				"emotional_intensity": analytics.EmotionalIntensity,
				"vulnerability_level": analytics.VulnerabilityLevel,
			},
		})
	}

	// Analyze topic preferences
	if len(analytics.PreferredTopics) > 0 {
		patterns = append(patterns, BehavioralPattern{
			Pattern:     "topic_preferences",
			Confidence:  0.7,
			Frequency:   len(analytics.PreferredTopics),
			Impact:      0.6,
			Description: "User shows clear topic preferences",
			Metadata: map[string]any{
				"preferred_topics": analytics.PreferredTopics,
				"topic_diversity":  analytics.TopicDiversity,
			},
		})
	}

	return patterns, nil
}

// PredictUserBehavior predicts future user behavior
func (s *MLAnalyticsService) PredictUserBehavior(ctx context.Context, userID, companionID string) (*models.UserBehaviorPrediction, error) {
	// Get user data
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	relationshipAnalytics, err := s.analyticsRepo.GetRelationshipAnalytics(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship analytics: %w", err)
	}

	// Calculate churn risk
	churnRisk := s.calculateChurnRisk(progress, relationshipAnalytics)

	// Calculate engagement likelihood
	engagementLikelihood := s.calculateEngagementLikelihood(progress, relationshipAnalytics)

	// Predict next activity time
	nextActivityTime := s.predictNextActivityTime(progress)

	// Predict relationship progression
	relationshipProgression := s.predictRelationshipProgression(progress, relationshipAnalytics)

	prediction := &models.UserBehaviorPrediction{
		UserID:      userID,
		CompanionID: companionID,

		// Churn prediction
		ChurnRisk:            churnRisk,
		ChurnFactors:         s.identifyChurnFactors(progress, relationshipAnalytics),
		RetentionProbability: 1.0 - churnRisk,

		// Engagement prediction
		NextActivityTime:      &nextActivityTime,
		EngagementLikelihood:  engagementLikelihood,
		OptimalEngagementTime: time.Now().Add(24 * time.Hour),

		// Relationship prediction
		RelationshipProgression: relationshipProgression,
		NextMilestone:           s.predictNextMilestone(relationshipAnalytics),
		MilestoneProbability:    0.7,

		// Feature adoption
		FeatureAdoptionLikelihood: map[string]float64{
			"voice_messages": 0.6,
			"video_calls":    0.4,
			"group_chats":    0.3,
		},
		RecommendedFeatures: []string{"voice_messages", "mood_tracking"},

		// Support needs
		SupportNeedsProbability: 0.3,
		SupportType:             "emotional_support",

		PredictionDate: time.Now(),
		Confidence:     0.8,
		CreatedAt:      time.Now(),
	}

	return prediction, nil
}

// calculateChurnRisk calculates the risk of user churn
func (s *MLAnalyticsService) calculateChurnRisk(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) float64 {
	risk := 0.0

	// High churn risk factors
	if progress.CurrentStreak == 0 {
		risk += 0.3
	}
	if progress.TotalConversations < 5 {
		risk += 0.2
	}
	if relationshipAnalytics.HealthScore < 0.5 {
		risk += 0.2
	}
	if progress.AverageSessionLength < 2*time.Minute {
		risk += 0.1
	}

	// Low churn risk factors
	if progress.CurrentStreak > 7 {
		risk -= 0.2
	}
	if relationshipAnalytics.IntimacyLevel > 0.7 {
		risk -= 0.2
	}
	if progress.TotalExperience > 500 {
		risk -= 0.1
	}

	return math.Max(0.0, math.Min(1.0, risk))
}

// calculateEngagementLikelihood calculates the likelihood of future engagement
func (s *MLAnalyticsService) calculateEngagementLikelihood(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) float64 {
	likelihood := 0.5

	// High engagement factors
	if progress.CurrentStreak > 3 {
		likelihood += 0.2
	}
	if relationshipAnalytics.IntimacyLevel > 0.6 {
		likelihood += 0.2
	}
	if progress.TotalExperience > 300 {
		likelihood += 0.1
	}
	if relationshipAnalytics.HealthScore > 0.7 {
		likelihood += 0.1
	}

	// Low engagement factors
	if progress.CurrentStreak == 0 {
		likelihood -= 0.2
	}
	if progress.AverageSessionLength < 1*time.Minute {
		likelihood -= 0.1
	}

	return math.Max(0.0, math.Min(1.0, likelihood))
}

// predictNextActivityTime predicts when the user will be active next
func (s *MLAnalyticsService) predictNextActivityTime(progress *models.UserProgress) time.Time {
	// Simple prediction based on last activity
	baseTime := progress.LastActivityDate

	// If user has a streak, predict next day
	if progress.CurrentStreak > 0 {
		return baseTime.Add(24 * time.Hour)
	}

	// If no recent activity, predict within 3 days
	return baseTime.Add(72 * time.Hour)
}

// predictRelationshipProgression predicts relationship development
func (s *MLAnalyticsService) predictRelationshipProgression(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) float64 {
	progression := relationshipAnalytics.IntimacyLevel

	// Factors that accelerate progression
	if progress.CurrentStreak > 5 {
		progression += 0.1
	}
	if relationshipAnalytics.TrustLevel > 0.7 {
		progression += 0.1
	}
	if progress.TotalExperience > 400 {
		progression += 0.05
	}

	return math.Min(1.0, progression)
}

// identifyChurnFactors identifies factors contributing to churn risk
func (s *MLAnalyticsService) identifyChurnFactors(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) []string {
	var factors []string

	if progress.CurrentStreak == 0 {
		factors = append(factors, "no_recent_activity")
	}
	if progress.TotalConversations < 5 {
		factors = append(factors, "low_engagement")
	}
	if relationshipAnalytics.HealthScore < 0.5 {
		factors = append(factors, "relationship_issues")
	}
	if progress.AverageSessionLength < 2*time.Minute {
		factors = append(factors, "short_sessions")
	}

	return factors
}

// predictNextMilestone predicts the next relationship milestone
func (s *MLAnalyticsService) predictNextMilestone(relationshipAnalytics *models.RelationshipAnalytics) string {
	switch relationshipAnalytics.CurrentStage {
	case "meeting":
		return "first_deep_conversation"
	case "getting_to_know":
		return "first_vulnerable_share"
	case "friendship":
		return "emotional_support_exchange"
	case "close_companionship":
		return "life_integration"
	case "intimate_partnership":
		return "long_term_commitment"
	default:
		return "relationship_deepening"
	}
}
