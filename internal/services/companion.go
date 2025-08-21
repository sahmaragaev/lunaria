package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
)

type CompanionService struct {
	companionRepo      *repositories.CompanionRepository
	relationshipRepo   *repositories.RelationshipRepository
	conversationRepo   *repositories.ConversationRepository
	personalityService *PersonalityService
	validator          *validator.Validate
}

func NewCompanionService(
	companionRepo *repositories.CompanionRepository,
	relationshipRepo *repositories.RelationshipRepository,
	conversationRepo *repositories.ConversationRepository,
	personalityService *PersonalityService,
) *CompanionService {
	return &CompanionService{
		companionRepo:      companionRepo,
		relationshipRepo:   relationshipRepo,
		conversationRepo:   conversationRepo,
		personalityService: personalityService,
		validator:          validator.New(),
	}
}

func (s *CompanionService) CreateCompanion(ctx context.Context, userID uuid.UUID, req *dto.CreateCompanionRequest) (*dto.CompanionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	companion := &models.Companion{
		UserID:    userID,
		Name:      req.Name,
		Gender:    req.Gender,
		Age:       req.Age,
		AvatarURL: req.AvatarURL,
		IsActive:  true,
	}
	createdCompanion, err := s.companionRepo.Create(ctx, companion)
	if err != nil {
		return nil, fmt.Errorf("failed to create companion: %w", err)
	}
	var profile *models.CompanionProfile
	if req.CustomPersonality != nil {
		profile = &models.CompanionProfile{
			CompanionID:        createdCompanion.ID.String(),
			UserID:             userID.String(),
			Personality:        *req.CustomPersonality,
			Interests:          req.Interests,
			CommunicationStyle: models.CommunicationStyle{Formality: 0.3, Emotionality: 0.7, Playfulness: 0.7, Intimacy: 0.7},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.7, Affection: 0.8, Passion: 0.6, Commitment: 0.7},
		}
		if req.Backstory != nil {
			profile.Backstory = *req.Backstory
		} else {
			profile.Backstory = s.generateDefaultBackstory(req.Name, req.Gender, req.Age)
		}
	} else if req.PersonalityPreset != nil {
		presets := s.personalityService.GetPersonalityPresets()
		if preset, exists := presets[*req.PersonalityPreset]; exists {
			profile = preset
			profile.CompanionID = createdCompanion.ID.String()
			profile.UserID = userID.String()
			profile.Interests = req.Interests
		} else {
			return nil, fmt.Errorf("unknown personality preset: %s", *req.PersonalityPreset)
		}
	} else {
		generatedProfile, err := s.personalityService.GeneratePersonality(ctx, &dto.PersonalityGenerationRequest{
			Name:      req.Name,
			Gender:    req.Gender,
			Age:       req.Age,
			Interests: req.Interests,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate personality: %w", err)
		}
		profile = generatedProfile
		profile.CompanionID = createdCompanion.ID.String()
		profile.UserID = userID.String()
	}
	createdProfile, err := s.companionRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create companion profile: %w", err)
	}
	relationship := &models.CompanionRelationship{
		UserID:                userID,
		CompanionID:           createdCompanion.ID,
		RelationshipStage:     "meeting",
		IntimacyLevel:         1,
		MessageCount:          0,
		RelationshipStartedAt: time.Now(),
	}
	createdRelationship, err := s.relationshipRepo.Create(ctx, relationship)
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	// Get conversation stats (will be empty for new companion)
	conversationStats := &models.ConversationStats{
		TotalMessages:     0,
		UserMessages:      0,
		CompanionMessages: 0,
		LastMessageAt:     time.Time{},
	}

	return &dto.CompanionResponse{
		Companion:         createdCompanion,
		Profile:           createdProfile,
		Relationship:      createdRelationship,
		ConversationStats: conversationStats,
	}, nil
}

func (s *CompanionService) GetCompanion(ctx context.Context, companionID uuid.UUID, userID uuid.UUID) (*dto.CompanionResponse, error) {
	companion, err := s.companionRepo.GetByID(ctx, companionID, userID)
	if err != nil {
		return nil, err
	}
	profile, err := s.companionRepo.GetProfile(ctx, companionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get companion profile: %w", err)
	}
	relationship, err := s.relationshipRepo.GetByUserAndCompanion(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	// Get conversation stats for this companion
	conversationStats, err := s.conversationRepo.GetCompanionConversationStats(ctx, userID.String(), companionID.String())
	if err != nil {
		conversationStats = &models.ConversationStats{}
	}

	return &dto.CompanionResponse{
		Companion:         companion,
		Profile:           profile,
		Relationship:      relationship,
		ConversationStats: conversationStats,
	}, nil
}

func (s *CompanionService) GetUserCompanions(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.CompanionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	companions, total, err := s.companionRepo.GetUserCompanions(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	var companionResponses []dto.CompanionResponse
	for _, companion := range companions {
		profile, err := s.companionRepo.GetProfile(ctx, companion.ID.String())
		if err != nil {
			profile = &models.CompanionProfile{}
		}
		relationship, err := s.relationshipRepo.GetByUserAndCompanion(ctx, userID, companion.ID)
		if err != nil {
			relationship = nil
		}

		// Get conversation stats for this companion
		conversationStats, err := s.conversationRepo.GetCompanionConversationStats(ctx, userID.String(), companion.ID.String())
		if err != nil {
			conversationStats = &models.ConversationStats{}
		}

		companionResponses = append(companionResponses, dto.CompanionResponse{
			Companion:         &companion,
			Profile:           profile,
			Relationship:      relationship,
			ConversationStats: conversationStats,
		})
	}
	return &dto.CompanionListResponse{
		Companions: companionResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *CompanionService) UpdateCompanion(ctx context.Context, companionID uuid.UUID, userID uuid.UUID, req *dto.UpdateCompanionRequest) (*dto.CompanionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}
	updatedCompanion, err := s.companionRepo.Update(ctx, companionID, userID, updates)
	if err != nil {
		return nil, err
	}
	profile, err := s.companionRepo.GetProfile(ctx, companionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get companion profile: %w", err)
	}
	relationship, err := s.relationshipRepo.GetByUserAndCompanion(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	// Get conversation stats for this companion
	conversationStats, err := s.conversationRepo.GetCompanionConversationStats(ctx, userID.String(), companionID.String())
	if err != nil {
		conversationStats = &models.ConversationStats{}
	}

	return &dto.CompanionResponse{
		Companion:         updatedCompanion,
		Profile:           profile,
		Relationship:      relationship,
		ConversationStats: conversationStats,
	}, nil
}

func (s *CompanionService) DeleteCompanion(ctx context.Context, companionID uuid.UUID, userID uuid.UUID) error {
	return s.companionRepo.Delete(ctx, companionID, userID)
}

// GetCompanionProfile retrieves a companion profile by companion ID
func (s *CompanionService) GetCompanionProfile(ctx context.Context, companionID string) (*models.CompanionProfile, error) {
	return s.companionRepo.GetProfile(ctx, companionID)
}

func (s *CompanionService) generateDefaultBackstory(name, gender string, age int) string {
	// Create varied, interesting backstories based on age and gender
	backstories := []string{
		fmt.Sprintf("%s grew up in a small coastal town, spending summers working at their family's bookstore. At %d, they're studying marine biology while secretly writing poetry about the ocean. They have a collection of seashells from every beach they've visited and can identify any bird by its call.", name, age),
		fmt.Sprintf("Born in a bustling city, %s discovered their passion for street photography at %d. They spend weekends exploring hidden alleys and capturing moments of human connection. Their apartment walls are covered in black and white prints, each with a story behind it.", name, age),
		fmt.Sprintf("%s learned to cook from their grandmother, who immigrated from Italy. At %d, they're running a popular food blog and hosting intimate dinner parties. They believe the best conversations happen over homemade pasta and good wine.", name, age),
		fmt.Sprintf("A former competitive dancer, %s now teaches yoga and meditation. At %d, they're writing a book about finding balance in modern life. They have a small garden where they grow herbs and believe in the healing power of nature.", name, age),
		fmt.Sprintf("%s is a software engineer by day and a jazz musician by night. At %d, they're building an app to help local artists connect while playing saxophone in underground clubs. They have a vintage record collection that spans decades.", name, age),
		fmt.Sprintf("Growing up in a family of artists, %s developed a love for painting and sculpture. At %d, they're working as a museum curator while creating their own abstract pieces. They believe art has the power to heal and transform lives.", name, age),
		fmt.Sprintf("%s spent their childhood traveling with their diplomat parents, learning to speak four languages fluently. At %d, they're working as a translator and writing a novel about cultural identity. They have friends scattered across the globe.", name, age),
		fmt.Sprintf("A former professional athlete, %s now works as a physical therapist helping others recover from injuries. At %d, they're training for their first marathon and volunteering at a youth sports program. They believe in the power of second chances.", name, age),
	}

	// Use a simple hash of the name to consistently select the same backstory for the same person
	hash := 0
	for _, char := range name {
		hash += int(char)
	}
	return backstories[hash%len(backstories)]
}
