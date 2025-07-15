//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/sync/errgroup"
)

// Project structure constants
const (
	cliBinaryName = "curseforge-autoupdate" // CLI binary name
	webBinaryName = "webserver"             // Web server binary name
	distDir       = "dist"
	templDir      = "./templates"
	cliDir        = "./cmd/cli"
	webDir        = "./cmd/web"
	webTmpDir     = "./cmd/web/tmp"
)

// Platforms for Release builds
var platforms = []struct {
	goos, goarch, ext string
}{
	{"windows", "amd64", ".exe"},
	{"linux", "amd64", ""},
	{"darwin", "amd64", ""},
	{"darwin", "arm64", ""},
}

// Dev runs both the web server (with hot reload) and the CLI in dev mode.
// It runs `templ generate --watch` for live template generation,
// starts the web server with Air if available (fallback to go run),
// and runs the CLI with modd if available (fallback to go run).
func Dev() error {
	fmt.Println("Starting development environment for both web and CLI...")

	// Check for required tools
	if _, err := exec.LookPath("templ"); err != nil {
		fmt.Println("❌ 'templ' not found. Please run 'mage install' first.")
		return err
	}

	// Always start templ generate --watch for live template generation
	var eg errgroup.Group
	eg.Go(func() error {
		fmt.Println("▶️  Running 'templ generate --watch' ...")
		return sh.RunV("templ", "generate", "--watch")
	})

	// Web server: prefer Air, fallback to go run
	eg.Go(func() error {
		if _, err := exec.LookPath("air"); err == nil {
			// Generate .air.toml if missing
			if _, err := os.Stat(".air.toml"); os.IsNotExist(err) {
				fmt.Println("⚠️  '.air.toml' not found. Generating default config...")
				airConfig := `
root = "."
tmp_dir = "cmd/web/tmp"
[build]
  bin = "cmd/web/tmp/main"
  cmd = "go build -o ./cmd/web/tmp/main ./cmd/web"
  include_ext = ["go", "templ", "html"]
  log = "build-errors.log"
[color]
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"
[screen]
  keep_scroll = true
`
				if err := os.WriteFile(".air.toml", []byte(airConfig), 0644); err != nil {
					return fmt.Errorf("failed to write .air.toml: %w", err)
				}
				fmt.Println("✅ .air.toml created.")
			}
			fmt.Println("▶️  Running 'air' for web hot reload ...")
			return sh.RunV("air", "-c", ".air.toml")
		}
		fmt.Println("⚠️  'air' not found. Falling back to 'go run' for web server.")
		return sh.RunV("go", "run", webDir)
	})

	// CLI: run in watch mode if modd is available, else just run once
	eg.Go(func() error {
		if _, err := exec.LookPath("modd"); err == nil {
			fmt.Println("▶️  Running 'modd' for CLI hot reload ...")
			moddConf := `
[mods]
**/*.go {
	prep: go run ./cmd/cli
}
`
			if _, err := os.Stat("modd.conf"); os.IsNotExist(err) {
				if err := os.WriteFile("modd.conf", []byte(moddConf), 0644); err != nil {
					return fmt.Errorf("failed to write modd.conf: %w", err)
				}
			}
			return sh.RunV("modd")
		}
		fmt.Println("⚠️  'modd' not found. Running CLI once with 'go run'.")
		return sh.RunV("go", "run", cliDir)
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Clean removes build artifacts, binaries, dist folder, and generated template files.
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	var errs []error

	// Remove local binaries (CLI and web)
	for _, f := range []string{cliBinaryName, cliBinaryName + ".exe", webBinaryName, webBinaryName + ".exe"} {
		if err := sh.Rm(f); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	// Remove dist folder
	if err := os.RemoveAll(distDir); err != nil && !os.IsNotExist(err) {
		errs = append(errs, err)
	}

	// Remove Templ-generated files (*.templ.go)
	err := filepath.Walk(templDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			if matched, _ := filepath.Match("*_templ.go", filepath.Base(path)); matched {
				if rmErr := os.Remove(path); rmErr != nil && !os.IsNotExist(rmErr) {
					errs = append(errs, fmt.Errorf("failed to remove templ file %s: %w", path, rmErr))
				}
			}
		}
		return nil
	})
	if err != nil {
		errs = append(errs, fmt.Errorf("walking templ dir: %w", err))
	}

	// Remove web build artifacts
	if err := os.RemoveAll(webTmpDir); err != nil && !os.IsNotExist(err) {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Println("Clean completed!")
	return nil
}

// Deps downloads Go module dependencies.
func Deps() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "tidy")
}

// Format runs all formatting tools (go fmt + goimports).
func Format() error {
	mg.Deps(Fmt, GoImports)
	return nil
}

// Fmt runs go fmt over the entire codebase.
func Fmt() error {
	return sh.RunV("go", "fmt", "./...")
}

// GoImports runs goimports to fix import grouping and order.
func GoImports() error {
	return sh.RunV("goimports", "-w", ".")
}

// Install installs all required dev tools (templ, air, linting, etc.).
func Install() error {
	tools := []struct {
		binary, url, version string
	}{
		{"golangci-lint", "github.com/golangci/golangci-lint/cmd/golangci-lint", "latest"},
		{"gosec", "github.com/securego/gosec/v2/cmd/gosec", "latest"},
		{"mage", "github.com/magefile/mage", "v1.14.0"},
		{"gofumpt", "mvdan.cc/gofumpt", "v0.4.0"},
		{"goimports", "golang.org/x/tools/cmd/goimports", "latest"},
		{"templ", "github.com/a-h/templ/cmd/templ", "latest"},
		{"air", "github.com/air-verse/air", "latest"},
	}

	var eg errgroup.Group
	for _, t := range tools {
		t := t // capture range variable
		eg.Go(func() error {
			if path, err := exec.LookPath(t.binary); err == nil {
				fmt.Printf("Tool %s already installed at %s\n", t.binary, path)
				return nil
			}

			fmt.Printf("Installing %s...\n", t.binary)
			cmd := []string{"go", "install", fmt.Sprintf("%s@%s", t.url, t.version)}
			if err := sh.RunV(cmd[0], cmd[1:]...); err != nil {
				return fmt.Errorf("failed to install %s: %w", t.binary, err)
			}
			fmt.Printf("Installed %s successfully\n", t.binary)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	fmt.Println("All tools installed successfully!")
	return nil
}

// InstallApp installs the CLI app to your GOPATH/bin.
func InstallApp() error {
	return sh.RunV("go", "install", cliDir)
}

// Lint runs static analysis using golangci-lint.
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Build compiles the CLI and webserver binaries to the dist directory.
// It first generates templ files, then builds the binaries.
func Build() error {
	fmt.Println("Building CLI and web server...")
	mg.Deps(Generate)

	// Create dist directory if missing
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	// Build CLI binary
	cliOut := filepath.Join(distDir, cliBinaryName)
	if err := sh.RunV("go", "build", "-o", cliOut, cliDir); err != nil {
		return fmt.Errorf("failed to build CLI: %w", err)
	}

	// Build web binary
	webOut := filepath.Join(distDir, webBinaryName)
	if err := sh.RunV("go", "build", "-o", webOut, webDir); err != nil {
		return fmt.Errorf("failed to build web: %w", err)
	}

	fmt.Println("Build completed!")
	return nil
}

// Release builds cross-platform CLI binaries into dist.
// Web is not cross-compiled yet.
func Release() error {
	fmt.Println("Building release for all platforms (CLI)...")
	mg.Deps(Generate)

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	var eg errgroup.Group
	for _, p := range platforms {
		p := p // capture range variable
		eg.Go(func() error {
			binary := fmt.Sprintf("%s-%s-%s%s", cliBinaryName, p.goos, p.goarch, p.ext)
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

// Run executes the CLI locally using go run.
func Run() error {
	return sh.RunV("go", "run", cliDir)
}

// RunWeb executes the web server locally after generating templates.
func RunWeb() error {
	mg.Deps(Generate)
	return sh.RunV("go", "run", webDir)
}

// Generate runs templ generate to compile .templ files once (no watch).
func Generate() error {
	fmt.Println("Generating templates with templ...")
	return sh.RunV("templ", "generate")
}

// Test runs unit tests with race detection and coverage.
func Test() error {
	return sh.RunV("go", "test", "-race", "-coverprofile=coverage.out", "./...")
}

// Check runs linting, tests, and format checks in sequence.
func Check() error {
	mg.Deps(Lint, Test, Format)
	return nil
}

// Security runs gosec for security static analysis.
func Security() error {
	return sh.RunV("gosec", "./...")
}

// Setup installs deps and all required dev tools.
func Setup() {
	mg.SerialDeps(Deps, Install)
}

// Tidy runs go mod tidy explicitly.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}


func CI() error {
    mg.Deps(Lint, Test, Format, Security, Build, Release)
    return nil
}