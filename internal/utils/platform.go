package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// Platform represents the operating system type
type Platform int

const (
	PlatformWindows Platform = iota
	PlatformDarwin
	PlatformLinux
	PlatformUnknown
)

// GetPlatform returns the current platform
func GetPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return PlatformWindows
	case "darwin":
		return PlatformDarwin
	case "linux":
		return PlatformLinux
	default:
		return PlatformUnknown
	}
}

// GetConfigDir returns the platform-specific configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch GetPlatform() {
	case PlatformWindows:
		// Windows: %APPDATA%/claude
		appData := os.Getenv("APPDATA")
		if appData == "" {
			// Fallback to user home + AppData\Roaming
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "claude"), nil
	
	case PlatformDarwin, PlatformLinux:
		// Unix-like: ~/.claude
		return filepath.Join(homeDir, ".claude"), nil
	
	default:
		// Unknown platform: fallback to user home + .claude
		return filepath.Join(homeDir, ".claude"), nil
	}
}

// GetConfigFilePath returns the platform-specific configuration file path
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "settings.json"), nil
}

// GetProfilesDir returns the platform-specific profiles directory
func GetProfilesDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "profiles"), nil
}

// GetFilePermissions returns platform-specific file permissions
func GetFilePermissions(isDir bool) os.FileMode {
	switch GetPlatform() {
	case PlatformWindows:
		// Windows: Go's file mode is mostly ignored, but we set reasonable defaults
		if isDir {
			return 0755 // rwxr-xr-x
		}
		return 0644 // rw-r--r--
	
	default:
		// Unix-like systems
		if isDir {
			return 0755 // rwxr-xr-x
		}
		return 0644 // rw-r--r--
	}
}

// EnsureDirectory creates a directory with platform-specific permissions
func EnsureDirectory(dirPath string) error {
	permissions := GetFilePermissions(true)
	return os.MkdirAll(dirPath, permissions)
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return GetPlatform() == PlatformWindows
}

// IsDarwin returns true if running on macOS
func IsDarwin() bool {
	return GetPlatform() == PlatformDarwin
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return GetPlatform() == PlatformLinux
}

// GetPlatformName returns human-readable platform name
func GetPlatformName() string {
	switch GetPlatform() {
	case PlatformWindows:
		return "Windows"
	case PlatformDarwin:
		return "macOS"
	case PlatformLinux:
		return "Linux"
	default:
		return "Unknown"
	}
}