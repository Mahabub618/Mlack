package services

import (
	"Mlack/models"
	"sync"
	"time"
)

type EventService struct {
	clients map[chan models.CommitEvent]bool
	mu      sync.Mutex
}

func NewEventService() *EventService {
	return &EventService{
		clients: make(map[chan models.CommitEvent]bool),
	}
}

func (s *EventService) Broadcast(event models.CommitEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		// Use non-blocking send with timeout
		select {
		case client <- event:
		case <-time.After(5 * time.Second):
			// Client is too slow - remove them
			delete(s.clients, client)
			close(client)
		}
	}
}

func (s *EventService) AddClient(client chan models.CommitEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = true
}

func (s *EventService) RemoveClient(client chan models.CommitEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
		close(client)
	}
}
