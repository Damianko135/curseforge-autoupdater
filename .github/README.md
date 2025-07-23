# GitHub Configuration

This directory contains GitHub-specific configuration files for the CurseForge Auto-Updater project.

## Directory Structure

```
.github/
├── ISSUE_TEMPLATE/          # Issue templates
│   ├── bug_report.yml       # Bug report template
│   └── feature_request.yml  # Feature request template
├── workflow/                # GitHub Actions workflows (renamed to prevent auto-execution)
│   ├── code-quality.yml     # Code quality checks
│   ├── dependency-update.yml # Automated dependency updates
│   ├── go-ci.yml           # Go CI/CD pipeline
│   ├── pr-labeller.yml     # Automatic PR labeling
│   ├── python-ci.yml       # Python CI/CD pipeline
│   ├── release.yml         # Release automation
│   ├── security.yml        # Security scanning
│   └── setup-labels.yml    # Label management
├── CODEOWNERS              # Code ownership and review assignments
├── labeller.yml            # PR labeling rules
├── labels.yml              # Repository label definitions
├── markdown-link-check.json # Markdown link validation config
├── pull_request_template.md # PR template
├── renovate.json           # Renovate dependency update config
└── README.md               # This file
```

## Workflows Overview

### CI/CD Workflows

1. **go-ci.yml** - Go Continuous Integration
   - Runs tests on Go 1.21 and 1.22
   - Performs static analysis with staticcheck
   - Builds CLI and Web binaries
   - Uploads coverage to Codecov

2. **python-ci.yml** - Python Continuous Integration
   - Tests on Python 3.9-3.12
   - Linting with flake8, black, isort
   - Type checking with mypy
   - Package building and installation testing

3. **code-quality.yml** - Code Quality Checks
   - golangci-lint for Go code
   - Multiple Python linters (flake8, pylint, mypy)
   - Markdown linting and link checking
   - Conventional commit validation

4. **security.yml** - Security Scanning
   - CodeQL analysis for Go and Python
   - Dependency vulnerability scanning with Trivy
   - Go security scanning with Gosec
   - Python security scanning with Safety and Bandit

### Automation Workflows

5. **pr-labeller.yml** - Automatic PR Labeling
   - Labels PRs based on changed files
   - Uses rules defined in `labeller.yml`

6. **setup-labels.yml** - Label Management
   - Creates/updates repository labels
   - Uses definitions from `labels.yml`
   - Runs on label file changes

7. **dependency-update.yml** - Dependency Updates
   - Weekly automated dependency updates
   - Separate PRs for Go, Python, and GitHub Actions
   - Uses Renovate for GitHub Actions updates

8. **release.yml** - Release Automation
   - Triggered on version tags (v*)
   - Builds binaries for multiple platforms
   - Creates GitHub releases with assets
   - Supports both Go binaries and Python packages

## Labels

The repository uses a comprehensive labeling system defined in `labels.yml`:

- **Priority**: `priority:high`, `priority:medium`, `priority:low`
- **Status**: `status:in-progress`, `status:blocked`
- **Area**: `area:golang`, `area:python`, `area:docs`
- **Type**: `type:ci`, `type:test`, `type:config`
- **Component**: `component:cli`, `component:web`, `component:api`
- **CurseForge Specific**: `curseforge:api`, `curseforge:download`, `curseforge:update`

## Issue Templates

Two structured issue templates are provided:

1. **Bug Report** (`bug_report.yml`)
   - Component selection
   - Reproduction steps
   - Environment details
   - Log output

2. **Feature Request** (`feature_request.yml`)
   - Component affected
   - Problem description
   - Proposed solution
   - Use case details

## Pull Request Template

The PR template (`pull_request_template.md`) includes:
- Change type classification
- Component checklist
- Testing requirements
- Review checklist
- Related issue linking

## Configuration Files

- **CODEOWNERS**: Defines code ownership and automatic review assignments
- **labeller.yml**: Rules for automatic PR labeling based on file changes
- **renovate.json**: Configuration for automated dependency updates
- **markdown-link-check.json**: Settings for markdown link validation

## Activating Workflows

Currently, workflows are stored in the `workflow/` directory to prevent automatic execution during setup. To activate them:

1. Rename the `workflow/` directory to `workflows/`:
   ```bash
   mv .github/workflow .github/workflows
   ```

2. Update the CODEOWNERS file to replace `@your-username` with your actual GitHub username

3. Commit and push the changes

## Customization

Before activating, consider customizing:

1. **CODEOWNERS**: Replace `@your-username` with actual GitHub usernames
2. **Go versions**: Update Go version matrix in `go-ci.yml` if needed
3. **Python versions**: Adjust Python version matrix in `python-ci.yml`
4. **Branch names**: Update branch names if you use different naming conventions
5. **Renovate config**: Adjust update schedules and grouping in `renovate.json`

## Security Considerations

- Workflows use `GITHUB_TOKEN` for most operations
- CodeQL and security scanning results are uploaded to GitHub Security tab
- Dependency vulnerability alerts are enabled
- All workflows follow security best practices with minimal permissions

## Monitoring

After activation, monitor:
- Workflow run status in the Actions tab
- Security alerts in the Security tab
- Dependency update PRs from Renovate
- Label application on new PRs and issues

## Troubleshooting

Common issues and solutions:

1. **Workflow not triggering**: Check file paths in trigger conditions
2. **Permission errors**: Verify `GITHUB_TOKEN` permissions in workflow files
3. **Build failures**: Check Go/Python versions and dependencies
4. **Label conflicts**: Review `labels.yml` for duplicate or conflicting labels