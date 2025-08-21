package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
)

type GamificationService struct {
	analyticsRepo *repositories.AnalyticsRepository
	convRepo      *repositories.ConversationRepository
}

func NewGamificationService(analyticsRepo *repositories.AnalyticsRepository, convRepo *repositories.ConversationRepository) *GamificationService {
	return &GamificationService{
		analyticsRepo: analyticsRepo,
		convRepo:      convRepo,
	}
}

// InitializeAchievementDefinitions initializes the default achievement definitions
func (s *GamificationService) InitializeAchievementDefinitions(ctx context.Context) error {
	definitions := []models.AchievementDefinition{
		// Conversation Achievements
		{
			ID:          "first_conversation",
			Title:       "First Steps",
			Description: "Have your first conversation",
			Category:    "conversation",
			Type:        "milestone",
			Points:      50,
			Rarity:      "common",
			IconURL:     "/icons/first-conversation.png",
			Criteria: models.AchievementCriteria{
				Type:        "total_sessions",
				Target:      1,
				Measurement: "count",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "conversation_streak_7",
			Title:       "Week Warrior",
			Description: "Have conversations for 7 days in a row",
			Category:    "conversation",
			Type:        "streak",
			Points:      100,
			Rarity:      "rare",
			IconURL:     "/icons/week-warrior.png",
			Criteria: models.AchievementCriteria{
				Type:        "streak",
				Target:      7,
				Measurement: "days",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "conversation_streak_30",
			Title:       "Monthly Master",
			Description: "Have conversations for 30 days in a row",
			Category:    "conversation",
			Type:        "streak",
			Points:      500,
			Rarity:      "epic",
			IconURL:     "/icons/monthly-master.png",
			Criteria: models.AchievementCriteria{
				Type:        "streak",
				Target:      30,
				Measurement: "days",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "deep_conversation",
			Title:       "Deep Dive",
			Description: "Have a conversation lasting more than 30 minutes",
			Category:    "conversation",
			Type:        "quality",
			Points:      75,
			Rarity:      "common",
			IconURL:     "/icons/deep-dive.png",
			Criteria: models.AchievementCriteria{
				Type:        "session_duration",
				Target:      30,
				Measurement: "minutes",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "message_master",
			Title:       "Message Master",
			Description: "Send 100 messages",
			Category:    "conversation",
			Type:        "milestone",
			Points:      150,
			Rarity:      "common",
			IconURL:     "/icons/message-master.png",
			Criteria: models.AchievementCriteria{
				Type:        "total_messages",
				Target:      100,
				Measurement: "count",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},

		// Relationship Achievements
		{
			ID:          "trust_builder",
			Title:       "Trust Builder",
			Description: "Reach a high trust level in your relationship",
			Category:    "relationship",
			Type:        "milestone",
			Points:      200,
			Rarity:      "rare",
			IconURL:     "/icons/trust-builder.png",
			Criteria: models.AchievementCriteria{
				Type:        "trust_level",
				Target:      0.8,
				Measurement: "percentage",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "intimacy_explorer",
			Title:       "Intimacy Explorer",
			Description: "Reach a high intimacy level",
			Category:    "relationship",
			Type:        "milestone",
			Points:      300,
			Rarity:      "epic",
			IconURL:     "/icons/intimacy-explorer.png",
			Criteria: models.AchievementCriteria{
				Type:        "intimacy_level",
				Target:      0.8,
				Measurement: "percentage",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "vulnerability_champion",
			Title:       "Vulnerability Champion",
			Description: "Share vulnerable moments in conversations",
			Category:    "relationship",
			Type:        "milestone",
			Points:      250,
			Rarity:      "rare",
			IconURL:     "/icons/vulnerability-champion.png",
			Criteria: models.AchievementCriteria{
				Type:        "vulnerability_level",
				Target:      0.7,
				Measurement: "percentage",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},

		// Personal Growth Achievements
		{
			ID:          "emotional_processor",
			Title:       "Emotional Processor",
			Description: "Work through difficult emotions in conversations",
			Category:    "growth",
			Type:        "milestone",
			Points:      175,
			Rarity:      "rare",
			IconURL:     "/icons/emotional-processor.png",
			Criteria: models.AchievementCriteria{
				Type:        "emotional_processing",
				Target:      5,
				Measurement: "sessions",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "self_reflection_expert",
			Title:       "Self-Reflection Expert",
			Description: "Engage in deep self-reflection conversations",
			Category:    "growth",
			Type:        "milestone",
			Points:      200,
			Rarity:      "rare",
			IconURL:     "/icons/self-reflection-expert.png",
			Criteria: models.AchievementCriteria{
				Type:        "self_reflection",
				Target:      10,
				Measurement: "sessions",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "goal_setter",
			Title:       "Goal Setter",
			Description: "Set and discuss personal goals",
			Category:    "growth",
			Type:        "milestone",
			Points:      150,
			Rarity:      "common",
			IconURL:     "/icons/goal-setter.png",
			Criteria: models.AchievementCriteria{
				Type:        "goal_setting",
				Target:      3,
				Measurement: "goals",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},

		// Level Achievements
		{
			ID:          "level_5",
			Title:       "Getting Started",
			Description: "Reach level 5",
			Category:    "progression",
			Type:        "level",
			Points:      100,
			Rarity:      "common",
			IconURL:     "/icons/level-5.png",
			Criteria: models.AchievementCriteria{
				Type:        "level",
				Target:      5,
				Measurement: "level",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "level_10",
			Title:       "Dedicated Companion",
			Description: "Reach level 10",
			Category:    "progression",
			Type:        "level",
			Points:      250,
			Rarity:      "rare",
			IconURL:     "/icons/level-10.png",
			Criteria: models.AchievementCriteria{
				Type:        "level",
				Target:      10,
				Measurement: "level",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "level_25",
			Title:       "Relationship Expert",
			Description: "Reach level 25",
			Category:    "progression",
			Type:        "level",
			Points:      750,
			Rarity:      "epic",
			IconURL:     "/icons/level-25.png",
			Criteria: models.AchievementCriteria{
				Type:        "level",
				Target:      25,
				Measurement: "level",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "level_50",
			Title:       "Master Companion",
			Description: "Reach level 50",
			Category:    "progression",
			Type:        "level",
			Points:      2000,
			Rarity:      "legendary",
			IconURL:     "/icons/level-50.png",
			Criteria: models.AchievementCriteria{
				Type:        "level",
				Target:      50,
				Measurement: "level",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},

		// Special Event Achievements
		{
			ID:          "anniversary_1_month",
			Title:       "One Month Together",
			Description: "Celebrate one month of companionship",
			Category:    "special",
			Type:        "anniversary",
			Points:      300,
			Rarity:      "rare",
			IconURL:     "/icons/anniversary-1-month.png",
			Criteria: models.AchievementCriteria{
				Type:        "relationship_duration",
				Target:      30,
				Measurement: "days",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "anniversary_6_months",
			Title:       "Six Months Strong",
			Description: "Celebrate six months of companionship",
			Category:    "special",
			Type:        "anniversary",
			Points:      1000,
			Rarity:      "epic",
			IconURL:     "/icons/anniversary-6-months.png",
			Criteria: models.AchievementCriteria{
				Type:        "relationship_duration",
				Target:      180,
				Measurement: "days",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "anniversary_1_year",
			Title:       "One Year Together",
			Description: "Celebrate one year of companionship",
			Category:    "special",
			Type:        "anniversary",
			Points:      2500,
			Rarity:      "legendary",
			IconURL:     "/icons/anniversary-1-year.png",
			Criteria: models.AchievementCriteria{
				Type:        "relationship_duration",
				Target:      365,
				Measurement: "days",
			},
			Active:    true,
			CreatedAt: time.Now(),
		},
	}

	// Insert achievement definitions
	for _, definition := range definitions {
		err := s.insertAchievementDefinition(ctx, &definition)
		if err != nil {
			return fmt.Errorf("failed to insert achievement definition %s: %w", definition.ID, err)
		}
	}

	return nil
}

// insertAchievementDefinition inserts an achievement definition
func (s *GamificationService) insertAchievementDefinition(ctx context.Context, definition *models.AchievementDefinition) error {
	// Check if already exists
	existing, err := s.analyticsRepo.GetAchievementDefinition(ctx, definition.ID)
	if err == nil && existing != nil {
		// Already exists, skip
		return nil
	}

	// Set creation timestamp
	definition.CreatedAt = time.Now()
	definition.Active = true

	// Insert new definition using the repository method
	err = s.analyticsRepo.InsertAchievementDefinition(ctx, definition)
	if err != nil {
		return fmt.Errorf("failed to insert achievement definition: %w", err)
	}

	return nil
}

// CheckAndAwardAchievements checks for and awards achievements based on user activity
func (s *GamificationService) CheckAndAwardAchievements(ctx context.Context, userID, companionID string, activityData *ActivityData) error {
	// Get user progress
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// Get achievement definitions
	definitions, err := s.analyticsRepo.GetAchievementDefinitions(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get achievement definitions: %w", err)
	}

	// Check each achievement
	for _, definition := range definitions {
		// Skip if already earned
		earned, err := s.analyticsRepo.CheckAchievementEarned(ctx, userID, companionID, definition.ID)
		if err != nil {
			continue
		}
		if earned {
			continue
		}

		// Check if achievement criteria are met
		if s.checkAchievementCriteria(ctx, &definition, progress, activityData) {
			err = s.awardAchievement(ctx, userID, companionID, &definition, activityData)
			if err != nil {
				return fmt.Errorf("failed to award achievement %s: %w", definition.ID, err)
			}
		}
	}

	return nil
}

// ActivityData represents user activity data for achievement checking
type ActivityData struct {
	SessionDuration    time.Duration
	MessageCount       int
	ConversationDepth  float64
	EmotionalIntensity float64
	VulnerabilityLevel float64
	TrustLevel         float64
	IntimacyLevel      float64
	RelationshipAge    time.Duration
}

// checkAchievementCriteria checks if achievement criteria are met
func (s *GamificationService) checkAchievementCriteria(ctx context.Context, definition *models.AchievementDefinition, progress *models.UserProgress, activityData *ActivityData) bool {
	switch definition.Criteria.Type {
	case "total_sessions":
		return float64(progress.TotalConversations) >= definition.Criteria.Target
	case "total_messages":
		return float64(progress.TotalMessages) >= definition.Criteria.Target
	case "streak":
		return float64(progress.CurrentStreak) >= definition.Criteria.Target
	case "level":
		return float64(progress.CurrentLevel) >= definition.Criteria.Target
	case "session_duration":
		return activityData.SessionDuration.Minutes() >= definition.Criteria.Target
	case "trust_level":
		return activityData.TrustLevel >= definition.Criteria.Target
	case "intimacy_level":
		return activityData.IntimacyLevel >= definition.Criteria.Target
	case "vulnerability_level":
		return activityData.VulnerabilityLevel >= definition.Criteria.Target
	case "relationship_duration":
		return activityData.RelationshipAge.Hours()/24 >= definition.Criteria.Target
	default:
		return false
	}
}

// awardAchievement awards an achievement to a user
func (s *GamificationService) awardAchievement(ctx context.Context, userID, companionID string, definition *models.AchievementDefinition, activityData *ActivityData) error {
	achievement := &models.UserAchievement{
		UserID:          userID,
		CompanionID:     companionID,
		AchievementID:   definition.ID,
		AchievementType: definition.Type,
		Title:           definition.Title,
		Description:     definition.Description,
		IconURL:         definition.IconURL,
		Points:          definition.Points,
		Rarity:          definition.Rarity,
		EarnedAt:        time.Now(),
		Context: map[string]any{
			"session_duration":    activityData.SessionDuration.String(),
			"message_count":       activityData.MessageCount,
			"conversation_depth":  activityData.ConversationDepth,
			"emotional_intensity": activityData.EmotionalIntensity,
			"vulnerability_level": activityData.VulnerabilityLevel,
			"trust_level":         activityData.TrustLevel,
			"intimacy_level":      activityData.IntimacyLevel,
		},
	}

	// Save achievement
	err := s.analyticsRepo.InsertUserAchievement(ctx, achievement)
	if err != nil {
		return fmt.Errorf("failed to insert achievement: %w", err)
	}

	// Update user progress
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	progress.TotalAchievements++
	if definition.Rarity == "rare" || definition.Rarity == "epic" || definition.Rarity == "legendary" {
		progress.RareAchievements++
	}

	// Add bonus experience points
	progress.TotalExperience += definition.Points * 10

	// Recalculate level
	progress.CurrentLevel = s.calculateLevel(progress.TotalExperience)
	progress.LevelProgress = s.calculateLevelProgress(progress.TotalExperience)
	progress.ExperienceToNext = s.calculateExperienceToNext(progress.TotalExperience)

	// Save updated progress
	err = s.analyticsRepo.UpsertUserProgress(ctx, progress)
	if err != nil {
		return fmt.Errorf("failed to update user progress: %w", err)
	}

	return nil
}

// calculateLevel calculates user level based on experience
func (s *GamificationService) calculateLevel(experience int) int {
	level := int(float64(experience)/100.0) + 1
	if level < 1 {
		level = 1
	}
	return level
}

// calculateLevelProgress calculates progress within current level
func (s *GamificationService) calculateLevelProgress(experience int) float64 {
	currentLevel := s.calculateLevel(experience)
	experienceForCurrentLevel := (currentLevel - 1) * 100
	experienceForNextLevel := currentLevel * 100

	if experienceForNextLevel == experienceForCurrentLevel {
		return 1.0
	}

	progress := float64(experience-experienceForCurrentLevel) / float64(experienceForNextLevel-experienceForCurrentLevel)
	if progress > 1.0 {
		progress = 1.0
	}
	return progress
}

// calculateExperienceToNext calculates experience needed for next level
func (s *GamificationService) calculateExperienceToNext(experience int) int {
	currentLevel := s.calculateLevel(experience)
	experienceForNextLevel := currentLevel * 100
	return experienceForNextLevel - experience
}

// GetUserAchievements gets achievements for a user
func (s *GamificationService) GetUserAchievements(ctx context.Context, userID, companionID string, limit int) ([]models.UserAchievement, error) {
	return s.analyticsRepo.GetUserAchievements(ctx, userID, companionID, limit)
}

// GetAchievementProgress gets progress for all achievements
func (s *GamificationService) GetAchievementProgress(ctx context.Context, userID, companionID string) (map[string]float64, error) {
	// Get user progress
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	// Get achievement definitions
	definitions, err := s.analyticsRepo.GetAchievementDefinitions(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement definitions: %w", err)
	}

	// Calculate progress for each achievement
	achievementProgress := make(map[string]float64)
	for _, definition := range definitions {
		// Check if already earned
		earned, err := s.analyticsRepo.CheckAchievementEarned(ctx, userID, companionID, definition.ID)
		if err != nil {
			continue
		}
		if earned {
			achievementProgress[definition.ID] = 1.0
			continue
		}

		// Calculate progress based on criteria
		progressValue := s.calculateAchievementProgress(&definition, progress)
		achievementProgress[definition.ID] = progressValue
	}

	return achievementProgress, nil
}

// calculateAchievementProgress calculates progress for an achievement
func (s *GamificationService) calculateAchievementProgress(definition *models.AchievementDefinition, progress *models.UserProgress) float64 {
	var currentValue float64

	switch definition.Criteria.Type {
	case "total_sessions":
		currentValue = float64(progress.TotalConversations)
	case "total_messages":
		currentValue = float64(progress.TotalMessages)
	case "streak":
		currentValue = float64(progress.CurrentStreak)
	case "level":
		currentValue = float64(progress.CurrentLevel)
	default:
		currentValue = 0.0
	}

	if definition.Criteria.Target <= 0 {
		return 0.0
	}

	progressValue := currentValue / definition.Criteria.Target
	if progressValue > 1.0 {
		progressValue = 1.0
	}

	return progressValue
}

// GetStreakInformation gets detailed streak information
func (s *GamificationService) GetStreakInformation(ctx context.Context, userID, companionID string) (*models.StreakInformation, error) {
	return s.analyticsRepo.GetStreakInformation(ctx, userID, companionID)
}

// UpdateStreak updates user streak based on activity
func (s *GamificationService) UpdateStreak(ctx context.Context, userID, companionID string) error {
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	lastActivity := progress.LastActivityDate.Truncate(24 * time.Hour)

	if today.Equal(lastActivity) {
		// Already updated today
		return nil
	}

	if today.Sub(lastActivity) == 24*time.Hour {
		// Consecutive day
		progress.CurrentStreak++
		if progress.CurrentStreak > progress.LongestStreak {
			progress.LongestStreak = progress.CurrentStreak
		}
	} else if today.Sub(lastActivity) > 24*time.Hour {
		// Streak broken
		progress.CurrentStreak = 1
	}

	progress.LastActivityDate = today

	return s.analyticsRepo.UpsertUserProgress(ctx, progress)
}

// GetLevelRewards gets rewards for reaching a specific level
func (s *GamificationService) GetLevelRewards(level int) map[string]any {
	rewards := make(map[string]any)

	switch level {
	case 5:
		rewards["title"] = "Getting Started"
		rewards["description"] = "You're making great progress!"
		rewards["bonus_points"] = 100
	case 10:
		rewards["title"] = "Dedicated Companion"
		rewards["description"] = "You're building a strong connection!"
		rewards["bonus_points"] = 250
	case 25:
		rewards["title"] = "Relationship Expert"
		rewards["description"] = "You've mastered the art of companionship!"
		rewards["bonus_points"] = 750
	case 50:
		rewards["title"] = "Master Companion"
		rewards["description"] = "You're a true master of relationships!"
		rewards["bonus_points"] = 2000
	default:
		rewards["title"] = "Level Up!"
		rewards["description"] = "Congratulations on reaching level " + fmt.Sprintf("%d", level)
		rewards["bonus_points"] = level * 10
	}

	return rewards
}

// GetAchievementCategories gets all achievement categories
func (s *GamificationService) GetAchievementCategories(ctx context.Context) ([]string, error) {
	definitions, err := s.analyticsRepo.GetAchievementDefinitions(ctx, "")
	if err != nil {
		return nil, err
	}

	categories := make(map[string]bool)
	for _, definition := range definitions {
		categories[definition.Category] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}

	return result, nil
}

// GetAchievementsByCategory gets achievements by category
func (s *GamificationService) GetAchievementsByCategory(ctx context.Context, category string) ([]models.AchievementDefinition, error) {
	return s.analyticsRepo.GetAchievementDefinitions(ctx, category)
}

// GetUserProgress gets user progress
func (s *GamificationService) GetUserProgress(ctx context.Context, userID, companionID string) (*models.UserProgress, error) {
	return s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
}

// GetAchievementDefinitions gets achievement definitions
func (s *GamificationService) GetAchievementDefinitions(ctx context.Context, category string) ([]models.AchievementDefinition, error) {
	return s.analyticsRepo.GetAchievementDefinitions(ctx, category)
}
