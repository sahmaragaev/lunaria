package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsService struct {
	grokService *GrokService
	repo        *repositories.AnalyticsRepository
	convRepo    *repositories.ConversationRepository
}

func NewAnalyticsService(grokService *GrokService, repo *repositories.AnalyticsRepository, convRepo *repositories.ConversationRepository) *AnalyticsService {
	return &AnalyticsService{
		grokService: grokService,
		repo:        repo,
		convRepo:    convRepo,
	}
}

// TrackUserEngagement tracks comprehensive user engagement metrics
func (s *AnalyticsService) TrackUserEngagement(ctx context.Context, userID, companionID string, conversationID primitive.ObjectID, sessionData *SessionData) error {
	// Get existing analytics or create new
	analytics, err := s.repo.GetUserEngagementAnalytics(ctx, userID, companionID, conversationID)
	if err != nil {
		// Create new analytics record
		analytics = &models.UserEngagementAnalytics{
			UserID:         userID,
			CompanionID:    companionID,
			ConversationID: conversationID,
			CreatedAt:      time.Now(),
		}
	}

	// Update session metrics
	analytics.SessionDuration = sessionData.Duration
	analytics.MessagesPerSession = sessionData.MessageCount
	analytics.ResponseTime = sessionData.AverageResponseTime
	analytics.PeakActivityTime = sessionData.PeakActivityTime
	analytics.UpdatedAt = time.Now()

	// Analyze conversation quality
	qualityMetrics, err := s.analyzeConversationQuality(ctx, conversationID, sessionData)
	if err != nil {
		return fmt.Errorf("failed to analyze conversation quality: %w", err)
	}

	analytics.ConversationDepth = qualityMetrics.Depth
	analytics.EmotionalIntensity = qualityMetrics.EmotionalIntensity
	analytics.TopicDiversity = qualityMetrics.TopicDiversity
	analytics.VulnerabilityLevel = qualityMetrics.VulnerabilityLevel
	analytics.EngagementScore = qualityMetrics.EngagementScore

	// Analyze behavioral patterns
	behavioralPatterns, err := s.analyzeBehavioralPatterns(ctx, userID, companionID)
	if err != nil {
		return fmt.Errorf("failed to analyze behavioral patterns: %w", err)
	}

	analytics.SessionFrequency = behavioralPatterns.SessionFrequency
	analytics.PreferredTopics = behavioralPatterns.PreferredTopics
	analytics.InteractionStyle = behavioralPatterns.InteractionStyle

	// Analyze relationship progression
	relationshipMetrics, err := s.analyzeRelationshipProgression(ctx, userID, companionID)
	if err != nil {
		return fmt.Errorf("failed to analyze relationship progression: %w", err)
	}

	analytics.IntimacyGrowth = relationshipMetrics.IntimacyGrowth
	analytics.TrustBuilding = relationshipMetrics.TrustBuilding
	analytics.RelationshipStage = relationshipMetrics.Stage
	analytics.MilestoneProgress = relationshipMetrics.MilestoneProgress

	// Analyze emotional intelligence
	emotionalMetrics, err := s.analyzeEmotionalIntelligence(ctx, conversationID, sessionData)
	if err != nil {
		return fmt.Errorf("failed to analyze emotional intelligence: %w", err)
	}

	analytics.SentimentTrend = emotionalMetrics.SentimentTrend
	analytics.EmotionalRegulation = emotionalMetrics.EmotionalRegulation
	analytics.EmpathyResponse = emotionalMetrics.EmpathyResponse
	analytics.MoodImpact = emotionalMetrics.MoodImpact

	// Save analytics
	return s.repo.UpsertUserEngagementAnalytics(ctx, analytics)
}

// SessionData represents session information for analytics
type SessionData struct {
	Duration            time.Duration
	MessageCount        int
	AverageResponseTime time.Duration
	PeakActivityTime    time.Time
	Messages            []*models.Message
	ResponseQuality     float64
}

// ConversationQualityMetrics represents conversation quality analysis
type ConversationQualityMetrics struct {
	Depth              float64
	EmotionalIntensity float64
	TopicDiversity     float64
	VulnerabilityLevel float64
	EngagementScore    float64
}

// analyzeConversationQuality analyzes the quality of a conversation
func (s *AnalyticsService) analyzeConversationQuality(ctx context.Context, conversationID primitive.ObjectID, sessionData *SessionData) (*ConversationQualityMetrics, error) {
	// Get recent messages for analysis
	messages, _, _, err := s.convRepo.ListMessages(ctx, conversationID, 50, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze the quality of this conversation:

CONVERSATION:
%s

Analyze and rate from 0.0 to 1.0:
1. Conversation depth (meaningfulness and substance)
2. Emotional intensity (emotional engagement and expression)
3. Topic diversity (variety of topics discussed)
4. Vulnerability level (openness and authenticity)
5. Overall engagement score

Respond with JSON:
{
  "depth": 0.0-1.0,
  "emotional_intensity": 0.0-1.0,
  "topic_diversity": 0.0-1.0,
  "vulnerability_level": 0.0-1.0,
  "engagement_score": 0.0-1.0,
  "analysis": "detailed analysis"
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are a conversation quality analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return &ConversationQualityMetrics{
			Depth:              0.5,
			EmotionalIntensity: 0.5,
			TopicDiversity:     0.5,
			VulnerabilityLevel: 0.5,
			EngagementScore:    0.5,
		}, nil
	}

	var metrics ConversationQualityMetrics
	if err := json.Unmarshal([]byte(response), &metrics); err != nil {
		return &ConversationQualityMetrics{
			Depth:              0.5,
			EmotionalIntensity: 0.5,
			TopicDiversity:     0.5,
			VulnerabilityLevel: 0.5,
			EngagementScore:    0.5,
		}, nil
	}

	return &metrics, nil
}

// BehavioralPatterns represents user behavioral analysis
type BehavioralPatterns struct {
	SessionFrequency int
	PreferredTopics  []string
	InteractionStyle string
}

// analyzeBehavioralPatterns analyzes user behavioral patterns
func (s *AnalyticsService) analyzeBehavioralPatterns(ctx context.Context, userID, companionID string) (*BehavioralPatterns, error) {
	// Get user progress to analyze patterns
	progress, err := s.repo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return &BehavioralPatterns{
			SessionFrequency: 1,
			PreferredTopics:  []string{},
			InteractionStyle: "balanced",
		}, nil
	}

	// Analyze session frequency based on total conversations and time
	sessionFrequency := 1
	if progress.TotalConversations > 0 {
		daysSinceCreation := time.Since(progress.CreatedAt).Hours() / 24
		if daysSinceCreation > 0 {
			sessionFrequency = int(float64(progress.TotalConversations) / daysSinceCreation)
		}
	}

	// Get preferred topics from recent conversations
	preferredTopics := []string{}
	// This would be enhanced with actual topic analysis from conversations

	// Determine interaction style based on message patterns
	interactionStyle := "balanced"
	if progress.AverageSessionLength > 30*time.Minute {
		interactionStyle = "deep"
	} else if progress.AverageSessionLength < 5*time.Minute {
		interactionStyle = "brief"
	}

	return &BehavioralPatterns{
		SessionFrequency: sessionFrequency,
		PreferredTopics:  preferredTopics,
		InteractionStyle: interactionStyle,
	}, nil
}

// RelationshipMetrics represents relationship progression analysis
type RelationshipMetrics struct {
	IntimacyGrowth    float64
	TrustBuilding     float64
	Stage             string
	MilestoneProgress map[string]float64
}

// analyzeRelationshipProgression analyzes relationship development
func (s *AnalyticsService) analyzeRelationshipProgression(ctx context.Context, userID, companionID string) (*RelationshipMetrics, error) {
	// Get relationship analytics
	relationshipAnalytics, err := s.repo.GetRelationshipAnalytics(ctx, userID, companionID)
	if err != nil {
		return &RelationshipMetrics{
			IntimacyGrowth:    0.0,
			TrustBuilding:     0.0,
			Stage:             "meeting",
			MilestoneProgress: make(map[string]float64),
		}, nil
	}

	return &RelationshipMetrics{
		IntimacyGrowth:    relationshipAnalytics.IntimacyGrowth,
		TrustBuilding:     relationshipAnalytics.TrustLevel,
		Stage:             relationshipAnalytics.CurrentStage,
		MilestoneProgress: make(map[string]float64), // Would be populated from actual milestones
	}, nil
}

// EmotionalMetrics represents emotional intelligence analysis
type EmotionalMetrics struct {
	SentimentTrend      []models.SentimentPoint
	EmotionalRegulation float64
	EmpathyResponse     float64
	MoodImpact          float64
}

// analyzeEmotionalIntelligence analyzes emotional aspects of conversations
func (s *AnalyticsService) analyzeEmotionalIntelligence(ctx context.Context, conversationID primitive.ObjectID, sessionData *SessionData) (*EmotionalMetrics, error) {
	// Get recent messages for sentiment analysis
	messages, _, _, err := s.convRepo.ListMessages(ctx, conversationID, 20, nil)
	if err != nil {
		return &EmotionalMetrics{
			SentimentTrend:      []models.SentimentPoint{},
			EmotionalRegulation: 0.5,
			EmpathyResponse:     0.5,
			MoodImpact:          0.5,
		}, nil
	}

	// Analyze sentiment trend
	sentimentTrend := s.analyzeSentimentTrend(messages)

	// Analyze emotional regulation and empathy
	emotionalAnalysis, err := s.analyzeEmotionalPatterns(ctx, messages)
	if err != nil {
		return &EmotionalMetrics{
			SentimentTrend:      sentimentTrend,
			EmotionalRegulation: 0.5,
			EmpathyResponse:     0.5,
			MoodImpact:          0.5,
		}, nil
	}

	return &EmotionalMetrics{
		SentimentTrend:      sentimentTrend,
		EmotionalRegulation: emotionalAnalysis.Regulation,
		EmpathyResponse:     emotionalAnalysis.Empathy,
		MoodImpact:          emotionalAnalysis.MoodImpact,
	}, nil
}

// analyzeSentimentTrend analyzes sentiment over time
func (s *AnalyticsService) analyzeSentimentTrend(messages []*models.Message) []models.SentimentPoint {
	var sentimentPoints []models.SentimentPoint

	for i, msg := range messages {
		if msg.Text == nil {
			continue
		}

		// Simple sentiment analysis (would be enhanced with AI)
		sentiment := s.calculateSimpleSentiment(*msg.Text)

		point := models.SentimentPoint{
			Timestamp: msg.CreatedAt,
			Score:     sentiment.Score,
			Intensity: sentiment.Intensity,
			Dominant:  sentiment.Dominant,
		}
		sentimentPoints = append(sentimentPoints, point)

		// Limit to last 10 points for performance
		if i >= 9 {
			break
		}
	}

	return sentimentPoints
}

// SimpleSentiment represents basic sentiment analysis
type SimpleSentiment struct {
	Score     float64
	Intensity float64
	Dominant  string
}

// calculateSimpleSentiment performs basic sentiment analysis
func (s *AnalyticsService) calculateSimpleSentiment(text string) SimpleSentiment {
	text = strings.ToLower(text)

	// Multi-language sentiment dictionaries
	sentimentWords := map[string]map[string][]string{
		"positive": {
			"en": {"love", "happy", "great", "wonderful", "amazing", "good", "excellent", "fantastic", "beautiful", "perfect", "joy", "excited", "grateful", "blessed", "awesome", "incredible", "outstanding", "brilliant", "splendid", "magnificent"},
			"es": {"amor", "feliz", "genial", "maravilloso", "increíble", "bueno", "excelente", "fantástico", "hermoso", "perfecto", "alegría", "emocionado", "agradecido", "bendecido", "asombroso", "increíble", "sobresaliente", "brillante", "espléndido", "magnífico"},
			"fr": {"amour", "heureux", "génial", "merveilleux", "incroyable", "bon", "excellent", "fantastique", "beau", "parfait", "joie", "excité", "reconnaissant", "béni", "formidable", "incroyable", "exceptionnel", "brillant", "splendide", "magnifique"},
			"de": {"liebe", "glücklich", "großartig", "wunderbar", "unglaublich", "gut", "ausgezeichnet", "fantastisch", "schön", "perfekt", "freude", "aufgeregt", "dankbar", "gesegnet", "großartig", "unglaublich", "hervorragend", "brillant", "prächtig", "magnifik"},
			"it": {"amore", "felice", "fantastico", "meraviglioso", "incredibile", "buono", "eccellente", "fantastico", "bello", "perfetto", "gioia", "eccitato", "grato", "benedetto", "fantastico", "incredibile", "eccezionale", "brillante", "splendido", "magnifico"},
			"pt": {"amor", "feliz", "ótimo", "maravilhoso", "incrível", "bom", "excelente", "fantástico", "bonito", "perfeito", "alegria", "empolgado", "grato", "abençoado", "incrível", "inacreditável", "excepcional", "brilhante", "esplêndido", "magnífico"},
			"ru": {"любовь", "счастливый", "отличный", "чудесный", "удивительный", "хороший", "отличный", "фантастический", "красивый", "идеальный", "радость", "взволнованный", "благодарный", "благословенный", "потрясающий", "невероятный", "выдающийся", "блестящий", "великолепный", "величественный"},
			"ja": {"愛", "幸せ", "素晴らしい", "素敵", "信じられない", "良い", "優秀", "素晴らしい", "美しい", "完璧", "喜び", "興奮", "感謝", "祝福", "素晴らしい", "信じられない", "卓越", "輝かしい", "華麗", "壮大"},
			"ko": {"사랑", "행복", "훌륭한", "멋진", "놀라운", "좋은", "훌륭한", "환상적인", "아름다운", "완벽한", "기쁨", "흥분", "감사한", "축복받은", "놀라운", "믿을 수 없는", "뛰어난", "빛나는", "화려한", "장엄한"},
			"zh": {"爱", "快乐", "伟大", "精彩", "惊人", "好", "优秀", "精彩", "美丽", "完美", "喜悦", "兴奋", "感激", "祝福", "惊人", "难以置信", "杰出", "辉煌", "华丽", "宏伟"},
		},
		"negative": {
			"en": {"sad", "angry", "terrible", "awful", "bad", "horrible", "disappointed", "frustrated", "upset", "worried", "depressed", "anxious", "scared", "lonely", "hurt", "pain", "suffering", "miserable", "hopeless", "desperate"},
			"es": {"triste", "enojado", "terrible", "horrible", "malo", "horrible", "decepcionado", "frustrado", "molesto", "preocupado", "deprimido", "ansioso", "asustado", "solo", "herido", "dolor", "sufrimiento", "miserable", "desesperado", "desesperado"},
			"fr": {"triste", "fâché", "terrible", "affreux", "mauvais", "horrible", "déçu", "frustré", "contrarié", "inquiet", "déprimé", "anxieux", "effrayé", "seul", "blessé", "douleur", "souffrance", "misérable", "désespéré", "désespéré"},
			"de": {"traurig", "wütend", "schrecklich", "furchtbar", "schlecht", "schrecklich", "enttäuscht", "frustriert", "verärgert", "besorgt", "deprimiert", "ängstlich", "verängstigt", "einsam", "verletzt", "schmerz", "leiden", "elend", "hoffnungslos", "verzweifelt"},
			"it": {"triste", "arrabbiato", "terribile", "orribile", "cattivo", "orribile", "deluso", "frustrato", "turbato", "preoccupato", "depresso", "ansioso", "spaventato", "solo", "ferito", "dolore", "sofferenza", "miserabile", "disperato", "disperato"},
			"pt": {"triste", "irritado", "terrível", "horrível", "ruim", "horrível", "decepcionado", "frustrado", "chateado", "preocupado", "deprimido", "ansioso", "assustado", "sozinho", "machucado", "dor", "sofrimento", "miserável", "desesperado", "desesperado"},
			"ru": {"грустный", "злой", "ужасный", "ужасный", "плохой", "ужасный", "разочарованный", "разочарованный", "расстроенный", "обеспокоенный", "подавленный", "тревожный", "испуганный", "одинокий", "раненый", "боль", "страдание", "несчастный", "безнадежный", "отчаянный"},
			"ja": {"悲しい", "怒った", "ひどい", "恐ろしい", "悪い", "恐ろしい", "失望", "イライラ", "動揺", "心配", "落ち込んだ", "不安", "怖い", "孤独", "傷ついた", "痛み", "苦しみ", "惨め", "絶望的", "絶望的"},
			"ko": {"슬픈", "화난", "끔찍한", "무서운", "나쁜", "끔찍한", "실망한", "좌절한", "화난", "걱정하는", "우울한", "불안한", "무서워하는", "외로운", "상처받은", "고통", "고통", "비참한", "절망적인", "절망적인"},
			"zh": {"悲伤", "愤怒", "可怕", "可怕", "坏", "可怕", "失望", "沮丧", "心烦", "担心", "沮丧", "焦虑", "害怕", "孤独", "受伤", "痛苦", "痛苦", "悲惨", "绝望", "绝望"},
		},
	}

	// Detect language (simplified - in production, use a proper language detection library)
	detectedLang := s.detectLanguage(text)

	// Get sentiment words for detected language, fallback to English
	positiveWords, ok := sentimentWords["positive"][detectedLang]
	if !ok {
		positiveWords = sentimentWords["positive"]["en"]
	}

	negativeWords, ok := sentimentWords["negative"][detectedLang]
	if !ok {
		negativeWords = sentimentWords["negative"]["en"]
	}

	positiveCount := 0
	negativeCount := 0

	// Count sentiment words
	for _, word := range positiveWords {
		if strings.Contains(text, word) {
			positiveCount++
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(text, word) {
			negativeCount++
		}
	}

	// Calculate sentiment score
	total := positiveCount + negativeCount
	if total == 0 {
		return SimpleSentiment{
			Score:     0.5,
			Intensity: 0.1,
			Dominant:  "neutral",
		}
	}

	score := float64(positiveCount) / float64(total)
	intensity := float64(total) / 10.0 // Normalize intensity

	dominant := "neutral"
	if score > 0.6 {
		dominant = "positive"
	} else if score < 0.4 {
		dominant = "negative"
	}

	return SimpleSentiment{
		Score:     score,
		Intensity: math.Min(intensity, 1.0),
		Dominant:  dominant,
	}
}

// detectLanguage performs basic language detection
func (s *AnalyticsService) detectLanguage(text string) string {
	// Simple language detection based on character sets
	// In production, use a proper language detection library like "github.com/bbalet/stopwords"

	// Check for Chinese characters
	if strings.ContainsAny(text, "的一是在不了有和人这中大为上个国我以要他时来用们生到作地于出就分对成会可主发年动同工也能下过子说产种面而方后多定行学法所民得经十三之进着等部度家电力里如水化高自二理起小物现实加量都两体制机当使点从业本去把性好应开它合还因由其些然前外天政四日那社义事平形相全表间样与关各重新线内数正心反你明看原又么利比或但质气第向道命此变条只没结解问意建月公无系军很情者最立代想已通并提直题党程展五果料象员革位入常文总次品式活设及管特件长求老头基资边流路级少图山统接知较将组见计别她手角期根论运农指几九区强放决西被干做必战先回则任取据处队南给色光门即保治北造百规热领七海口东导器压志世金增争济阶油思术极交受联什认六共权收证改清己美再采转更单风切打白教速花带安场身车例真务具万每目至达走积示议声报斗完类八离华名确才科张信马节话米整空元况今集温传土许步群广石记需段研界拉林律叫且究观越织装影算低持音众书布复容儿须际商非验连断深难近矿千周委素技备半办青省列习响约支般史感劳便团往酸历市克何除消构府称太准精值号率族维划选标写存候毛亲快效斯院查江型眼王按格养易置派层片始却专状育厂京识适属圆包火住调满县局照参红细引听该铁价严龙飞") {
		return "zh"
	}

	// Check for Japanese characters
	if strings.ContainsAny(text, "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんアイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン") {
		return "ja"
	}

	// Check for Korean characters
	if strings.ContainsAny(text, "가나다라마바사아자차카타파하거너더러머버서어저처커터퍼허기니디리미비시이지치키티피히구누두루무부수우주추쿠투푸후그느드르므브스으즈츠크트프흐긔늬듸리미비시이지치키티피히") {
		return "ko"
	}

	// Check for Russian characters
	if strings.ContainsAny(text, "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ") {
		return "ru"
	}

	// Check for common European language patterns
	// Spanish: common words like "el", "la", "de", "que", "y"
	if strings.Contains(text, " el ") || strings.Contains(text, " la ") || strings.Contains(text, " de ") || strings.Contains(text, " que ") || strings.Contains(text, " y ") {
		return "es"
	}

	// French: common words like "le", "la", "de", "et", "que"
	if strings.Contains(text, " le ") || strings.Contains(text, " la ") || strings.Contains(text, " de ") || strings.Contains(text, " et ") || strings.Contains(text, " que ") {
		return "fr"
	}

	// German: common words like "der", "die", "das", "und", "in"
	if strings.Contains(text, " der ") || strings.Contains(text, " die ") || strings.Contains(text, " das ") || strings.Contains(text, " und ") || strings.Contains(text, " in ") {
		return "de"
	}

	// Italian: common words like "il", "la", "di", "e", "che"
	if strings.Contains(text, " il ") || strings.Contains(text, " la ") || strings.Contains(text, " di ") || strings.Contains(text, " e ") || strings.Contains(text, " che ") {
		return "it"
	}

	// Portuguese: common words like "o", "a", "de", "e", "que"
	if strings.Contains(text, " o ") || strings.Contains(text, " a ") || strings.Contains(text, " de ") || strings.Contains(text, " e ") || strings.Contains(text, " que ") {
		return "pt"
	}

	// Default to English
	return "en"
}

// EmotionalAnalysis represents emotional pattern analysis
type EmotionalAnalysis struct {
	Regulation float64
	Empathy    float64
	MoodImpact float64
}

// analyzeEmotionalPatterns analyzes emotional patterns in conversations
func (s *AnalyticsService) analyzeEmotionalPatterns(ctx context.Context, messages []*models.Message) (*EmotionalAnalysis, error) {
	conversationText := s.formatConversationForAnalysis(messages)

	prompt := fmt.Sprintf(`Analyze the emotional patterns in this conversation:

CONVERSATION:
%s

Analyze and rate from 0.0 to 1.0:
1. Emotional regulation (how well emotions are managed)
2. Empathy response (how empathetic the responses are)
3. Mood impact (how the conversation affects mood)

Respond with JSON:
{
  "regulation": 0.0-1.0,
  "empathy": 0.0-1.0,
  "mood_impact": 0.0-1.0,
  "analysis": "detailed analysis"
}`,
		conversationText)

	llmMessages := []LLMMessage{
		{Role: "system", Content: "You are an emotional pattern analyzer. Respond only with valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := s.grokService.SendMiniMessage(ctx, llmMessages)
	if err != nil {
		return &EmotionalAnalysis{
			Regulation: 0.5,
			Empathy:    0.5,
			MoodImpact: 0.5,
		}, nil
	}

	var analysis EmotionalAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return &EmotionalAnalysis{
			Regulation: 0.5,
			Empathy:    0.5,
			MoodImpact: 0.5,
		}, nil
	}

	return &analysis, nil
}

// formatConversationForAnalysis formats conversation for AI analysis
func (s *AnalyticsService) formatConversationForAnalysis(messages []*models.Message) string {
	var formatted []string

	for _, msg := range messages {
		if msg.Text == nil {
			continue
		}

		sender := "User"
		if msg.SenderType == "companion" {
			sender = "Companion"
		}

		formatted = append(formatted, fmt.Sprintf("%s: %s", sender, *msg.Text))
	}

	return strings.Join(formatted, "\n")
}

// Gamification Methods

// ProcessUserProgress processes and updates user progress
func (s *AnalyticsService) ProcessUserProgress(ctx context.Context, userID, companionID string, sessionData *SessionData) error {
	// Get current progress
	progress, err := s.repo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		// Create new progress
		progress = &models.UserProgress{
			UserID:      userID,
			CompanionID: companionID,
			CreatedAt:   time.Now(),
		}
	}

	// Calculate experience points
	experienceGained := s.calculateExperiencePoints(sessionData)
	progress.TotalExperience += experienceGained

	// Update level
	progress.CurrentLevel = s.calculateLevel(progress.TotalExperience)
	progress.LevelProgress = s.calculateLevelProgress(progress.TotalExperience)
	progress.ExperienceToNext = s.calculateExperienceToNext(progress.TotalExperience)

	// Update session statistics
	progress.TotalConversations++
	progress.TotalMessages += sessionData.MessageCount
	progress.TotalTimeSpent += sessionData.Duration

	// Update average session length
	if progress.TotalConversations > 0 {
		progress.AverageSessionLength = progress.TotalTimeSpent / time.Duration(progress.TotalConversations)
	}

	// Update streak
	s.updateStreak(progress)

	// Update achievement progress
	s.updateAchievementProgress(ctx, progress, sessionData)

	// Save progress
	return s.repo.UpsertUserProgress(ctx, progress)
}

// calculateExperiencePoints calculates experience points for a session
func (s *AnalyticsService) calculateExperiencePoints(sessionData *SessionData) int {
	basePoints := 10

	// Bonus for session duration
	durationBonus := int(sessionData.Duration.Minutes()) / 5 // 1 point per 5 minutes

	// Bonus for message count
	messageBonus := sessionData.MessageCount / 2 // 1 point per 2 messages

	// Bonus for conversation quality
	qualityBonus := int(sessionData.ResponseQuality * 20) // Up to 20 points for quality

	// Bonus for engagement
	engagementBonus := 0
	if sessionData.Duration > 10*time.Minute {
		engagementBonus = 5
	}
	if sessionData.MessageCount > 10 {
		engagementBonus += 5
	}

	return basePoints + durationBonus + messageBonus + qualityBonus + engagementBonus
}

// calculateLevel calculates user level based on experience
func (s *AnalyticsService) calculateLevel(experience int) int {
	// Level formula: level = sqrt(experience / 100) + 1
	level := int(math.Sqrt(float64(experience)/100.0)) + 1
	if level < 1 {
		level = 1
	}
	return level
}

// calculateLevelProgress calculates progress within current level
func (s *AnalyticsService) calculateLevelProgress(experience int) float64 {
	currentLevel := s.calculateLevel(experience)
	experienceForCurrentLevel := (currentLevel - 1) * (currentLevel - 1) * 100
	experienceForNextLevel := currentLevel * currentLevel * 100

	if experienceForNextLevel == experienceForCurrentLevel {
		return 1.0
	}

	progress := float64(experience-experienceForCurrentLevel) / float64(experienceForNextLevel-experienceForCurrentLevel)
	return math.Min(progress, 1.0)
}

// calculateExperienceToNext calculates experience needed for next level
func (s *AnalyticsService) calculateExperienceToNext(experience int) int {
	currentLevel := s.calculateLevel(experience)
	experienceForNextLevel := currentLevel * currentLevel * 100
	return experienceForNextLevel - experience
}

// updateStreak updates user streak information
func (s *AnalyticsService) updateStreak(progress *models.UserProgress) {
	today := time.Now().Truncate(24 * time.Hour)
	lastActivity := progress.LastActivityDate.Truncate(24 * time.Hour)

	if today.Equal(lastActivity) {
		// Already updated today
		return
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
}

// updateAchievementProgress updates achievement progress
func (s *AnalyticsService) updateAchievementProgress(ctx context.Context, progress *models.UserProgress, sessionData *SessionData) {
	// Get achievement definitions
	definitions, err := s.repo.GetAchievementDefinitions(ctx, "")
	if err != nil {
		return
	}

	// Update progress for each achievement
	for _, definition := range definitions {
		if progress.AchievementProgress == nil {
			progress.AchievementProgress = make(map[string]float64)
		}

		// Check if achievement is already earned
		earned, err := s.repo.CheckAchievementEarned(ctx, progress.UserID, progress.CompanionID, definition.ID)
		if err != nil || earned {
			continue
		}

		// Calculate progress based on achievement criteria
		progressValue := s.calculateAchievementProgress(definition, progress, sessionData)
		progress.AchievementProgress[definition.ID] = progressValue

		// Check if achievement is completed
		if progressValue >= definition.Criteria.Target {
			s.awardAchievement(ctx, progress, &definition)
		}
	}
}

// calculateAchievementProgress calculates progress for an achievement
func (s *AnalyticsService) calculateAchievementProgress(definition models.AchievementDefinition, progress *models.UserProgress, sessionData *SessionData) float64 {
	switch definition.Criteria.Type {
	case "total_messages":
		return float64(progress.TotalMessages)
	case "total_sessions":
		return float64(progress.TotalConversations)
	case "total_time":
		return progress.TotalTimeSpent.Hours()
	case "streak":
		return float64(progress.CurrentStreak)
	case "level":
		return float64(progress.CurrentLevel)
	default:
		return 0.0
	}
}

// awardAchievement awards an achievement to a user
func (s *AnalyticsService) awardAchievement(ctx context.Context, progress *models.UserProgress, definition *models.AchievementDefinition) {
	achievement := &models.UserAchievement{
		UserID:          progress.UserID,
		CompanionID:     progress.CompanionID,
		AchievementID:   definition.ID,
		AchievementType: definition.Type,
		Title:           definition.Title,
		Description:     definition.Description,
		IconURL:         definition.IconURL,
		Points:          definition.Points,
		Rarity:          definition.Rarity,
		Context:         make(map[string]any),
	}

	// Save achievement
	err := s.repo.InsertUserAchievement(ctx, achievement)
	if err != nil {
		return
	}

	// Update progress
	progress.TotalAchievements++
	if definition.Rarity == "rare" || definition.Rarity == "epic" || definition.Rarity == "legendary" {
		progress.RareAchievements++
	}

	// Add bonus experience
	progress.TotalExperience += definition.Points * 10
}

// GetUserDashboardData gets comprehensive dashboard data for a user
func (s *AnalyticsService) GetUserDashboardData(ctx context.Context, userID, companionID string) (*models.UserDashboardData, error) {
	// Get user progress
	progress, err := s.repo.GetUserProgress(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	// Get recent achievements
	achievements, err := s.repo.GetUserAchievements(ctx, userID, companionID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}

	// Get relationship analytics
	relationshipAnalytics, err := s.repo.GetRelationshipAnalytics(ctx, userID, companionID)
	if err != nil {
		// Create empty analytics if not found
		relationshipAnalytics = &models.RelationshipAnalytics{
			UserID:      userID,
			CompanionID: companionID,
		}
	}

	// Get engagement trends
	trends, err := s.repo.GetEngagementTrends(ctx, userID, companionID, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement trends: %w", err)
	}

	// Get user statistics
	statistics, err := s.repo.GetUserStatistics(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	// Get streak information
	streakInfo, err := s.repo.GetStreakInformation(ctx, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak information: %w", err)
	}

	// Generate recommendations
	recommendations := s.generateRecommendations(progress, relationshipAnalytics, statistics)

	// Get next milestones
	nextMilestones := s.getNextMilestones(progress, relationshipAnalytics)

	dashboard := &models.UserDashboardData{
		UserID:                userID,
		CompanionID:           companionID,
		Progress:              progress,
		RecentAchievements:    achievements,
		RelationshipAnalytics: relationshipAnalytics,
		EngagementTrends:      trends,
		Recommendations:       recommendations,
		NextMilestones:        nextMilestones,
		Statistics:            statistics,
		StreakInfo:            streakInfo,
		LastUpdated:           time.Now(),
	}

	return dashboard, nil
}

// generateRecommendations generates personalized recommendations
func (s *AnalyticsService) generateRecommendations(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics, statistics *models.UserStatistics) []models.Recommendation {
	var recommendations []models.Recommendation

	// Recommendation based on session frequency
	if statistics.TotalSessions < 5 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "engagement",
			Title:       "Start Your Journey",
			Description: "Try having a conversation every day to build a stronger connection.",
			Priority:    1,
			Confidence:  0.9,
			Action:      "start_daily_conversation",
		})
	}

	// Recommendation based on session length
	if statistics.AverageSessionLength < 5*time.Minute {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "quality",
			Title:       "Deeper Conversations",
			Description: "Try spending more time in conversations to explore topics more deeply.",
			Priority:    2,
			Confidence:  0.8,
			Action:      "extend_conversation_time",
		})
	}

	// Recommendation based on streak
	if progress.CurrentStreak < 3 {
		recommendations = append(recommendations, models.Recommendation{
			Type:        "consistency",
			Title:       "Build a Streak",
			Description: "Try to have conversations for 3 days in a row to build momentum.",
			Priority:    1,
			Confidence:  0.9,
			Action:      "build_streak",
		})
	}

	return recommendations
}

// getNextMilestones gets upcoming milestones for the user
func (s *AnalyticsService) getNextMilestones(progress *models.UserProgress, relationshipAnalytics *models.RelationshipAnalytics) []models.StageMilestone {
	var milestones []models.StageMilestone

	// Level milestone
	nextLevel := progress.CurrentLevel + 1
	experienceForNextLevel := nextLevel * nextLevel * 100
	levelProgress := float64(progress.TotalExperience) / float64(experienceForNextLevel)

	milestones = append(milestones, models.StageMilestone{
		ID:          fmt.Sprintf("level_%d", nextLevel),
		Title:       fmt.Sprintf("Reach Level %d", nextLevel),
		Description: fmt.Sprintf("Gain %d more experience points to reach level %d", experienceForNextLevel-progress.TotalExperience, nextLevel),
		Completed:   false,
		Progress:    levelProgress,
	})

	// Streak milestone
	nextStreakMilestone := ((progress.CurrentStreak / 7) + 1) * 7
	streakProgress := float64(progress.CurrentStreak%7) / 7.0

	milestones = append(milestones, models.StageMilestone{
		ID:          fmt.Sprintf("streak_%d", nextStreakMilestone),
		Title:       fmt.Sprintf("%d-Day Streak", nextStreakMilestone),
		Description: fmt.Sprintf("Maintain a %d-day conversation streak", nextStreakMilestone),
		Completed:   false,
		Progress:    streakProgress,
	})

	return milestones
}

// GetEngagementTrends gets engagement trends for a user
func (s *AnalyticsService) GetEngagementTrends(ctx context.Context, userID, companionID string, days int) ([]models.EngagementTrendPoint, error) {
	return s.repo.GetEngagementTrends(ctx, userID, companionID, days)
}

// GetUserStatistics gets user statistics
func (s *AnalyticsService) GetUserStatistics(ctx context.Context, userID, companionID string) (*models.UserStatistics, error) {
	return s.repo.GetUserStatistics(ctx, userID, companionID)
}

// GetRelationshipAnalytics gets relationship analytics
func (s *AnalyticsService) GetRelationshipAnalytics(ctx context.Context, userID, companionID string) (*models.RelationshipAnalytics, error) {
	return s.repo.GetRelationshipAnalytics(ctx, userID, companionID)
}

// GetPlatformAnalytics gets platform-wide analytics
func (s *AnalyticsService) GetPlatformAnalytics(ctx context.Context, days int) (map[string]any, error) {
	return s.repo.GetPlatformAnalytics(ctx, days)
}
