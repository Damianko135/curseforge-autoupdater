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

// Constants for directories and binary names — adjust these to your project structure
const (
	binaryName = "myapp" // Replace with your binary name
	distDir    = "dist"
	templDir   = "./templates" // Replace with your templ directory
	cliDir     = "./cmd/myapp" // Replace with your CLI directory
	webDir     = "./cmd/web"   // Replace with your web server directory
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

// Dev runs the dev web server with hot reload using Air and templ.
func Dev() error {
	fmt.Println("Running in development mode with Templ and Air...")

	// Check dependencies
	if _, err := exec.LookPath("air"); err != nil {
		fmt.Println("⚠️  'air' not found. Falling back to 'go run'.")
		return sh.RunV("go", "run", webDir)
	}
	if _, err := exec.LookPath("templ"); err != nil {
		fmt.Println("⚠️  'templ' not found. Please install it.")
		return err
	}

	// Generate .air.toml if missing
	if _, err := os.Stat(".air.toml"); os.IsNotExist(err) {
		fmt.Println("⚠️  '.air.toml' not found. Generating custom config...")

		airConfig := `
root = "."
testdata_dir = "testdata"
tmp_dir = "cmd/web/tmp"

[build]
  args_bin = []
  bin = "cmd/web/tmp/main"
  cmd = "go build -o ./cmd/web/tmp/main ./cmd/web"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  silent = false
  time = false

[misc]
  clean_on_exit = false

[proxy]
  app_port = 0
  enabled = false
  proxy_port = 0

[screen]
  clear_on_rebuild = false
  keep_scroll = true
`
		if err := os.WriteFile(".air.toml", []byte(airConfig), 0644); err != nil {
			return fmt.Errorf("failed to write .air.toml: %w", err)
		}
		fmt.Println("✅ Custom .air.toml created.")
	}

	// Run Air in a goroutine
	var eg errgroup.Group
	eg.Go(func() error {
		fmt.Println("🚀 Launching hot-reload dev server with Air...")
		return sh.RunV("air", "-c", ".air.toml")
	})

	// Run templ generate with watch in parallel
	eg.Go(func() error {
		fmt.Println("🚀 Running templ generate with --watch...")
		return sh.RunV("templ", "generate", "--watch")
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// Clean removes build artifacts, including Templ-generated files.
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	var errs []error

	// Remove local binaries
	for _, f := range []string{binaryName, binaryName + ".exe"} {
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

// Fmt formats the codebase using `go fmt`.
func Fmt() error {
	return sh.RunV("go", "fmt", "./...")
}

// GoImports formats the codebase using `goimports`.
func GoImports() error {
	return sh.RunV("goimports", "-w", ".")
}

// Install installs all required dev tools concurrently.
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

// InstallApp installs the application into GOPATH/bin.
func InstallApp() error {
	return sh.RunV("go", "install", cliDir)
}

// Lint runs the linter (requires golangci-lint).
func Lint() error {
	return sh.RunV("golangci-lint", "run")
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

// Run executes the application locally.
func Run() error {
	return sh.RunV("go", "run", cliDir)
}

// Security runs security analysis (requires gosec).
func Security() error {
	return sh.RunV("gosec", "./...")
}

// Setup prepares the development environment (deps + tools).
func Setup() {
	mg.SerialDeps(Deps, Install)
}

// Tidy runs `go mod tidy` explicitly.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}
