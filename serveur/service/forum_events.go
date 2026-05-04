package service

import "sync"

type ForumEvent struct {
	Type    string `json:"type"`
	Action  string `json:"action"`
	SujetID string `json:"sujetId,omitempty"`
	ThreadID string `json:"threadId,omitempty"`
}

type ForumEventHub struct {
	mu          sync.Mutex
	subscribers map[chan ForumEvent]struct{}
}

func NewForumEventHub() *ForumEventHub {
	return &ForumEventHub{
		subscribers: make(map[chan ForumEvent]struct{}),
	}
}

func (h *ForumEventHub) Subscribe() chan ForumEvent {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan ForumEvent, 16)
	h.subscribers[ch] = struct{}{}
	return ch
}

func (h *ForumEventHub) Unsubscribe(ch chan ForumEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.subscribers[ch]; ok {
		delete(h.subscribers, ch)
		close(ch)
	}
}

func (h *ForumEventHub) Publish(event ForumEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}