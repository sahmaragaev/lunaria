package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationRepository struct {
	db *mongo.Database
}

func NewConversationRepository(db *mongo.Database) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) CreateConversation(ctx context.Context, conv *models.Conversation) (*models.Conversation, error) {
	conv.ID = primitive.NewObjectID()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()

	_, err := r.db.Collection("conversations").InsertOne(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

func (r *ConversationRepository) GetConversationByID(ctx context.Context, id primitive.ObjectID) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Collection("conversations").FindOne(ctx, bson.M{"_id": id}).Decode(&conv)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}
	return &conv, nil
}

func (r *ConversationRepository) ListUserConversations(ctx context.Context, userID string, archived bool, limit, offset int) ([]*models.Conversation, error) {
	filter := bson.M{"user_id": userID, "archived": archived}
	opts := options.Find().SetSort(bson.M{"last_activity": -1}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cur, err := r.db.Collection("conversations").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer cur.Close(ctx)
	var conversations []*models.Conversation
	for cur.Next(ctx) {
		var conv models.Conversation
		if err := cur.Decode(&conv); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}
	return conversations, nil
}

// ListConversations lists conversations between a user and companion
func (r *ConversationRepository) ListConversations(ctx context.Context, userID, companionID string, limit int, cursor any) ([]*models.Conversation, error) {
	filter := bson.M{
		"user_id":      userID,
		"companion_id": companionID,
	}

	// Add cursor-based pagination if provided
	if cursor != nil {
		if cursorID, ok := cursor.(primitive.ObjectID); ok {
			filter["_id"] = bson.M{"$lt": cursorID}
		}
	}

	opts := options.Find().
		SetSort(bson.M{"last_activity": -1}).
		SetLimit(int64(limit))

	cur, err := r.db.Collection("conversations").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer cur.Close(ctx)

	var conversations []*models.Conversation
	for cur.Next(ctx) {
		var conv models.Conversation
		if err := cur.Decode(&conv); err != nil {
			return nil, fmt.Errorf("failed to decode conversation: %w", err)
		}
		conversations = append(conversations, &conv)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return conversations, nil
}

// ListConversationsWithFilter lists all conversations with optional filtering
func (r *ConversationRepository) ListConversationsWithFilter(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Conversation, error) {
	opts := options.Find().
		SetSort(bson.M{"last_activity": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cur, err := r.db.Collection("conversations").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer cur.Close(ctx)

	var conversations []*models.Conversation
	for cur.Next(ctx) {
		var conv models.Conversation
		if err := cur.Decode(&conv); err != nil {
			return nil, fmt.Errorf("failed to decode conversation: %w", err)
		}
		conversations = append(conversations, &conv)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return conversations, nil
}

func (r *ConversationRepository) ArchiveConversation(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("conversations").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"archived": true, "updated_at": time.Now()}})
	return err
}

func (r *ConversationRepository) ReactivateConversation(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("conversations").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"archived": false, "updated_at": time.Now()}})
	return err
}

func (r *ConversationRepository) CreateMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	msg.ID = primitive.NewObjectID()
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	_, err := r.db.Collection("messages").InsertOne(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	return msg, nil
}

func (r *ConversationRepository) GetMessageByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error) {
	var msg models.Message
	err := r.db.Collection("messages").FindOne(ctx, bson.M{"_id": id}).Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	return &msg, nil
}

func (r *ConversationRepository) ListMessages(ctx context.Context, conversationID primitive.ObjectID, limit int, cursor *primitive.ObjectID) ([]*models.Message, *primitive.ObjectID, bool, error) {
	filter := bson.M{"conversation_id": conversationID}
	if cursor != nil {
		filter["_id"] = bson.M{"$lt": *cursor}
	}
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(int64(limit))
	cur, err := r.db.Collection("messages").Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to list messages: %w", err)
	}
	defer cur.Close(ctx)
	var messages []*models.Message
	var lastID *primitive.ObjectID
	for cur.Next(ctx) {
		var msg models.Message
		if err := cur.Decode(&msg); err != nil {
			return nil, nil, false, err
		}
		lastID = &msg.ID
		messages = append(messages, &msg)
	}
	hasMore := len(messages) == limit
	return messages, lastID, hasMore, nil
}

func (r *ConversationRepository) UpdateMessage(ctx context.Context, msg *models.Message) error {
	collection := r.db.Collection("messages")
	filter := bson.M{"_id": msg.ID}
	update := bson.M{"$set": bson.M{"read": msg.Read, "updated_at": msg.UpdatedAt}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

func (r *ConversationRepository) CreateMediaMetadata(ctx context.Context, media *models.MediaMetadata) (*models.MediaMetadata, error) {
	media.ID = primitive.NewObjectID()
	media.CreatedAt = time.Now()
	media.UpdatedAt = time.Now()
	_, err := r.db.Collection("media_metadata").InsertOne(ctx, media)
	if err != nil {
		return nil, fmt.Errorf("failed to create media metadata: %w", err)
	}
	return media, nil
}

func (r *ConversationRepository) GetMediaMetadataByID(ctx context.Context, id primitive.ObjectID) (*models.MediaMetadata, error) {
	var media models.MediaMetadata
	err := r.db.Collection("media_metadata").FindOne(ctx, bson.M{"_id": id}).Decode(&media)
	if err != nil {
		return nil, fmt.Errorf("media metadata not found: %w", err)
	}
	return &media, nil
}

// SaveConversationContext saves or updates conversation context
func (r *ConversationRepository) SaveConversationContext(ctx context.Context, context *models.ConversationContext) error {
	collection := r.db.Collection("conversation_contexts")

	// Use upsert to create or update
	filter := bson.M{"conversation_id": context.ConversationID}
	update := bson.M{"$set": context}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save conversation context: %w", err)
	}

	return nil
}

// GetConversationContext retrieves conversation context by conversation ID
func (r *ConversationRepository) GetConversationContext(ctx context.Context, conversationID primitive.ObjectID) (*models.ConversationContext, error) {
	collection := r.db.Collection("conversation_contexts")
	var context models.ConversationContext

	err := collection.FindOne(ctx, bson.M{"conversation_id": conversationID}).Decode(&context)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("conversation context not found")
		}
		return nil, fmt.Errorf("failed to get conversation context: %w", err)
	}

	return &context, nil
}

// SaveMemories stores AI-enhanced memories for a conversation
func (r *ConversationRepository) SaveMemories(ctx context.Context, conversationID primitive.ObjectID, memories []models.AIEnhancedMemoryEntry) error {
	collection := r.db.Collection("ai_memories")

	// Convert memories to documents for insertion
	var documents []any
	for _, memory := range memories {
		memory.ConversationID = conversationID
		documents = append(documents, memory)
	}

	if len(documents) > 0 {
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			return fmt.Errorf("failed to save memories: %w", err)
		}
	}

	return nil
}

// GetMemories retrieves AI-enhanced memories for a conversation
func (r *ConversationRepository) GetMemories(ctx context.Context, conversationID primitive.ObjectID, limit int) ([]models.AIEnhancedMemoryEntry, error) {
	collection := r.db.Collection("ai_memories")

	filter := bson.M{"conversation_id": conversationID}
	opts := options.Find().
		SetSort(bson.M{"importance": -1, "last_referenced": -1}).
		SetLimit(int64(limit))

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories: %w", err)
	}
	defer cur.Close(ctx)

	var memories []models.AIEnhancedMemoryEntry
	for cur.Next(ctx) {
		var memory models.AIEnhancedMemoryEntry
		if err := cur.Decode(&memory); err != nil {
			return nil, fmt.Errorf("failed to decode memory: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// UpdateMemoryReference updates the last referenced time and frequency of a memory
func (r *ConversationRepository) UpdateMemoryReference(ctx context.Context, memoryID primitive.ObjectID) error {
	collection := r.db.Collection("ai_memories")

	filter := bson.M{"_id": memoryID}
	update := bson.M{
		"$set": bson.M{
			"last_referenced": time.Now(),
			"updated_at":      time.Now(),
		},
		"$inc": bson.M{
			"frequency": 1,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update memory reference: %w", err)
	}

	return nil
}

// DeleteUserConversations deletes all conversations for a specific user
func (r *ConversationRepository) DeleteUserConversations(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}

	// Delete conversations
	_, err := r.db.Collection("conversations").DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user conversations: %w", err)
	}

	// Delete messages for this user's conversations
	// First get all conversation IDs for this user
	convFilter := bson.M{"user_id": userID}
	convCursor, err := r.db.Collection("conversations").Find(ctx, convFilter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return fmt.Errorf("failed to find user conversations: %w", err)
	}
	defer convCursor.Close(ctx)

	var conversationIDs []primitive.ObjectID
	for convCursor.Next(ctx) {
		var conv struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := convCursor.Decode(&conv); err != nil {
			return fmt.Errorf("failed to decode conversation ID: %w", err)
		}
		conversationIDs = append(conversationIDs, conv.ID)
	}

	if len(conversationIDs) > 0 {
		// Delete messages
		msgFilter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
		_, err = r.db.Collection("messages").DeleteMany(ctx, msgFilter)
		if err != nil {
			return fmt.Errorf("failed to delete user messages: %w", err)
		}

		// Delete conversation contexts
		ctxFilter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
		_, err = r.db.Collection("conversation_contexts").DeleteMany(ctx, ctxFilter)
		if err != nil {
			return fmt.Errorf("failed to delete conversation contexts: %w", err)
		}

		// Delete AI memories
		memoryFilter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
		_, err = r.db.Collection("ai_memories").DeleteMany(ctx, memoryFilter)
		if err != nil {
			return fmt.Errorf("failed to delete AI memories: %w", err)
		}
	}

	return nil
}

// GetConversationStats gets statistics about conversations
func (r *ConversationRepository) GetConversationStats(ctx context.Context, userID string) (map[string]any, error) {
	stats := make(map[string]any)

	// Total conversations
	totalFilter := bson.M{"user_id": userID}
	totalCount, err := r.db.Collection("conversations").CountDocuments(ctx, totalFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count total conversations: %w", err)
	}
	stats["total_conversations"] = totalCount

	// Active conversations (not archived)
	activeFilter := bson.M{"user_id": userID, "archived": false}
	activeCount, err := r.db.Collection("conversations").CountDocuments(ctx, activeFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count active conversations: %w", err)
	}
	stats["active_conversations"] = activeCount

	// Archived conversations
	archivedFilter := bson.M{"user_id": userID, "archived": true}
	archivedCount, err := r.db.Collection("conversations").CountDocuments(ctx, archivedFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count archived conversations: %w", err)
	}
	stats["archived_conversations"] = archivedCount

	// Total messages
	convFilter := bson.M{"user_id": userID}
	convCursor, err := r.db.Collection("conversations").Find(ctx, convFilter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find user conversations: %w", err)
	}
	defer convCursor.Close(ctx)

	var conversationIDs []primitive.ObjectID
	for convCursor.Next(ctx) {
		var conv struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := convCursor.Decode(&conv); err != nil {
			return nil, fmt.Errorf("failed to decode conversation ID: %w", err)
		}
		conversationIDs = append(conversationIDs, conv.ID)
	}

	if len(conversationIDs) > 0 {
		msgFilter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
		messageCount, err := r.db.Collection("messages").CountDocuments(ctx, msgFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to count messages: %w", err)
		}
		stats["total_messages"] = messageCount
	} else {
		stats["total_messages"] = int64(0)
	}

	return stats, nil
}

// GetCompanionConversationStats gets conversation statistics for a specific companion
func (r *ConversationRepository) GetCompanionConversationStats(ctx context.Context, userID, companionID string) (*models.ConversationStats, error) {
	// Get conversations between user and companion
	convFilter := bson.M{"user_id": userID, "companion_id": companionID}
	convCursor, err := r.db.Collection("conversations").Find(ctx, convFilter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations: %w", err)
	}
	defer convCursor.Close(ctx)

	var conversationIDs []primitive.ObjectID
	for convCursor.Next(ctx) {
		var conv struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := convCursor.Decode(&conv); err != nil {
			return nil, fmt.Errorf("failed to decode conversation ID: %w", err)
		}
		conversationIDs = append(conversationIDs, conv.ID)
	}

	stats := &models.ConversationStats{
		TotalMessages:     0,
		UserMessages:      0,
		CompanionMessages: 0,
		LastMessageAt:     time.Time{},
	}

	if len(conversationIDs) > 0 {
		// Get total messages
		msgFilter := bson.M{"conversation_id": bson.M{"$in": conversationIDs}}
		totalMessages, err := r.db.Collection("messages").CountDocuments(ctx, msgFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to count total messages: %w", err)
		}
		stats.TotalMessages = int(totalMessages)

		// Get user messages count
		userMsgFilter := bson.M{
			"conversation_id": bson.M{"$in": conversationIDs},
			"sender_type":     "user",
		}
		userMessages, err := r.db.Collection("messages").CountDocuments(ctx, userMsgFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to count user messages: %w", err)
		}
		stats.UserMessages = int(userMessages)

		// Get companion messages count
		companionMsgFilter := bson.M{
			"conversation_id": bson.M{"$in": conversationIDs},
			"sender_type":     "companion",
		}
		companionMessages, err := r.db.Collection("messages").CountDocuments(ctx, companionMsgFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to count companion messages: %w", err)
		}
		stats.CompanionMessages = int(companionMessages)

		// Get last message timestamp
		lastMsgOpts := options.FindOne().SetSort(bson.M{"created_at": -1})
		var lastMessage models.Message
		err = r.db.Collection("messages").FindOne(ctx, msgFilter, lastMsgOpts).Decode(&lastMessage)
		if err == nil {
			stats.LastMessageAt = lastMessage.CreatedAt
		}
	}

	return stats, nil
}

// GetConversationsByDateRange gets conversations within a date range
func (r *ConversationRepository) GetConversationsByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.Conversation, error) {
	filter := bson.M{
		"user_id": userID,
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cur, err := r.db.Collection("conversations").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations by date range: %w", err)
	}
	defer cur.Close(ctx)

	var conversations []*models.Conversation
	for cur.Next(ctx) {
		var conv models.Conversation
		if err := cur.Decode(&conv); err != nil {
			return nil, fmt.Errorf("failed to decode conversation: %w", err)
		}
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}
