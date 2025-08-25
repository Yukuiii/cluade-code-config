package services

import (
	"encoding/json"
	"fmt"
	"os"

	"claude-code-config/internal/models"
	"claude-code-config/internal/utils"
)

// ConfigService handles Claude Code configuration operations
type ConfigService struct{}

// NewConfigService creates a new ConfigService instance
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfigPath returns Claude Code configuration file path
func (cs *ConfigService) GetConfigPath() string {
	configPath, _ := utils.GetConfigFilePath()
	return configPath
}

// Load reads and parses Claude Code configuration
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

// Save saves Claude Code configuration
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

	// Create config directory if it doesn't exist
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get config directory: %v", err),
		}
	}
	
	if err := utils.EnsureDirectory(configDir); err != nil {
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

	// Write to file with platform-specific permissions
	permissions := utils.GetFilePermissions(false)
	if err := os.WriteFile(configPath, jsonData, permissions); err != nil {
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