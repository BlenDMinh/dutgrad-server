package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type KVStorage interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Entry struct {
	Value      interface{} `json:"value"`
	Expiration int64       `json:"expiration"`
}

type InMemoryStorage struct {
	store map[string]Entry
	mu    sync.RWMutex
}

func NewInMemoryStorage() KVStorage {
	return &InMemoryStorage{
		store: make(map[string]Entry),
		mu:    sync.RWMutex{},
	}
}

func (ms *InMemoryStorage) Set(key string, value interface{}, expiration time.Duration) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var expirationTime int64 = 0
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration).Unix()
	}

	ms.store[key] = Entry{
		Value:      value,
		Expiration: expirationTime,
	}

	return nil
}

func (ms *InMemoryStorage) Get(key string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	entry, exists := ms.store[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	// Check if the entry has expired
	if entry.Expiration > 0 && entry.Expiration < time.Now().Unix() {
		// Delete the expired entry (will require a write lock)
		ms.mu.RUnlock()
		_ = ms.Delete(key)
		ms.mu.RLock()
		return "", fmt.Errorf("key expired: %s", key)
	}

	// Convert value to string
	valueStr, ok := entry.Value.(string)
	if !ok {
		// Convert to JSON string if not already a string
		valueBytes, err := json.Marshal(entry.Value)
		if err != nil {
			return "", fmt.Errorf("failed to convert value to string: %w", err)
		}
		valueStr = string(valueBytes)
	}

	return valueStr, nil
}

func (ms *InMemoryStorage) Delete(key string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.store, key)
	return nil
}

// Optional: Add cleanup routine to periodically remove expired entries
func (ms *InMemoryStorage) StartCleanupRoutine(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			ms.cleanupExpiredEntries()
		}
	}()
}

func (ms *InMemoryStorage) cleanupExpiredEntries() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now().Unix()
	for key, entry := range ms.store {
		if entry.Expiration > 0 && entry.Expiration < now {
			delete(ms.store, key)
		}
	}
}
