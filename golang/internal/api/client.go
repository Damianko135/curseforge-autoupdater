package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client wraps configuration and HTTP client for CurseForge API
type Client struct {
	APIKey     string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

// NewClient creates a new CurseForge API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:    apiKey,
		BaseURL:   "https://api.curseforge.com/v1",
		UserAgent: "CurseForge Auto-Updater/1.0",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// addHeaders sets required headers for each request
func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(method, path string, params map[string]string) (*http.Response, error) {
	// Build URL with parameters
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if params != nil {
		query := u.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		u.RawQuery = query.Encode()
	}

	// Create request
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addHeaders(req)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// GetMod retrieves information about a specific mod
func (c *Client) GetMod(modID int) (*ModInfo, error) {
	path := fmt.Sprintf("/mods/%d", modID)

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("mod with ID %d not found", modID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[ModInfo]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}

// GetModFiles retrieves files for a specific mod
func (c *Client) GetModFiles(modID int, gameVersion string, modLoaderType int, pageSize int, index int) ([]ModFile, error) {
	path := fmt.Sprintf("/mods/%d/files", modID)

	params := make(map[string]string)
	if gameVersion != "" {
		params["gameVersion"] = gameVersion
	}
	if modLoaderType > 0 {
		params["modLoaderType"] = strconv.Itoa(modLoaderType)
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}

	resp, err := c.doRequest("GET", path, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[[]ModFile]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// GetModFile retrieves a specific mod file
func (c *Client) GetModFile(modID, fileID int) (*ModFile, error) {
	path := fmt.Sprintf("/mods/%d/files/%d", modID, fileID)

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file with ID %d not found for mod %d", fileID, modID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[ModFile]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}

// GetModFileDownloadURL retrieves the download URL for a specific mod file
func (c *Client) GetModFileDownloadURL(modID, fileID int) (string, error) {
	path := fmt.Sprintf("/mods/%d/files/%d/download-url", modID, fileID)

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// SearchMods searches for mods based on various criteria
func (c *Client) SearchMods(gameID int, categoryID int, searchFilter string, sortField int, sortOrder string, gameVersion string, pageSize int, index int) ([]ModInfo, error) {
	path := "/mods/search"

	params := make(map[string]string)
	if gameID > 0 {
		params["gameId"] = strconv.Itoa(gameID)
	}
	if categoryID > 0 {
		params["categoryId"] = strconv.Itoa(categoryID)
	}
	if searchFilter != "" {
		params["searchFilter"] = searchFilter
	}
	if sortField > 0 {
		params["sortField"] = strconv.Itoa(sortField)
	}
	if sortOrder != "" {
		params["sortOrder"] = sortOrder
	}
	if gameVersion != "" {
		params["gameVersion"] = gameVersion
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}

	resp, err := c.doRequest("GET", path, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[[]ModInfo]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// GetGameVersions retrieves available game versions
func (c *Client) GetGameVersions(gameID int) ([]GameVersion, error) {
	path := fmt.Sprintf("/games/%d/versions", gameID)

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[[]GameVersion]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// CheckIfModExists checks if a mod with the given ID exists
func (c *Client) CheckIfModExists(modID int) (bool, error) {
	_, err := c.GetMod(modID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetLatestModFile retrieves the latest file for a mod based on game version and release type
func (c *Client) GetLatestModFile(modID int, gameVersion string, releaseType int) (*ModFile, error) {
	files, err := c.GetModFiles(modID, gameVersion, 0, 50, 0)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found for mod %d", modID)
	}

	// Filter by release type if specified
	var filteredFiles []ModFile
	if releaseType > 0 {
		for _, file := range files {
			if file.ReleaseType == releaseType {
				filteredFiles = append(filteredFiles, file)
			}
		}
	} else {
		filteredFiles = files
	}

	if len(filteredFiles) == 0 {
		return nil, fmt.Errorf("no files found for mod %d with release type %d", modID, releaseType)
	}

	// Return the most recent file (files are usually sorted by date)
	return &filteredFiles[0], nil
}

// DownloadFile downloads a file from the given URL
func (c *Client) DownloadFile(url string, writer io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write downloaded data: %w", err)
	}

	return nil
}
