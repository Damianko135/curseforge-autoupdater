package curseforge

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/pkg/models"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	BaseURL = "https://api.curseforge.com/v1"
)

// Client represents a CurseForge API client
type Client struct {
	client *resty.Client
	logger *logrus.Logger
}

// NewClient creates a new CurseForge API client
func NewClient(apiKey string, logger *logrus.Logger) *Client {
	client := resty.New()
	client.SetBaseURL(BaseURL)
	client.SetHeader("Accept", "application/json")
	client.SetHeader("x-api-key", apiKey)
	client.SetHeader("User-Agent", "CurseForge Auto-Updater/1.0")

	return &Client{
		client: client,
		logger: logger,
	}
}

// GetModInfo retrieves basic information about a mod
func (c *Client) GetModInfo(modID string) (*models.ModInfo, error) {
	c.logger.Debugf("Fetching mod info for ID: %s", modID)

	var response models.ModInfoResponse
	resp, err := c.client.R().
		SetResult(&response).
		Get(fmt.Sprintf("/mods/%s", modID))

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	c.logger.Infof("Found mod: %s by %s", response.Data.Name,
		func() string {
			if len(response.Data.Authors) > 0 {
				return response.Data.Authors[0].Name
			}
			return "Unknown"
		}())

	return &response.Data, nil
}

// GetModFiles retrieves all files for a mod
func (c *Client) GetModFiles(modID string, gameID int) ([]models.CurseForgeFile, error) {
	c.logger.Debugf("Fetching files for mod ID: %s", modID)

	var response models.FilesResponse
	req := c.client.R().SetResult(&response)

	if gameID > 0 {
		req.SetQueryParam("gameId", strconv.Itoa(gameID))
	}

	resp, err := req.Get(fmt.Sprintf("/mods/%s/files", modID))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	c.logger.Infof("Found %d files for mod", len(response.Data))

	// Log pagination info if available
	if response.Pagination.TotalCount > 0 {
		c.logger.Debugf("Pagination: %d/%d results", response.Pagination.ResultCount, response.Pagination.TotalCount)
	}

	return response.Data, nil
}

// GetLatestFile returns the latest file from a list of files
func (c *Client) GetLatestFile(files []models.CurseForgeFile) *models.CurseForgeFile {
	if len(files) == 0 {
		return nil
	}

	// Sort by file date descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].FileDate > files[j].FileDate
	})

	latest := &files[0]
	c.logger.Infof("Latest file: %s (%s)", latest.FileName, latest.FileDate)
	return latest
}

// DownloadFile downloads a file to the specified directory
func (c *Client) DownloadFile(file *models.CurseForgeFile, downloadPath string) error {
	if file.DownloadURL == "" {
		return fmt.Errorf("no download URL available for file %s", file.FileName)
	}

	// Ensure download directory exists
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	filePath := filepath.Join(downloadPath, file.FileName)
	c.logger.Infof("Downloading %s to %s", file.FileName, filePath)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download the file
	resp, err := c.client.R().
		SetDoNotParseResponse(true).
		Get(file.DownloadURL)

	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode())
	}

	// Copy the response body to the file
	_, err = io.Copy(out, resp.RawBody())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	c.logger.Infof("Successfully downloaded %s (%d bytes)", file.FileName, file.FileLength)
	return nil
}

// LoadDownloadMetadata loads metadata about previously downloaded files
func LoadDownloadMetadata(downloadPath string) (map[string]models.DownloadMetadata, error) {
	metadataFile := filepath.Join(downloadPath, "download_metadata.json")

	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		return make(map[string]models.DownloadMetadata), nil
	}

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata map[string]models.DownloadMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	return metadata, nil
}

// SaveDownloadMetadata saves metadata about downloaded files
func SaveDownloadMetadata(downloadPath string, metadata map[string]models.DownloadMetadata) error {
	metadataFile := filepath.Join(downloadPath, "download_metadata.json")

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// IsDownloadNeeded checks if a file needs to be downloaded
func IsDownloadNeeded(file *models.CurseForgeFile, downloadPath string, metadata map[string]models.DownloadMetadata, logger *logrus.Logger) (bool, string) {
	fileID := strconv.Itoa(file.ID)
	fileName := file.FileName

	// Check if file exists locally
	localFilePath := filepath.Join(downloadPath, fileName)
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		logger.Debugf("File not found locally: %s", fileName)
		return true, "File not downloaded yet"
	}

	// Check metadata
	localMetadata, exists := metadata[fileID]
	if !exists {
		logger.Debugf("No metadata found for file ID %s", fileID)
		return true, "No metadata for this file"
	}

	// Check date
	if localMetadata.FileDate != file.FileDate {
		logger.Debugf("Date mismatch - Local: %s, Remote: %s", localMetadata.FileDate, file.FileDate)
		return true, fmt.Sprintf("File updated (was: %s, now: %s)", localMetadata.FileDate, file.FileDate)
	}

	// Check hash if available
	remoteHash := ""
	for _, hash := range file.Hashes {
		if hash.Algo == 1 { // SHA-1
			remoteHash = hash.Value
			break
		}
	}

	if remoteHash != "" && localMetadata.Hash != remoteHash {
		logger.Debugf("Hash mismatch - Local: %s, Remote: %s", localMetadata.Hash, remoteHash)
		return true, "File hash changed"
	}

	logger.Debugf("File up to date: %s", fileName)
	return false, "File is current"
}

// RecordDownload records a successful download in metadata
func RecordDownload(file *models.CurseForgeFile, downloadPath string, metadata map[string]models.DownloadMetadata, logger *logrus.Logger) error {
	fileID := strconv.Itoa(file.ID)

	// Get hash
	fileHash := ""
	for _, hash := range file.Hashes {
		if hash.Algo == 1 { // SHA-1
			fileHash = hash.Value
			break
		}
	}

	metadata[fileID] = models.DownloadMetadata{
		FileName:     file.FileName,
		FileDate:     file.FileDate,
		DownloadedAt: time.Now(),
		Hash:         fileHash,
		FileLength:   file.FileLength,
	}

	if err := SaveDownloadMetadata(downloadPath, metadata); err != nil {
		return fmt.Errorf("failed to save download metadata: %w", err)
	}

	logger.Debugf("Recorded download metadata for %s", file.FileName)
	return nil
}
