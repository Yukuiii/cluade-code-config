package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"claude-code-config/internal/models"
	"claude-code-config/internal/utils"
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

// GetProfilesDir returns directory for storing configuration profiles
func (ps *ProfileService) GetProfilesDir() string {
	profilesDir, _ := utils.GetProfilesDir()
	return profilesDir
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
	if err := utils.EnsureDirectory(profilesDir); err != nil {
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
	profilePath := utils.JoinPaths(profilesDir, name+".json")
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

	// Write with platform-specific permissions
	permissions := utils.GetFilePermissions(false)
	if err := os.WriteFile(profilePath, jsonData, permissions); err != nil {
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
			profilePath := utils.JoinPaths(profilesDir, file.Name())
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

// Apply applies a saved profile to current Claude Code configuration
func (ps *ProfileService) Apply(profileName string) models.ConfigResponse {
	profilesDir := ps.GetProfilesDir()
	profilePath := utils.JoinPaths(profilesDir, profileName+".json")
	
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
	profilePath := utils.JoinPaths(profilesDir, profileName+".json")
	
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