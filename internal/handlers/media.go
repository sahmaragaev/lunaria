package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/response"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaHandler struct {
	mediaService *services.MediaService
}

func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func (h *MediaHandler) GenerateUploadURL(c *gin.Context) {
	var req dto.PresignedURLRequest
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
	url, fileID, err := h.mediaService.GeneratePresignedUploadURL(c.Request.Context(), user.ID.String(), req.Type, req.Format)
	if err != nil {
		response.InternalServerError(c, err, nil)
		return
	}
	response.Success(c, dto.PresignedURLResponse{UploadURL: url, FileID: fileID}, "Presigned URL generated")
}

func (h *MediaHandler) ValidateMedia(c *gin.Context) {
	var req dto.ValidateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, nil)
		return
	}
	mediaID, err := primitive.ObjectIDFromHex(req.FileID)
	if err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid file ID"})
		return
	}
	meta, err := h.mediaService.GetMediaMetadataByID(c.Request.Context(), mediaID)
	if err != nil {
		response.NotFound(c, err, gin.H{"error": "Media not found"})
		return
	}
	if meta.Status != "validated" {
		response.BadRequest(c, nil, gin.H{"error": "Media not validated yet"})
		return
	}
	response.Success(c, dto.ValidateMediaResponse{Media: meta}, "Media validated")
}

func (h *MediaHandler) GetMediaFile(c *gin.Context) {
	fileID := c.Param("file_id")
	userInterface, exists := c.Get("user")
	if !exists {
		response.Error(c, 401, nil, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(*models.User)
	mediaID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		response.BadRequest(c, err, gin.H{"error": "Invalid file ID"})
		return
	}
	meta, err := h.mediaService.GetMediaMetadataByID(c.Request.Context(), mediaID)
	if err != nil {
		response.NotFound(c, err, gin.H{"error": "Media not found"})
		return
	}
	if meta == nil || meta.UserID != user.ID.String() {
		response.Forbidden(c, nil, gin.H{"error": "Access denied"})
		return
	}
	response.Success(c, gin.H{"file_id": fileID, "media": meta}, "Media file access granted")
}
