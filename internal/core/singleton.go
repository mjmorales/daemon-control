package core

import (
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton core config manager
func GetManager() *Manager {
	once.Do(func() {
		instance = NewManager()
		if err := instance.Init(); err != nil {
			log.Warn().Err(err).Msg("Failed to initialize core config, using defaults")
			instance.config = DefaultConfig()
		}
	})
	return instance
}

// MustGetConfig returns the core config or panics
func MustGetConfig() *CoreConfig {
	manager := GetManager()
	if manager.config == nil {
		panic("core config not initialized")
	}
	return manager.config
}
