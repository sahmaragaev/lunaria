package services

import (
	"sync"
	"time"
)

// TypingState holds the transient typing status for a conversation
type TypingState struct {
	IsTyping      bool
	MessageIndex  int
	TotalMessages int
	LastUpdate    time.Time
}

// TypingTracker manages in-memory typing states
type TypingTracker struct {
	mu sync.RWMutex
	m  map[string]TypingState
}

func NewTypingTracker() *TypingTracker {
	return &TypingTracker{m: make(map[string]TypingState)}
}

func (t *TypingTracker) SetStart(convID string) {
	t.mu.Lock()
	t.m[convID] = TypingState{IsTyping: true, MessageIndex: 0, TotalMessages: 0, LastUpdate: time.Now()}
	t.mu.Unlock()
}

func (t *TypingTracker) SetTotal(convID string, total int) {
	t.mu.Lock()
	state := t.m[convID]
	state.IsTyping = true
	state.TotalMessages = total
	state.LastUpdate = time.Now()
	t.m[convID] = state
	t.mu.Unlock()
}

func (t *TypingTracker) Update(convID string, index int, total int) {
	t.mu.Lock()
	state := t.m[convID]
	state.IsTyping = index < total-1
	state.MessageIndex = index
	state.TotalMessages = total
	state.LastUpdate = time.Now()
	t.m[convID] = state
	t.mu.Unlock()
}

func (t *TypingTracker) Stop(convID string) {
	t.mu.Lock()
	delete(t.m, convID)
	t.mu.Unlock()
}

func (t *TypingTracker) Get(convID string) (TypingState, bool) {
	t.mu.RLock()
	state, ok := t.m[convID]
	t.mu.RUnlock()
	return state, ok
}

var globalTypingTracker = NewTypingTracker()

// GetTypingTracker returns the process-wide typing tracker
func GetTypingTracker() *TypingTracker {
	return globalTypingTracker
}
