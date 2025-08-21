package dto

import (
	"github.com/sahmaragaev/lunaria-backend/internal/models"
)

type CreateMessageRequest struct {
	Type        string              `json:"type" binding:"required,oneof=text photo voice sticker system"`
	Text        *string             `json:"text,omitempty"`
	MediaID     *string             `json:"media_id,omitempty"`
	Sticker     *models.StickerInfo `json:"sticker,omitempty"`
	SystemEvent *models.SystemEvent `json:"system_event,omitempty"`
}

type CreateMessageResponse struct {
	Message *models.Message `json:"message"`
}

type GetMessagesResponse struct {
	Messages   []*models.Message `json:"messages"`
	NextCursor *string           `json:"next_cursor,omitempty"`
	HasMore    bool              `json:"has_more"`
}

type PresignedURLRequest struct {
	Type   string `json:"type" binding:"required,oneof=photo voice"`
	Format string `json:"format" binding:"required"`
}

type PresignedURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
}

type ValidateMediaRequest struct {
	FileID string `json:"file_id" binding:"required"`
}

type ValidateMediaResponse struct {
	Media *models.MediaMetadata `json:"media"`
}
