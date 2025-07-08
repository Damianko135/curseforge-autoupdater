//go:build mage

// vscode:ignoreFileError
package main

import (
	"errors"
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

var platforms = []struct {
	goos, goarch, ext string
}{
	{"windows", "amd64", ".exe"},
	{"linux", "amd64", ""},
	{"darwin", "amd64", ""},
	{"darwin", "arm64", ""},
}

// Build builds the application for the current platform.
func Build() error {
	fmt.Println("Building application...")

	binary := binaryName
	goos, err := sh.Output("go", "env", "GOOS")
	if err != nil {
		return fmt.Errorf("detecting GOOS: %w", err)
	}
	if goos == "windows" {
		binary += ".exe"
	}

	return sh.RunV("go", "build", "-trimpath", "-o", binary, cmdDir)
}

// BuildAll builds the application for all target platforms.
func BuildAll() error {
	fmt.Println("Building for all platforms...")

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	for _, p := range platforms {
		binary := fmt.Sprintf("%s-%s-%s%s", binaryName, p.goos, p.goarch, p.ext)
		binaryPath := filepath.Join(distDir, binary)

		fmt.Printf("Building %s...\n", binary)
		env := map[string]string{
			"GOOS":   p.goos,
			"GOARCH": p.goarch,
		}

		if err := sh.RunWith(env, "go", "build", "-trimpath", "-o", binaryPath, cmdDir); err != nil {
			return fmt.Errorf("failed to build %s: %w", binary, err)
		}
	}

	fmt.Println("All builds completed successfully!")
	return nil
}

// Deps downloads and tidies Go module dependencies.
func Deps() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "tidy")
}

// Test runs unit tests.
func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
}

// Run executes the application locally.
func Run() error {
	return sh.RunV("go", "run", cmdDir)
}

// Clean removes build artifacts.
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	var errs []error

	for _, f := range []string{binaryName, binaryName + ".exe"} {
		if err := sh.Rm(f); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	if err := os.RemoveAll(distDir); err != nil && !os.IsNotExist(err) {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Println("Clean completed!")
	return nil
}

// InstallApp installs the application into GOPATH/bin.
func InstallApp() error {
	return sh.RunV("go", "install", cmdDir)
}

// Fmt formats the codebase.
func Fmt() error {
	return sh.RunV("go", "fmt", "./...")
}

// Lint runs the linter (requires golangci-lint).
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Security runs security analysis (requires gosec).
func Security() error {
	return sh.RunV("gosec", "./...")
}

// Install installs all required dev tools.
func Install() error {
	tools := []struct {
		name, url, version string
	}{
		{"golangci-lint", "github.com/golangci/golangci-lint/cmd/golangci-lint", "latest"},
		{"gosec", "github.com/securego/gosec/v2/cmd/gosec", "latest"},
	}

	for _, t := range tools {
		fmt.Printf("Installing %s...\n", t.name)
		if err := sh.RunV("go", "install", t.url+"@"+t.version); err != nil {
			return fmt.Errorf("failed to install %s: %w", t.name, err)
		}
	}

	fmt.Println("All tools installed successfully!")
	return nil
}

// Setup prepares the development environment (deps + tools).
func Setup() {
	mg.SerialDeps(Deps, Install)
}

// All runs the full build pipeline: Deps, Fmt, Test, Build.
func All() {
	mg.SerialDeps(Deps, Fmt, Test, Build)
}

// CI runs the CI pipeline: Deps, Fmt, Lint, Test, BuildAll.
func CI() {
	mg.SerialDeps(Deps, Fmt, Lint, Test, BuildAll)
}
