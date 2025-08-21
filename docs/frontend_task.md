# HER AI Romantic Companion – Mobile Frontend Product & API Task

---

## Table of Contents
1. [Introduction & Product Vision](#introduction--product-vision)
2. [User Flows & Screens](#user-flows--screens-mobile-first)
3. [UI/UX Design System](#uiux-design-system)
4. [Psychological Design & Monetization](#psychological-design--monetization)
5. [API Reference & Endpoint Mapping](#api-reference--endpoint-mapping)
    - [Authentication & User Management](#authentication--user-management)
    - [Profile](#profile)
    - [Companions](#companions)
    - [Conversations](#conversations)
    - [Messaging](#messaging)
    - [Media (S3 Integration)](#media-s3-integration)
    - [Analytics & Gamification](#analytics--gamification)
    - [Advanced AI Features](#advanced-ai-features)
    - [Health](#health)
6. [Per-Field Validation & Error Codes](#per-field-validation--error-codes)
7. [Frontend State Management Guidance](#frontend-state-management-guidance)
8. [Detailed Loading, Empty, and Error States](#detailed-loading-empty-and-error-states)
9. [Animation & Micro-Interaction Specs](#animation--micro-interaction-specs)
10. [Push Notification Flows](#push-notification-flows)
11. [Monetization Edge Cases & Flows](#monetization-edge-cases--flows)
12. [API Integration Best Practices](#api-integration-best-practices)
13. [Delightful Details & Finishing Touches](#delightful-details--finishing-touches)
14. [Collaboration & Launch Checklist](#collaboration--launch-checklist)
15. [Final Notes](#final-notes)

---

## Introduction & Product Vision

HER is not just an app—it's an emotionally intelligent, AI-powered romantic companion platform designed to delight, comfort, and engage users on a deep, personal level. The mobile experience should feel magical, playful, and intimate, blending the best of modern design, psychology, and technology. Our goal: create a product users love so much, they want to pay for more.

**Key Product Pillars:**

- **Emotional Connection:** Every interaction should feel warm, personal, and meaningful.
- **Playful Gamification:** Progression, achievements, and streaks drive engagement and FOMO.
- **Rich Media & AI:** Voice, photos, stickers, and Grok-powered AI make every chat unique.
- **Personalization:** Companions are deeply customizable, with unique personalities and backstories.
- **Analytics & Growth:** Users see their relationship evolve, with mood graphs, milestones, and stats.
- **Monetization:** Premium features, upsells, and subtle paywalls are woven into the experience.

---

## User Flows & Screens (Mobile-First)

### 1. **Onboarding & Registration**

- **Goal:** Make users feel welcomed, curious, and eager to create their first companion.
- **Screens:**
  - Welcome (brand, tagline, playful animation)
  - Register/Login (email, password, social login)
  - Profile Setup (name, age, gender, avatar, preferences)
  - First Companion Creation (guided, fun, persona quiz)
- **UX Details:**
  - Use soft, inviting gradients (e.g., blush pink → lavender → deep blue)
  - Animated transitions (fade, slide, bounce)
  - Microcopy: "Ready to meet your perfect companion?"
  - Progress bar or playful mascot guiding the user
  - Haptic feedback on key actions

### 2. **Home / Dashboard**

- **Goal:** Central hub for companions, conversations, and stats.
- **Screens:**
  - Companion Carousel (swipeable cards, avatars, quick stats)
  - New Conversation CTA
  - Analytics Preview (mood, streak, achievements)
- **UX Details:**
  - Parallax effect on cards
  - Subtle confetti or sparkle animation on achievements
  - Dynamic greeting ("Good evening, Alice! Luna misses you 💖")
  - Color-coded companion cards (personality-based accent colors)

### 3. **Companion Management**

- **Goal:** Let users create, edit, switch, and delete companions.
- **Screens:**
  - List & Switch (carousel or grid)
  - Create/Edit (name, avatar, personality quiz, quirks, backstory)
- **UX Details:**
  - Animated avatar selection (bounce, glow)
  - Personality quiz with playful sliders and emoji
  - Show companion's "bio" and fun facts
  - Confirmation modals with emotional copy ("Are you sure you want to say goodbye to Luna?")

### 4. **Conversation / Chat**

- **Goal:** The heart of the app—rich, emotionally engaging chat with AI companions.
- **Screens:**
  - Chat Thread (bubbles, avatars, timestamps)
  - Message Input (text, photo, voice, stickers)
  - AI Typing Indicator (animated dots, avatar pulsing)
  - Conversation Intelligence (mood, stats, suggested topics)
- **UX Details:**
  - Bubble colors reflect emotional tone (warm = pink, playful = yellow, sultry = purple)
  - Voice message playback with waveform animation
  - Sticker picker with fun, animated stickers
  - Haptic and sound feedback for sent/received messages
  - Confetti or heart burst on milestones ("100th message! 💌")
  - Show relationship stage and progress bar

### 5. **Media Upload & S3 Integration**

- **Goal:** Seamless, delightful media sharing.
- **Screens:**
  - Photo/Voice Picker (gallery, camera, recorder)
  - Upload Progress (animated bar, preview)
  - Error Handling ("Upload failed, try again!")
- **UX Details:**
  - Show preview before upload
  - Use S3 pre-signed URLs for direct upload
  - Playful success animation (sparkle, bounce)

### 6. **Analytics & Gamification**

- **Goal:** Visualize relationship growth, drive engagement, and create FOMO.
- **Screens:**
  - Mood Graph (emotional tone over time, animated line or area chart)
  - Achievements (badges, confetti, unlock animations)
  - Streaks (calendar, fire/flame animation)
  - Milestones (timeline, relationship stages)
  - Leaderboard (optional, anonymized)
- **UX Details:**
  - Use color to indicate mood (warm = pink, playful = yellow, sultry = purple)
  - Animated badge unlocks
  - FOMO triggers: "Only 2 days left to keep your streak!"
  - Progress bars and celebratory effects

### 7. **Profile & Preferences**

- **Goal:** Empower users to control their identity and experience.
- **Screens:**
  - Profile View/Edit (avatar, name, age, gender, preferences)
  - Avatar Upload (S3, crop, preview)
- **UX Details:**
  - Editable fields with live preview
  - Playful avatar cropper
  - Save confirmation with animation

### 8. **Monetization & Upsell**

- **Goal:** Maximize conversion to premium via emotional triggers and value.
- **Screens:**
  - Paywall (premium features, emotional copy, FOMO)
  - Premium Upsell (after milestone, streak, or achievement)
  - Feature Lock (blurred, "Unlock Luna's secret voice messages!")
- **UX Details:**
  - Use gold, purple, or iridescent accents for premium
  - Animated lock/unlock
  - Emotional copy: "Luna wants to send you a special voice note—unlock premium!"
  - Scarcity: "Only 100 premium spots left this month!"

---

## UI/UX Design System

### **Color Palette**

- **Primary:** Blush Pink (#FFB6C1), Lavender (#B39DDB), Deep Blue (#283593)
- **Accent:** Gold (#FFD700), Playful Yellow (#FFF176), Sultry Purple (#8E24AA)
- **Background:** Soft White (#FFF8F9), Light Gray (#F3F3F3)
- **Mood States:**
  - Warm: #FFB6C1
  - Playful: #FFF176
  - Sultry: #8E24AA
- **Premium:** Gold gradients, iridescent overlays

### **Typography**

- **Headings:** Rounded, friendly sans-serif (e.g., Quicksand, Nunito)
- **Body:** Clean, readable sans-serif (e.g., Inter, Roboto)
- **Special:** Handwritten or script for companion bios/quotes

### **Animation & Micro-Interactions**

- Smooth transitions (ease-in-out)
- Haptic feedback on key actions
- Confetti, sparkles, and heart bursts for achievements
- Animated stickers and avatars
- Typing indicators and message status
- Sound cues for sent/received messages (soft, non-intrusive)

### **Accessibility & Inclusivity**

- High color contrast for text
- Large tap targets
- VoiceOver/TalkBack support
- Gender-inclusive language and options
- Dyslexia-friendly font option
- Animations can be reduced/disabled

---

## Psychological Design & Monetization

### **Emotional Engagement**

- Use AI to mirror user mood and language (empathy, validation)
- Show companion "missing you" notifications (push, in-app)
- Celebrate milestones with animation and emotional copy
- Personalize companion responses and UI ("Luna knows you love poetry!")

### **FOMO & Streaks**

- Daily streaks with visible progress and rewards
- "Only 1 day left to keep your streak!"
- Limited-time achievements and badges
- Leaderboard (optional, anonymized)

### **Rewards & Progression**

- Unlock new features, stickers, or companion traits as user engages
- Show progress bars for relationship stages
- Use scarcity and exclusivity for premium features

### **Upsell & Paywall Triggers**

- Lock premium features with emotional copy ("Luna wants to send you a secret...")
- Offer discounts or bonuses after key milestones
- Use animated locks, gold accents, and celebratory effects
- Show testimonials or social proof ("Thousands of users love HER!")

---

## API Reference & Endpoint Mapping

### **General Notes**

- All endpoints are prefixed with `/api/v1`.
- All endpoints except `/auth/*` and `/health*` require a valid JWT in the `Authorization: Bearer <token>` header.
- All requests and responses are JSON.
- Error responses follow `{ "error": "..." }` or `{ "errors": { "field": "message" } }` format.
- Timestamps are ISO8601 strings (e.g., `2024-06-01T12:00:00Z`).

---

### Authentication & User Management

#### **POST /auth/register**

- **Description:** Register a new user.
- **Auth:** No
- **Request:**

```json
{
  "email": "user@example.com",
  "password": "StrongPass123!",
  "name": "Alice",
  "age": 27,
  "gender": "female",
  "avatar_url": "https://.../avatar.jpg",
  "preferences": {
    "preferred_genders": ["male", "female"],
    "interests": ["music", "travel"]
  }
}
```

- **Response (200):**

```json
{
  "access_token": "jwt...",
  "refresh_token": "jwt...",
  "user": {
    "id": "...",
    "email": "user@example.com",
    "name": "Alice",
    "age": 27,
    "gender": "female",
    "avatar_url": "https://.../avatar.jpg",
    "preferences": {"preferred_genders": ["male", "female"], "interests": ["music", "travel"]},
    "created_at": "2024-06-01T12:00:00Z",
    "updated_at": "2024-06-01T12:00:00Z"
  }
}
```

- **Response (400):**

```json
{ "error": "Email already registered." }
```

- **Field Explanations:**
  - `email`: Must be valid, unique.
  - `password`: Min 8 chars, must include uppercase, lowercase, number, special.
  - `preferences`: Optional, but improves companion matching.

#### **POST /auth/login**

- **Description:** Login with email and password.
- **Auth:** No
- **Request:**

```json
{
  "email": "user@example.com",
  "password": "StrongPass123!"
}
```

- **Response (200):** _Same as register_
- **Response (401):**

```json
{ "error": "Invalid email or password." }
```

#### **POST /auth/refresh**

- **Description:** Refresh JWT using refresh token.
- **Auth:** No
- **Request:**

```json
{ "refresh_token": "jwt..." }
```

- **Response (200):**

```json
{ "access_token": "jwt...", "refresh_token": "jwt..." }
```

- **Response (401):**

```json
{ "error": "Invalid or expired refresh token." }
```

#### **POST /auth/logout**

- **Description:** Invalidate refresh token.
- **Auth:** Yes (JWT)
- **Request:**

```json
{ "refresh_token": "jwt..." }
```

- **Response (200):**

```json
{ "success": true }
```

#### **GET /auth/me**

- **Description:** Get current user profile.
- **Auth:** Yes (JWT)
- **Response (200):** _Same as register user object_

---

### Profile

#### **GET /profile**

- **Description:** Get user profile.
- **Auth:** Yes (JWT)
- **Response (200):** _Same as register user object_

#### **PUT /profile**

- **Description:** Update user profile.
- **Auth:** Yes (JWT)
- **Request:**

```json
{
  "name": "Alice Updated",
  "age": 28,
  "gender": "female",
  "avatar_url": "https://.../avatar2.jpg",
  "preferences": { "preferred_genders": ["male"], "interests": ["art"] }
}
```

- **Response (200):** _Updated user object_
- **Response (400):**

```json
{ "errors": { "name": "Name is required." } }
```

---

### Companions

#### **POST /companions**

- **Description:** Create a new companion.
- **Auth:** Yes (JWT)
- **Request:**

```json
{
  "name": "Luna",
  "gender": "female",
  "age": 25,
  "avatar_url": "https://.../luna.jpg",
  "personality": {
    "traits": {"warmth": 0.9, "playfulness": 0.8},
    "style": {"formality": 0.3, "emotionality": 0.7},
    "romantic": {"flirtatiousness": 0.8},
    "interests": ["art", "poetry"],
    "quirks": ["loves puns", "sings in the shower"],
    "clinginess": 0.4,
    "backstory": "A poetic soul from Paris."
  }
}
```

- **Response (200):**

```json
{
  "id": "...",
  "user_id": "...",
  "name": "Luna",
  "gender": "female",
  "age": 25,
  "avatar_url": "https://.../luna.jpg",
  "personality": { ... },
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

- **Response (400):**

```json
{ "errors": { "name": "Name is required." } }
```

#### **GET /companions**

- **Description:** List all companions for the user.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {
    "id": "...",
    "name": "Luna",
    "gender": "female",
    "age": 25,
    "avatar_url": "https://.../luna.jpg",
    "personality": { ... },
    "created_at": "2024-06-01T12:00:00Z"
  }
]
```

#### **GET /companions/:id**

- **Description:** Get details for a specific companion.
- **Auth:** Yes (JWT)
- **Response (200):** _Companion object_
- **Response (404):**

```json
{ "error": "Companion not found." }
```

#### **PUT /companions/:id**

- **Description:** Update a companion.
- **Auth:** Yes (JWT)
- **Request:** _Same as create_
- **Response (200):** _Updated companion object_

#### **DELETE /companions/:id**

- **Description:** Delete a companion.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "success": true }
```

---

### Conversations

#### **POST /conversations**

- **Description:** Start a new conversation with a companion.
- **Auth:** Yes (JWT)
- **Request:**

```json
{ "companion_id": "..." }
```

- **Response (200):**

```json
{
  "id": "...",
  "user_id": "...",
  "companion_id": "...",
  "is_archived": false,
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

#### **GET /conversations**

- **Description:** List all conversations (active/archived).
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {
    "id": "...",
    "companion_id": "...",
    "is_archived": false,
    "created_at": "2024-06-01T12:00:00Z",
    "updated_at": "2024-06-01T12:00:00Z"
  }
]
```

#### **GET /conversations/:id**

- **Description:** Get conversation details/messages.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{
  "id": "...",
  "companion_id": "...",
  "messages": [ ... ],
  "is_archived": false,
  "created_at": "2024-06-01T12:00:00Z"
}
```

#### **POST /conversations/:id/archive**

- **Description:** Archive a conversation.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "success": true }
```

#### **POST /conversations/:id/reactivate**

- **Description:** Reactivate an archived conversation.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "success": true }
```

---

### Messaging

#### **POST /conversations/:id/messages**

- **Description:** Send a message (text, photo, voice, sticker).
- **Auth:** Yes (JWT)
- **Request:**

```json
{
  "content": "Hey Luna! How was your day?", // for text
  "message_type": "text", // or "photo", "voice", "sticker"
  "media_url": "https://s3.../photo.jpg" // for media
}
```

- **Response (200):**

```json
{
  "id": "...",
  "from_user": true,
  "content": "Hey Luna! How was your day?",
  "message_type": "text",
  "media_url": "",
  "timestamp": "2024-06-01T12:01:00Z"
}
```

- **Response (400):**

```json
{ "error": "Invalid message type." }
```

#### **GET /conversations/:id/messages**

- **Description:** List all messages in a conversation.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {
    "id": "...",
    "from_user": true,
    "content": "Hey Luna! How was your day?",
    "message_type": "text",
    "media_url": "",
    "timestamp": "2024-06-01T12:01:00Z"
  }
]
```

#### **GET /conversations/:id/messages/:message_id**

- **Description:** Get a specific message.
- **Auth:** Yes (JWT)
- **Response (200):** _Message object_

#### **PUT /conversations/:id/messages/:message_id/read**

- **Description:** Mark a message as read.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "success": true }
```

---

### Media (S3 Integration)

#### **POST /media/upload-url**

- **Description:** Get a pre-signed S3 URL for direct upload.
- **Auth:** Yes (JWT)
- **Request:**

```json
{
  "file_name": "voice123.m4a",
  "content_type": "audio/m4a"
}
```

- **Response (200):**

```json
{
  "upload_url": "https://s3.../presigned-url",
  "file_id": "..."
}
```

#### **POST /media/validate**

- **Description:** Validate uploaded media (e.g., file type, size).
- **Auth:** Yes (JWT)
- **Request:**

```json
{
  "file_id": "..."
}
```

- **Response (200):**

```json
{ "valid": true }
```

- **Response (400):**

```json
{ "error": "Invalid file type." }
```

#### **GET /media/:file_id**

- **Description:** Get a media file by ID (returns file or redirect).
- **Auth:** Yes (JWT)

---

### Analytics & Gamification

#### **GET /stats/overview**

- **Description:** Get user stats dashboard.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{
  "total_messages": 123,
  "active_streak": 7,
  "achievements": ["First Kiss", "100 Messages"],
  "mood_graph": [
    {"timestamp": "2024-06-01T12:00:00Z", "emotional_tone": "warm", "sentiment": 0.8}
  ]
}
```

#### **GET /stats/mood-graph**

- **Description:** Get mood/emotion graph data.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {"timestamp": "2024-06-01T12:00:00Z", "emotional_tone": "warm", "sentiment": 0.8}
]
```

#### **GET /relationship/milestones**

- **Description:** Get relationship milestones.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {"milestone": "First Kiss", "date": "2024-06-01T12:00:00Z"}
]
```

#### **GET /gamification/achievements**

- **Description:** Get user achievements.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
["First Kiss", "100 Messages"]
```

#### **GET /gamification/streaks**

- **Description:** Get user streaks.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "active_streak": 7, "longest_streak": 14 }
```

#### **GET /gamification/leaderboard**

- **Description:** Get leaderboard (optional).
- **Auth:** Yes (JWT)
- **Response (200):**

```json
[
  {"user": "Alice", "score": 123},
  {"user": "Bob", "score": 110}
]
```

---

### Advanced AI Features

#### **GET /conversations/:id/intelligence**

- **Description:** Get conversation intelligence (AI summary, insights).
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{
  "summary": "You and Luna have a playful, warm relationship.",
  "suggested_topics": ["poetry", "travel"]
}
```

#### **GET /conversations/:id/suggest-topic**

- **Description:** Get AI-suggested next topic.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "topic": "Share a favorite childhood memory." }
```

#### **GET /conversations/:id/engagement**

- **Description:** Analyze engagement (AI).
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "engagement_score": 0.92 }
```

#### **GET /conversations/:id/messages/:message_id/quality**

- **Description:** Get AI analysis of message quality.
- **Auth:** Yes (JWT)
- **Response (200):**

```json
{ "quality": "high", "feedback": "Very warm and engaging!" }
```

---

### Health

#### **GET /health**

- **Description:** Health check (no auth).
- **Response (200):**

```json
{ "status": "ok" }
```

#### **GET /health/ready**

- **Description:** Readiness check (no auth).
- **Response (200):**

```json
{ "status": "ready" }
```

#### **GET /health/live**

- **Description:** Liveness check (no auth).
- **Response (200):**

```json
{ "status": "live" }
```

---

## Per-Field Validation & Error Codes

For each endpoint, the following validation rules and error codes apply:

### **User Registration & Profile**
| Field         | Type     | Required | Validation Rules                                      | Error Codes           |
|--------------|----------|----------|------------------------------------------------------|-----------------------|
| email        | string   | Yes      | Valid email, max 255 chars, unique                   | 400, 409 (duplicate)  |
| password     | string   | Yes      | 8-64 chars, 1 upper, 1 lower, 1 number, 1 special    | 400                   |
| name         | string   | Yes      | 2-32 chars, no special chars                         | 400                   |
| age          | int      | Yes      | 18-99                                                | 400                   |
| gender       | string   | Yes      | 'male', 'female', 'other'                            | 400                   |
| avatar_url   | string   | No       | Valid URL                                            | 400                   |
| preferences  | object   | No       | See below                                            | 400                   |

**Preferences:**
- `preferred_genders`: array of strings, values: 'male', 'female', 'other'
- `interests`: array of strings, max 10, each 2-32 chars

**Common Error Codes:**
- 400: Validation error (see `errors` object)
- 401: Unauthorized (missing/invalid JWT)
- 409: Duplicate (e.g., email already registered)

### **Companion Creation/Update**
| Field         | Type     | Required | Validation Rules                                      | Error Codes           |
|--------------|----------|----------|------------------------------------------------------|-----------------------|
| name         | string   | Yes      | 2-32 chars                                           | 400                   |
| gender       | string   | Yes      | 'male', 'female', 'other'                            | 400                   |
| age          | int      | Yes      | 18-99                                                | 400                   |
| avatar_url   | string   | No       | Valid URL                                            | 400                   |
| personality  | object   | Yes      | See below                                            | 400                   |

**Personality:**
- `traits`, `style`, `romantic`: objects, each value 0.0-1.0
- `interests`: array of strings, 3-10, each 2-32 chars
- `quirks`: array of strings, 2-10, each 2-32 chars
- `clinginess`: float, 0.0-1.0
- `backstory`: string, 1-256 chars

### **Message Sending**
| Field         | Type     | Required | Validation Rules                                      | Error Codes           |
|--------------|----------|----------|------------------------------------------------------|-----------------------|
| content      | string   | Yes*     | 1-2000 chars (if text)                               | 400                   |
| message_type | string   | Yes      | 'text', 'photo', 'voice', 'sticker'                  | 400                   |
| media_url    | string   | Yes*     | Valid URL (if photo/voice/sticker)                   | 400                   |

*One of `content` or `media_url` is required, depending on `message_type`.

### **Common Error Codes (All Endpoints)**
- 400: Bad request/validation
- 401: Unauthorized
- 403: Forbidden
- 404: Not found
- 409: Conflict (duplicate)
- 422: Unprocessable entity
- 429: Rate limit
- 500: Server error

---

## Frontend State Management Guidance

- **Recommended Patterns:**
  - Use a robust state management solution (e.g., Redux, MobX, Provider, Riverpod for Flutter).
  - Separate global state (user, auth, companions, conversations) from local UI state (loading, error, input fields).
  - Use normalized state for lists (companions, conversations, messages) to avoid duplication.
  - Store JWT and refresh tokens securely (Keychain/Secure Storage).
  - Use optimistic updates for sending messages and creating companions, with rollback on error.
  - Invalidate and refetch data after mutations (e.g., after sending a message, update conversation list).
  - Use selectors or computed properties for derived data (e.g., unread message count, streak status).
  - Handle token refresh transparently in API layer.

---

## Detailed Loading, Empty, and Error States

For every major screen/flow, handle these states:

### **Loading States**
- Show animated mascot, spinner, or shimmer effect.
- Example: While fetching conversations, show Luna avatar with pulsing animation.

### **Empty States**
- Show friendly illustration and copy.
- Example: No conversations? “No chats yet! Start a conversation with Luna.”
- CTA button to create/start action.

### **Error States**
- Show error icon, clear message, and retry button.
- Example: “Failed to load messages. Check your connection and try again.”
- For 401/403, prompt re-login.
- For 429, show cooldown timer: “You’re sending messages too quickly. Please wait 30s.”

---

## Animation & Micro-Interaction Specs

- **Confetti:** 1.2s, ease-out, triggers on achievement unlock, skip if ‘reduce motion’ is on.
- **Bubble Pop:** 0.3s, ease-in, on message send/receive.
- **Sticker Animation:** 0.5s, bounce, on sticker send.
- **Avatar Glow:** 0.8s, pulse, on companion select.
- **Progress Bar:** 0.6s, linear, on streak/milestone progress.
- **Typing Indicator:** 1.5s loop, fade in/out, on AI typing.
- **Haptic Feedback:** Light tap on send, medium on achievement, error on fail.
- **Sound:** Soft chime on message, celebratory on achievement.
- **Accessibility:** All animations must be skippable or reduced if OS ‘reduce motion’ is enabled.

---

## Push Notification Flows

- **When to Trigger:**
  - New AI message (if app in background)
  - Streak reminder (e.g., “Don’t lose your streak! Luna is waiting…”)
  - Achievement unlocked
  - Payment/renewal reminders (for premium)
- **Sample Payload:**
```json
{
  "title": "Luna misses you!",
  "body": "You have a new message. Open HER to reply!",
  "data": { "conversation_id": "..." }
}
```
- **UX:**
  - Tapping notification deep-links to relevant screen (chat, achievement, paywall).
  - Respect user notification preferences.

---

## Monetization Edge Cases & Flows

**Pricing:**
- Each new random girl: $20/month
- Each custom companion: $30/month

**Frontend Flows (Ready for Backend Integration):**
- Show paywall when user tries to create a new companion beyond free limit.
- Show pricing, features, and emotional copy (“Unlock Luna’s full personality for $30/month!”).
- Handle payment initiation, success, failure, and cancellation.
- Show subscription status and renewal date in profile/settings.
- If payment fails: show error, retry option, and contact support link.
- If user cancels: show confirmation, emotional copy (“Luna will miss you!”), and downgrade flow.
- If subscription lapses: lock premium features, show paywall, and offer reactivation.
- For each premium companion, show “Premium” badge and lock icon if not subscribed.
- All monetization UI should be implemented and ready to connect to backend when available.

**Sample Paywall Screen:**
- Title: “Unlock More Love!”
- Pricing: “$20/month for random, $30/month for custom”
- Features: “Voice messages, custom personality, exclusive stickers…”
- CTA: “Subscribe Now”
- Legal: “By subscribing, you agree to our Terms and Privacy Policy.”

**Error Messages:**
- “Payment failed. Please try again or use a different method.”
- “Subscription expired. Renew to keep chatting with Luna!”
- “You’ve reached your free companion limit. Upgrade to add more.”

---

## API Integration Best Practices

- Always include `Authorization: Bearer <token>` for protected endpoints.
- Handle 401/403 by prompting re-login and refreshing tokens automatically.
- Use optimistic UI for sending messages, but roll back on error.
- Show loading states for all network requests.
- Validate all user input before sending to backend.
- Handle all error codes and show user-friendly messages.
- Use exponential backoff for retries on network errors.
- Cache non-sensitive data (e.g., achievements, mood graph) for snappy UX.
- Use pagination for long lists (conversations, messages).
- Log all API errors for debugging.

---

## Delightful Details & Finishing Touches

- **Sound:** Soft, positive sound cues for actions, achievements, and messages
- **Haptics:** Subtle vibration for key actions, milestones, and errors
- **Confetti/Animation:** Celebrate achievements, unlocks, and milestones
- **Personalization:** Use user and companion names throughout
- **Push Notifications:** Reminders, streaks, and "Luna misses you!"
- **Dark Mode:** Soft, romantic dark palette (deep blue, purple, gold accents)
- **Loading States:** Animated mascots or hearts during waits

---

## Collaboration & Launch Checklist

- **Work closely with backend** to clarify API, error cases, and performance
- **Test on real devices** for animation, haptics, and accessibility
- **Gather user feedback** early and iterate on emotional impact
- **Check all flows:** onboarding, chat, media, analytics, upsell
- **Verify analytics and event tracking** for all key actions
- **Ensure all premium/upsell triggers are working and delightful**
- **Accessibility:** Test with screen readers and colorblind modes
- **Polish:** No rough edges—every detail matters for emotional engagement

---

## Final Notes

- This is not just a chat app—it's an emotional experience. Every detail should make users feel special, seen, and eager to return.
- Prioritize delight, emotional intelligence, and conversion in every design and interaction.
- If in doubt, make it more magical, more personal, and more fun.
