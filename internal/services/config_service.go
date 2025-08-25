package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"claude-code-config/internal/models"
)

// ConfigService handles Claude Code configuration operations
type ConfigService struct{}

// NewConfigService creates a new ConfigService instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfigPath returns the Claude Code configuration file path
func (cs *ConfigService) GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".claude", "settings.json")
}

// Load reads and parses the Claude Code configuration
func (cs *ConfigService) Load() models.ConfigResponse {
	configPath := cs.GetConfigPath()
	
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return models.ConfigResponse{
			Success: false,
			Message: "Claude Code configuration file not found",
		}
	}

	// Read file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to read configuration file: %v", err),
		}
	}

	// Parse JSON
	var config models.ClaudeConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to parse configuration file: %v", err),
		}
	}

	return models.ConfigResponse{
		Success: true,
		Message: "Configuration loaded successfully",
		Data:    config,
	}
}

// Save saves the Claude Code configuration
func (cs *ConfigService) Save(authToken, baseURL string) models.ConfigResponse {
	configPath := cs.GetConfigPath()
	
	// Load existing config to preserve permissions
	existing := cs.Load()
	var config models.ClaudeConfig
	
	if existing.Success {
		// Use existing config as base
		if existingData, ok := existing.Data.(models.ClaudeConfig); ok {
			config = existingData
		}
	}
	
	// Update env variables
	config.Env.AnthropicAuthToken = authToken
	config.Env.AnthropicBaseURL = baseURL
	
	// Ensure permissions exist (initialize if empty)
	if config.Permissions.Allow == nil {
		config.Permissions.Allow = []string{}
	}
	if config.Permissions.Deny == nil {
		config.Permissions.Deny = []string{}
	}

	// Create .claude directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create configuration directory: %v", err),
		}
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to serialize configuration: %v", err),
		}
	}

	// Write to file
	if err := os.WriteFile(configPath, jsonData, 0644); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to write configuration file: %v", err),
		}
	}

	return models.ConfigResponse{
		Success: true,
		Message: "Configuration saved successfully",
		Data:    config,
	}
}