package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunMigrations(db *mongo.Database) error {
	ctx := context.Background()

	// Conversations
	_, err := db.Collection("conversations").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "companion_id", Value: 1}, {Key: "last_activity", Value: -1}},
			Options: options.Index().SetName("idx_conversations_user_companion"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_conversations_created_at"),
		},
	})
	if err != nil {
		log.Printf("MongoDB migration (conversations) failed: %v", err)
		return err
	}

	// Messages
	_, err = db.Collection("messages").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "conversation_id", Value: 1}, {Key: "created_at", Value: -1}},
		Options: options.Index().SetName("idx_messages_conversation_created"),
	})
	if err != nil {
		log.Printf("MongoDB migration (messages) failed: %v", err)
		return err
	}

	// User engagement analytics
	_, err = db.Collection("user_engagement_analytics").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "companion_id", Value: 1}, {Key: "conversation_id", Value: 1}, {Key: "created_at", Value: -1}},
		Options: options.Index().SetName("idx_analytics_user_companion_conversation_created"),
	})
	if err != nil {
		log.Printf("MongoDB migration (analytics) failed: %v", err)
		return err
	}

	log.Println("MongoDB migrations applied successfully.")
	return nil
}
