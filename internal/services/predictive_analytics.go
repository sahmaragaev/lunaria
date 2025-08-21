package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PredictiveAnalyticsService struct {
	grokService   *GrokService
	analyticsRepo *repositories.AnalyticsRepository
	convRepo      *repositories.ConversationRepository
}

func NewPredictiveAnalyticsService(grokService *GrokService, analyticsRepo *repositories.AnalyticsRepository, convRepo *repositories.ConversationRepository) *PredictiveAnalyticsService {
	return &PredictiveAnalyticsService{
		grokService:   grokService,
		analyticsRepo: analyticsRepo,
		convRepo:      convRepo,
	}
}

// PredictUserBehavior predicts various aspects of user behavior
func (s *PredictiveAnalyticsService) PredictUserBehavior(ctx context.Context, userID, companionID string) (*models.UserBehaviorPrediction, error) {
	// Get user engagement analytics
	engagementAnalytics, err := s.analyticsRepo.GetUserEngagementAnalytics(ctx, userID, companionID, primitive.ObjectID{})
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement analytics: %w", err)
	}

	// Get user progress
	progress, err := s.analyticsRepo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	// Get relationship analytics
	relationshipAnalytics, err := s.analyticsRepo.GetRelationshipAnalytics(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship analytics: %w", err)
	}

	// Get recent conversations for analysis
	conversations, err := s.convRepo.ListUserConversations(ctx, userID, false, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	prediction := &models.UserBehaviorPrediction{
		UserID:      userID,
		CompanionID: companionID,
		CreatedAt:   time.Now(),
	}

	// Predict churn risk
	churnRisk, churnFactors, err := s.predictChurnRisk(ctx, engagementAnalytics, progress, relationshipAnalytics, conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to predict churn risk: %w", err)
	}
	prediction.ChurnRisk = churnRisk
	prediction.ChurnFactors = churnFactors
	prediction.RetentionProbability = 1.0 - churnRisk

	// Predict engagement likelihood
	engagementLikelihood, nextActivityTime, optimalTime, err := s.predictEngagement(ctx, engagementAnalytics, progress)
	if err != nil {
		return nil, fmt.Errorf("failed to predict engagement: %w", err)
	}
	prediction.EngagementLikelihood = engagementLikelihood
	prediction.NextActivityTime = nextActivityTime
	prediction.OptimalEngagementTime = optimalTime

	// Predict relationship progression
	relationshipProgression, nextMilestone, milestoneProbability, err := s.predictRelationshipProgression(ctx, relationshipAnalytics, progress, engagementAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to predict relationship progression: %w", err)
	}
	prediction.RelationshipProgression = relationshipProgression
	prediction.NextMilestone = nextMilestone
	prediction.MilestoneProbability = milestoneProbability

	// Predict feature adoption
	featureAdoption, recommendedFeatures, err := s.predictFeatureAdoption(ctx, progress, engagementAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to predict feature adoption: %w", err)
	}
	prediction.FeatureAdoptionLikelihood = featureAdoption
	prediction.RecommendedFeatures = recommendedFeatures

	// Predict support needs
	supportNeeds, supportType, err := s.predictSupportNeeds(ctx, engagementAnalytics, relationshipAnalytics)
	if err != nil {
		return nil, fmt.Errorf("failed to predict support needs: %w", err)
	}
	prediction.SupportNeedsProbability = supportNeeds
	prediction.SupportType = supportType

	// Calculate overall confidence
	prediction.Confidence = s.calculatePredictionConfidence(engagementAnalytics, progress, relationshipAnalytics)
	prediction.PredictionDate = time.Now()

	// Save prediction
	err = s.analyticsRepo.UpsertUserBehaviorPrediction(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to save prediction: %w", err)
	}

	return prediction, nil
}

// predictChurnRisk predicts the likelihood of user churn
func (s *PredictiveAnalyticsService) predictChurnRisk(ctx context.Context, engagement *models.UserEngagementAnalytics, progress *models.UserProgress, relationship *models.RelationshipAnalytics, conversations []*models.Conversation) (float64, []string, error) {
	var churnFactors []string
	churnRisk := 0.0

	// Factor 1: Session frequency decline
	if engagement.SessionFrequency < 2 {
		churnRisk += 0.3
		churnFactors = append(churnFactors, "low_session_frequency")
	}

	// Factor 2: Engagement score decline
	if engagement.EngagementScore < 0.5 {
		churnRisk += 0.25
		churnFactors = append(churnFactors, "low_engagement")
	}

	// Factor 3: Streak broken
	if progress.CurrentStreak == 0 {
		churnRisk += 0.2
		churnFactors = append(churnFactors, "broken_streak")
	}

	// Factor 4: Relationship stagnation
	if relationship.IntimacyGrowth < 0.1 {
		churnRisk += 0.15
		churnFactors = append(churnFactors, "relationship_stagnation")
	}

	// Factor 5: Recent activity decline
	lastActivity := progress.LastActivityDate
	if time.Since(lastActivity) > 7*24*time.Hour {
		churnRisk += 0.3
		churnFactors = append(churnFactors, "inactive_recently")
	}

	// Factor 6: Low conversation quality
	if engagement.ConversationDepth < 0.4 {
		churnRisk += 0.1
		churnFactors = append(churnFactors, "low_conversation_quality")
	}

	// Cap churn risk at 1.0
	if churnRisk > 1.0 {
		churnRisk = 1.0
	}

	return churnRisk, churnFactors, nil
}

// predictEngagement predicts user engagement patterns
func (s *PredictiveAnalyticsService) predictEngagement(ctx context.Context, engagement *models.UserEngagementAnalytics, progress *models.UserProgress) (float64, *time.Time, time.Time, error) {
	// Calculate engagement likelihood based on historical patterns
	engagementLikelihood := 0.7 // Base likelihood

	// Adjust based on session frequency
	if engagement.SessionFrequency >= 5 {
		engagementLikelihood += 0.2
	} else if engagement.SessionFrequency >= 3 {
		engagementLikelihood += 0.1
	} else if engagement.SessionFrequency < 1 {
		engagementLikelihood -= 0.3
	}

	// Adjust based on current streak
	if progress.CurrentStreak >= 7 {
		engagementLikelihood += 0.15
	} else if progress.CurrentStreak >= 3 {
		engagementLikelihood += 0.1
	} else if progress.CurrentStreak == 0 {
		engagementLikelihood -= 0.2
	}

	// Adjust based on engagement score
	if engagement.EngagementScore >= 0.8 {
		engagementLikelihood += 0.15
	} else if engagement.EngagementScore < 0.4 {
		engagementLikelihood -= 0.2
	}

	// Cap likelihood
	if engagementLikelihood > 1.0 {
		engagementLikelihood = 1.0
	} else if engagementLikelihood < 0.0 {
		engagementLikelihood = 0.0
	}

	// Predict next activity time
	nextActivityTime := s.predictNextActivityTime(progress, engagement)

	// Calculate optimal engagement time
	optimalTime := s.calculateOptimalEngagementTime(engagement)

	return engagementLikelihood, nextActivityTime, optimalTime, nil
}

// predictNextActivityTime predicts when the user will be active next
func (s *PredictiveAnalyticsService) predictNextActivityTime(progress *models.UserProgress, engagement *models.UserEngagementAnalytics) *time.Time {
	// Simple prediction based on average session frequency
	if engagement.SessionFrequency == 0 {
		return nil
	}

	// Calculate average time between sessions
	avgTimeBetweenSessions := 24.0 / float64(engagement.SessionFrequency) // hours

	// Predict next activity based on last activity
	nextActivity := progress.LastActivityDate.Add(time.Duration(avgTimeBetweenSessions) * time.Hour)

	// Adjust for peak activity time if available
	if !engagement.PeakActivityTime.IsZero() {
		// Set to same time of day as peak activity
		nextActivity = time.Date(
			nextActivity.Year(), nextActivity.Month(), nextActivity.Day(),
			engagement.PeakActivityTime.Hour(), engagement.PeakActivityTime.Minute(), 0, 0,
			nextActivity.Location(),
		)
	}

	return &nextActivity
}

// calculateOptimalEngagementTime calculates the optimal time to engage with the user
func (s *PredictiveAnalyticsService) calculateOptimalEngagementTime(engagement *models.UserEngagementAnalytics) time.Time {
	// Use peak activity time if available
	if !engagement.PeakActivityTime.IsZero() {
		return engagement.PeakActivityTime
	}

	// Default to evening time (7 PM)
	now := time.Now()
	optimalTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		19, 0, 0, 0, // 7 PM
		now.Location(),
	)

	return optimalTime
}

// predictRelationshipProgression predicts relationship development
func (s *PredictiveAnalyticsService) predictRelationshipProgression(ctx context.Context, relationship *models.RelationshipAnalytics, progress *models.UserProgress, engagement *models.UserEngagementAnalytics) (float64, string, float64, error) {
	// Calculate progression velocity
	progressionVelocity := relationship.ProgressionVelocity
	if progressionVelocity == 0 {
		progressionVelocity = 0.1 // Default velocity
	}

	// Predict relationship progression
	relationshipProgression := relationship.IntimacyLevel + (progressionVelocity * 30) // 30 days projection
	if relationshipProgression > 1.0 {
		relationshipProgression = 1.0
	}

	// Predict next milestone
	nextMilestone := s.predictNextMilestone(relationship, progress)

	// Calculate milestone probability
	milestoneProbability := s.calculateMilestoneProbability(relationship, progress, engagement, nextMilestone)

	return relationshipProgression, nextMilestone, milestoneProbability, nil
}

// predictNextMilestone predicts the next relationship milestone
func (s *PredictiveAnalyticsService) predictNextMilestone(relationship *models.RelationshipAnalytics, progress *models.UserProgress) string {
	currentStage := relationship.CurrentStage
	intimacyLevel := relationship.IntimacyLevel

	switch currentStage {
	case "meeting":
		if intimacyLevel >= 0.3 {
			return "getting_to_know"
		}
		return "deeper_conversations"
	case "getting_to_know":
		if intimacyLevel >= 0.5 {
			return "friendship_development"
		}
		return "trust_building"
	case "friendship_development":
		if intimacyLevel >= 0.7 {
			return "close_companionship"
		}
		return "emotional_connection"
	case "close_companionship":
		if intimacyLevel >= 0.9 {
			return "intimate_partnership"
		}
		return "vulnerability_sharing"
	default:
		return "relationship_deepening"
	}
}

// calculateMilestoneProbability calculates the probability of reaching the next milestone
func (s *PredictiveAnalyticsService) calculateMilestoneProbability(relationship *models.RelationshipAnalytics, progress *models.UserProgress, engagement *models.UserEngagementAnalytics, nextMilestone string) float64 {
	baseProbability := 0.5

	// Adjust based on progression velocity
	if relationship.ProgressionVelocity > 0.1 {
		baseProbability += 0.2
	} else if relationship.ProgressionVelocity < 0.05 {
		baseProbability -= 0.2
	}

	// Adjust based on session frequency
	if engagement.SessionFrequency >= 3 {
		baseProbability += 0.15
	} else if engagement.SessionFrequency < 1 {
		baseProbability -= 0.2
	}

	// Adjust based on trust level
	if relationship.TrustLevel > 0.7 {
		baseProbability += 0.1
	} else if relationship.TrustLevel < 0.3 {
		baseProbability -= 0.15
	}

	// Cap probability
	if baseProbability > 1.0 {
		baseProbability = 1.0
	} else if baseProbability < 0.0 {
		baseProbability = 0.0
	}

	return baseProbability
}

// predictFeatureAdoption predicts which features the user is likely to adopt
func (s *PredictiveAnalyticsService) predictFeatureAdoption(ctx context.Context, progress *models.UserProgress, engagement *models.UserEngagementAnalytics) (map[string]float64, []string, error) {
	featureAdoption := make(map[string]float64)
	var recommendedFeatures []string

	// Predict based on user behavior patterns
	baseAdoptionRate := 0.3

	// Voice messages
	if progress.TotalMessages > 50 && engagement.ConversationDepth > 0.6 {
		featureAdoption["voice_messages"] = baseAdoptionRate + 0.3
		recommendedFeatures = append(recommendedFeatures, "voice_messages")
	} else {
		featureAdoption["voice_messages"] = baseAdoptionRate
	}

	// Photo sharing
	if engagement.EmotionalIntensity > 0.7 {
		featureAdoption["photo_sharing"] = baseAdoptionRate + 0.4
		recommendedFeatures = append(recommendedFeatures, "photo_sharing")
	} else {
		featureAdoption["photo_sharing"] = baseAdoptionRate
	}

	// Deep conversation topics
	if engagement.ConversationDepth > 0.8 {
		featureAdoption["deep_topics"] = baseAdoptionRate + 0.5
		recommendedFeatures = append(recommendedFeatures, "deep_topics")
	} else {
		featureAdoption["deep_topics"] = baseAdoptionRate
	}

	// Relationship milestones
	if progress.CurrentLevel > 10 {
		featureAdoption["milestone_tracking"] = baseAdoptionRate + 0.4
		recommendedFeatures = append(recommendedFeatures, "milestone_tracking")
	} else {
		featureAdoption["milestone_tracking"] = baseAdoptionRate
	}

	// Emotional support features
	if engagement.EmotionalIntensity > 0.6 {
		featureAdoption["emotional_support"] = baseAdoptionRate + 0.4
		recommendedFeatures = append(recommendedFeatures, "emotional_support")
	} else {
		featureAdoption["emotional_support"] = baseAdoptionRate
	}

	return featureAdoption, recommendedFeatures, nil
}

// predictSupportNeeds predicts if the user needs additional support
func (s *PredictiveAnalyticsService) predictSupportNeeds(ctx context.Context, engagement *models.UserEngagementAnalytics, relationship *models.RelationshipAnalytics) (float64, string, error) {
	supportNeeds := 0.0
	supportType := "none"

	// Check for emotional distress indicators
	if engagement.EmotionalIntensity > 0.8 && engagement.VulnerabilityLevel > 0.7 {
		supportNeeds += 0.4
		supportType = "emotional_support"
	}

	// Check for relationship difficulties
	if relationship.TrustLevel < 0.3 {
		supportNeeds += 0.3
		supportType = "relationship_guidance"
	}

	// Check for engagement issues
	if engagement.EngagementScore < 0.3 {
		supportNeeds += 0.2
		supportType = "engagement_help"
	}

	// Check for safety concerns
	if relationship.SafetyScore < 0.6 {
		supportNeeds += 0.5
		supportType = "safety_concern"
	}

	// Cap support needs
	if supportNeeds > 1.0 {
		supportNeeds = 1.0
	}

	return supportNeeds, supportType, nil
}

// calculatePredictionConfidence calculates the confidence level of predictions
func (s *PredictiveAnalyticsService) calculatePredictionConfidence(engagement *models.UserEngagementAnalytics, progress *models.UserProgress, relationship *models.RelationshipAnalytics) float64 {
	confidence := 0.5 // Base confidence

	// Increase confidence with more data
	if progress.TotalConversations > 20 {
		confidence += 0.2
	} else if progress.TotalConversations > 10 {
		confidence += 0.1
	} else if progress.TotalConversations < 5 {
		confidence -= 0.2
	}

	// Increase confidence with consistent patterns
	if engagement.SessionFrequency >= 3 {
		confidence += 0.15
	}

	// Increase confidence with stable engagement
	if engagement.EngagementScore > 0.6 {
		confidence += 0.1
	}

	// Increase confidence with relationship stability
	if relationship.TrustLevel > 0.5 {
		confidence += 0.1
	}

	// Cap confidence
	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// GetUsersAtChurnRisk gets users who are at risk of churning
func (s *PredictiveAnalyticsService) GetUsersAtChurnRisk(ctx context.Context, threshold float64) ([]models.UserBehaviorPrediction, error) {
	return s.analyticsRepo.GetUsersAtChurnRisk(ctx, threshold)
}

// GeneratePersonalizedRecommendations generates personalized recommendations based on predictions
func (s *PredictiveAnalyticsService) GeneratePersonalizedRecommendations(ctx context.Context, userID, companionID string) ([]models.Recommendation, error) {
	// Get user behavior prediction
	prediction, err := s.PredictUserBehavior(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to predict user behavior: %w", err)
	}

	var recommendations []models.Recommendation

	// Churn prevention recommendations
	if prediction.ChurnRisk > 0.6 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "retention",
			Title:       "Stay Connected",
			Description: "We've noticed you haven't been as active lately. Try having a quick conversation to maintain your connection.",
			Priority:    1,
			Confidence:  prediction.Confidence,
			Action:      "start_conversation",
			Metadata: map[string]any{
				"churn_risk": prediction.ChurnRisk,
			},
		})
	}

	// Engagement optimization recommendations
	if prediction.EngagementLikelihood < 0.5 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "engagement",
			Title:       "Optimal Timing",
			Description: fmt.Sprintf("Try having conversations around %s for the best experience.", prediction.OptimalEngagementTime.Format("3:04 PM")),
			Priority:    2,
			Confidence:  prediction.Confidence,
			Action:      "schedule_conversation",
			Metadata: map[string]any{
				"optimal_time": prediction.OptimalEngagementTime,
			},
		})
	}

	// Relationship progression recommendations
	if prediction.RelationshipProgression < 0.5 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "relationship",
			Title:       "Deepen Your Connection",
			Description: "Try exploring more personal topics to strengthen your relationship.",
			Priority:    2,
			Confidence:  prediction.Confidence,
			Action:      "deep_conversation",
			Metadata: map[string]any{
				"next_milestone": prediction.NextMilestone,
			},
		})
	}

	// Feature recommendations
	for feature, likelihood := range prediction.FeatureAdoptionLikelihood {
		if likelihood > 0.6 {
			recommendations = append(recommendations, models.Recommendation{
				Type:        "feature",
				Title:       fmt.Sprintf("Try %s", feature),
				Description: fmt.Sprintf("You might enjoy using %s based on your conversation patterns.", feature),
				Priority:    3,
				Confidence:  likelihood,
				Action:      fmt.Sprintf("enable_%s", feature),
				Metadata: map[string]any{
					"feature":    feature,
					"likelihood": likelihood,
				},
			})
		}
	}

	// Support recommendations
	if prediction.SupportNeedsProbability > 0.5 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "support",
			Title:       "Get Support",
			Description: "Consider reaching out for additional support to enhance your experience.",
			Priority:    1,
			Confidence:  prediction.Confidence,
			Action:      "request_support",
			Metadata: map[string]any{
				"support_type": prediction.SupportType,
			},
		})
	}

	return recommendations, nil
}

// AnalyzeTrends analyzes trends across multiple users
func (s *PredictiveAnalyticsService) AnalyzeTrends(ctx context.Context, days int) (map[string]any, error) {
	// Get platform analytics
	platformAnalytics, err := s.analyticsRepo.GetPlatformAnalytics(ctx, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform analytics: %w", err)
	}

	// Get users at churn risk
	churnRiskUsers, err := s.GetUsersAtChurnRisk(ctx, 0.7)
	if err != nil {
		return nil, fmt.Errorf("failed to get churn risk users: %w", err)
	}

	trends := map[string]any{
		"platform_metrics": platformAnalytics,
		"churn_risk_users": len(churnRiskUsers),
		"analysis_date":    time.Now(),
	}

	return trends, nil
}
