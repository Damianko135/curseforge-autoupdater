## Potential Issues in the Codebase (as of July 2025)


### 1. **Error Handling**
- Error wrapping with `fmt.Errorf("context: %w", err)` is now consistently applied throughout the codebase.
- Most functions now provide contextual error messages, improving traceability and debugging.
- User-facing CLI commands print errors with context to `os.Stderr`.
- **Remaining improvement:** Review all error messages for sufficient context (e.g., include file paths, operation names) and ensure errors are logged (not just returned) in background operations or goroutines.

### 2. **Configuration Validation**
- The config validation is present, but some fields (like notification URLs, webhook methods, etc.) could use stricter validation (e.g., URL format, allowed HTTP methods).
- Default values are set, but missing or invalid config files may not always be caught early.

### 3. **Atomic File Operations**
- `SafeWriteFile` uses temp files, but if `os.Chmod` fails, the temp file may not be cleaned up. Ensure all error paths clean up temp files.
- File and directory existence checks (`FileExists`, `DirExists`) do not distinguish between permission errors and non-existence.

### 4. **Testing Coverage**
- There are some unit tests (e.g., version compare), but many core modules (API, server, notification, CLI) lack comprehensive tests.
- No integration or end-to-end tests for update/backup/restore workflows.

### 5. **Logging and Observability**
- Logging is not consistently used across all modules. Some errors are returned but not logged, making debugging harder.
- No structured logging or log levels (info, warn, error) are evident.

### 6. **Security Practices**
- API keys and sensitive config are handled via config/env, but there is no mention of secure storage or redaction in logs.
- No input sanitization for user-supplied paths or config values (potential path traversal or injection risk).
- No rate limiting or retry logic for API/network operations.

### 7. **Concurrency and Race Conditions**
- No explicit locking or concurrency control for file operations, backups, or server state changes. Potential for race conditions if run in parallel.

### 8. **Worst Practices / Code Smells**
- Some functions are very large or do too much (e.g., backup/restore logic could be further modularized).
- Some repeated logic (e.g., error aggregation, file checks) could be refactored into helpers.
- Magic strings and numbers (e.g., HTTP methods, file extensions) are scattered—should be constants/enums.

### 9. **Documentation and Comments**
- Some exported functions lack doc comments (Go best practice).
- No clear documentation for environment variables, config file precedence, or error codes.

### 10. **Feature Gaps / TODOs**
- No scheduling/automation yet (planned, but not implemented).
- No rollback/transactional update support.
- No manifest-based mod version tracking (planned, not implemented).
- No multi-server or multi-modpack support yet.

---

## Steps to Fix and Improve

1. **Improve Error Handling**
   - Use `fmt.Errorf("context: %w", err)` everywhere an error is returned.
   - Ensure all errors are logged with context before returning or exiting.
   - Add user-facing error messages for CLI commands.

2. **Strengthen Configuration Validation**
   - Add stricter validation for URLs, HTTP methods, and required fields.
   - Fail fast on invalid or missing config.

3. **Enhance File Operation Safety**
   - Ensure all temp files are cleaned up on error in `SafeWriteFile` and similar functions.
   - Distinguish between permission errors and non-existence in file/dir checks.

4. **Increase Test Coverage**
   - Add unit tests for all major modules (API, server, notification, CLI).
   - Add integration tests for update, backup, and restore workflows.

5. **Add Structured Logging**
   - Use a logging library with levels (info, warn, error).
   - Add logs for all major operations and errors.

6. **Harden Security**
   - Redact API keys and sensitive data in logs.
   - Validate and sanitize all user/config input.
   - Add retry logic and rate limiting for network operations.

7. **Address Concurrency**
   - Add mutexes or other concurrency controls for file and server operations.
   - Document thread-safety assumptions.

8. **Refactor for Maintainability**
   - Break up large functions into smaller helpers.
   - Move magic strings/numbers to constants.
   - Add doc comments to all exported functions.

9. **Fill Feature Gaps**
   - Implement planned features: scheduling, rollback, manifest tracking, multi-server support.
   - Add CLI commands for all planned operations.

10. **Improve Documentation**
    - Document all config options, environment variables, and error codes.
    - Add usage examples and troubleshooting guides.

---

**Note:**
This list is not exhaustive. For a full audit, run:
- `golangci-lint run` (linting)
- `gosec ./...` (security)
- `go test -v ./...` (unit tests)
and review the output for actionable issues.
