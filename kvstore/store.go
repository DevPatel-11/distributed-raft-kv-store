package kvstore

import (
	"fmt"
	"sync"
)

// Store represents the key-value store state machine
type Store struct {
	mu   sync.RWMutex
	data map[string][]byte
	meta map[string]*Meta
}

// Meta stores metadata about a key
type Meta struct {
	Version   int64
	Timestamp int64
}

// New creates a new key-value store
func New() *Store {
	return &Store{
		data: make(map[string][]byte),
		meta: make(map[string]*Meta),
	}
}

// Set sets a key-value pair
func (s *Store) Set(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	s.data[key] = value
	s.meta[key] = &Meta{
		Version: int64(len(s.data)),
	}

	return nil
}

// Get retrieves a value by key
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Delete removes a key
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key not found")
	}

	delete(s.data, key)
	delete(s.meta, key)

	return nil
}

// Keys returns all keys in the store
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}

// Clear clears all data from the store
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string][]byte)
	s.meta = make(map[string]*Meta)
}

// Size returns the number of keys in the store
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}
