package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client wraps configuration and HTTP client
type Client struct {
	APIKey     string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:    apiKey,
		BaseURL:   "https://api.curseforge.com/v1",
		UserAgent: "CurseForge Auto-Updater PoC/1.0",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// addHeaders sets headers for each request
func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)
}

// ModInfo is a basic struct for mod data (simplified, adjust per API spec)
type ModInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// APIResponse is a generic wrapper used by CurseForge
type APIResponse[T any] struct {
	Data T `json:"data"`
}

// CheckIfExists checks if a mod with a given ID exists
func (c *Client) CheckIfExists(id int) (bool, error) {
	url := fmt.Sprintf("%s/mods/%d", c.BaseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request failed: %w", err)
	}

	c.addHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result APIResponse[ModInfo]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Optional sanity check
	if result.Data.ID != id {
		return false, errors.New("mod ID mismatch in response")
	}

	return true, nil
}
