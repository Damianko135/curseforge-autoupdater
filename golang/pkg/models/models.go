package models

import (
	"time"
)

// CurseForgeFile represents a file from the CurseForge API
type CurseForgeFile struct {
	ID          int    `json:"id"`
	FileName    string `json:"fileName"`
	DisplayName string `json:"displayName"`
	FileDate    string `json:"fileDate"`
	FileLength  int64  `json:"fileLength"`
	DownloadURL string `json:"downloadUrl"`
	Hashes      []Hash `json:"hashes"`
}

// Hash represents a file hash
type Hash struct {
	Algo  int    `json:"algo"`
	Value string `json:"value"`
}

// FilesResponse represents the API response for mod files
type FilesResponse struct {
	Data       []CurseForgeFile `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// Pagination represents pagination info
type Pagination struct {
	Index       int `json:"index"`
	PageSize    int `json:"pageSize"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}

// ModInfo represents basic mod information
type ModInfo struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	GameID  int      `json:"gameId"`
	ClassID int      `json:"classId"`
	Authors []Author `json:"authors"`
}

// Author represents a mod author
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ModInfoResponse represents the API response for mod info
type ModInfoResponse struct {
	Data ModInfo `json:"data"`
}

// DownloadMetadata represents metadata about downloaded files
type DownloadMetadata struct {
	FileName     string    `json:"fileName"`
	FileDate     string    `json:"fileDate"`
	DownloadedAt time.Time `json:"downloadedAt"`
	Hash         string    `json:"hash,omitempty"`
	FileLength   int64     `json:"fileLength"`
}

// Config represents the application configuration
type Config struct {
	APIKey       string `koanf:"api_key"`
	ModID        string `koanf:"mod_id"`
	DownloadPath string `koanf:"download_path"`
	GameID       int    `koanf:"game_id"`
	LogLevel     string `koanf:"log_level"`
}
