package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/damianko135/curseforge-autoupdate/golang/helper/filesystem"
	"github.com/klauspost/compress/zip"
)

// BackupManager handles server backups
type BackupManager struct {
	serverPath  string
	backupPath  string
	compression bool
	retention   int // days
}

// NewBackupManager creates a new backup manager
func NewBackupManager(serverPath, backupPath string, compression bool, retention int) *BackupManager {
	return &BackupManager{
		serverPath:  serverPath,
		backupPath:  backupPath,
		compression: compression,
		retention:   retention,
	}
}

// BackupInfo represents information about a backup
type BackupInfo struct {
	Name         string
	Path         string
	Size         int64
	Created      time.Time
	IsCompressed bool
	Type         string // full, incremental, pre-update, etc.
}

// CreateBackup creates a new backup
func (bm *BackupManager) CreateBackup(name string, backupType string) (*BackupInfo, error) {
	// Ensure backup directory exists
	if err := filesystem.EnsureDir(bm.backupPath); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup name if not provided
	if name == "" {
		name = fmt.Sprintf("backup_%s", time.Now().Format("20060102_150405"))
	}

	// Add type suffix if provided
	if backupType != "" {
		name = fmt.Sprintf("%s_%s", name, backupType)
	}

	var backupFilePath string
	var err error

	if bm.compression {
		backupFilePath = filepath.Join(bm.backupPath, name+".zip")
		err = bm.createCompressedBackup(backupFilePath)
	} else {
		backupFilePath = filepath.Join(bm.backupPath, name)
		err = bm.createUncompressedBackup(backupFilePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Get backup size
	size, err := bm.getBackupSize(backupFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup size: %w", err)
	}

	return &BackupInfo{
		Name:         name,
		Path:         backupFilePath,
		Size:         size,
		Created:      time.Now(),
		IsCompressed: bm.compression,
		Type:         backupType,
	}, nil
}

// createCompressedBackup creates a compressed backup
func (bm *BackupManager) createCompressedBackup(backupPath string) error {
	// Create zip file
	// #nosec G304 -- backupPath is constructed internally
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through server directory and add files to zip
	err = filepath.Walk(bm.serverPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk error at %s: %w", path, walkErr)
		}

		// Skip certain files/directories
		if bm.shouldSkipFile(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, relErr := filepath.Rel(bm.serverPath, path)
		if relErr != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, relErr)
		}

		if info.IsDir() {
			// Create directory entry
			header := &zip.FileHeader{
				Name:     relPath + "/",
				Method:   zip.Store,
				Modified: info.ModTime(),
			}
			_, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("failed to create zip dir header for %s: %w", relPath, err)
			}
			return nil
		}

		// Create file entry
		header := &zip.FileHeader{
			Name:     relPath,
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip file header for %s: %w", relPath, err)
		}

		// Copy file content to zip
		// #nosec G304 -- path is validated by Walk
		file, openErr := os.Open(path)
		if openErr != nil {
			return fmt.Errorf("failed to open file %s: %w", path, openErr)
		}
		defer file.Close()

		if _, copyErr := io.Copy(writer, file); copyErr != nil {
			return fmt.Errorf("failed to copy file %s to zip: %w", path, copyErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("backup zip creation failed: %w", err)
	}
	return nil
}

// createUncompressedBackup creates an uncompressed backup
func (bm *BackupManager) createUncompressedBackup(backupPath string) error {
	return filesystem.CopyDir(bm.serverPath, backupPath)
}

// shouldSkipFile determines if a file should be skipped during backup
func (bm *BackupManager) shouldSkipFile(path string, info os.FileInfo) bool {
	// Skip lock files
	if strings.HasSuffix(path, ".lock") {
		return true
	}

	// Skip log files (optional)
	if strings.HasSuffix(path, ".log") || strings.HasSuffix(path, ".log.gz") {
		return true
	}

	// Skip temporary files
	if strings.HasPrefix(info.Name(), "tmp_") || strings.HasPrefix(info.Name(), "temp_") {
		return true
	}

	// Skip certain directories
	baseName := filepath.Base(path)
	skipDirs := []string{"logs", "cache", "tmp", "temp"}
	for _, skipDir := range skipDirs {
		if baseName == skipDir {
			return true
		}
	}

	return false
}

// getBackupSize calculates the size of a backup
func (bm *BackupManager) getBackupSize(backupPath string) (int64, error) {
	if filesystem.DirExists(backupPath) {
		return filesystem.GetDirSize(backupPath)
	}
	if filesystem.FileExists(backupPath) {
		size, err := filesystem.GetFileSize(backupPath)
		if err != nil {
			return 0, fmt.Errorf("getBackupSize: failed to get file size for %q: %w", backupPath, err)
		}
		return size, nil
	}
	return 0, fmt.Errorf("getBackupSize: path does not exist: %q", backupPath)
}

// ListBackups lists all available backups
func (bm *BackupManager) ListBackups() ([]BackupInfo, error) {
	if !filesystem.DirExists(bm.backupPath) {
		return []BackupInfo{}, nil
	}

	entries, err := os.ReadDir(bm.backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".zip") {
			backupPath := filepath.Join(bm.backupPath, entry.Name())
			var createdTime time.Time
			size, err := bm.getBackupSize(backupPath)
			if err != nil {
				size = 0
			}
			if filesystem.FileExists(backupPath) || filesystem.DirExists(backupPath) {
				info, err := os.Stat(backupPath)
				if err == nil {
					createdTime = info.ModTime()
				}
			}
			backupInfo := BackupInfo{
				Name:         entry.Name(),
				Path:         backupPath,
				Size:         size,
				Created:      createdTime,
				IsCompressed: strings.HasSuffix(entry.Name(), ".zip"),
			}

			// Try to determine backup type from name
			if strings.Contains(entry.Name(), "_pre_update") {
				backupInfo.Type = "pre-update"
			} else if strings.Contains(entry.Name(), "_post_update") {
				backupInfo.Type = "post-update"
			} else if strings.Contains(entry.Name(), "_manual") {
				backupInfo.Type = "manual"
			} else {
				backupInfo.Type = "automatic"
			}

			backups = append(backups, backupInfo)
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Created.After(backups[j].Created)
	})

	return backups, nil
}

// RestoreBackup restores a backup
func (bm *BackupManager) RestoreBackup(backupName string) error {
	// Find backup
	backups, err := bm.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	var targetBackup *BackupInfo
	for _, backup := range backups {
		if backup.Name == backupName {
			targetBackup = &backup
			break
		}
	}

	if targetBackup == nil {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Create temporary restore directory
	tempDir := filepath.Join(bm.backupPath, "temp_restore_"+time.Now().Format("20060102_150405"))
	if err := filesystem.EnsureDir(tempDir); err != nil {
		return fmt.Errorf("failed to create temp restore directory: %w", err)
	}
	// Clean up tempDir after restore, log error if any
	defer func() {
		if err := filesystem.RemoveDir(tempDir); err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] failed to remove temp restore dir %s: %v\n", tempDir, err)
		}
	}()

	// Extract/copy backup to temp directory
	if targetBackup.IsCompressed {
		if err := bm.extractBackup(targetBackup.Path, tempDir); err != nil {
			return fmt.Errorf("failed to extract backup: %w", err)
		}
	} else {
		if err := filesystem.CopyDir(targetBackup.Path, tempDir); err != nil {
			return fmt.Errorf("failed to copy backup: %w", err)
		}
	}

	// Remove current server directory
	if err := filesystem.RemoveDir(bm.serverPath); err != nil {
		return fmt.Errorf("failed to remove current server directory: %w", err)
	}

	// Move temp directory to server path
	if err := filesystem.MoveFile(tempDir, bm.serverPath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return nil
}

// extractBackup extracts a compressed backup
func (bm *BackupManager) extractBackup(backupPath, targetPath string) error {
	// Open zip file
	reader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer reader.Close()

	// Extract files
	for _, file := range reader.File {
		filePath := filepath.Join(targetPath, file.Name)

		if file.FileInfo().IsDir() {
			// Create directory
			if err := filesystem.EnsureDir(filePath); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filePath, err)
			}
			continue
		}

		// Create file
		if err := filesystem.EnsureDir(filepath.Dir(filePath)); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", filePath, err)
		}

		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in backup: %w", err)
		}
		defer fileReader.Close()

		// #nosec G304 -- filePath is constructed internally
		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", filePath, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, fileReader)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}

// DeleteBackup deletes a backup
func (bm *BackupManager) DeleteBackup(backupName string) error {
	backupPath := filepath.Join(bm.backupPath, backupName)

	if !filesystem.FileExists(backupPath) && !filesystem.DirExists(backupPath) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	if filesystem.DirExists(backupPath) {
		return filesystem.RemoveDir(backupPath)
	}

	return filesystem.RemoveFile(backupPath)
}

// CleanupOldBackups removes old backups based on retention policy
func (bm *BackupManager) CleanupOldBackups() error {
	if bm.retention <= 0 {
		return nil // No retention policy
	}

	backups, err := bm.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -bm.retention)

	for _, backup := range backups {
		if backup.Created.Before(cutoffTime) {
			if err := bm.DeleteBackup(backup.Name); err != nil {
				return fmt.Errorf("failed to delete old backup %s: %w", backup.Name, err)
			}
		}
	}

	return nil
}

// GetBackupInfo gets information about a specific backup
func (bm *BackupManager) GetBackupInfo(backupName string) (*BackupInfo, error) {
	backups, err := bm.ListBackups()
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	for _, backup := range backups {
		if backup.Name == backupName {
			return &backup, nil
		}
	}

	return nil, fmt.Errorf("backup not found: %s", backupName)
}

// ValidateBackup validates a backup file
func (bm *BackupManager) ValidateBackup(backupName string) error {
	backup, err := bm.GetBackupInfo(backupName)
	if err != nil {
		return err
	}

	// Check if backup file exists
	if !filesystem.FileExists(backup.Path) && !filesystem.DirExists(backup.Path) {
		return fmt.Errorf("backup file does not exist: %s", backup.Path)
	}

	// If compressed, try to open the zip file
	if backup.IsCompressed {
		reader, err := zip.OpenReader(backup.Path)
		if err != nil {
			return fmt.Errorf("failed to open backup zip file: %w", err)
		}
		defer reader.Close()

		// Check if zip contains expected files
		hasServerProperties := false
		for _, file := range reader.File {
			if filepath.Base(file.Name) == "server.properties" {
				hasServerProperties = true
				break
			}
		}

		if !hasServerProperties {
			return fmt.Errorf("backup appears to be invalid: missing server.properties")
		}
	}

	return nil
}

// GetBackupSpace returns the total space used by backups
func (bm *BackupManager) GetBackupSpace() (int64, error) {
	if !filesystem.DirExists(bm.backupPath) {
		return 0, nil
	}

	return filesystem.GetDirSize(bm.backupPath)
}

// CreatePreUpdateBackup creates a backup before updating
func (bm *BackupManager) CreatePreUpdateBackup(version string) (*BackupInfo, error) {
	name := fmt.Sprintf("pre_update_%s_%s", version, time.Now().Format("20060102_150405"))
	return bm.CreateBackup(name, "pre-update")
}

// CreatePostUpdateBackup creates a backup after updating
func (bm *BackupManager) CreatePostUpdateBackup(version string) (*BackupInfo, error) {
	name := fmt.Sprintf("post_update_%s_%s", version, time.Now().Format("20060102_150405"))
	return bm.CreateBackup(name, "post-update")
}

// CreateManualBackup creates a manual backup
func (bm *BackupManager) CreateManualBackup(name string) (*BackupInfo, error) {
	if name == "" {
		name = fmt.Sprintf("manual_%s", time.Now().Format("20060102_150405"))
	} else {
		name = fmt.Sprintf("manual_%s_%s", name, time.Now().Format("20060102_150405"))
	}
	return bm.CreateBackup(name, "manual")
}

// GetLatestBackup returns the most recent backup
func (bm *BackupManager) GetLatestBackup() (*BackupInfo, error) {
	backups, err := bm.ListBackups()
	if err != nil {
		return nil, err
	}

	if len(backups) == 0 {
		return nil, fmt.Errorf("no backups found")
	}

	return &backups[0], nil
}

// UpdateRetentionPolicy updates the retention policy
func (bm *BackupManager) UpdateRetentionPolicy(days int) {
	bm.retention = days
}

// EnableCompression enables or disables compression
func (bm *BackupManager) EnableCompression(enabled bool) {
	bm.compression = enabled
}
