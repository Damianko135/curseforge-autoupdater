package api

import (
	"fmt"
	"strings"
	"time"
)

// ModpackInfo represents modpack-specific information
type ModpackInfo struct {
	*ModInfo
	LatestVersion   string
	CurrentVersion  string
	HasUpdate       bool
	UpdateAvailable *ModFile
	Changelog       string
}

// ModLoaderType constants
const (
	ModLoaderTypeAny    int = 0
	ModLoaderTypeForge  int = 1
	ModLoaderTypeFabric int = 4
	ModLoaderTypeQuilt  int = 5
)

// GameID constants
const (
	GameIDMinecraft int = 432
)

// GetModpackInfo retrieves comprehensive information about a modpack
func (c *Client) GetModpackInfo(modpackID int, gameVersion string, currentVersion string, releaseChannel string) (*ModpackInfo, error) {
	// Get basic mod info
	modInfo, err := c.GetMod(modpackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack info: %w", err)
	}

	// Check if this is actually a modpack
	if modInfo.ClassID != 4471 { // 4471 is the class ID for modpacks
		return nil, fmt.Errorf("mod %d is not a modpack (class ID: %d)", modpackID, modInfo.ClassID)
	}

	// Get latest file based on release channel
	releaseType := getReleaseTypeFromChannel(releaseChannel)
	latestFile, err := c.GetLatestModFile(modpackID, gameVersion, releaseType)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest modpack file: %w", err)
	}

	// Create modpack info
	modpackInfo := &ModpackInfo{
		ModInfo:         modInfo,
		LatestVersion:   latestFile.DisplayName,
		CurrentVersion:  currentVersion,
		UpdateAvailable: latestFile,
	}

	// Check if update is available
	if currentVersion != "" {
		modpackInfo.HasUpdate = !isVersionEqual(currentVersion, latestFile.DisplayName)
	} else {
		modpackInfo.HasUpdate = true // No current version means update is available
	}

	return modpackInfo, nil
}

// GetModpackVersions retrieves all available versions for a modpack
func (c *Client) GetModpackVersions(modpackID int, gameVersion string) ([]ModFile, error) {
	// Get all files for the modpack
	files, err := c.GetModFiles(modpackID, gameVersion, 0, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack files: %w", err)
	}

	// Filter for server files if available
	var serverFiles []ModFile
	var allFiles []ModFile

	for _, file := range files {
		if file.IsServerPack {
			serverFiles = append(serverFiles, file)
		}
		allFiles = append(allFiles, file)
	}

	// Prefer server files if available, otherwise use all files
	if len(serverFiles) > 0 {
		return serverFiles, nil
	}

	return allFiles, nil
}

// GetModpackChangelog retrieves the changelog for a modpack version
func (c *Client) GetModpackChangelog(modpackID int, fileID int) (string, error) {
	// Note: CurseForge API doesn't directly provide changelog in the file info
	// This would need to be implemented by parsing the file description or
	// using additional API endpoints if available

	file, err := c.GetModFile(modpackID, fileID)
	if err != nil {
		return "", fmt.Errorf("failed to get modpack file: %w", err)
	}

	// For now, return basic information about the file
	changelog := fmt.Sprintf("Version: %s\n", file.DisplayName)
	changelog += fmt.Sprintf("Release Date: %s\n", file.FileDate.Format("2006-01-02 15:04:05"))
	changelog += fmt.Sprintf("File Size: %d bytes\n", file.FileLength)
	changelog += fmt.Sprintf("Downloads: %d\n", file.DownloadCount)

	if len(file.GameVersions) > 0 {
		changelog += fmt.Sprintf("Game Versions: %s\n", strings.Join(file.GameVersions, ", "))
	}

	return changelog, nil
}

// CompareModpackVersions compares two modpack versions
func (c *Client) CompareModpackVersions(modpackID int, currentFileID int, latestFileID int) (*VersionComparison, error) {
	if currentFileID == latestFileID {
		return &VersionComparison{
			IsNewer:      false,
			CurrentFile:  nil,
			LatestFile:   nil,
			ChangelogURL: "",
		}, nil
	}

	// Get current file info
	var currentFile *ModFile
	var err error
	if currentFileID > 0 {
		currentFile, err = c.GetModFile(modpackID, currentFileID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current modpack file: %w", err)
		}
	}

	// Get latest file info
	latestFile, err := c.GetModFile(modpackID, latestFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest modpack file: %w", err)
	}

	// Compare file dates
	isNewer := false
	if currentFile == nil || latestFile.FileDate.After(currentFile.FileDate) {
		isNewer = true
	}

	return &VersionComparison{
		IsNewer:      isNewer,
		CurrentFile:  currentFile,
		LatestFile:   latestFile,
		ChangelogURL: "", // Would need to be implemented based on available API
	}, nil
}

// VersionComparison represents a comparison between two modpack versions
type VersionComparison struct {
	IsNewer      bool
	CurrentFile  *ModFile
	LatestFile   *ModFile
	ChangelogURL string
}

// IsModpackCompatible checks if a modpack is compatible with the specified game version
func (c *Client) IsModpackCompatible(modpackID int, gameVersion string) (bool, error) {
	files, err := c.GetModFiles(modpackID, gameVersion, 0, 10, 0)
	if err != nil {
		return false, fmt.Errorf("failed to get modpack files: %w", err)
	}

	// Check if any files are compatible with the specified game version
	for _, file := range files {
		for _, version := range file.GameVersions {
			if version == gameVersion {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetModpackDownloadURL retrieves the download URL for a modpack file
func (c *Client) GetModpackDownloadURL(modpackID int, fileID int) (string, error) {
	return c.GetModFileDownloadURL(modpackID, fileID)
}

// GetModpackServerFile retrieves the server file for a modpack if available
func (c *Client) GetModpackServerFile(modpackID int, gameVersion string, releaseChannel string) (*ModFile, error) {
	releaseType := getReleaseTypeFromChannel(releaseChannel)

	files, err := c.GetModFiles(modpackID, gameVersion, 0, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack files: %w", err)
	}

	// Look for server files first
	for _, file := range files {
		if file.IsServerPack && (releaseType == 0 || file.ReleaseType == releaseType) {
			return &file, nil
		}
	}

	// If no server file found, look for regular files
	for _, file := range files {
		if releaseType == 0 || file.ReleaseType == releaseType {
			return &file, nil
		}
	}

	return nil, fmt.Errorf("no suitable file found for modpack %d", modpackID)
}

// getReleaseTypeFromChannel converts a release channel string to release type int
func getReleaseTypeFromChannel(channel string) int {
	switch strings.ToLower(channel) {
	case "stable", "release":
		return ReleaseTypeRelease
	case "beta":
		return ReleaseTypeBeta
	case "alpha":
		return ReleaseTypeAlpha
	default:
		return 0 // Any release type
	}
}

// isVersionEqual compares two version strings for equality
func isVersionEqual(version1, version2 string) bool {
	// Simple string comparison for now
	// This could be enhanced with semantic version comparison
	return strings.TrimSpace(version1) == strings.TrimSpace(version2)
}

// GetModpackDependencies retrieves dependencies for a modpack
func (c *Client) GetModpackDependencies(modpackID int, fileID int) ([]ModDependency, error) {
	file, err := c.GetModFile(modpackID, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack file: %w", err)
	}

	return file.Dependencies, nil
}

// ModpackUpdateInfo represents information about a modpack update
type ModpackUpdateInfo struct {
	HasUpdate      bool
	CurrentVersion string
	LatestVersion  string
	CurrentFileID  int
	LatestFileID   int
	UpdateSize     int64
	ReleaseDate    time.Time
	GameVersions   []string
	IsServerPack   bool
	DownloadURL    string
	Changelog      string
	Dependencies   []ModDependency
	IsCompatible   bool
	RequiredMods   []ModDependency
	OptionalMods   []ModDependency
}

// GetModpackUpdateInfo retrieves comprehensive update information
func (c *Client) GetModpackUpdateInfo(modpackID int, currentVersion string, currentFileID int, gameVersion string, releaseChannel string) (*ModpackUpdateInfo, error) {
	// Get modpack info
	modpackInfo, err := c.GetModpackInfo(modpackID, gameVersion, currentVersion, releaseChannel)
	if err != nil {
		return nil, fmt.Errorf("failed to get modpack info: %w", err)
	}

	// Get download URL
	downloadURL, err := c.GetModpackDownloadURL(modpackID, modpackInfo.UpdateAvailable.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get download URL: %w", err)
	}

	// Get changelog
	changelog, err := c.GetModpackChangelog(modpackID, modpackInfo.UpdateAvailable.ID)
	if err != nil {
		// Don't fail if changelog is not available
		changelog = "Changelog not available"
	}

	// Check compatibility
	isCompatible, err := c.IsModpackCompatible(modpackID, gameVersion)
	if err != nil {
		// Don't fail if compatibility check fails
		isCompatible = false
	}

	// Separate dependencies into required and optional
	var requiredMods []ModDependency
	var optionalMods []ModDependency

	for _, dep := range modpackInfo.UpdateAvailable.Dependencies {
		if dep.RelationType == RelationTypeRequiredDependency {
			requiredMods = append(requiredMods, dep)
		} else if dep.RelationType == RelationTypeOptionalDependency {
			optionalMods = append(optionalMods, dep)
		}
	}

	return &ModpackUpdateInfo{
		HasUpdate:      modpackInfo.HasUpdate,
		CurrentVersion: modpackInfo.CurrentVersion,
		LatestVersion:  modpackInfo.LatestVersion,
		CurrentFileID:  currentFileID,
		LatestFileID:   modpackInfo.UpdateAvailable.ID,
		UpdateSize:     modpackInfo.UpdateAvailable.FileLength,
		ReleaseDate:    modpackInfo.UpdateAvailable.FileDate,
		GameVersions:   modpackInfo.UpdateAvailable.GameVersions,
		IsServerPack:   modpackInfo.UpdateAvailable.IsServerPack,
		DownloadURL:    downloadURL,
		Changelog:      changelog,
		Dependencies:   modpackInfo.UpdateAvailable.Dependencies,
		IsCompatible:   isCompatible,
		RequiredMods:   requiredMods,
		OptionalMods:   optionalMods,
	}, nil
}
