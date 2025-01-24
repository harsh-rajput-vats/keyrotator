package keyrotator

import (
    "os"
    "testing"
)

func TestNewAPIKeyConfig(t *testing.T) {
    // Create a temporary config file
    configContent := `{
        "tts": {
            "1": {
                "accessKey": "your_access_key",
                "secretKey": "your_secret__key"
            }
        },
        "llm": {
            "1": "ind_LYUiDAiCptU7rt6HrdPalUzXTZtjJIITK"
        }
    }`
    configPath := "test_config.json"
    err := os.WriteFile(configPath, []byte(configContent), 0644)
    if err != nil {
        t.Fatalf("Failed to create test config file: %v", err)
    }
    defer os.Remove(configPath)

    // Test NewAPIKeyConfig
    config, err := NewAPIKeyConfig(configPath)
    if err != nil {
        t.Fatalf("Failed to create APIKeyConfig: %v", err)
    }

    if len(config.Keys) != 2 {
        t.Errorf("Expected 2 categories, got %d", len(config.Keys))
    }
}

func TestGetRandomAPIKey(t *testing.T) {
    configContent := `{
        "llm": {
            "1": "key1",
            "2": "key2"
        }
    }`
    configPath := "test_config.json"
    err := os.WriteFile(configPath, []byte(configContent), 0644)
    if err != nil {
        t.Fatalf("Failed to create test config file: %v", err)
    }
    defer os.Remove(configPath)

    config, err := NewAPIKeyConfig(configPath)
    if err != nil {
        t.Fatalf("Failed to create APIKeyConfig: %v", err)
    }

    keyID, keyValue, err := config.GetRandomAPIKey("llm")
    if err != nil {
        t.Fatalf("Failed to get random API key: %v", err)
    }

    if keyID == "" || keyValue == nil {
        t.Errorf("Expected valid keyID and keyValue, got keyID: %s, keyValue: %v", keyID, keyValue)
    }
}

func TestMarkKeyAsExhausted(t *testing.T) {
    configContent := `{
        "llm": {
            "1": "key1"
        }
    }`
    configPath := "test_config.json"
    err := os.WriteFile(configPath, []byte(configContent), 0644)
    if err != nil {
        t.Fatalf("Failed to create test config file: %v", err)
    }
    defer os.Remove(configPath)

    config, err := NewAPIKeyConfig(configPath)
    if err != nil {
        t.Fatalf("Failed to create APIKeyConfig: %v", err)
    }

    err = config.MarkKeyAsExhausted("llm", "1")
    if err != nil {
        t.Fatalf("Failed to mark key as exhausted: %v", err)
    }

    if !config.ExhaustedKeys["llm"]["1"] {
        t.Errorf("Expected key to be marked as exhausted")
    }
}