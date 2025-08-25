package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"claude-code-config/internal/models"
)

// ProfileService handles configuration profile operations
type ProfileService struct {
	configService *ConfigService
}

// NewProfileService creates a new ProfileService instance
func NewProfileService(configService *ConfigService) *ProfileService {
	return &ProfileService{
		configService: configService,
	}
}

// GetProfilesDir returns the directory for storing configuration profiles
func (ps *ProfileService) GetProfilesDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".claude", "profiles")
}

// Save saves a configuration as a reusable profile
func (ps *ProfileService) Save(name, authToken, baseURL, description string) models.ConfigResponse {
	// Validate inputs
	name = strings.TrimSpace(name)
	if name == "" {
		return models.ConfigResponse{
			Success: false,
			Message: "Profile name cannot be empty",
		}
	}

	// Create profiles directory
	profilesDir := ps.GetProfilesDir()
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create profiles directory: %v", err),
		}
	}

	// Create profile
	profile := models.ConfigProfile{
		Name:        name,
		AuthToken:   strings.TrimSpace(authToken),
		BaseURL:     strings.TrimSpace(baseURL),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Description: strings.TrimSpace(description),
	}

	// Check if profile already exists
	profilePath := filepath.Join(profilesDir, name+".json")
	if _, err := os.Stat(profilePath); err == nil {
		// Profile exists, update it
		if existingProfile, err := ps.loadFromFile(profilePath); err == nil {
			profile.CreatedAt = existingProfile.CreatedAt
		}
	}

	// Save profile to file
	jsonData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to serialize profile: %v", err),
		}
	}

	if err := os.WriteFile(profilePath, jsonData, 0644); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save profile: %v", err),
		}
	}

	return models.ConfigResponse{
		Success: true,
		Message: "Profile saved successfully",
		Data:    profile,
	}
}

// LoadAll loads all saved configuration profiles
func (ps *ProfileService) LoadAll() models.ConfigResponse {
	profilesDir := ps.GetProfilesDir()
	
	// Check if profiles directory exists
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		return models.ConfigResponse{
			Success: true,
			Message: "No profiles found",
			Data:    []models.ConfigProfile{},
		}
	}

	// Read all profile files
	files, err := os.ReadDir(profilesDir)
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to read profiles directory: %v", err),
		}
	}

	var profiles []models.ConfigProfile
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			profilePath := filepath.Join(profilesDir, file.Name())
			if profile, err := ps.loadFromFile(profilePath); err == nil {
				profiles = append(profiles, profile)
			}
		}
	}

	return models.ConfigResponse{
		Success: true,
		Message: fmt.Sprintf("Loaded %d profiles", len(profiles)),
		Data:    profiles,
	}
}

// Apply applies a saved profile to the current Claude Code configuration
func (ps *ProfileService) Apply(profileName string) models.ConfigResponse {
	profilesDir := ps.GetProfilesDir()
	profilePath := filepath.Join(profilesDir, profileName+".json")
	
	// Load profile
	profile, err := ps.loadFromFile(profilePath)
	if err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to load profile '%s': %v", profileName, err),
		}
	}

	// Apply profile to Claude Code configuration
	return ps.configService.Save(profile.AuthToken, profile.BaseURL)
}

// Delete deletes a saved configuration profile
func (ps *ProfileService) Delete(profileName string) models.ConfigResponse {
	profilesDir := ps.GetProfilesDir()
	profilePath := filepath.Join(profilesDir, profileName+".json")
	
	// Check if profile exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Profile '%s' not found", profileName),
		}
	}

	// Delete profile file
	if err := os.Remove(profilePath); err != nil {
		return models.ConfigResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete profile '%s': %v", profileName, err),
		}
	}

	return models.ConfigResponse{
		Success: true,
		Message: fmt.Sprintf("Profile '%s' deleted successfully", profileName),
	}
}

// loadFromFile loads a profile from a JSON file
func (ps *ProfileService) loadFromFile(filePath string) (models.ConfigProfile, error) {
	var profile models.ConfigProfile
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return profile, err
	}

	if err := json.Unmarshal(content, &profile); err != nil {
		return profile, err
	}

	return profile, nil
}