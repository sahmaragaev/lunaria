package models

import (
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/enums/mediastatus"
	"github.com/sahmaragaev/lunaria-backend/internal/enums/mediatype"
	"github.com/sahmaragaev/lunaria-backend/internal/enums/messagetype"
	"github.com/sahmaragaev/lunaria-backend/internal/enums/sendertype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	CompanionID    string             `bson:"companion_id" json:"companion_id"`
	RecentMessages []Message          `bson:"recent_messages" json:"recent_messages"`
	Archived       bool               `bson:"archived" json:"archived"`
	Relationship   string             `bson:"relationship" json:"relationship"`
	LastActivity   time.Time          `bson:"last_activity" json:"last_activity"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`
	SenderID       string             `bson:"sender_id" json:"sender_id"`
	SenderType     sendertype.Type    `bson:"sender_type" json:"sender_type"` // user, companion, system
	Type           messagetype.Type   `bson:"type" json:"type"`               // text, photo, voice, sticker, system
	Text           *string            `bson:"text,omitempty" json:"text,omitempty"`
	Media          *MediaMetadata     `bson:"media,omitempty" json:"media,omitempty"`
	Sticker        *StickerInfo       `bson:"sticker,omitempty" json:"sticker,omitempty"`
	SystemEvent    *SystemEvent       `bson:"system_event,omitempty" json:"system_event,omitempty"`
	Read           bool               `bson:"read" json:"read"`
	IsTyping       bool               `bson:"is_typing" json:"is_typing"`           // Indicates if this message is part of a typing sequence
	MessageIndex   int                `bson:"message_index" json:"message_index"`   // Index of this message in a sequence (0-based)
	TotalMessages  int                `bson:"total_messages" json:"total_messages"` // Total number of messages in the sequence
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type MediaMetadata struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Type         mediatype.Type     `bson:"type" json:"type"` // photo, voice
	S3URL        string             `bson:"s3_url" json:"s3_url"`
	ThumbnailURL *string            `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"`
	Format       string             `bson:"format" json:"format"`
	Size         int64              `bson:"size" json:"size"`
	Width        *int               `bson:"width,omitempty" json:"width,omitempty"`
	Height       *int               `bson:"height,omitempty" json:"height,omitempty"`
	Duration     *float64           `bson:"duration,omitempty" json:"duration,omitempty"`
	Bitrate      *int               `bson:"bitrate,omitempty" json:"bitrate,omitempty"`
	EXIF         map[string]string  `bson:"exif,omitempty" json:"exif,omitempty"`
	Status       mediastatus.Type   `bson:"status" json:"status"` // pending, validated, rejected
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type StickerInfo struct {
	Pack string `bson:"pack" json:"pack"`
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
}

type SystemEvent struct {
	EventType string `bson:"event_type" json:"event_type"`
	Details   string `bson:"details" json:"details"`
}

// ConversationStats represents conversation statistics for a companion
type ConversationStats struct {
	TotalMessages     int       `json:"total_messages"`
	UserMessages      int       `json:"user_messages"`
	CompanionMessages int       `json:"companion_messages"`
	LastMessageAt     time.Time `json:"last_message_at"`
}
