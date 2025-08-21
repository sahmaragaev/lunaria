package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
)

type CompanionHandler struct {
	companionService *services.CompanionService
	validator        *validator.Validate
}

func NewCompanionHandler(companionService *services.CompanionService) *CompanionHandler {
	return &CompanionHandler{
		companionService: companionService,
		validator:        validator.New(),
	}
}

func (h *CompanionHandler) CreateCompanion(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	var req dto.CreateCompanionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}
	companion, err := h.companionService.CreateCompanion(c.Request.Context(), user.ID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation error") {
			response.BadRequest(c, err, nil)
			return
		}
		response.InternalServerError(c, err, gin.H{"error": "Failed to create companion"})
		return
	}
	response.Created(c, companion, "Companion created successfully")
}

func (h *CompanionHandler) GetCompanion(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	companionIDStr := c.Param("id")
	companionID, err := uuid.Parse(companionIDStr)
	if err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid companion ID"})
		return
	}
	companion, err := h.companionService.GetCompanion(c.Request.Context(), companionID, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err, nil)
			return
		}
		response.InternalServerError(c, err, gin.H{"error": "Failed to get companion"})
		return
	}
	response.Success(c, companion, "Companion retrieved successfully")
}

func (h *CompanionHandler) GetUserCompanions(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	page := 1
	pageSize := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}
	companions, err := h.companionService.GetUserCompanions(c.Request.Context(), user.ID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, err, gin.H{"error": "Failed to get companions"})
		return
	}
	response.Success(c, companions, "Companions retrieved successfully")
}

func (h *CompanionHandler) UpdateCompanion(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	companionIDStr := c.Param("id")
	companionID, err := uuid.Parse(companionIDStr)
	if err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid companion ID"})
		return
	}
	var req dto.UpdateCompanionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid request body"})
		return
	}
	companion, err := h.companionService.UpdateCompanion(c.Request.Context(), companionID, user.ID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err, nil)
			return
		}
		if strings.Contains(err.Error(), "validation error") {
			response.BadRequest(c, err, nil)
			return
		}
		response.InternalServerError(c, err, gin.H{"error": "Failed to update companion"})
		return
	}
	response.Success(c, companion, "Companion updated successfully")
}

func (h *CompanionHandler) DeleteCompanion(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	companionIDStr := c.Param("id")
	companionID, err := uuid.Parse(companionIDStr)
	if err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid companion ID"})
		return
	}
	err = h.companionService.DeleteCompanion(c.Request.Context(), companionID, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err, nil)
			return
		}
		response.InternalServerError(c, err, gin.H{"error": "Failed to delete companion"})
		return
	}
	response.Success(c, nil, "Companion deleted successfully")
}
