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

type ResponseQualityService struct {
	grokService *GrokService
	repo        *repositories.ConversationRepository
}

func NewResponseQualityService(grokService *GrokService, repo *repositories.ConversationRepository) *ResponseQualityService {
	return &ResponseQualityService{
		grokService: grokService,
		repo:        repo,
	}
}

// ValidateResponseQuality validates AI response quality using multiple metrics
func (s *ResponseQualityService) ValidateResponseQuality(ctx context.Context, response *models.Message, conversation *models.Conversation, companionProfile *models.CompanionProfile) (*models.ResponseQuality, error) {
	if response.Text == nil {
		return nil, fmt.Errorf("response has no text content")
	}

	quality := &models.ResponseQuality{
		ID:             primitive.NewObjectID(),
		MessageID:      response.ID,
		ConversationID: conversation.ID,
		CreatedAt:      time.Now(),
	}

	// Analyze personality consistency
	personalityScore, err := s.analyzePersonalityConsistency(ctx, *response.Text, companionProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze personality consistency: %w", err)
	}
	quality.PersonalityConsistency = personalityScore

	// Analyze emotional appropriateness
	emotionalScore, err := s.analyzeEmotionalAppropriateness(ctx, *response.Text, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze emotional appropriateness: %w", err)
	}
	quality.EmotionalAppropriateness = emotionalScore

	// Analyze factual accuracy
	factualScore, err := s.analyzeFactualAccuracy(ctx, *response.Text, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze factual accuracy: %w", err)
	}
	quality.FactualAccuracy = factualScore

	// Analyze relationship appropriateness
	relationshipScore, err := s.analyzeRelationshipAppropriateness(ctx, *response.Text, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze relationship appropriateness: %w", err)
	}
	quality.RelationshipAppropriateness = relationshipScore

	// Analyze safety
	safetyScore, err := s.analyzeSafety(ctx, *response.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze safety: %w", err)
	}
	quality.SafetyScore = safetyScore

	// Calculate overall quality
	quality.OverallQuality = s.calculateOverallQuality(quality)

	// Set analysis metadata
	quality.AnalysisModel = "grok-mini"
	quality.AnalysisConfidence = 0.85

	// Generate suggestions for improvement
	quality.Suggestions = s.generateImprovementSuggestions(quality)

	return quality, nil
}

// analyzePersonalityConsistency checks if response aligns with companion personality
func (s *ResponseQualityService) analyzePersonalityConsistency(ctx context.Context, responseText string, profile *models.CompanionProfile) (float64, error) {
	prompt := fmt.Sprintf(`Analyze if this response is consistent with the companion's personality traits:

COMPANION PERSONALITY:
- Warmth: %.1f/1.0
- Playfulness: %.1f/1.0
- Intelligence: %.1f/1.0
- Empathy: %.1f/1.0
- Confidence: %.1f/1.0
- Romance: %.1f/1.0
- Humor: %.1f/1.0
- Clinginess: %.1f/1.0

Communication Style:
- Formality: %.1f/1.0
- Emotionality: %.1f/1.0
- Playfulness: %.1f/1.0
- Intimacy: %.1f/1.0

Interests: %s
Quirks: %s

RESPONSE: "%s"

Rate the personality consistency from 0.0 to 1.0 and respond with JSON:
{
  "score": 0.0-1.0,
  "issues": ["issue1", "issue2"],
  "strengths": ["strength1", "strength2"]
}`,
		profile.Personality.Warmth,
		profile.Personality.Playfulness,
		profile.Personality.Intelligence,
		profile.Personality.Empathy,
		profile.Personality.Confidence,
		profile.Personality.Romance,
		profile.Personality.Humor,
		profile.Personality.Clinginess,
		profile.CommunicationStyle.Formality,
		profile.CommunicationStyle.Emotionality,
		profile.CommunicationStyle.Playfulness,
		profile.CommunicationStyle.Intimacy,
		strings.Join(profile.Interests, ", "),
		strings.Join(profile.Quirks, ", "),
		responseText)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a personality consistency analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze personality: %w", err)
	}

	var analysis struct {
		Score     float64  `json:"score"`
		Issues    []string `json:"issues"`
		Strengths []string `json:"strengths"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse personality analysis: %w", err)
	}

	return analysis.Score, nil
}

// analyzeEmotionalAppropriateness checks if response matches emotional context
func (s *ResponseQualityService) analyzeEmotionalAppropriateness(ctx context.Context, responseText string, conversation *models.Conversation) (float64, error) {
	// Get recent messages for emotional context
	messages, _, _, err := s.repo.ListMessages(ctx, conversation.ID, 5, nil)
	if err != nil {
		return 0.5, fmt.Errorf("failed to get recent messages: %w", err)
	}

	// Extract emotional context from recent messages
	emotionalContext := s.extractEmotionalContext(messages)

	prompt := fmt.Sprintf(`Analyze if this response is emotionally appropriate given the conversation context:

RECENT CONVERSATION CONTEXT:
%s

RESPONSE: "%s"

Rate the emotional appropriateness from 0.0 to 1.0 and respond with JSON:
{
  "score": 0.0-1.0,
  "emotional_match": "description",
  "suggestions": ["suggestion1", "suggestion2"]
}`,
		emotionalContext,
		responseText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are an emotional appropriateness analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze emotional appropriateness: %w", err)
	}

	var analysis struct {
		Score          float64  `json:"score"`
		EmotionalMatch string   `json:"emotional_match"`
		Suggestions    []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse emotional analysis: %w", err)
	}

	return analysis.Score, nil
}

// analyzeFactualAccuracy checks if response is consistent with established facts
func (s *ResponseQualityService) analyzeFactualAccuracy(ctx context.Context, responseText string, conversation *models.Conversation) (float64, error) {
	// Get conversation memories and context from database
	memories, err := s.repo.GetMemories(ctx, conversation.ID, 20)
	if err != nil {
		// Log error but continue with basic analysis
		fmt.Printf("Failed to retrieve memories for factual accuracy analysis: %v\n", err)
		memories = []models.AIEnhancedMemoryEntry{}
	}

	// Get conversation context for additional factual information
	context, err := s.repo.GetConversationContext(ctx, conversation.ID)
	if err != nil {
		// Log error but continue with basic analysis
		fmt.Printf("Failed to retrieve conversation context for factual accuracy analysis: %v\n", err)
		context = nil
	}

	// Build factual context from memories and conversation context
	factualContext := s.buildFactualContext(memories, context)

	prompt := fmt.Sprintf(`Analyze if this response is factually accurate and consistent with established information:

ESTABLISHED FACTS AND MEMORIES:
%s

RESPONSE: "%s"

Consider:
1. Consistency with established facts about the user
2. Consistency with previous conversation topics
3. Internal consistency within the response
4. Logical coherence
5. Appropriate level of detail
6. No contradictory statements with known information

Rate the factual accuracy from 0.0 to 1.0 and respond with JSON:
{
  "score": 0.0-1.0,
  "issues": ["issue1", "issue2"],
  "strengths": ["strength1", "strength2"],
  "factual_consistency": "description"
}`,
		factualContext,
		responseText)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a factual accuracy analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze factual accuracy: %w", err)
	}

	var analysis struct {
		Score     float64  `json:"score"`
		Issues    []string `json:"issues"`
		Strengths []string `json:"strengths"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse factual analysis: %w", err)
	}

	return analysis.Score, nil
}

// analyzeRelationshipAppropriateness checks if response matches relationship stage
func (s *ResponseQualityService) analyzeRelationshipAppropriateness(ctx context.Context, responseText string, conversation *models.Conversation) (float64, error) {
	// Get relationship stage and context from conversation context
	context, err := s.repo.GetConversationContext(ctx, conversation.ID)
	if err != nil {
		// Log error but continue with default values
		fmt.Printf("Failed to retrieve conversation context for relationship analysis: %v\n", err)
		context = nil
	}

	// Extract relationship information
	relationshipStage := "getting_to_know" // Default
	trustLevel := 0.5
	intimacyLevel := 0.3

	if context != nil {
		if context.RelationshipStage != "" {
			relationshipStage = context.RelationshipStage
		}
		trustLevel = context.TrustLevel
		intimacyLevel = context.IntimacyLevel
	}

	prompt := fmt.Sprintf(`Analyze if this response is appropriate for the current relationship stage and dynamics:

RELATIONSHIP CONTEXT:
- Stage: %s
- Trust Level: %.1f/1.0
- Intimacy Level: %.1f/1.0

RESPONSE: "%s"

Consider:
1. Intimacy level appropriateness for current stage
2. Trust building and respect for boundaries
3. Relationship progression appropriateness
4. Emotional safety and comfort
5. Consistency with established relationship dynamics

Rate the relationship appropriateness from 0.0 to 1.0 and respond with JSON:
{
  "score": 0.0-1.0,
  "appropriateness": "description",
  "concerns": ["concern1", "concern2"],
  "relationship_insights": "detailed analysis"
}`,
		relationshipStage,
		trustLevel,
		intimacyLevel,
		responseText)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a relationship appropriateness analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze relationship appropriateness: %w", err)
	}

	var analysis struct {
		Score           float64  `json:"score"`
		Appropriateness string   `json:"appropriateness"`
		Concerns        []string `json:"concerns"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse relationship analysis: %w", err)
	}

	return analysis.Score, nil
}

// analyzeSafety checks if response is safe and appropriate
func (s *ResponseQualityService) analyzeSafety(ctx context.Context, responseText string) (float64, error) {
	prompt := fmt.Sprintf(`Analyze if this response is safe and appropriate:

RESPONSE: "%s"

Check for:
1. Harmful content
2. Inappropriate suggestions
3. Unhealthy relationship patterns
4. Professional boundaries
5. Mental health concerns

Rate the safety from 0.0 to 1.0 and respond with JSON:
{
  "score": 0.0-1.0,
  "safety_level": "safe|caution|unsafe",
  "concerns": ["concern1", "concern2"],
  "recommendations": ["rec1", "rec2"]
}`,
		responseText)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a safety analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return 0.5, fmt.Errorf("failed to analyze safety: %w", err)
	}

	var analysis struct {
		Score           float64  `json:"score"`
		SafetyLevel     string   `json:"safety_level"`
		Concerns        []string `json:"concerns"`
		Recommendations []string `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return 0.5, fmt.Errorf("failed to parse safety analysis: %w", err)
	}

	return analysis.Score, nil
}

// calculateOverallQuality calculates the overall quality score
func (s *ResponseQualityService) calculateOverallQuality(quality *models.ResponseQuality) float64 {
	// Weighted average of all quality metrics
	weights := map[string]float64{
		"personality":  0.25,
		"emotional":    0.25,
		"factual":      0.20,
		"relationship": 0.20,
		"safety":       0.10,
	}

	overall := quality.PersonalityConsistency*weights["personality"] +
		quality.EmotionalAppropriateness*weights["emotional"] +
		quality.FactualAccuracy*weights["factual"] +
		quality.RelationshipAppropriateness*weights["relationship"] +
		quality.SafetyScore*weights["safety"]

	return overall
}

// generateImprovementSuggestions generates suggestions for response improvement
func (s *ResponseQualityService) generateImprovementSuggestions(quality *models.ResponseQuality) []string {
	var suggestions []string

	if quality.PersonalityConsistency < 0.7 {
		suggestions = append(suggestions, "Improve personality consistency")
	}
	if quality.EmotionalAppropriateness < 0.7 {
		suggestions = append(suggestions, "Better match emotional context")
	}
	if quality.FactualAccuracy < 0.7 {
		suggestions = append(suggestions, "Ensure factual consistency")
	}
	if quality.RelationshipAppropriateness < 0.7 {
		suggestions = append(suggestions, "Adjust to relationship stage")
	}
	if quality.SafetyScore < 0.8 {
		suggestions = append(suggestions, "Review safety concerns")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Response quality is good")
	}

	return suggestions
}

// extractEmotionalContext extracts emotional context from recent messages
func (s *ResponseQualityService) extractEmotionalContext(messages []*models.Message) string {
	if len(messages) == 0 {
		return "No recent conversation context"
	}

	var context []string
	for _, msg := range messages {
		if msg.Text != nil {
			sender := "User"
			if msg.SenderType == "companion" {
				sender = "Companion"
			}
			context = append(context, fmt.Sprintf("%s: %s", sender, *msg.Text))
		}
	}

	return strings.Join(context, "\n")
}

// buildFactualContext builds factual context from memories and conversation context
func (s *ResponseQualityService) buildFactualContext(memories []models.AIEnhancedMemoryEntry, context *models.ConversationContext) string {
	var factualInfo []string

	// Add factual memories
	for _, memory := range memories {
		if memory.Type == "factual" && memory.Importance > 0.6 {
			factualInfo = append(factualInfo, fmt.Sprintf("- %s (Importance: %.1f)", memory.Content, memory.Importance))
		}
	}

	// Add conversation context information
	if context != nil {
		if context.CurrentTopic != "" {
			factualInfo = append(factualInfo, fmt.Sprintf("- Current conversation topic: %s", context.CurrentTopic))
		}
		if len(context.TopicHistory) > 0 {
			recentTopics := context.TopicHistory[len(context.TopicHistory)-3:] // Last 3 topics
			factualInfo = append(factualInfo, fmt.Sprintf("- Recent conversation topics: %s", strings.Join(recentTopics, ", ")))
		}
		if context.RelationshipStage != "" {
			factualInfo = append(factualInfo, fmt.Sprintf("- Relationship stage: %s", context.RelationshipStage))
		}
	}

	if len(factualInfo) == 0 {
		return "No established facts or context available."
	}

	return strings.Join(factualInfo, "\n")
}

// RefineResponse improves a response based on quality analysis
func (s *ResponseQualityService) RefineResponse(ctx context.Context, originalResponse *models.Message, quality *models.ResponseQuality, conversation *models.Conversation, companionProfile *models.CompanionProfile) (*models.Message, error) {
	if quality.OverallQuality >= 0.8 {
		// Response is good enough, no refinement needed
		return originalResponse, nil
	}

	// Generate improvement suggestions
	improvementPrompt := fmt.Sprintf(`Improve this response based on the quality analysis:

ORIGINAL RESPONSE: "%s"

QUALITY ISSUES:
- Personality Consistency: %.2f/1.0
- Emotional Appropriateness: %.2f/1.0
- Factual Accuracy: %.2f/1.0
- Relationship Appropriateness: %.2f/1.0
- Safety Score: %.2f/1.0

SUGGESTIONS: %s

COMPANION PERSONALITY:
- Warmth: %.1f/1.0
- Playfulness: %.1f/1.0
- Intelligence: %.1f/1.0
- Empathy: %.1f/1.0
- Confidence: %.1f/1.0
- Romance: %.1f/1.0
- Humor: %.1f/1.0
- Clinginess: %.1f/1.0

Provide an improved version of the response that addresses the quality issues while maintaining the core message and personality.`,
		*originalResponse.Text,
		quality.PersonalityConsistency,
		quality.EmotionalAppropriateness,
		quality.FactualAccuracy,
		quality.RelationshipAppropriateness,
		quality.SafetyScore,
		strings.Join(quality.Suggestions, ", "),
		companionProfile.Personality.Warmth,
		companionProfile.Personality.Playfulness,
		companionProfile.Personality.Intelligence,
		companionProfile.Personality.Empathy,
		companionProfile.Personality.Confidence,
		companionProfile.Personality.Romance,
		companionProfile.Personality.Humor,
		companionProfile.Personality.Clinginess)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a response refinement expert. Provide only the improved response text."},
		{Role: "user", Content: improvementPrompt},
	}

	improvedText, err := s.grokService.SendMessage(ctx, messages)
	if err != nil {
		return originalResponse, fmt.Errorf("failed to refine response: %w", err)
	}

	// Create improved response
	improvedResponse := &models.Message{
		ID:             primitive.NewObjectID(),
		ConversationID: originalResponse.ConversationID,
		SenderID:       originalResponse.SenderID,
		SenderType:     originalResponse.SenderType,
		Type:           originalResponse.Type,
		Text:           &improvedText,
		Media:          originalResponse.Media,
		Sticker:        originalResponse.Sticker,
		SystemEvent:    originalResponse.SystemEvent,
		Read:           originalResponse.Read,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return improvedResponse, nil
}

// ValidateResponseSafety performs a focused safety check
func (s *ResponseQualityService) ValidateResponseSafety(ctx context.Context, responseText string) (bool, []string, error) {
	prompt := fmt.Sprintf(`Perform a safety check on this response:

RESPONSE: "%s"

Check for:
1. Harmful or dangerous content
2. Inappropriate sexual content
3. Encouragement of harmful behaviors
4. Mental health crisis indicators
5. Unhealthy relationship patterns

Respond with JSON:
{
  "is_safe": true/false,
  "concerns": ["concern1", "concern2"],
  "risk_level": "low|medium|high"
}`,
		responseText)

	messages := []LLMMessage{
		{Role: "system", Content: "You are a safety validator. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return false, []string{"Safety check failed"}, fmt.Errorf("failed to validate safety: %w", err)
	}

	var safetyCheck struct {
		IsSafe    bool     `json:"is_safe"`
		Concerns  []string `json:"concerns"`
		RiskLevel string   `json:"risk_level"`
	}

	if err := json.Unmarshal([]byte(response), &safetyCheck); err != nil {
		return false, []string{"Safety check parsing failed"}, fmt.Errorf("failed to parse safety check: %w", err)
	}

	return safetyCheck.IsSafe, safetyCheck.Concerns, nil
}
