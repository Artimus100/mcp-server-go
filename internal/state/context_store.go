package state

import (
	"sync"
)

// ClientContext represents the context data for a client connection
type ClientContext struct {
	Values map[string]string
}

// ContextStore provides a thread-safe store for client context information
type ContextStore struct {
	contexts map[string]*ClientContext
	mu       sync.RWMutex
}

// NewContextStore creates a new empty context store
func NewContextStore() *ContextStore {
	return &ContextStore{
		contexts: make(map[string]*ClientContext),
	}
}

// Get retrieves a specific context value for a client
func (s *ContextStore) Get(clientID, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.contexts[clientID]
	if !exists {
		return "", false
	}

	val, exists := client.Values[key]
	return val, exists
}

// GetAll returns all context values for a client
func (s *ContextStore) GetAll(clientID string) (map[string]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.contexts[clientID]
	if !exists {
		return nil, false
	}

	// Copy the map to avoid external modification
	result := make(map[string]string)
	for k, v := range client.Values {
		result[k] = v
	}

	return result, true
}

// Set updates a context value for a client
func (s *ContextStore) Set(clientID, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.contexts[clientID]
	if !exists {
		client = &ClientContext{
			Values: make(map[string]string),
		}
		s.contexts[clientID] = client
	}

	client.Values[key] = value
}

// SetMultiple updates multiple context values for a client
func (s *ContextStore) SetMultiple(clientID string, values map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.contexts[clientID]
	if !exists {
		client = &ClientContext{
			Values: make(map[string]string),
		}
		s.contexts[clientID] = client
	}

	for k, v := range values {
		client.Values[k] = v
	}
}

// Remove deletes a context value for a client
func (s *ContextStore) Remove(clientID, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.contexts[clientID]
	if !exists {
		return
	}

	delete(client.Values, key)
}

// Clear removes all context values for a client
func (s *ContextStore) Clear(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.contexts, clientID)
}

// ListClients returns a list of all client IDs in the store
func (s *ContextStore) ListClients() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]string, 0, len(s.contexts))
	for id := range s.contexts {
		clients = append(clients, id)
	}

	return clients
}

// QueryClients finds clients that match a given key-value condition
// TODO: Implement more sophisticated query capabilities
func (s *ContextStore) QueryClients(key, value string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []string

	for clientID, ctx := range s.contexts {
		if v, exists := ctx.Values[key]; exists && v == value {
			matches = append(matches, clientID)
		}
	}

	return matches
}

// TODO: Add more advanced context operations:
// - Context expiration/TTL
// - Context snapshots/history
// - Subscription to context changes
// - Context serialization/persistence
