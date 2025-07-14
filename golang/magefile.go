//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binaryName = "cf-updater"
	distDir    = "dist"
	cliDir     = "./cmd/cli"
	webDir     = "./cmd/web"
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

	return sh.RunV("go", "build", "-trimpath", "-o", binary, cliDir)
}

// Release builds the application for all target platforms concurrently.
func Release() error {
	fmt.Println("Building release for all platforms...")

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	var eg errgroup.Group

	for _, p := range platforms {
		p := p // capture range variable
		eg.Go(func() error {
			binary := fmt.Sprintf("%s-%s-%s%s", binaryName, p.goos, p.goarch, p.ext)
			binaryPath := filepath.Join(distDir, binary)

			fmt.Printf("Building %s...\n", binary)
			env := map[string]string{
				"GOOS":   p.goos,
				"GOARCH": p.goarch,
			}

			if err := sh.RunWith(env, "go", "build", "-trimpath", "-o", binaryPath, cliDir); err != nil {
				return fmt.Errorf("failed to build %s: %w", binary, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	fmt.Println("All builds completed successfully!")
	return nil
}

// Deps downloads Go module dependencies.
func Deps() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "tidy")
}

// Tidy runs `go mod tidy` explicitly.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// Fmt formats the codebase using `go fmt`.
func Fmt() error {
	return sh.RunV("go", "fmt", "./...")
}

// GoImports formats the codebase using `goimports`.
func GoImports() error {
	return sh.RunV("goimports", "-w", ".")
}

// Format runs all formatting tools (go fmt + goimports).
func Format() error {
	mg.Deps(Fmt, GoImports)
	return nil
}

// Test runs unit tests.
func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
}

// Run executes the application locally.
func Run() error {
	return sh.RunV("go", "run", cliDir)
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
	return sh.RunV("go", "install", cliDir)
}

// Lint runs the linter (requires golangci-lint).
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Security runs security analysis (requires gosec).
func Security() error {
	return sh.RunV("gosec", "./...")
}

// Install installs all required dev tools concurrently.
func Install() error {
	tools := []struct {
		name, url, version string
	}{
		{"golangci-lint", "github.com/golangci/golangci-lint/cmd/golangci-lint", "latest"},
		{"gosec", "github.com/securego/gosec/v2/cmd/gosec", "latest"},
		{"go-toml", "github.com/pelletier/go-toml/v2/cmd/tomlv", "latest"},
		{"mage", "github.com/magefile/mage", "v1.14.0"},
		{"gofumpt", "mvdan.cc/gofumpt", "v0.4.0"},
		{"goimports", "golang.org/x/tools/cmd/goimports", "latest"},
	}

	var eg errgroup.Group
	for _, t := range tools {
		t := t // capture variable
		eg.Go(func() error {
			fmt.Printf("Installing %s...\n", t.name)
			if err := sh.RunV("go", "install", t.url+"@"+t.version); err != nil {
				return fmt.Errorf("failed to install %s: %w", t.name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	fmt.Println("All tools installed successfully!")
	return nil
}

// Setup prepares the development environment (deps + tools).
func Setup() {
	mg.SerialDeps(Deps, Install)
}

// All runs the full build pipeline: Clean, Deps, Format, Test, Build.
func All() {
	mg.SerialDeps(Clean, Deps, Format, Test, Build)
}

// CI runs the CI pipeline: Clean, Deps, Format, Test, then Lint + Security in parallel, finally Release.
func CI() error {
	// Run initial tasks serially
	mg.SerialDeps(Clean, Deps, Format, Test)

	// Run Lint and Security concurrently
	var eg errgroup.Group
	eg.Go(Lint)
	eg.Go(Security)
	if err := eg.Wait(); err != nil {
		return err
	}

	// Final release build
	return Release()
}
