package api

import (
	"time"
)

// APIResponse is a generic wrapper used by CurseForge API
type APIResponse[T any] struct {
	Data       T           `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Index       int `json:"index"`
	PageSize    int `json:"pageSize"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}

// ModInfo represents basic mod information
type ModInfo struct {
	ID                   int         `json:"id"`
	GameID               int         `json:"gameId"`
	Name                 string      `json:"name"`
	Slug                 string      `json:"slug"`
	Summary              string      `json:"summary"`
	Status               int         `json:"status"`
	DownloadCount        int64       `json:"downloadCount"`
	IsFeatured           bool        `json:"isFeatured"`
	PrimaryCategoryID    int         `json:"primaryCategoryId"`
	Categories           []Category  `json:"categories"`
	ClassID              int         `json:"classId"`
	Authors              []Author    `json:"authors"`
	Logo                 *ModAsset   `json:"logo"`
	Screenshots          []ModAsset  `json:"screenshots"`
	MainFileID           int         `json:"mainFileId"`
	LatestFiles          []ModFile   `json:"latestFiles"`
	LatestFilesIndexes   []FileIndex `json:"latestFilesIndexes"`
	DateCreated          time.Time   `json:"dateCreated"`
	DateModified         time.Time   `json:"dateModified"`
	DateReleased         time.Time   `json:"dateReleased"`
	AllowModDistribution bool        `json:"allowModDistribution"`
	GamePopularityRank   int         `json:"gamePopularityRank"`
	IsAvailable          bool        `json:"isAvailable"`
	ThumbsUpCount        int         `json:"thumbsUpCount"`
	Rating               float64     `json:"rating"`
}

// Category represents a mod category
type Category struct {
	ID               int       `json:"id"`
	GameID           int       `json:"gameId"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	URL              string    `json:"url"`
	IconURL          string    `json:"iconUrl"`
	DateModified     time.Time `json:"dateModified"`
	IsClass          bool      `json:"isClass"`
	ClassID          int       `json:"classId"`
	ParentCategoryID int       `json:"parentCategoryId"`
}

// Author represents a mod author
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ModAsset represents a mod asset (logo, screenshot, etc.)
type ModAsset struct {
	ID           int    `json:"id"`
	ModID        int    `json:"modId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	URL          string `json:"url"`
}

// ModFile represents a mod file
type ModFile struct {
	ID                   int                   `json:"id"`
	GameID               int                   `json:"gameId"`
	ModID                int                   `json:"modId"`
	IsAvailable          bool                  `json:"isAvailable"`
	DisplayName          string                `json:"displayName"`
	FileName             string                `json:"fileName"`
	ReleaseType          int                   `json:"releaseType"`
	FileStatus           int                   `json:"fileStatus"`
	Hashes               []FileHash            `json:"hashes"`
	FileDate             time.Time             `json:"fileDate"`
	FileLength           int64                 `json:"fileLength"`
	DownloadCount        int64                 `json:"downloadCount"`
	FileSizeOnDisk       int64                 `json:"fileSizeOnDisk"`
	DownloadURL          string                `json:"downloadUrl"`
	GameVersions         []string              `json:"gameVersions"`
	SortableGameVersions []SortableGameVersion `json:"sortableGameVersions"`
	Dependencies         []ModDependency       `json:"dependencies"`
	ExposeAsAlternative  bool                  `json:"exposeAsAlternative"`
	ParentProjectFileID  int                   `json:"parentProjectFileId"`
	AlternateFileID      int                   `json:"alternateFileId"`
	IsServerPack         bool                  `json:"isServerPack"`
	ServerPackFileID     int                   `json:"serverPackFileId"`
	IsEarlyAccessContent bool                  `json:"isEarlyAccessContent"`
	EarlyAccessEndDate   *time.Time            `json:"earlyAccessEndDate"`
	FileFingerprint      int64                 `json:"fileFingerprint"`
	Modules              []ModModule           `json:"modules"`
}

// FileHash represents a file hash
type FileHash struct {
	Value string `json:"value"`
	Algo  int    `json:"algo"`
}

// SortableGameVersion represents a sortable game version
type SortableGameVersion struct {
	GameVersionName        string    `json:"gameVersionName"`
	GameVersionPadded      string    `json:"gameVersionPadded"`
	GameVersion            string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionTypeID      int       `json:"gameVersionTypeId"`
}

// ModDependency represents a mod dependency
type ModDependency struct {
	ModID        int `json:"modId"`
	RelationType int `json:"relationType"`
	FileID       int `json:"fileId"`
}

// ModModule represents a mod module
type ModModule struct {
	Name        string `json:"name"`
	Fingerprint int64  `json:"fingerprint"`
}

// FileIndex represents a file index
type FileIndex struct {
	GameVersion       string `json:"gameVersion"`
	FileID            int    `json:"fileId"`
	Filename          string `json:"filename"`
	ReleaseType       int    `json:"releaseType"`
	GameVersionTypeID int    `json:"gameVersionTypeId"`
	ModLoader         int    `json:"modLoader"`
}

// GameVersion represents a game version
type GameVersion struct {
	ID                int       `json:"id"`
	GameID            int       `json:"gameId"`
	VersionString     string    `json:"versionString"`
	DateModified      time.Time `json:"dateModified"`
	GameVersionTypeID int       `json:"gameVersionTypeId"`
}

// ReleaseType constants
const (
	ReleaseTypeRelease int = 1
	ReleaseTypeBeta    int = 2
	ReleaseTypeAlpha   int = 3
)

// RelationType constants for dependencies
const (
	RelationTypeEmbeddedLibrary    int = 1
	RelationTypeOptionalDependency int = 2
	RelationTypeRequiredDependency int = 3
	RelationTypeTool               int = 4
	RelationTypeIncompatible       int = 5
	RelationTypeInclude            int = 6
)

// FileStatus constants
const (
	FileStatusProcessing         int = 1
	FileStatusChangesRequired    int = 2
	FileStatusUnderReview        int = 3
	FileStatusApproved           int = 4
	FileStatusRejected           int = 5
	FileStatusMalwareDetected    int = 6
	FileStatusDeleted            int = 7
	FileStatusArchived           int = 8
	FileStatusTesting            int = 9
	FileStatusReleased           int = 10
	FileStatusReadyForReview     int = 11
	FileStatusDeprecated         int = 12
	FileStatusBaking             int = 13
	FileStatusAwaitingPublishing int = 14
	FileStatusFailedPublishing   int = 15
)

// ModStatus constants
const (
	ModStatusNew             int = 1
	ModStatusChangesRequired int = 2
	ModStatusUnderSoftReview int = 3
	ModStatusApproved        int = 4
	ModStatusRejected        int = 5
	ModStatusChangesMade     int = 6
	ModStatusInactive        int = 7
	ModStatusAbandoned       int = 8
	ModStatusDeleted         int = 9
	ModStatusUnderReview     int = 10
)
