package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx            context.Context
	configService  *ConfigService
	profileService *ProfileService
}

// NewApp creates a new App application struct
func NewApp() *App {
	configService := NewConfigService()
	profileService := NewProfileService(configService)
	
	return &App{
		configService:  configService,
		profileService: profileService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetClaudeConfigPath returns the Claude Code configuration file path
func (a *App) GetClaudeConfigPath() string {
	return a.configService.GetConfigPath()
}

// GetProfilesDir returns the directory for storing configuration profiles
func (a *App) GetProfilesDir() string {
	return a.profileService.GetProfilesDir()
}

// LoadClaudeConfig reads and parses the Claude Code configuration
func (a *App) LoadClaudeConfig() ConfigResponse {
	return a.configService.Load()
}

// SaveClaudeConfig saves the Claude Code configuration
func (a *App) SaveClaudeConfig(authToken, baseURL string) ConfigResponse {
	return a.configService.Save(authToken, baseURL)
}

// SaveConfigProfile saves a configuration as a reusable profile
func (a *App) SaveConfigProfile(name, authToken, baseURL, description string) ConfigResponse {
	return a.profileService.Save(name, authToken, baseURL, description)
}

// LoadConfigProfiles loads all saved configuration profiles
func (a *App) LoadConfigProfiles() ConfigResponse {
	return a.profileService.LoadAll()
}

// ApplyConfigProfile applies a saved profile to the current Claude Code configuration
func (a *App) ApplyConfigProfile(profileName string) ConfigResponse {
	return a.profileService.Apply(profileName)
}

// DeleteConfigProfile deletes a saved configuration profile
func (a *App) DeleteConfigProfile(profileName string) ConfigResponse {
	return a.profileService.Delete(profileName)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}