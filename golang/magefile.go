//go:build mage

// vscode:ignoreFileError
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binaryName = "curseforge-updater"
	distDir    = "dist"
	cmdDir     = "./cmd"
)

// Build builds the application for the current platform
func Build() error {
	fmt.Println("Building application...")
	binary := binaryName

	// Check if we're on Windows
	goos, err := sh.Output("go", "env", "GOOS")
	if err != nil {
		return err
	}
	if goos == "windows" {
		binary += ".exe"
	}

	return sh.Run("go", "build", "-o", binary, cmdDir)
}

// BuildAll builds the application for multiple platforms
func BuildAll() error {
	fmt.Println("Building for all platforms...")

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	platforms := []struct {
		goos   string
		goarch string
		ext    string
	}{
		{"windows", "amd64", ".exe"},
		{"linux", "amd64", ""},
		{"darwin", "amd64", ""},
		{"darwin", "arm64", ""},
	}

	for _, platform := range platforms {
		binary := fmt.Sprintf("%s-%s-%s%s", binaryName, platform.goos, platform.goarch, platform.ext)
		binaryPath := filepath.Join(distDir, binary)

		fmt.Printf("Building %s...\n", binary)
		env := map[string]string{
			"GOOS":   platform.goos,
			"GOARCH": platform.goarch,
		}

		if err := sh.RunWith(env, "go", "build", "-o", binaryPath, cmdDir); err != nil {
			return fmt.Errorf("failed to build %s: %w", binary, err)
		}
	}

	fmt.Println("All builds completed successfully!")
	return nil
}

// Deps downloads and tidies dependencies
func Deps() error {
	fmt.Println("Downloading dependencies...")
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	fmt.Println("Tidying dependencies...")
	return sh.Run("go", "mod", "tidy")
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "-v", "./...")
}

// Run runs the application
func Run() error {
	fmt.Println("Running application...")
	return sh.Run("go", "run", cmdDir)
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	// Remove binary from current directory
	if err := sh.Rm(binaryName); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := sh.Rm(binaryName + ".exe"); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Remove dist directory
	if err := os.RemoveAll(distDir); err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Println("Clean completed!")
	return nil
}

// Install installs the application to GOPATH/bin
func Install() error {
	fmt.Println("Installing application...")
	return sh.Run("go", "install", cmdDir)
}

// Fmt formats the code
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.Run("go", "fmt", "./...")
}

// Lint runs golangci-lint (requires golangci-lint to be installed)
func Lint() error {
	fmt.Println("Running linter...")
	return sh.Run("golangci-lint", "run")
}

// Security runs security checks (requires gosec to be installed)
func Security() error {
	fmt.Println("Running security checks...")
	return sh.Run("gosec", "./...")
}

// Dev runs the application with air for live reload (requires air to be installed)
func Dev() error {
	fmt.Println("Starting development server with live reload...")
	return sh.Run("air")
}

// All runs a complete build pipeline: deps, fmt, test, build
func All() {
	mg.SerialDeps(Deps, Fmt, Test, Build)
}

// CI runs the CI pipeline: deps, fmt, lint, test, build-all
func CI() {
	mg.SerialDeps(Deps, Fmt, Lint, Test, BuildAll)
}
