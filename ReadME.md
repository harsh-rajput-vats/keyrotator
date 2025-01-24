# keyrotator

## Overview
`keyrotator` is a flexible Go package for managing and rotating API keys across multiple categories with thread-safe operations.

## Features
- Multiple key categories support
- Random key selection
- Usage tracking
- Automatic key reset at configurable time
- Thread-safe operations


## Example Limited Resources
```json
{
  "tts": {
    "1": {
      "accessKey": "your_access_key",
      "secretKey": "your_secret__key"
    },
    "2": {
      "accessKey": "another-access-key",
      "secretKey": "another-secret-key"
    }
  },
  "llm": {
    "1": "ind_LYUiDAiCptU7rt6HrdPalUzXTZtjJIITK",
    "2": "another-valid-key-retrieved-from-diff-userAcc",
    "3": "yet-another-valid-llm-key"
  },
}
```

## Installation
```bash
go get github.com/harsh-rajput-vats/keyrotator
```

## Quick Start
```go
config, _ := keyrotator.NewAPIKeyConfig("config.json")

// Start midnight reset for UTC timezone
config.StartMidnightReset(time.UTC)

// Get a random API key
keyID, apiKey, _ := config.GetRandomAPIKey("llm")
```

## Example Usage

```go
package main

import (
	"fmt"
	"log"
	"time"
    "keyrotator"
)

func main() {
	// Load configuration from a specific path
	config, _ := keyrotator.NewAPIKeyConfig("resources/limited.json")

	// Start midnight reset for UTC timezone
	config.StartMidnightReset(time.UTC)

	// Example of getting a random API key from LLM category
	keyID, apiKey, _ := config.GetRandomAPIKey("llm")

	// Example of getting a specific key from TTS category
	ttsKey, _ := config.GetAPIKey("tts", "1")

	// Mark a key as exhausted like on getting HTTP response of 429
	config.MarkKeyAsExhausted("llm", keyID)

  
}
```
