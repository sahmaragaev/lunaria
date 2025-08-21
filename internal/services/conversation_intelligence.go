package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationIntelligenceService struct {
	grokService *GrokService
	repo        *repositories.ConversationRepository
}

func NewConversationIntelligenceService(grokService *GrokService, repo *repositories.ConversationRepository) *ConversationIntelligenceService {
	return &ConversationIntelligenceService{
		grokService: grokService,
		repo:        repo,
	}
}

// AnalyzeConversationFlow analyzes the current conversation flow and provides insights
func (s *ConversationIntelligenceService) AnalyzeConversationFlow(ctx context.Context, conversationID primitive.ObjectID) (*models.ConversationIntelligence, error) {
	// Get recent messages for analysis
	messages, _, _, err := s.repo.ListMessages(ctx, conversationID, 20, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	intelligence := &models.ConversationIntelligence{
		ID:             primitive.NewObjectID(),
		ConversationID: conversationID,
		UpdatedAt:      time.Now(),
	}

	// Analyze topics
	topics, topicDepth, err := s.analyzeTopics(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze topics: %w", err)
	}
	intelligence.CurrentTopics = topics
	intelligence.TopicDepth = topicDepth

	// Analyze conversation flow
	flowAnalysis, err := s.analyzeFlow(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze flow: %w", err)
	}
	intelligence.ConversationPacing = flowAnalysis.Pacing
	intelligence.QuestionBalance = flowAnalysis.QuestionBalance
	intelligence.InformationExchange = flowAnalysis.InformationExchange
	intelligence.EngagementLevel = flowAnalysis.EngagementLevel

	// Analyze relationship dynamics
	relationshipAnalysis, err := s.analyzeRelationshipDynamics(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze relationship dynamics: %w", err)
	}
	intelligence.IntimacyProgression = relationshipAnalysis.IntimacyProgression
	intelligence.TrustBuilding = relationshipAnalysis.TrustBuilding
	intelligence.VulnerabilityLevel = relationshipAnalysis.VulnerabilityLevel

	// Generate recommendations
	recommendations, err := s.generateRecommendations(ctx, intelligence, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}
	intelligence.SuggestedTopics = recommendations.Topics
	intelligence.PacingRecommendations = recommendations.Pacing
	intelligence.RelationshipSuggestions = recommendations.Relationship

	// Analyze topic transitions
	transitions, err := s.analyzeTopicTransitions(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze topic transitions: %w", err)
	}
	intelligence.TopicTransitions = transitions

	return intelligence, nil
}

// analyzeTopics identifies current topics and their depth
func (s *ConversationIntelligenceService) analyzeTopics(ctx context.Context, messages []*models.Message) ([]string, map[string]float64, error) {
	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze this conversation and identify the main topics being discussed:

CONVERSATION:
%s

Identify:
1. Current active topics (max 5)
2. Depth of each topic (0.0-1.0 scale)
3. Topic categories (personal, work, interests, emotions, etc.)

Respond with JSON:
{
  "topics": ["topic1", "topic2", "topic3"],
  "topic_depth": {
    "topic1": 0.8,
    "topic2": 0.5,
    "topic3": 0.3
  },
  "categories": {
    "topic1": "personal",
    "topic2": "work",
    "topic3": "interests"
  }
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a conversation topic analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to analyze topics: %w", err)
	}

	var analysis struct {
		Topics     []string           `json:"topics"`
		TopicDepth map[string]float64 `json:"topic_depth"`
		Categories map[string]string  `json:"categories"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, nil, fmt.Errorf("failed to parse topic analysis: %w", err)
	}

	return analysis.Topics, analysis.TopicDepth, nil
}

// analyzeFlow analyzes conversation flow characteristics
func (s *ConversationIntelligenceService) analyzeFlow(ctx context.Context, messages []*models.Message) (*FlowAnalysis, error) {
	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze the conversation flow characteristics:

CONVERSATION:
%s

Analyze:
1. Conversation pacing (slow, normal, fast)
2. Question balance (ratio of questions to statements)
3. Information exchange (how much each person shares)
4. Engagement level (how engaged both parties seem)

Respond with JSON:
{
  "pacing": "slow|normal|fast",
  "question_balance": 0.0-1.0,
  "information_exchange": 0.0-1.0,
  "engagement_level": 0.0-1.0,
  "analysis": "description of flow characteristics"
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a conversation flow analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze flow: %w", err)
	}

	var analysis FlowAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse flow analysis: %w", err)
	}

	return &analysis, nil
}

// analyzeRelationshipDynamics analyzes relationship progression
func (s *ConversationIntelligenceService) analyzeRelationshipDynamics(ctx context.Context, messages []*models.Message) (*RelationshipAnalysis, error) {
	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze the relationship dynamics in this conversation:

CONVERSATION:
%s

Analyze:
1. Intimacy progression (how intimate the conversation is)
2. Trust building (signs of trust development)
3. Vulnerability level (how vulnerable both parties are)

Respond with JSON:
{
  "intimacy_progression": 0.0-1.0,
  "trust_building": 0.0-1.0,
  "vulnerability_level": 0.0-1.0,
  "dynamics": "description of relationship dynamics"
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a relationship dynamics analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze relationship dynamics: %w", err)
	}

	var analysis RelationshipAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse relationship analysis: %w", err)
	}

	return &analysis, nil
}

// generateRecommendations generates conversation recommendations
func (s *ConversationIntelligenceService) generateRecommendations(ctx context.Context, intelligence *models.ConversationIntelligence, messages []*models.Message) (*Recommendations, error) {
	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Based on this conversation analysis, generate recommendations:

CONVERSATION:
%s

CURRENT ANALYSIS:
- Topics: %s
- Pacing: %s
- Engagement: %.2f
- Intimacy: %.2f
- Trust: %.2f

Generate recommendations for:
1. Suggested topics to explore
2. Pacing adjustments
3. Relationship building opportunities

Respond with JSON:
{
  "topics": ["topic1", "topic2", "topic3"],
  "pacing": ["rec1", "rec2"],
  "relationship": ["rec1", "rec2"]
}`,
		conversationText,
		strings.Join(intelligence.CurrentTopics, ", "),
		intelligence.ConversationPacing,
		intelligence.EngagementLevel,
		intelligence.IntimacyProgression,
		intelligence.TrustBuilding)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a conversation recommendation generator. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	var recommendations Recommendations
	if err := json.Unmarshal([]byte(response), &recommendations); err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %w", err)
	}

	return &recommendations, nil
}

// analyzeTopicTransitions analyzes how topics change throughout the conversation
func (s *ConversationIntelligenceService) analyzeTopicTransitions(ctx context.Context, messages []*models.Message) ([]models.TopicTransition, error) {
	if len(messages) < 2 {
		return []models.TopicTransition{}, nil
	}

	// Group messages into chunks for topic analysis
	chunks := s.createMessageChunks(messages, 5)
	var transitions []models.TopicTransition

	for i := 1; i < len(chunks); i++ {
		fromChunk := chunks[i-1]
		toChunk := chunks[i]

		transition, err := s.analyzeChunkTransition(ctx, fromChunk, toChunk, messages[i*5-1].CreatedAt)
		if err != nil {
			continue // Skip failed transitions
		}
		transitions = append(transitions, *transition)
	}

	return transitions, nil
}

// createMessageChunks groups messages into chunks for analysis
func (s *ConversationIntelligenceService) createMessageChunks(messages []*models.Message, chunkSize int) [][]*models.Message {
	var chunks [][]*models.Message
	for i := 0; i < len(messages); i += chunkSize {
		end := i + chunkSize
		if end > len(messages) {
			end = len(messages)
		}
		chunks = append(chunks, messages[i:end])
	}
	return chunks
}

// analyzeChunkTransition analyzes transition between two message chunks
func (s *ConversationIntelligenceService) analyzeChunkTransition(ctx context.Context, fromChunk, toChunk []*models.Message, timestamp time.Time) (*models.TopicTransition, error) {
	fromText := s.formatConversationForAnalysis(fromChunk)
	toText := s.formatConversationForAnalysis(toChunk)

	prompt := fmt.Sprintf(`Analyze the topic transition between these conversation chunks:

FROM CHUNK:
%s

TO CHUNK:
%s

Identify:
1. From topic
2. To topic
3. Transition type (natural, forced, user-initiated)
4. Smoothness (0.0-1.0)

Respond with JSON:
{
  "from_topic": "topic name",
  "to_topic": "topic name",
  "transition_type": "natural|forced|user-initiated",
  "smoothness": 0.0-1.0
}`,
		fromText,
		toText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a topic transition analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze transition: %w", err)
	}

	var transition models.TopicTransition
	if err := json.Unmarshal([]byte(response), &transition); err != nil {
		return nil, fmt.Errorf("failed to parse transition: %w", err)
	}

	transition.Timestamp = timestamp
	return &transition, nil
}

// formatConversationForAnalysis formats messages for AI analysis
func (s *ConversationIntelligenceService) formatConversationForAnalysis(messages []*models.Message) string {
	var formatted []string
	for _, msg := range messages {
		if msg.Text != nil {
			sender := "User"
			if msg.SenderType == "companion" {
				sender = "Companion"
			}
			formatted = append(formatted, fmt.Sprintf("%s: %s", sender, *msg.Text))
		}
	}
	return strings.Join(formatted, "\n")
}

// SuggestNextTopic suggests the next topic based on conversation context
func (s *ConversationIntelligenceService) SuggestNextTopic(ctx context.Context, conversationID primitive.ObjectID, companionProfile *models.CompanionProfile) (string, error) {
	// Get conversation intelligence
	intelligence, err := s.AnalyzeConversationFlow(ctx, conversationID)
	if err != nil {
		return "", fmt.Errorf("failed to analyze conversation: %w", err)
	}

	// Get recent messages for context
	messages, _, _, err := s.repo.ListMessages(ctx, conversationID, 10, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get messages: %w", err)
	}

	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Suggest the next topic for this conversation:

RECENT CONVERSATION:
%s

CURRENT TOPICS: %s
COMPANION INTERESTS: %s
ENGAGEMENT LEVEL: %.2f
INTIMACY LEVEL: %.2f

Suggest a natural next topic that:
1. Builds on current conversation
2. Matches companion interests
3. Maintains appropriate engagement
4. Supports relationship progression

Respond with just the topic name.`,
		conversationText,
		strings.Join(intelligence.CurrentTopics, ", "),
		strings.Join(companionProfile.Interests, ", "),
		intelligence.EngagementLevel,
		intelligence.IntimacyProgression)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a topic suggestion expert. Respond with only the topic name."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return "", fmt.Errorf("failed to suggest topic: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// AnalyzeEngagementLevel analyzes how engaged the user is in the conversation
func (s *ConversationIntelligenceService) AnalyzeEngagementLevel(ctx context.Context, conversationID primitive.ObjectID) (float64, error) {
	// Get recent messages
	messages, _, _, err := s.repo.ListMessages(ctx, conversationID, 15, nil)
	if err != nil {
		return 0.5, fmt.Errorf("failed to get messages: %w", err)
	}

	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze the user's engagement level in this conversation:

CONVERSATION:
%s

Consider:
1. Response length and detail
2. Question asking
3. Emotional investment
4. Topic contribution
5. Response time patterns

Rate engagement from 0.0 to 1.0 and respond with JSON:
{
  "engagement_level": 0.0-1.0,
  "indicators": ["indicator1", "indicator2"],
  "suggestions": ["suggestion1", "suggestion2"]
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are an engagement analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze engagement: %w", err)
	}

	var analysis struct {
		EngagementLevel float64  `json:"engagement_level"`
		Indicators      []string `json:"indicators"`
		Suggestions     []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse engagement analysis: %w", err)
	}

	return analysis.EngagementLevel, nil
}

// Helper structs for analysis results
type FlowAnalysis struct {
	Pacing              string  `json:"pacing"`
	QuestionBalance     float64 `json:"question_balance"`
	InformationExchange float64 `json:"information_exchange"`
	EngagementLevel     float64 `json:"engagement_level"`
	Analysis            string  `json:"analysis"`
}

type RelationshipAnalysis struct {
	IntimacyProgression float64 `json:"intimacy_progression"`
	TrustBuilding       float64 `json:"trust_building"`
	VulnerabilityLevel  float64 `json:"vulnerability_level"`
	Dynamics            string  `json:"dynamics"`
}

type Recommendations struct {
	Topics       []string `json:"topics"`
	Pacing       []string `json:"pacing"`
	Relationship []string `json:"relationship"`
}
