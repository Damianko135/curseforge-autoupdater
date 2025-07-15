package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	// #nosec G304 -- src is validated by caller
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Ensure destination directory exists
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("failed to ensure destination directory for %s: %w", dst, err)
	}

	// #nosec G304 -- dst is validated by caller
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %w", src, dst, err)
	}

	return nil
}

// MoveFile moves a file from src to dst
func MoveFile(src, dst string) error {
	// Ensure destination directory exists
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("failed to ensure destination directory for %s: %w", dst, err)
	}

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", src, dst, err)
	}

	return nil
}

// RemoveFile removes a file
func RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", path, err)
	}
	return nil
}

// RemoveDir removes a directory and all its contents
func RemoveDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", path, err)
	}
	return nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileSize returns the size of a file in bytes, with improved error context
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("GetFileSize: failed to stat %q: %w", path, err)
	}
	if info.IsDir() {
		return 0, fmt.Errorf("GetFileSize: %q is a directory, not a file", path)
	}
	return info.Size(), nil
}

// IsEmpty checks if a directory is empty
func IsEmpty(path string) (bool, error) {
	// #nosec G304 -- path is validated by caller
	dir, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("failed to open directory %s: %w", path, err)
	}
	defer dir.Close()

	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// CopyDir copies a directory from src to dst recursively
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory %s: %w", src, err)
	}

	// #nosec G301 -- dst permissions are inherited from src
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory %s: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListFiles lists all files in a directory (non-recursively)
func ListFiles(path string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	return files, nil
}

// ListDirs lists all directories in a directory (non-recursively)
func ListDirs(path string) ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(path, entry.Name()))
		}
	}

	return dirs, nil
}

// FindFiles finds files matching a pattern in a directory (recursively)
func FindFiles(root, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				matches = append(matches, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", root, err)
	}

	return matches, nil
}

// GetDirSize calculates the total size of a directory and all its contents
func GetDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate directory size for %s: %w", path, err)
	}

	return size, nil
}

// CleanPath cleans and normalizes a file path
func CleanPath(path string) string {
	// Convert to forward slashes and clean
	path = strings.ReplaceAll(path, "\\", "/")
	return filepath.Clean(path)
}

// IsSubPath checks if child is a subdirectory of parent
func IsSubPath(parent, child string) bool {
	parent = CleanPath(parent)
	child = CleanPath(child)

	if parent == child {
		return true
	}

	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}

	return strings.HasPrefix(child, parent)
}

// SafeWriteFile writes data to a file atomically by writing to a temporary file first
func SafeWriteFile(path string, data []byte, perm os.FileMode) error {
	// Create temporary file in the same directory
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "tmp_*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write data to temporary file
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Set permissions
	if err := os.Chmod(tmpPath, perm); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions on temporary file: %w", err)
	}

	// Atomically move temporary file to final location
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to move temporary file to final location: %w", err)
	}

	return nil
}
