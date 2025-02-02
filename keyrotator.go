package keyrotator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var resetInstance int = 1;

// APIKeyConfig represents the configuration for API keys
type APIKeyConfig struct {
	// Nested map to support multiple categories and keys
	// First key is the category (e.g., "llm", "tts")
	// Second key is the key identifier
	// Value depends on the category structure
	Keys map[string]map[string]interface{} `json:"keys"`

	// Tracks how many times each key has been used
	UsageCount map[string]map[string]int `json:"usageCount"`

	// Tracks exhausted keys per category
	ExhaustedKeys map[string]map[string]bool `json:"exhaustedKeys"`

	// Mutex to ensure thread-safety
	mu sync.RWMutex

	// ConfigPath stores the path to the configuration file
	ConfigPath string
}

// NewAPIKeyConfig creates a new APIKeyConfig from a JSON file
func NewAPIKeyConfig(configPath string) (*APIKeyConfig, error) {
	// Enter absolute path to resource/configuration file
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Read config file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal JSON
	var config APIKeyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	// Store the config path
	config.ConfigPath = absPath

	// Initialize usage count and exhausted keys if not present
	if config.UsageCount == nil {
		config.UsageCount = make(map[string]map[string]int)
	}
	if config.ExhaustedKeys == nil {
		config.ExhaustedKeys = make(map[string]map[string]bool)
	}

	// Populate usage count and exhausted keys for each category and key
	for category, keys := range config.Keys {
		if config.UsageCount[category] == nil {
			config.UsageCount[category] = make(map[string]int)
		}
		if config.ExhaustedKeys[category] == nil {
			config.ExhaustedKeys[category] = make(map[string]bool)
		}

		for keyID := range keys {
			config.UsageCount[category][keyID] = 0
			config.ExhaustedKeys[category][keyID] = false
		}
	}

	return &config, nil
}

// GetRandomAPIKey retrieves a random available API key from a specified category
func (c *APIKeyConfig) GetRandomAPIKey(category string) (string, interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if category exists
	categoryKeys, ok := c.Keys[category]
	if !ok {
		return "", nil, fmt.Errorf("category %s not found", category)
	}

	// Find available keys
	availableKeys := []string{}
	for keyID, exhausted := range c.ExhaustedKeys[category] {
		if !exhausted {
			availableKeys = append(availableKeys, keyID)
		}
	}

	// No available keys
	if len(availableKeys) == 0 {
		return "", nil, fmt.Errorf("no available keys in category %s", category)
	}

	// Randomly select a key
	rand.Seed(time.Now().UnixNano())
	selectedKeyID := availableKeys[rand.Intn(len(availableKeys))]

	// Increment usage count
	c.UsageCount[category][selectedKeyID]++

	return selectedKeyID, categoryKeys[selectedKeyID], nil
}

// GetAPIKey retrieves a specific API key from a category
func (c *APIKeyConfig) GetAPIKey(category, keyID string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if category exists
	categoryKeys, ok := c.Keys[category]
	if !ok {
		return nil, fmt.Errorf("category %s not found", category)
	}

	// Check if key exists
	key, ok := categoryKeys[keyID]
	if !ok {
		return nil, fmt.Errorf("key %s not found in category %s", keyID, category)
	}

	return key, nil
}

// MarkKeyAsExhausted marks a specific key as exhausted
func (c *APIKeyConfig) MarkKeyAsExhausted(category, keyID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Validate category and key
	if _, ok := c.Keys[category]; !ok {
		return fmt.Errorf("category %s not found", category)
	}
	if _, ok := c.Keys[category][keyID]; !ok {
		return fmt.Errorf("key %s not found in category %s", keyID, category)
	}

	// Mark key as exhausted
	c.ExhaustedKeys[category][keyID] = true
	return nil
}

// StartMidnightReset starts a goroutine to reset keys at specified timezone midnight
func (c *APIKeyConfig) StartMidnightReset(timezone *time.Location) {
	go func() {
		for {
			now := time.Now().In(timezone)
			nextMidnight := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
			sleepDuration := time.Until(nextMidnight)
			// Log the scheduled reset time
			log.Printf("Next key reset scheduled for %v (%s)", nextMidnight, timezone.String())
			time.Sleep(sleepDuration)
			c.ResetExhaustedKeys()
		}
	}()
}

// ResetExhaustedKeys resets all exhausted keys across all categories
func (c *APIKeyConfig) ResetExhaustedKeys() {
	c.mu.Lock()
	defer c.mu.Unlock()
	resetCount := 0
	for category := range c.ExhaustedKeys {
		for keyID, wasExhausted := range c.ExhaustedKeys[category] {
			if wasExhausted {
				resetCount++
			}
			c.ExhaustedKeys[category][keyID] = false
		}
	}
	log.Printf("%d reset completed: %d exhausted keys were reset across all categories", resetInstance, resetCount)
	resetInstance++
}

// GetUsageStatistics returns usage statistics for all keys
func (c *APIKeyConfig) GetUsageStatistics() map[string]map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a deep copy to prevent external modifications
	stats := make(map[string]map[string]int)
	for category, keys := range c.UsageCount {
		stats[category] = make(map[string]int)
		for keyID, count := range keys {
			stats[category][keyID] = count
		}
	}
	return stats
}

// Reload reloads the configuration from the original file
func (c *APIKeyConfig) Reload() error {
	newConfig, err := NewAPIKeyConfig(c.ConfigPath)
	if err != nil {
		return err
	}

	// Replace the current config with the new one
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Keys = newConfig.Keys
	c.UsageCount = newConfig.UsageCount
	c.ExhaustedKeys = newConfig.ExhaustedKeys

	return nil
}

// SaveConfig saves the current configuration back to a JSON file
func (c *APIKeyConfig) SaveConfig(configPath ...string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Use provided path or original config path
	path := c.ConfigPath
	if len(configPath) > 0 {
		path = configPath[0]
	}

	data, err := json.MarshalIndent(c.Keys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}