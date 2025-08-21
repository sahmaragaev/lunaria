package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/models/dto"
)

type PersonalityService struct {
	grokService *GrokService
	validator   *validator.Validate
}

func NewPersonalityService(grokService *GrokService) *PersonalityService {
	return &PersonalityService{
		grokService: grokService,
		validator:   validator.New(),
	}
}

func (p *PersonalityService) GeneratePersonality(ctx context.Context, req *dto.PersonalityGenerationRequest) (*models.CompanionProfile, error) {
	if err := p.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	prompt := p.buildPersonalityPrompt(req)
	messages := []LLMMessage{{Role: "system", Content: prompt}}
	response, err := p.grokService.SendMiniMessage(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate personality: %w", err)
	}

	profile, err := p.parsePersonalityResponse(response, req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse personality response: %w", err)
	}
	return profile, nil
}

func (p *PersonalityService) buildPersonalityPrompt(req *dto.PersonalityGenerationRequest) string {
	interestsStr := ""
	if len(req.Interests) > 0 {
		interestsStr = fmt.Sprintf("They should be interested in: %s. ", strings.Join(req.Interests, ", "))
	}
	styleStr := ""
	if req.Style != "" {
		styleStr = fmt.Sprintf("Their overall personality style should be: %s. ", req.Style)
	}
	return fmt.Sprintf(`Create a realistic, diverse personality for %s, a %d-year-old %s person. %s%s

Create a unique individual with their own flaws, strengths, and realistic personality traits. Make them feel like a real human being who is open to romantic connections and intimate conversations.

Return JSON with this structure:
{
  "personality": {
    "warmth": 0.6-0.9,
    "playfulness": 0.5-0.9, 
    "intelligence": 0.6-0.9,
    "empathy": 0.6-0.9,
    "confidence": 0.6-0.9,
    "romance": 0.6-0.9,
    "humor": 0.5-0.9,
    "clinginess": 0.2-0.5
  },
  "communication_style": {
    "formality": 0.1-0.5,
    "emotionality": 0.5-0.9,
    "playfulness": 0.5-0.9,
    "intimacy": 0.7-0.9
  },
  "romantic_behavior": {
    "flirtatiousness": 0.6-0.9,
    "affection": 0.6-0.9,
    "passion": 0.6-0.9,
    "commitment": 0.6-0.9
  },
  "interests": ["realistic hobbies and interests"],
  "quirks": ["authentic personality quirks, not forced"],
  "backstory": "Write a realistic life story for %s with specific details, challenges, achievements, and personality development. Include real-life experiences, relationships, career/education, and what shaped their personality. Make it authentic and human, not idealized. Include their experiences with romance and relationships.",
  "preferences": {
    "preferred_topics": ["genuine interests", "romance", "relationships", "personal growth"],
    "avoided_topics": ["sensitive topics"],
    "response_length": "short|medium|long",
    "emoji_usage": "rare|moderate|frequent",
    "conversation_pacing": "slow|balanced|fast"
  }
}

IMPORTANT: Create diverse, realistic personalities who are open to romantic connections and intimate conversations. They should be comfortable with emotional vulnerability and romantic topics. Values 0.0-1.0. JSON only.`, req.Name, req.Age, req.Gender, interestsStr, styleStr, req.Name)
}

func (p *PersonalityService) parsePersonalityResponse(response string, req *dto.PersonalityGenerationRequest) (*models.CompanionProfile, error) {
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	}
	var personalityData struct {
		Personality        models.PersonalityTraits    `json:"personality"`
		CommunicationStyle models.CommunicationStyle   `json:"communication_style"`
		RomanticBehavior   models.RomanticBehavior     `json:"romantic_behavior"`
		Interests          []string                    `json:"interests"`
		Quirks             []string                    `json:"quirks"`
		Backstory          string                      `json:"backstory"`
		Preferences        models.CompanionPreferences `json:"preferences"`
	}
	if err := json.Unmarshal([]byte(response), &personalityData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal personality data: %w", err)
	}
	profile := &models.CompanionProfile{
		Personality:        personalityData.Personality,
		CommunicationStyle: personalityData.CommunicationStyle,
		RomanticBehavior:   personalityData.RomanticBehavior,
		Interests:          personalityData.Interests,
		Quirks:             personalityData.Quirks,
		Backstory:          personalityData.Backstory,
		Preferences:        personalityData.Preferences,
		MemoryContext:      []models.MemoryEntry{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := p.validator.Struct(profile); err != nil {
		return nil, fmt.Errorf("generated personality validation failed: %w", err)
	}
	return profile, nil
}

func (p *PersonalityService) GetPersonalityPresets() map[string]*models.CompanionProfile {
	return map[string]*models.CompanionProfile{
		"warm": {
			Personality:        models.PersonalityTraits{Warmth: 0.8, Playfulness: 0.6, Intelligence: 0.7, Empathy: 0.8, Confidence: 0.7, Romance: 0.7, Humor: 0.7, Clinginess: 0.3},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.3, Emotionality: 0.7, Playfulness: 0.6, Intimacy: 0.7},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.6, Affection: 0.8, Passion: 0.6, Commitment: 0.7},
		},
		"reserved": {
			Personality:        models.PersonalityTraits{Warmth: 0.5, Playfulness: 0.4, Intelligence: 0.8, Empathy: 0.7, Confidence: 0.7, Romance: 0.6, Humor: 0.5, Clinginess: 0.2},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.6, Emotionality: 0.5, Playfulness: 0.3, Intimacy: 0.6},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.4, Affection: 0.6, Passion: 0.5, Commitment: 0.6},
		},
		"confident": {
			Personality:        models.PersonalityTraits{Warmth: 0.7, Playfulness: 0.8, Intelligence: 0.8, Empathy: 0.6, Confidence: 0.9, Romance: 0.8, Humor: 0.8, Clinginess: 0.3},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.2, Emotionality: 0.6, Playfulness: 0.8, Intimacy: 0.7},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.7, Affection: 0.7, Passion: 0.7, Commitment: 0.7},
		},
		"creative": {
			Personality:        models.PersonalityTraits{Warmth: 0.8, Playfulness: 0.9, Intelligence: 0.7, Empathy: 0.8, Confidence: 0.7, Romance: 0.7, Humor: 0.9, Clinginess: 0.4},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.1, Emotionality: 0.8, Playfulness: 0.9, Intimacy: 0.7},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.6, Affection: 0.7, Passion: 0.6, Commitment: 0.6},
		},
		"busy": {
			Personality:        models.PersonalityTraits{Warmth: 0.6, Playfulness: 0.5, Intelligence: 0.8, Empathy: 0.7, Confidence: 0.8, Romance: 0.6, Humor: 0.6, Clinginess: 0.2},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.5, Emotionality: 0.5, Playfulness: 0.4, Intimacy: 0.6},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.5, Affection: 0.6, Passion: 0.5, Commitment: 0.7},
		},
		"romantic": {
			Personality:        models.PersonalityTraits{Warmth: 0.9, Playfulness: 0.7, Intelligence: 0.6, Empathy: 0.9, Confidence: 0.6, Romance: 0.9, Humor: 0.6, Clinginess: 0.5},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.2, Emotionality: 0.9, Playfulness: 0.6, Intimacy: 0.9},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.8, Affection: 0.9, Passion: 0.8, Commitment: 0.8},
		},
		"passionate": {
			Personality:        models.PersonalityTraits{Warmth: 0.8, Playfulness: 0.8, Intelligence: 0.7, Empathy: 0.8, Confidence: 0.8, Romance: 0.9, Humor: 0.7, Clinginess: 0.4},
			CommunicationStyle: models.CommunicationStyle{Formality: 0.1, Emotionality: 0.9, Playfulness: 0.7, Intimacy: 0.9},
			RomanticBehavior:   models.RomanticBehavior{Flirtatiousness: 0.9, Affection: 0.8, Passion: 0.9, Commitment: 0.7},
		},
	}
}
