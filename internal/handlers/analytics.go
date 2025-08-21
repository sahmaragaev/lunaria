package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsHandler struct {
	analyticsService           *services.AnalyticsService
	gamificationService        *services.GamificationService
	predictiveAnalyticsService *services.PredictiveAnalyticsService
}

func NewAnalyticsHandler(
	analyticsService *services.AnalyticsService,
	gamificationService *services.GamificationService,
	predictiveAnalyticsService *services.PredictiveAnalyticsService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService:           analyticsService,
		gamificationService:        gamificationService,
		predictiveAnalyticsService: predictiveAnalyticsService,
	}
}

// GetUserDashboard gets comprehensive dashboard data for a user
func (h *AnalyticsHandler) GetUserDashboard(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	dashboard, err := h.analyticsService.GetUserDashboardData(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetUserProgress gets user progress and gamification data
func (h *AnalyticsHandler) GetUserProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	progress, err := h.gamificationService.GetUserProgress(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetUserAchievements gets user achievements
func (h *AnalyticsHandler) GetUserAchievements(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	achievements, err := h.gamificationService.GetUserAchievements(c.Request.Context(), userID, companionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// GetAchievementProgress gets progress for all achievements
func (h *AnalyticsHandler) GetAchievementProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	progress, err := h.gamificationService.GetAchievementProgress(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get achievement progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetAchievementDefinitions gets available achievement definitions
func (h *AnalyticsHandler) GetAchievementDefinitions(c *gin.Context) {
	category := c.Query("category")

	definitions, err := h.gamificationService.GetAchievementDefinitions(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get achievement definitions"})
		return
	}

	c.JSON(http.StatusOK, definitions)
}

// GetAchievementCategories gets all achievement categories
func (h *AnalyticsHandler) GetAchievementCategories(c *gin.Context) {
	categories, err := h.gamificationService.GetAchievementCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get achievement categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetStreakInformation gets user streak information
func (h *AnalyticsHandler) GetStreakInformation(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	streakInfo, err := h.gamificationService.GetStreakInformation(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get streak information"})
		return
	}

	c.JSON(http.StatusOK, streakInfo)
}

// GetEngagementTrends gets user engagement trends
func (h *AnalyticsHandler) GetEngagementTrends(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trends, err := h.analyticsService.GetEngagementTrends(c.Request.Context(), userID, companionID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get engagement trends"})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetUserStatistics gets user statistics
func (h *AnalyticsHandler) GetUserStatistics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	statistics, err := h.analyticsService.GetUserStatistics(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user statistics"})
		return
	}

	c.JSON(http.StatusOK, statistics)
}

// GetRelationshipAnalytics gets relationship analytics
func (h *AnalyticsHandler) GetRelationshipAnalytics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	analytics, err := h.analyticsService.GetRelationshipAnalytics(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relationship analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetUserBehaviorPrediction gets user behavior prediction
func (h *AnalyticsHandler) GetUserBehaviorPrediction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	prediction, err := h.predictiveAnalyticsService.PredictUserBehavior(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get behavior prediction"})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// GetPersonalizedRecommendations gets personalized recommendations
func (h *AnalyticsHandler) GetPersonalizedRecommendations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	recommendations, err := h.predictiveAnalyticsService.GeneratePersonalizedRecommendations(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// TrackSessionActivity tracks user session activity
func (h *AnalyticsHandler) TrackSessionActivity(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request struct {
		CompanionID        string        `json:"companion_id" binding:"required"`
		ConversationID     string        `json:"conversation_id" binding:"required"`
		SessionDuration    time.Duration `json:"session_duration"`
		MessageCount       int           `json:"message_count"`
		ResponseQuality    float64       `json:"response_quality"`
		ConversationDepth  float64       `json:"conversation_depth"`
		EmotionalIntensity float64       `json:"emotional_intensity"`
		VulnerabilityLevel float64       `json:"vulnerability_level"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	conversationID, err := primitive.ObjectIDFromHex(request.ConversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Create session data
	sessionData := &services.SessionData{
		Duration:            request.SessionDuration,
		MessageCount:        request.MessageCount,
		ResponseQuality:     request.ResponseQuality,
		PeakActivityTime:    time.Now(),
		AverageResponseTime: 30 * time.Second, // Default value
	}

	// Track user engagement
	err = h.analyticsService.TrackUserEngagement(c.Request.Context(), userID, request.CompanionID, conversationID, sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track session activity"})
		return
	}

	// Process user progress
	err = h.analyticsService.ProcessUserProgress(c.Request.Context(), userID, request.CompanionID, sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user progress"})
		return
	}

	// Check and award achievements
	activityData := &services.ActivityData{
		SessionDuration:    request.SessionDuration,
		MessageCount:       request.MessageCount,
		ConversationDepth:  request.ConversationDepth,
		EmotionalIntensity: request.EmotionalIntensity,
		VulnerabilityLevel: request.VulnerabilityLevel,
		TrustLevel:         0.5,            // Default value
		IntimacyLevel:      0.5,            // Default value
		RelationshipAge:    24 * time.Hour, // Default value
	}

	err = h.gamificationService.CheckAndAwardAchievements(c.Request.Context(), userID, request.CompanionID, activityData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session activity tracked successfully"})
}

// UpdateStreak updates user streak
func (h *AnalyticsHandler) UpdateStreak(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companionID := c.Query("companion_id")
	if companionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companion_id is required"})
		return
	}

	err := h.gamificationService.UpdateStreak(c.Request.Context(), userID, companionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update streak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Streak updated successfully"})
}

// GetLevelRewards gets rewards for a specific level
func (h *AnalyticsHandler) GetLevelRewards(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level"})
		return
	}

	rewards := h.gamificationService.GetLevelRewards(level)
	c.JSON(http.StatusOK, rewards)
}

// Admin Analytics Endpoints

// GetPlatformAnalytics gets platform-wide analytics (admin only)
func (h *AnalyticsHandler) GetPlatformAnalytics(c *gin.Context) {
	// Check if user is admin (implement your admin check logic)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Add admin role check
	// if !isAdmin(userID) {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	analytics, err := h.analyticsService.GetPlatformAnalytics(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetUsersAtChurnRisk gets users at risk of churning (admin only)
func (h *AnalyticsHandler) GetUsersAtChurnRisk(c *gin.Context) {
	// Check if user is admin (implement your admin check logic)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Add admin role check
	// if !isAdmin(userID) {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	thresholdStr := c.DefaultQuery("threshold", "0.7")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		threshold = 0.7
	}

	users, err := h.predictiveAnalyticsService.GetUsersAtChurnRisk(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users at churn risk"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetAnalyticsTrends gets analytics trends (admin only)
func (h *AnalyticsHandler) GetAnalyticsTrends(c *gin.Context) {
	// Check if user is admin (implement your admin check logic)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Add admin role check
	// if !isAdmin(userID) {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trends, err := h.predictiveAnalyticsService.AnalyzeTrends(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics trends"})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// InitializeAchievements initializes achievement definitions (admin only)
func (h *AnalyticsHandler) InitializeAchievements(c *gin.Context) {
	// Check if user is admin (implement your admin check logic)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Add admin role check
	// if !isAdmin(userID) {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	err := h.gamificationService.InitializeAchievementDefinitions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievements initialized successfully"})
}
