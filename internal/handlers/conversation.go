package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationHandler struct {
	service *services.ConversationService
}

func NewConversationHandler(service *services.ConversationService) *ConversationHandler {
	return &ConversationHandler{service: service}
}

func (h *ConversationHandler) StartConversation(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}

	user := userInterface.(*models.User)
	companionID := c.Query("companion_id")
	relationship := c.Query("relationship")

	conv, err := h.service.StartConversation(c.Request.Context(), user.ID.String(), companionID, relationship)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}

	response.Created(c, conv, "Conversation started")
}

func (h *ConversationHandler) ListConversations(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	archived := c.Query("archived") == "true"
	convs, err := h.service.ListConversations(c.Request.Context(), user.ID.String(), archived, 20, 0)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}
	response.Success(c, convs, "Conversations listed")
}

func (h *ConversationHandler) GetConversation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	conv, err := h.service.GetConversation(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err, nil)
		return
	}
	response.Success(c, conv, "Conversation details")
}

func (h *ConversationHandler) ArchiveConversation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	if err := h.service.ArchiveConversation(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, nil)
		return
	}
	response.Success(c, nil, "Conversation archived")
}

func (h *ConversationHandler) ReactivateConversation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	if err := h.service.ReactivateConversation(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, nil)
		return
	}
	response.Success(c, nil, "Conversation reactivated")
}
