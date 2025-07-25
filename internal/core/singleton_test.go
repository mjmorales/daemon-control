package core

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetManager(t *testing.T) {
	// Reset singleton for testing
	once = sync.Once{}
	instance = nil
	
	// First call should create instance
	manager1 := GetManager()
	assert.NotNil(t, manager1)
	assert.NotNil(t, manager1.config)
	
	// Second call should return same instance
	manager2 := GetManager()
	assert.Same(t, manager1, manager2)
	
	// Concurrent calls should all get same instance
	var wg sync.WaitGroup
	managers := make([]*Manager, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			managers[idx] = GetManager()
		}(i)
	}
	
	wg.Wait()
	
	// All should be the same instance
	for i := 0; i < 10; i++ {
		assert.Same(t, manager1, managers[i])
	}
}

func TestMustGetConfig(t *testing.T) {
	// Reset singleton for testing
	once = sync.Once{}
	instance = nil
	
	// Test normal operation - should not panic
	assert.NotPanics(t, func() {
		config := MustGetConfig()
		assert.NotNil(t, config)
	})
}