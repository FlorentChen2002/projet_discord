package service

import "sync"

// Structure représentant un événement dans le forum
type ForumEvent struct {
    Type string `json:"type"`
    Action string `json:"action"`
    SujetID string `json:"sujetId,omitempty"`
    ThreadID string `json:"threadId,omitempty"`
}

// Structure représentant un hub d'événements pour le forum
type ForumEventHub struct {
    mu sync.Mutex
    subs map[chan ForumEvent]struct{}
}

// fonction qui permet de créer un nouveau hub d'événements pour le forum
func NewForumEventHub() *ForumEventHub {
    return &ForumEventHub{subs: make(map[chan ForumEvent]struct{}),}
}

// fonction qui permet de s'abonner à un hub d'événements pour le forum
func (h *ForumEventHub) Subscribe() chan ForumEvent {
    h.mu.Lock()
    defer h.mu.Unlock()
    ch := make(chan ForumEvent, 16)
    h.subs[ch] = struct{}{}
    return ch
}

// fonction qui permet de se désabonner d'un hub d'événements pour le forum
func (h *ForumEventHub) Unsubscribe(ch chan ForumEvent) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if _, ok := h.subs[ch]; ok {
        delete(h.subs, ch)
        close(ch)
    }
}

// fonction qui permet de publier un événement dans le hub d'événements pour le forum
func (h *ForumEventHub) Publish(event ForumEvent) {
    h.mu.Lock()
    defer h.mu.Unlock()
    for ch := range h.subs {
        select {
        case ch <- event:
        default:
        }
    }
}