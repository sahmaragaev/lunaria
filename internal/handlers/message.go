package handlers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/enums/messagetype"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHandler struct {
	service             *services.MessageService
	conversationService *services.ConversationService
	companionService    *services.CompanionService
	pendingResponses    map[string]*time.Timer
	responseMutex       sync.RWMutex
	generatingResponses map[string]bool
	aggregationStart    map[string]time.Time
	aggregationWindow   time.Duration
	aggregationMax      time.Duration
}

func NewMessageHandler(service *services.MessageService, conversationService *services.ConversationService, companionService *services.CompanionService) *MessageHandler {
	return &MessageHandler{
		service:             service,
		conversationService: conversationService,
		companionService:    companionService,
		pendingResponses:    make(map[string]*time.Timer),
		responseMutex:       sync.RWMutex{},
		generatingResponses: make(map[string]bool),
		aggregationStart:    make(map[string]time.Time),
		aggregationWindow:   2 * time.Second,
		aggregationMax:      6 * time.Second,
	}
}

func MessageFromDTO(req dto.CreateMessageRequest, convID primitive.ObjectID, userID string, media *models.MediaMetadata) *models.Message {
	return &models.Message{
		ConversationID: convID,
		SenderID:       userID,
		SenderType:     "user",
		Type:           messagetype.Type(req.Type),
		Text:           req.Text,
		Media:          media,
		Sticker:        req.Sticker,
		SystemEvent:    req.SystemEvent,
		Read:           false,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	convIDStr := c.Param("id")
	convID, _ := primitive.ObjectIDFromHex(convIDStr)

	var media *models.MediaMetadata

	if req.MediaID != nil {
		mediaID, _ := primitive.ObjectIDFromHex(*req.MediaID)
		media, _ = h.service.GetMediaByID(c.Request.Context(), mediaID)
	}
	msg := MessageFromDTO(req, convID, user.ID.String(), media)
	storedMsg, err := h.service.SendMessage(c.Request.Context(), msg)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	h.responseMutex.Lock()

	// Check if already generating a response for this conversation
	if h.generatingResponses[convIDStr] {
		h.responseMutex.Unlock()
		response.Created(c, storedMsg, "Message sent")
		return
	}

	// Stop any existing timer (we will reschedule below)
	if existingTimer, exists := h.pendingResponses[convIDStr]; exists {
		existingTimer.Stop()
	}

	// Aggregation window to collect multiple user messages
	now := time.Now()
	if _, ok := h.aggregationStart[convIDStr]; !ok {
		h.aggregationStart[convIDStr] = now
	}
	elapsed := now.Sub(h.aggregationStart[convIDStr])
	delay := h.aggregationWindow
	if elapsed >= h.aggregationMax {
		delay = 200 * time.Millisecond
	}

	timer := time.AfterFunc(delay, func() {
		h.responseMutex.Lock()
		if h.generatingResponses[convIDStr] {
			h.responseMutex.Unlock()
			return
		}
		h.generatingResponses[convIDStr] = true
		h.responseMutex.Unlock()

		h.generateBotResponse(convID, storedMsg)

		h.responseMutex.Lock()
		delete(h.pendingResponses, convIDStr)
		delete(h.generatingResponses, convIDStr)
		delete(h.aggregationStart, convIDStr)
		h.responseMutex.Unlock()
	})

	h.pendingResponses[convIDStr] = timer
	h.responseMutex.Unlock()

	response.Created(c, storedMsg, "Message sent")
}

func (h *MessageHandler) generateBotResponse(convID primitive.ObjectID, userMsg *models.Message) {
	conversation, err := h.conversationService.GetConversation(context.Background(), convID)
	if err != nil {
		fmt.Printf("Failed to get conversation: %v\n", err)
		return
	}

	companionProfile, err := h.companionService.GetCompanionProfile(context.Background(), conversation.CompanionID)
	if err != nil {
		fmt.Printf("Failed to get companion profile: %v\n", err)
		return
	}

	botResponse, err := h.service.GenerateAIResponse(context.Background(), conversation, userMsg, companionProfile)
	if err != nil {
		fmt.Printf("Failed to generate AI response: %v\n", err)
		return
	}

	fmt.Printf("Bot response generated and stored: %s\n", botResponse.ID.Hex())
}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	convIDStr := c.Param("id")
	convID, _ := primitive.ObjectIDFromHex(convIDStr)
	msgs, _, _, err := h.service.ListMessages(c.Request.Context(), convID, 20, nil)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, msgs, "Messages listed")
}

func (h *MessageHandler) GetMessage(c *gin.Context) {
	msgIDStr := c.Param("message_id")
	msgID, _ := primitive.ObjectIDFromHex(msgIDStr)
	msg, err := h.service.GetMessageByID(c.Request.Context(), msgID)
	if err != nil {
		response.NotFound(c, err, nil)
		return
	}

	response.Success(c, msg, "Message details")
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	msgIDStr := c.Param("message_id")
	msgID, _ := primitive.ObjectIDFromHex(msgIDStr)
	if err := h.service.MarkMessageAsRead(c.Request.Context(), msgID); err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, nil, "Message marked as read")
}

// GetConversationIntelligence retrieves conversation intelligence insights
func (h *MessageHandler) GetConversationIntelligence(c *gin.Context) {
	convIDStr := c.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	intelligence, err := h.service.GetConversationIntelligence(c.Request.Context(), convID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, intelligence, "Conversation intelligence retrieved")
}

// SuggestNextTopic suggests the next conversation topic
func (h *MessageHandler) SuggestNextTopic(c *gin.Context) {
	convIDStr := c.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	// Get conversation to retrieve companion ID
	conversation, err := h.conversationService.GetConversation(c.Request.Context(), convID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	// Get companion profile using the companion ID from conversation
	companionProfile, err := h.companionService.GetCompanionProfile(c.Request.Context(), conversation.CompanionID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	topic, err := h.service.SuggestNextTopic(c.Request.Context(), convID, companionProfile)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, gin.H{"suggested_topic": topic}, "Topic suggestion generated")
}

// AnalyzeEngagement analyzes user engagement in the conversation
func (h *MessageHandler) AnalyzeEngagement(c *gin.Context) {
	convIDStr := c.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	engagement, err := h.service.AnalyzeEngagement(c.Request.Context(), convID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, gin.H{"engagement_level": engagement}, "Engagement analysis completed")
}

// GetResponseQuality retrieves quality metrics for a specific response
func (h *MessageHandler) GetResponseQuality(c *gin.Context) {
	msgIDStr := c.Param("message_id")
	msgID, err := primitive.ObjectIDFromHex(msgIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	convIDStr := c.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	// Get conversation to retrieve companion ID
	conversation, err := h.conversationService.GetConversation(c.Request.Context(), convID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	// Get companion profile using the companion ID from conversation
	companionProfile, err := h.companionService.GetCompanionProfile(c.Request.Context(), conversation.CompanionID)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	quality, err := h.service.GetResponseQuality(c.Request.Context(), msgID, conversation, companionProfile)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Success(c, quality, "Response quality analysis completed")
}

// CheckTypingStatus checks if the AI is still typing for a conversation
func (h *MessageHandler) CheckTypingStatus(c *gin.Context) {
	convIDStr := c.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		response.BadRequest(c, err, nil)
		return
	}

	// Prefer in-memory live tracker for accuracy
	if state, ok := services.GetTypingTracker().Get(convIDStr); ok {
		response.Success(c, gin.H{
			"is_typing":       state.IsTyping,
			"message_index":   state.MessageIndex,
			"total_messages":  state.TotalMessages,
			"last_message_at": state.LastUpdate,
		}, "Typing status retrieved (live)")
		return
	}

	// Fallback to DB when tracker has no state
	messages, _, _, err := h.service.ListMessages(c.Request.Context(), convID, 10, nil)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	var latestCompanionMessage *models.Message
	for _, msg := range messages {
		if msg.SenderType == "companion" {
			if latestCompanionMessage == nil || msg.CreatedAt.After(latestCompanionMessage.CreatedAt) {
				latestCompanionMessage = msg
			}
		}
	}

	if latestCompanionMessage == nil {
		response.Success(c, gin.H{"is_typing": false, "message_index": 0, "total_messages": 0}, "No companion messages found")
		return
	}

	response.Success(c, gin.H{
		"is_typing":       latestCompanionMessage.IsTyping,
		"message_index":   latestCompanionMessage.MessageIndex,
		"total_messages":  latestCompanionMessage.TotalMessages,
		"last_message_at": latestCompanionMessage.CreatedAt,
	}, "Typing status retrieved")
}
