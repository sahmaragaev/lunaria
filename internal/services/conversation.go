package services

import (
	"context"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationService struct {
	repo      *repositories.ConversationRepository
	analytics *repositories.AnalyticsRepository
}

func NewConversationService(repo *repositories.ConversationRepository, analytics *repositories.AnalyticsRepository) *ConversationService {
	return &ConversationService{repo: repo, analytics: analytics}
}

func (s *ConversationService) StartConversation(ctx context.Context, userID, companionID string, relationship string) (*models.Conversation, error) {
	conv := &models.Conversation{
		UserID:         userID,
		CompanionID:    companionID,
		Relationship:   relationship,
		RecentMessages: []models.Message{},
		Archived:       false,
		LastActivity:   time.Now(),
	}

	return s.repo.CreateConversation(ctx, conv)
}

func (s *ConversationService) ListConversations(ctx context.Context, userID string, archived bool, limit, offset int) ([]*models.Conversation, error) {
	return s.repo.ListUserConversations(ctx, userID, archived, limit, offset)
}

func (s *ConversationService) GetConversation(ctx context.Context, id primitive.ObjectID) (*models.Conversation, error) {
	return s.repo.GetConversationByID(ctx, id)
}

func (s *ConversationService) ArchiveConversation(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.ArchiveConversation(ctx, id)
}

func (s *ConversationService) ReactivateConversation(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.ReactivateConversation(ctx, id)
}
