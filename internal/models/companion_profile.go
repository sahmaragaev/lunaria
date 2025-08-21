package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanionProfile struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CompanionID        string               `bson:"companion_id" json:"companion_id"`
	UserID             string               `bson:"user_id" json:"user_id"`
	Personality        PersonalityTraits    `bson:"personality" json:"personality"`
	Backstory          string               `bson:"backstory" json:"backstory"`
	Interests          []string             `bson:"interests" json:"interests"`
	Quirks             []string             `bson:"quirks" json:"quirks"`
	CommunicationStyle CommunicationStyle   `bson:"communication_style" json:"communication_style"`
	RomanticBehavior   RomanticBehavior     `bson:"romantic_behavior" json:"romantic_behavior"`
	Preferences        CompanionPreferences `bson:"preferences" json:"preferences"`
	MemoryContext      []MemoryEntry        `bson:"memory_context" json:"memory_context"`
	CreatedAt          time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time            `bson:"updated_at" json:"updated_at"`
}

type PersonalityTraits struct {
	Warmth       float64 `bson:"warmth" json:"warmth" validate:"min=0,max=1"`
	Playfulness  float64 `bson:"playfulness" json:"playfulness" validate:"min=0,max=1"`
	Intelligence float64 `bson:"intelligence" json:"intelligence" validate:"min=0,max=1"`
	Empathy      float64 `bson:"empathy" json:"empathy" validate:"min=0,max=1"`
	Confidence   float64 `bson:"confidence" json:"confidence" validate:"min=0,max=1"`
	Romance      float64 `bson:"romance" json:"romance" validate:"min=0,max=1"`
	Humor        float64 `bson:"humor" json:"humor" validate:"min=0,max=1"`
	Clinginess   float64 `bson:"clinginess" json:"clinginess" validate:"min=0,max=1"`
}

type CommunicationStyle struct {
	Formality    float64 `bson:"formality" json:"formality" validate:"min=0,max=1"`
	Emotionality float64 `bson:"emotionality" json:"emotionality" validate:"min=0,max=1"`
	Playfulness  float64 `bson:"playfulness" json:"playfulness" validate:"min=0,max=1"`
	Intimacy     float64 `bson:"intimacy" json:"intimacy" validate:"min=0,max=1"`
}

type RomanticBehavior struct {
	Flirtatiousness float64 `bson:"flirtatiousness" json:"flirtatiousness" validate:"min=0,max=1"`
	Affection       float64 `bson:"affection" json:"affection" validate:"min=0,max=1"`
	Passion         float64 `bson:"passion" json:"passion" validate:"min=0,max=1"`
	Commitment      float64 `bson:"commitment" json:"commitment" validate:"min=0,max=1"`
}

type CompanionPreferences struct {
	PreferredTopics    []string `bson:"preferred_topics" json:"preferred_topics"`
	AvoidedTopics      []string `bson:"avoided_topics" json:"avoided_topics"`
	ResponseLength     string   `bson:"response_length" json:"response_length"`
	EmojiUsage         string   `bson:"emoji_usage" json:"emoji_usage"`
	ConversationPacing string   `bson:"conversation_pacing" json:"conversation_pacing"`
}

type MemoryEntry struct {
	Type       string    `bson:"type" json:"type"`
	Content    string    `bson:"content" json:"content"`
	Importance int       `bson:"importance" json:"importance"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Tags       []string  `bson:"tags" json:"tags"`
}
