package router

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sahmaragaev/lunaria-backend/internal/config"
	"github.com/sahmaragaev/lunaria-backend/internal/database/mongodb"
	"github.com/sahmaragaev/lunaria-backend/internal/database/postgres"
	"github.com/sahmaragaev/lunaria-backend/internal/handlers"
	"github.com/sahmaragaev/lunaria-backend/internal/middleware"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"github.com/sahmaragaev/lunaria-backend/internal/services"
)

func SetupRouter(cfg *config.Config, pgDB *postgres.PostgresDB, mongoDB *mongodb.MongoDB) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Services
	redisService := services.NewRedisService(&cfg.Redis)
	jwtService := services.NewJWTService(&cfg.JWT, redisService)
	passwordService := services.NewPasswordService()
	grokService := services.NewGrokService(&cfg.Grok)
	personalityService := services.NewPersonalityService(grokService)

	// Repositories
	userRepo := repositories.NewUserRepository(pgDB.DB)
	companionRepo := repositories.NewCompanionRepository(pgDB.DB, mongoDB.Database)
	relationshipRepo := repositories.NewRelationshipRepository(pgDB.DB)
	conversationRepo := repositories.NewConversationRepository(mongoDB.Database)
	analyticsRepo := repositories.NewAnalyticsRepository(pgDB.DB, mongoDB.Database)

	// Services
	authService := services.NewAuthService(userRepo, jwtService, passwordService)
	companionService := services.NewCompanionService(companionRepo, relationshipRepo, conversationRepo, personalityService)

	// S3 custom config for Contabo or any S3-compatible storage
	s3cfg := cfg.S3
	awsCfg, _ := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(s3cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3cfg.AccessKeyID, s3cfg.SecretAccessKey, "")),
	)
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = &s3cfg.Endpoint
		o.UsePathStyle = s3cfg.UsePathStyle
	})
	mediaService := services.NewMediaServiceWithClient(s3Client, s3cfg.S3Bucket, conversationRepo, analyticsRepo, s3cfg.Endpoint)
	conversationService := services.NewConversationService(conversationRepo, analyticsRepo)

	// Initialize advanced AI services
	aiContextService := services.NewAIContextService(grokService, conversationRepo)
	responseQualityService := services.NewResponseQualityService(grokService, conversationRepo)
	conversationIntelligenceService := services.NewConversationIntelligenceService(grokService, conversationRepo)

	// Initialize message service with all AI components
	messageService := services.NewMessageService(conversationRepo, analyticsRepo, grokService, aiContextService, responseQualityService, conversationIntelligenceService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	healthHandler := handlers.NewHealthHandler(pgDB, mongoDB)
	companionHandler := handlers.NewCompanionHandler(companionService)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	messageHandler := handlers.NewMessageHandler(messageService, conversationService, companionService)

	// Routes
	v1 := router.Group("/api/v1")

	// Health checks
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.GET("/health/live", healthHandler.LivenessCheck)

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
		auth.GET("/me", authMiddleware.RequireAuth(), authHandler.GetProfile)
	}

	// Profile routes (protected)
	profile := v1.Group("/profile")
	profile.Use(authMiddleware.RequireAuth())
	{
		profile.GET("", authHandler.GetProfile)
		profile.PUT("", authHandler.UpdateProfile)
	}

	// Companion routes (protected)
	companions := v1.Group("/companions")
	companions.Use(authMiddleware.RequireAuth())
	{
		companions.POST("", companionHandler.CreateCompanion)
		companions.GET("", companionHandler.GetUserCompanions)
		companions.GET(":id", companionHandler.GetCompanion)
		companions.PUT(":id", companionHandler.UpdateCompanion)
		companions.DELETE(":id", companionHandler.DeleteCompanion)
	}

	// Media routes
	media := v1.Group("/media")
	media.Use(authMiddleware.RequireAuth())
	{
		media.POST("/upload-url", mediaHandler.GenerateUploadURL)
		media.POST("/validate", mediaHandler.ValidateMedia)
		media.GET(":file_id", mediaHandler.GetMediaFile)
	}

	// Conversation routes
	conversations := v1.Group("/conversations")
	conversations.Use(authMiddleware.RequireAuth())
	{
		conversations.POST("", conversationHandler.StartConversation)
		conversations.GET("", conversationHandler.ListConversations)
		conversations.GET(":id", conversationHandler.GetConversation)
		conversations.POST(":id/archive", conversationHandler.ArchiveConversation)
		conversations.POST(":id/reactivate", conversationHandler.ReactivateConversation)
		// Messaging routes
		conversations.POST(":id/messages", messageHandler.SendMessage)
		conversations.GET(":id/messages", messageHandler.ListMessages)
		conversations.GET(":id/messages/:message_id", messageHandler.GetMessage)
		conversations.PUT(":id/messages/:message_id/read", messageHandler.MarkAsRead)
		// Advanced AI routes
		conversations.GET(":id/intelligence", messageHandler.GetConversationIntelligence)
		conversations.GET(":id/suggest-topic", messageHandler.SuggestNextTopic)
		conversations.GET(":id/engagement", messageHandler.AnalyzeEngagement)
		conversations.GET(":id/messages/:message_id/quality", messageHandler.GetResponseQuality)
		conversations.GET(":id/typing-status", messageHandler.CheckTypingStatus)
	}

	return router
}
