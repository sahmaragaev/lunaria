package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RealTimeAnalyticsService handles real-time analytics processing
type RealTimeAnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
	convRepo      *repositories.ConversationRepository
	grokService   *GrokService

	// Real-time processing
	eventStream chan *AnalyticsEvent
	processors  map[string]EventProcessor
	mu          sync.RWMutex

	// Performance monitoring
	processingStats *ProcessingStats
}

// AnalyticsEvent represents a real-time analytics event
type AnalyticsEvent struct {
	Type           string             `json:"type"`
	UserID         string             `json:"user_id"`
	CompanionID    string             `json:"companion_id"`
	ConversationID primitive.ObjectID `json:"conversation_id"`
	Timestamp      time.Time          `json:"timestamp"`
	Data           map[string]any     `json:"data"`
	Priority       int                `json:"priority"` // 1=high, 2=medium, 3=low
}

// EventProcessor processes specific types of analytics events
type EventProcessor interface {
	Process(ctx context.Context, event *AnalyticsEvent) error
	GetEventType() string
}

// ProcessingStats tracks real-time processing performance
type ProcessingStats struct {
	EventsProcessed    int64         `json:"events_processed"`
	AverageProcessTime time.Duration `json:"average_process_time"`
	ErrorCount         int64         `json:"error_count"`
	LastProcessed      time.Time     `json:"last_processed"`
	mu                 sync.RWMutex
}

// NewRealTimeAnalyticsService creates a new real-time analytics service
func NewRealTimeAnalyticsService(analyticsRepo *repositories.AnalyticsRepository, convRepo *repositories.ConversationRepository, grokService *GrokService) *RealTimeAnalyticsService {
	service := &RealTimeAnalyticsService{
		analyticsRepo:   analyticsRepo,
		convRepo:        convRepo,
		grokService:     grokService,
		eventStream:     make(chan *AnalyticsEvent, 1000),
		processors:      make(map[string]EventProcessor),
		processingStats: &ProcessingStats{},
	}

	// Register event processors
	service.registerProcessors()

	// Start event processing
	go service.startEventProcessing()

	return service
}

// registerProcessors registers all event processors
func (s *RealTimeAnalyticsService) registerProcessors() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.processors["message_sent"] = &MessageSentProcessor{s}
	s.processors["session_started"] = &SessionStartedProcessor{s}
	s.processors["session_ended"] = &SessionEndedProcessor{s}
	s.processors["response_received"] = &ResponseReceivedProcessor{s}
	s.processors["user_activity"] = &UserActivityProcessor{s}
	s.processors["emotional_state"] = &EmotionalStateProcessor{s}
	s.processors["engagement_level"] = &EngagementLevelProcessor{s}
}

// startEventProcessing starts the event processing loop
func (s *RealTimeAnalyticsService) startEventProcessing() {
	for event := range s.eventStream {
		s.processEvent(context.Background(), event)
	}
}

// processEvent processes a single analytics event
func (s *RealTimeAnalyticsService) processEvent(ctx context.Context, event *AnalyticsEvent) {
	start := time.Now()

	s.mu.RLock()
	processor, exists := s.processors[event.Type]
	s.mu.RUnlock()

	if !exists {
		s.updateStats(time.Since(start), true)
		return
	}

	err := processor.Process(ctx, event)
	s.updateStats(time.Since(start), err != nil)

	if err != nil {
		// Log error but don't block processing
		fmt.Printf("Error processing event %s: %v\n", event.Type, err)
	}
}

// updateStats updates processing statistics
func (s *RealTimeAnalyticsService) updateStats(processTime time.Duration, hasError bool) {
	s.processingStats.mu.Lock()
	defer s.processingStats.mu.Unlock()

	s.processingStats.EventsProcessed++
	s.processingStats.LastProcessed = time.Now()

	if hasError {
		s.processingStats.ErrorCount++
	}

	// Update average processing time
	if s.processingStats.EventsProcessed == 1 {
		s.processingStats.AverageProcessTime = processTime
	} else {
		total := s.processingStats.AverageProcessTime * time.Duration(s.processingStats.EventsProcessed-1)
		s.processingStats.AverageProcessTime = (total + processTime) / time.Duration(s.processingStats.EventsProcessed)
	}
}

// EmitEvent emits an analytics event for processing
func (s *RealTimeAnalyticsService) EmitEvent(event *AnalyticsEvent) {
	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Send to event stream (non-blocking)
	select {
	case s.eventStream <- event:
	default:
		// Channel full, log warning
		fmt.Printf("Warning: Event stream full, dropping event %s\n", event.Type)
	}
}

// TrackMessageSent tracks when a message is sent
func (s *RealTimeAnalyticsService) TrackMessageSent(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, messageData map[string]any) {
	event := &AnalyticsEvent{
		Type:           "message_sent",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           messageData,
		Priority:       1,
	}

	s.EmitEvent(event)
}

// TrackSessionStarted tracks when a session starts
func (s *RealTimeAnalyticsService) TrackSessionStarted(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID) {
	event := &AnalyticsEvent{
		Type:           "session_started",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           make(map[string]any),
		Priority:       1,
	}

	s.EmitEvent(event)
}

// TrackSessionEnded tracks when a session ends
func (s *RealTimeAnalyticsService) TrackSessionEnded(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, sessionData map[string]any) {
	event := &AnalyticsEvent{
		Type:           "session_ended",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           sessionData,
		Priority:       1,
	}

	s.EmitEvent(event)
}

// TrackResponseReceived tracks when a response is received
func (s *RealTimeAnalyticsService) TrackResponseReceived(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, responseData map[string]any) {
	event := &AnalyticsEvent{
		Type:           "response_received",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           responseData,
		Priority:       2,
	}

	s.EmitEvent(event)
}

// TrackUserActivity tracks general user activity
func (s *RealTimeAnalyticsService) TrackUserActivity(ctx context.Context, userID, companionID string, activityType string, activityData map[string]any) {
	event := &AnalyticsEvent{
		Type:        "user_activity",
		UserID:      userID,
		CompanionID: companionID,
		Timestamp:   time.Now(),
		Data: map[string]any{
			"activity_type": activityType,
			"activity_data": activityData,
		},
		Priority: 3,
	}

	s.EmitEvent(event)
}

// TrackEmotionalState tracks user emotional state
func (s *RealTimeAnalyticsService) TrackEmotionalState(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, emotionalData map[string]any) {
	event := &AnalyticsEvent{
		Type:           "emotional_state",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           emotionalData,
		Priority:       2,
	}

	s.EmitEvent(event)
}

// TrackEngagementLevel tracks user engagement level
func (s *RealTimeAnalyticsService) TrackEngagementLevel(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, engagementData map[string]any) {
	event := &AnalyticsEvent{
		Type:           "engagement_level",
		UserID:         userID,
		CompanionID:    companionID,
		ConversationID: conversationID,
		Timestamp:      time.Now(),
		Data:           engagementData,
		Priority:       2,
	}

	s.EmitEvent(event)
}

// GetProcessingStats returns current processing statistics
func (s *RealTimeAnalyticsService) GetProcessingStats() *ProcessingStats {
	s.processingStats.mu.RLock()
	defer s.processingStats.mu.RUnlock()

	stats := *s.processingStats
	return &stats
}

// GetActiveSessions gets currently active sessions
func (s *RealTimeAnalyticsService) GetActiveSessions(ctx context.Context) ([]models.RealTimeMetrics, error) {
	// This would query the real-time metrics collection for active sessions
	// For now, return empty slice
	return []models.RealTimeMetrics{}, nil
}

// GetSystemHealth gets system health metrics
func (s *RealTimeAnalyticsService) GetSystemHealth(ctx context.Context) map[string]any {
	stats := s.GetProcessingStats()

	return map[string]any{
		"events_processed":     stats.EventsProcessed,
		"average_process_time": stats.AverageProcessTime.String(),
		"error_rate":           float64(stats.ErrorCount) / float64(stats.EventsProcessed),
		"last_processed":       stats.LastProcessed,
		"queue_size":           len(s.eventStream),
		"active_processors":    len(s.processors),
	}
}

// Event Processors

// MessageSentProcessor processes message sent events
type MessageSentProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *MessageSentProcessor) GetEventType() string {
	return "message_sent"
}

func (p *MessageSentProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Update real-time metrics
	metrics := &models.RealTimeMetrics{
		UserID:            event.UserID,
		CompanionID:       event.CompanionID,
		ConversationID:    event.ConversationID,
		IsActive:          true,
		MessagesInSession: 1, // Increment message count
		LastResponseTime:  event.Timestamp,
		Timestamp:         event.Timestamp,
	}

	// Update message count from data if available
	if messageCount, ok := event.Data["message_count"].(int); ok {
		metrics.MessagesInSession = messageCount
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}

// SessionStartedProcessor processes session started events
type SessionStartedProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *SessionStartedProcessor) GetEventType() string {
	return "session_started"
}

func (p *SessionStartedProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	metrics := &models.RealTimeMetrics{
		UserID:            event.UserID,
		CompanionID:       event.CompanionID,
		ConversationID:    event.ConversationID,
		IsActive:          true,
		SessionStartTime:  event.Timestamp,
		MessagesInSession: 0,
		Timestamp:         event.Timestamp,
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}

// SessionEndedProcessor processes session ended events
type SessionEndedProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *SessionEndedProcessor) GetEventType() string {
	return "session_ended"
}

func (p *SessionEndedProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Get session duration from data
	sessionDuration := time.Duration(0)
	if duration, ok := event.Data["session_duration"].(time.Duration); ok {
		sessionDuration = duration
	}

	metrics := &models.RealTimeMetrics{
		UserID:                 event.UserID,
		CompanionID:            event.CompanionID,
		ConversationID:         event.ConversationID,
		IsActive:               false,
		CurrentSessionDuration: sessionDuration,
		Timestamp:              event.Timestamp,
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}

// ResponseReceivedProcessor processes response received events
type ResponseReceivedProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *ResponseReceivedProcessor) GetEventType() string {
	return "response_received"
}

func (p *ResponseReceivedProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Update response time metrics
	responseTime := time.Duration(0)
	if rt, ok := event.Data["response_time"].(time.Duration); ok {
		responseTime = rt
	}

	metrics := &models.RealTimeMetrics{
		UserID:           event.UserID,
		CompanionID:      event.CompanionID,
		ConversationID:   event.ConversationID,
		IsActive:         true,
		LastResponseTime: event.Timestamp,
		AIResponseTime:   responseTime,
		Timestamp:        event.Timestamp,
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}

// UserActivityProcessor processes user activity events
type UserActivityProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *UserActivityProcessor) GetEventType() string {
	return "user_activity"
}

func (p *UserActivityProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Process user activity for engagement tracking
	// This could trigger engagement level updates, streak tracking, etc.
	return nil
}

// EmotionalStateProcessor processes emotional state events
type EmotionalStateProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *EmotionalStateProcessor) GetEventType() string {
	return "emotional_state"
}

func (p *EmotionalStateProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Update emotional state in real-time metrics
	emotionalState := "neutral"
	if state, ok := event.Data["emotional_state"].(string); ok {
		emotionalState = state
	}

	moodIndicator := "neutral"
	if mood, ok := event.Data["mood_indicator"].(string); ok {
		moodIndicator = mood
	}

	metrics := &models.RealTimeMetrics{
		UserID:         event.UserID,
		CompanionID:    event.CompanionID,
		ConversationID: event.ConversationID,
		IsActive:       true,
		EmotionalState: emotionalState,
		MoodIndicator:  moodIndicator,
		Timestamp:      event.Timestamp,
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}

// EngagementLevelProcessor processes engagement level events
type EngagementLevelProcessor struct {
	service *RealTimeAnalyticsService
}

func (p *EngagementLevelProcessor) GetEventType() string {
	return "engagement_level"
}

func (p *EngagementLevelProcessor) Process(ctx context.Context, event *AnalyticsEvent) error {
	// Update engagement level in real-time metrics
	engagementLevel := 0.5
	if level, ok := event.Data["engagement_level"].(float64); ok {
		engagementLevel = level
	}

	responseQuality := 0.5
	if quality, ok := event.Data["response_quality"].(float64); ok {
		responseQuality = quality
	}

	metrics := &models.RealTimeMetrics{
		UserID:          event.UserID,
		CompanionID:     event.CompanionID,
		ConversationID:  event.ConversationID,
		IsActive:        true,
		EngagementLevel: engagementLevel,
		ResponseQuality: responseQuality,
		Timestamp:       event.Timestamp,
	}

	return p.service.analyticsRepo.UpsertRealTimeMetrics(ctx, metrics)
}
