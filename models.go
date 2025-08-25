package main

import "time"

// ClaudeConfig represents the Claude Code configuration structure
type ClaudeConfig struct {
	Env struct {
		AnthropicAuthToken string `json:"ANTHROPIC_AUTH_TOKEN"`
		AnthropicBaseURL   string `json:"ANTHROPIC_BASE_URL"`
	} `json:"env"`
	Permissions struct {
		Allow []string `json:"allow"`
		Deny  []string `json:"deny"`
	} `json:"permissions"`
}

// ConfigProfile represents a saved configuration profile
type ConfigProfile struct {
	Name        string    `json:"name"`
	AuthToken   string    `json:"authToken"`
	BaseURL     string    `json:"baseURL"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description,omitempty"`
}

// ConfigResponse represents the response structure for API calls
type ConfigResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}