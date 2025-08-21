package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	repo                     *repositories.ConversationRepository
	analytics                *repositories.AnalyticsRepository
	grok                     *GrokService
	aiContext                *AIContextService
	responseQuality          *ResponseQualityService
	conversationIntelligence *ConversationIntelligenceService
}

func NewMessageService(repo *repositories.ConversationRepository, analytics *repositories.AnalyticsRepository, grok *GrokService, aiContext *AIContextService, responseQuality *ResponseQualityService, conversationIntelligence *ConversationIntelligenceService) *MessageService {
	return &MessageService{
		repo:                     repo,
		analytics:                analytics,
		grok:                     grok,
		aiContext:                aiContext,
		responseQuality:          responseQuality,
		conversationIntelligence: conversationIntelligence,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	if err := s.validateMessage(msg); err != nil {
		return nil, err
	}

	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	storedMsg, err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	analytics := &models.MessageAnalytics{
		ID:             uuid.New(),
		ConversationID: msg.ConversationID.Hex(),
		SenderID:       uuid.MustParse(msg.SenderID),
		Type:           string(msg.Type),
		CreatedAt:      msg.CreatedAt,
	}
	s.analytics.InsertMessageAnalytics(ctx, analytics)

	return storedMsg, nil
}

func (s *MessageService) validateMessage(msg *models.Message) error {
	switch msg.Type {
	case "text":
		if msg.Text == nil || len(*msg.Text) == 0 {
			return fmt.Errorf("text message cannot be empty")
		}
	case "photo", "voice":
		if msg.Media == nil {
			return fmt.Errorf("media required for %s message", msg.Type)
		}
	case "sticker":
		if msg.Sticker == nil {
			return fmt.Errorf("sticker required")
		}
	case "system":
		if msg.SystemEvent == nil {
			return fmt.Errorf("system event required")
		}
	default:
		return fmt.Errorf("unsupported message type")
	}
	return nil
}

func (s *MessageService) GenerateAIResponse(ctx context.Context, conversation *models.Conversation, userMsg *models.Message, companionProfile *models.CompanionProfile) (*models.Message, error) {
	// Get conversation context and build dynamic prompt
	dynamicPrompt, err := s.aiContext.BuildDynamicPrompt(ctx, conversation, userMsg, companionProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to build dynamic prompt: %w", err)
	}

	// Get recent messages for context
	msgs, _, _, err := s.repo.ListMessages(ctx, conversation.ID, 10, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	fmt.Printf("DEBUG: Retrieved %d recent messages for conversation %s\n", len(msgs), conversation.ID.Hex())

	// Build conversation history for AI
	llmMessages := s.buildConversationHistory(msgs, userMsg)

	// Add dynamic system prompt plus an additional style directive to reduce clichés/idioms
	styleDirective := "Write in a natural, down-to-earth tone. Avoid clichés and idioms. Keep sentences concise, warm, and conversational. Speak like a real person."
	llmMessages = append([]LLMMessage{{Role: "system", Content: dynamicPrompt}, {Role: "system", Content: styleDirective}}, llmMessages...)

	fmt.Printf("DEBUG: Sending %d messages to AI (including system prompt)\n", len(llmMessages))

	// Signal typing start immediately
	GetTypingTracker().SetStart(conversation.ID.Hex())

	// Generate multiple AI responses with realistic delays
	aiResponses, err := s.generateMultipleAIResponses(ctx, llmMessages, conversation, companionProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI responses: %w", err)
	}

	// Inform tracker about total messages
	GetTypingTracker().SetTotal(conversation.ID.Hex(), len(aiResponses))

	// Store all responses in database
	var finalResponse *models.Message
	for i, aiText := range aiResponses {
		// Create AI response message
		aiResponse := &models.Message{
			ConversationID: userMsg.ConversationID,
			SenderID:       conversation.CompanionID,
			SenderType:     "companion",
			Type:           "text",
			Text:           &aiText,
			Read:           false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			IsTyping:       i < len(aiResponses)-1, // Mark as typing except for the last message
			MessageIndex:   i,                      // Track message order
			TotalMessages:  len(aiResponses),       // Total messages in this response
		}

		// Update typing tracker for each chunk
		GetTypingTracker().Update(conversation.ID.Hex(), i, len(aiResponses))

		// Store the response
		storedResponse, err := s.repo.CreateMessage(ctx, aiResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to store AI response: %w", err)
		}

		// Add realistic delay between messages (except for the last one)
		if i < len(aiResponses)-1 {
			delay := s.calculateTypingDelay(aiText, companionProfile)
			time.Sleep(delay)
		}

		finalResponse = storedResponse
	}

	// Extract and store memories from the conversation
	go func() {
		// Run memory extraction in background
		allMessages := append(msgs, userMsg)
		for _, response := range aiResponses {
			msg := &models.Message{
				ConversationID: userMsg.ConversationID,
				SenderID:       conversation.CompanionID,
				SenderType:     "companion",
				Type:           "text",
				Text:           &response,
			}
			allMessages = append(allMessages, msg)
		}
		if err := s.aiContext.ExtractAndStoreMemory(context.Background(), conversation.ID, allMessages); err != nil {
			fmt.Printf("Memory extraction failed: %v\n", err)
		}
	}()

	// Update conversation intelligence in background
	go func() {
		if _, err := s.conversationIntelligence.AnalyzeConversationFlow(context.Background(), conversation.ID); err != nil {
			fmt.Printf("Conversation intelligence update failed: %v\n", err)
		}
	}()

	return finalResponse, nil
}

// buildConversationHistory builds the conversation history for AI context
func (s *MessageService) buildConversationHistory(messages []*models.Message, userMsg *models.Message) []LLMMessage {
	var llmMessages []LLMMessage

	// Add recent messages in reverse chronological order
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		if m.Text != nil {
			role := "user"
			if m.SenderType == "companion" {
				role = "assistant"
			}
			llmMessages = append(llmMessages, LLMMessage{Role: role, Content: *m.Text})
		}
	}

	// Check if the current user message is already in the messages array
	// to avoid duplication
	if userMsg.Text != nil {
		messageAlreadyIncluded := false
		for _, m := range messages {
			if m.ID == userMsg.ID && m.Text != nil && *m.Text == *userMsg.Text {
				messageAlreadyIncluded = true
				fmt.Printf("DEBUG: User message already included in conversation history, skipping duplication. Message: %s\n", *userMsg.Text)
				break
			}
		}

		// Only add the current user message if it's not already included
		if !messageAlreadyIncluded {
			llmMessages = append(llmMessages, LLMMessage{Role: "user", Content: *userMsg.Text})
		}
	}

	return llmMessages
}

func (s *MessageService) ListMessages(ctx context.Context, conversationID primitive.ObjectID, limit int, cursor *primitive.ObjectID) ([]*models.Message, *primitive.ObjectID, bool, error) {
	return s.repo.ListMessages(ctx, conversationID, limit, cursor)
}

func (s *MessageService) GetMessageByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error) {
	return s.repo.GetMessageByID(ctx, id)
}

func (s *MessageService) GetMediaByID(ctx context.Context, id primitive.ObjectID) (*models.MediaMetadata, error) {
	return s.repo.GetMediaMetadataByID(ctx, id)
}

func (s *MessageService) MarkMessageAsRead(ctx context.Context, id primitive.ObjectID) error {
	msg, err := s.repo.GetMessageByID(ctx, id)
	if err != nil {
		return err
	}
	msg.Read = true
	msg.UpdatedAt = time.Now()

	return s.repo.UpdateMessage(ctx, msg)
}

// GetConversationIntelligence retrieves conversation intelligence insights
func (s *MessageService) GetConversationIntelligence(ctx context.Context, conversationID primitive.ObjectID) (*models.ConversationIntelligence, error) {
	return s.conversationIntelligence.AnalyzeConversationFlow(ctx, conversationID)
}

// SuggestNextTopic suggests the next conversation topic
func (s *MessageService) SuggestNextTopic(ctx context.Context, conversationID primitive.ObjectID, companionProfile *models.CompanionProfile) (string, error) {
	return s.conversationIntelligence.SuggestNextTopic(ctx, conversationID, companionProfile)
}

// AnalyzeEngagement analyzes user engagement in the conversation
func (s *MessageService) AnalyzeEngagement(ctx context.Context, conversationID primitive.ObjectID) (float64, error) {
	return s.conversationIntelligence.AnalyzeEngagementLevel(ctx, conversationID)
}

// GetResponseQuality retrieves quality metrics for a specific response
func (s *MessageService) GetResponseQuality(ctx context.Context, messageID primitive.ObjectID, conversation *models.Conversation, companionProfile *models.CompanionProfile) (*models.ResponseQuality, error) {
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return s.responseQuality.ValidateResponseQuality(ctx, message, conversation, companionProfile)
}

// generateMultipleAIResponses generates multiple AI responses that simulate natural conversation flow
func (s *MessageService) generateMultipleAIResponses(ctx context.Context, llmMessages []LLMMessage, conversation *models.Conversation, companionProfile *models.CompanionProfile) ([]string, error) {
	// Generate the full response first
	fullResponse, err := s.grok.SendMessage(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Split the response into multiple messages based on natural breaks
	messages := s.splitResponseIntoMessages(fullResponse, companionProfile)

	return messages, nil
}

// splitResponseIntoMessages splits a long response into multiple natural messages
func (s *MessageService) splitResponseIntoMessages(response string, companionProfile *models.CompanionProfile) []string {
	var messages []string

	// Split by sentences and group them naturally
	sentences := s.splitIntoSentences(response)

	var currentMessage strings.Builder
	messageCount := 0
	maxMessages := 5 // Maximum number of messages to split into

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// Check if adding this sentence would make the message too long
		if currentMessage.Len()+len(sentence) > 150 && currentMessage.Len() > 0 {
			// Start a new message
			if messageCount < maxMessages-1 {
				messages = append(messages, currentMessage.String())
				currentMessage.Reset()
				messageCount++
			}
		}

		if currentMessage.Len() > 0 {
			currentMessage.WriteString(" ")
		}
		currentMessage.WriteString(sentence)
	}

	// Add the last message
	if currentMessage.Len() > 0 {
		messages = append(messages, currentMessage.String())
	}

	// If we only have one message, try to split it further
	if len(messages) == 1 && len(response) > 200 {
		messages = s.splitLongMessage(response)
	}

	return messages
}

// splitIntoSentences splits text into sentences
func (s *MessageService) splitIntoSentences(text string) []string {
	// Improved splitting: keep ., !, ? and split accordingly
	var result []string
	var current strings.Builder
	for _, r := range text {
		current.WriteRune(r)
		if r == '.' || r == '!' || r == '?' {
			snt := strings.TrimSpace(current.String())
			if snt != "" {
				result = append(result, snt)
			}
			current.Reset()
		}
	}
	tail := strings.TrimSpace(current.String())
	if tail != "" {
		result = append(result, tail)
	}
	return result
}

// splitLongMessage splits a very long message into smaller chunks
func (s *MessageService) splitLongMessage(text string) []string {
	var messages []string
	words := strings.Fields(text)

	var currentMessage strings.Builder
	wordCount := 0
	maxWordsPerMessage := 30

	for _, word := range words {
		if wordCount >= maxWordsPerMessage && currentMessage.Len() > 0 {
			messages = append(messages, currentMessage.String())
			currentMessage.Reset()
			wordCount = 0
		}

		if currentMessage.Len() > 0 {
			currentMessage.WriteString(" ")
		}
		currentMessage.WriteString(word)
		wordCount++
	}

	if currentMessage.Len() > 0 {
		messages = append(messages, currentMessage.String())
	}

	return messages
}

// calculateTypingDelay calculates realistic typing delay based on message length and personality
func (s *MessageService) calculateTypingDelay(message string, companionProfile *models.CompanionProfile) time.Duration {
	words := len(strings.Fields(message))

	// Base pause grows with chunk length
	baseMs := 750 + int(float64(words)*90.0)

	// Punctuation-based thinking pauses
	punctuationBonus := 0
	trimmed := strings.TrimSpace(message)
	if strings.HasSuffix(trimmed, ".") || strings.HasSuffix(trimmed, "!") || strings.HasSuffix(trimmed, "?") {
		punctuationBonus += 250
	}
	if strings.Contains(trimmed, "\n") {
		punctuationBonus += 350
	}

	totalMs := baseMs + punctuationBonus

	// Personality adjustments
	multiplier := 1.0
	if companionProfile.Personality.Confidence > 0.7 {
		multiplier *= 0.9
	}
	if companionProfile.Personality.Intelligence > 0.8 {
		multiplier *= 0.95
	}

	// Randomness 0.9 - 1.3
	randomFactor := 0.9 + (rand.Float64() * 0.4)
	total := time.Duration(float64(totalMs)*multiplier*randomFactor) * time.Millisecond

	// Clamp to sensible bounds
	minDelay := 1 * time.Second
	maxDelay := 6 * time.Second
	if total < minDelay {
		total = minDelay
	}
	if total > maxDelay {
		total = maxDelay
	}

	// Small easing for very short chunks
	if words <= 4 {
		if total > 1200*time.Millisecond {
			total -= 400 * time.Millisecond
		}
	}

	return total
}
