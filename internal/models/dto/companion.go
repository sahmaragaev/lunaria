package dto

import (
	"github.com/sahmaragaev/lunaria-backend/internal/models"
)

type CreateCompanionRequest struct {
	Name              string                    `json:"name" validate:"required,min=1,max=50"`
	Gender            string                    `json:"gender" validate:"required,oneof=male female other"`
	Age               int                       `json:"age" validate:"required,min=18,max=99"`
	AvatarURL         *string                   `json:"avatar_url,omitempty" validate:"omitempty,url"`
	PersonalityPreset *string                   `json:"personality_preset,omitempty"`
	CustomPersonality *models.PersonalityTraits `json:"custom_personality,omitempty"`
	Interests         []string                  `json:"interests,omitempty"`
	Backstory         *string                   `json:"backstory,omitempty"`
}

type UpdateCompanionRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Age       *int    `json:"age,omitempty" validate:"omitempty,min=18,max=99"`
}

type CompanionResponse struct {
	Companion         *models.Companion             `json:"companion"`
	Profile           *models.CompanionProfile      `json:"profile"`
	Relationship      *models.CompanionRelationship `json:"relationship,omitempty"`
	ConversationStats *models.ConversationStats     `json:"conversation_stats,omitempty"`
}

type CompanionListResponse struct {
	Companions []CompanionResponse `json:"companions"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
}

type PersonalityGenerationRequest struct {
	Name      string   `json:"name"`
	Gender    string   `json:"gender"`
	Age       int      `json:"age"`
	Interests []string `json:"interests,omitempty"`
	Style     string   `json:"style,omitempty"`
}
