//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "github.com/magefile/mage/mage"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Project structure constants
const (
	cliBinaryName = "curseforge-autoupdate"
	webBinaryName = "webserver"
	distDir       = "dist"
	templDir      = "./views"
	cliDir        = "./cmd/cli"
	webDir        = "./cmd/web"
	webTmpDir     = "./cmd/web/tmp"
)

// func main() {
// 	os.Exit(mage.Main())
// }

// Platforms for Release builds
var platforms = []struct {
	goos, goarch, ext string
}{
	{"windows", "amd64", ".exe"},
	{"linux", "amd64", ""},
	{"darwin", "amd64", ""},
	{"darwin", "arm64", ""},
}

func ensureAirToml() error {
	if _, err := os.Stat(".air.toml"); errors.Is(err, os.ErrNotExist) {
		airConfig := `
		root = "."
testdata_dir = "testdata"
tmp_dir = "cmd/web/tmp"

[build]
  args_bin = []
  bin = "cmd\\web\\tmp\\main.exe"
  cmd = "go build -o ./cmd/web/tmp/main.exe ./cmd/web"
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
		return os.WriteFile(".air.toml", []byte(airConfig), 0644)
	}
	return nil
}

// Dev starts the CLI and web server in development mode with hot-reload enabled.
func Dev() error {
	mg.Deps(Install, Generate)
	if err := ensureAirToml(); err != nil {
		return err
	}
	return sh.RunV("air")
}

// Clean removes generated files, binaries, and the dist directory.
func Clean() error {
	fmt.Println("Cleaning project...")
	_ = os.RemoveAll(distDir)
	_ = os.RemoveAll("tmp")
	return sh.RunV("go", "clean", "-modcache")
}

// Deps ensures Go module dependencies are downloaded.
func Deps() error {
	return sh.RunV("go", "mod", "download")
}

// Format runs go fmt and goimports on the entire codebase.
func Format() error {
	mg.Deps(FormatGo, FormatImports)
	return nil
}

// FormatGo formats Go source code using go fmt.
func FormatGo() error {
	return sh.RunV("go", "fmt", "./...")
}

// FormatImports organizes Go imports using goimports.
func FormatImports() error {
	return sh.RunV("goimports", "-w", ".")
}

// Install sets up development tools such as templ, air, and linters.
func Install() error {
	tools := []string{
		"github.com/air-verse/air@latest",
		"github.com/incu6us/goimports-reviser/v2@latest",
		"github.com/go-task/task/v3/cmd/task@latest",
		"github.com/securego/gosec/v2/cmd/gosec@latest",
		"github.com/a-h/templ/cmd/templ@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	}
	for _, tool := range tools {
		if err := sh.RunV("go", "install", tool); err != nil {
			return err
		}
	}
	return nil
}

// InstallApp compiles and installs the CLI application.
func InstallApp() error {
	return sh.RunV("go", "install", filepath.Join(cliDir))
}

// Lint runs static code analysis using golangci-lint.
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Build compiles the CLI and webserver into the dist directory.
func Build() error {
	mg.Deps(Generate)
	_ = os.MkdirAll(distDir, 0755)
	ldflags := buildLdflags()
	if err := sh.RunV("go", "build", "-ldflags", ldflags, "-o", filepath.Join(distDir, cliBinaryName), cliDir); err != nil {
		return err
	}
	return sh.RunV("go", "build", "-ldflags", ldflags, "-o", filepath.Join(distDir, webBinaryName), webDir)
}

// Release builds binaries for multiple OS/architectures.
func Release() error {
	mg.Deps(Generate)
	ldflags := buildLdflags()
	for _, platform := range platforms {
		output := fmt.Sprintf("%s-%s-%s%s", cliBinaryName, platform.goos, platform.goarch, platform.ext)
		env := map[string]string{"GOOS": platform.goos, "GOARCH": platform.goarch}
		if err := sh.RunWithV(env, "go", "build", "-ldflags", ldflags, "-o", filepath.Join(distDir, output), cliDir); err != nil {
			return err
		}
	}
	return nil
}

func buildLdflags() string {
	version := "dev"
	commit := "unknown"
	if out, err := sh.Output("git", "rev-parse", "--short", "HEAD"); err == nil {
		commit = strings.TrimSpace(out)
	}
	date := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf("-X main.version=%s -X main.commit=%s -X main.date=%s", version, commit, date)
}

// Run executes the CLI application.
func Run() error {
	return sh.RunV(filepath.Join(".", distDir, cliBinaryName))
}

// RunWeb starts the web server after compiling templates.
func RunWeb() error {
	mg.Deps(Generate)
	return sh.RunV(filepath.Join(".", distDir, webBinaryName))
}

// Generate compiles templ files into Go code.
func Generate() error {
	return sh.RunV("templ", "generate")
}

// Test runs unit tests with race detector and coverage analysis.
func Test() error {
	if os.Getenv("SHORT") == "1" {
		return sh.RunV("go", "test", "-short", "./...")
	}
	return sh.RunV("go", "test", "-race", "-coverprofile=coverage.out", "./...")
}

// Check runs all static checks: lint, format, and test.
func Check() error {
	mg.Deps(Lint, Format, Test)
	return nil
}

// Security performs security scans using gosec.
func Security() error {
	return sh.RunV("gosec", "-severity", "medium", "./...")
}

// Setup installs dependencies and development tools.
func Setup() {
	mg.Deps(Deps, Tidy, Install)
}

// Tidy runs go mod tidy to clean up go.mod and go.sum.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}
