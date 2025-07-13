package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"toby_launcher/core/version"
	"toby_launcher/utils/file_utils"
)

// PathConfig manages paths to configuration files.
type PathConfig struct {
	BaseDir    string
	FilesDir   string
	isPortable bool
}

// NewPathConfig creates a new PathConfig instance and validates the configuration directory.
func NewPathConfig() (*PathConfig, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	dataDir := filepath.Join(exeDir, "data")
	if file_utils.Exists(dataDir) {
		// Portable mode: use data as BaseDir
		filesDir := filepath.Join(dataDir, "files")
		if err := os.MkdirAll(filesDir, 0755); err != nil {
			return nil, err
		}
		return &PathConfig{
			BaseDir:    dataDir,
			FilesDir:   filesDir,
			isPortable: true,
		}, nil
	}

	// Installed mode: use system configuration directory
	configDir, err := OsConfigDir(runtime.GOOS)
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(configDir, version.AppName)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration directory does not exist: %s", baseDir)
	}
	filesDir := filepath.Join(baseDir, "files")
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		return nil, err
	}
	return &PathConfig{
		BaseDir:    baseDir,
		FilesDir:   filesDir,
		isPortable: false,
	}, nil
}

func (pc *PathConfig) ConfigFilePath() string {
	return filepath.Join(pc.BaseDir, "config.json")
}

func (pc *PathConfig) LogFilePath() string {
	logFile := version.AppName + "." + "log"
	return filepath.Join(pc.BaseDir, logFile)
}

func (pc *PathConfig) GamesPath() string {
	return filepath.Join(pc.BaseDir, "games.json")
}

func (pc *PathConfig) TextRulesPath() string {
	return filepath.Join(pc.BaseDir, "text_rules.json")
}

func (pc *PathConfig) GameFilePath(file string) string {
	return filepath.Join(pc.BaseDir, "files", file)
}

// GzdoomPath returns the path to the gzdoom executable.
func (pc *PathConfig) GzdoomPath() (string, error) {
	exeName := "gzdoom"
	if runtime.GOOS == "windows" {
		exeName = "gzdoom.exe"
	}

	// 1. Check $PATH
	if path, err := exec.LookPath(exeName); err == nil {
		return path, nil
	}

	// 2. Check system installation (e.g., /usr/local/share/TobyLauncher/gzdoom/gzdoom) or portable mode
	var systemPath string
	if runtime.GOOS == "darwin" {
		systemPath = filepath.Join(pc.BaseDir, "gzdoom", "GZDoom.app", "Contents", "MacOS", exeName)
	} else {
		systemPath = filepath.Join(pc.BaseDir, "gzdoom", exeName)
	}
	if file_utils.Exists(systemPath) {
		return systemPath, nil
	}

	// 3. Check-specific path (e.g., /Applications/Gzdoom.app/Contents/MacOS/gzdoom)
	if runtime.GOOS == "darwin" {
		macPath := filepath.Join("/Applications", "Gzdoom.app", "Contents", "MacOS", exeName)
		if file_utils.Exists(macPath) {
			return macPath, nil
		}
	}

	return "", fmt.Errorf("gzdoom executable not found in portable, system, or PATH")
}

func OsConfigDir(platform string) (string, error) {
	switch platform {
	case "linux":
		return "/usr/local/share", nil
	case "windows":
		return "C:\\ProgramData", nil
	case "darwin":
		return "/Library/Application Support", nil
	default:
		return "", fmt.Errorf("platform %s is not supported", platform)
	}
}
