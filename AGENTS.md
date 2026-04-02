# AGENTS.md - Coding Guidelines for gmcc

## Project Overview

A Go-based Minecraft client application (`gmcc`) providing headless/automated gameplay capabilities.

**Module**: `gmcc`  
**Go Version**: 1.25.1

## Build Commands

```bash
# Build the main binary
go build -o gmcc.exe ./cmd/gmcc

# Build with version info (used in releases)
go build -ldflags="-s -w -X main.Version=v1.0.0" -o gmcc.exe ./cmd/gmcc

# Build for production (stripped)
go build -ldflags="-s -w" -o gmcc.exe ./cmd/gmcc

# Cross-compile examples:
# Linux:   GOOS=linux GOARCH=amd64 go build -o gmcc-linux ./cmd/gmcc
# macOS:   GOOS=darwin GOARCH=amd64 go build -o gmcc-darwin ./cmd/gmcc
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./pkg/binutil

# Run a single test function
go test -v ./pkg/binutil -run TestReader_ReadVarInt

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Lint Commands

```bash
# Format code (must pass before committing)
go fmt ./...

# Vet code for issues
go vet ./...

# Run goimports (organize imports)
goimports -w .

# Recommended: Use golangci-lint for comprehensive checks
golangci-lint run
```

## Code Style Guidelines

### Imports
- Group imports: stdlib first, then external packages, then project packages
- Use `goimports` to automatically organize imports
- Avoid unused imports

```go
import (
    "bytes"
    "encoding/binary"
    "fmt"

    "golang.org/x/term"
    "gopkg.in/yaml.v3"

    "gmcc/internal/constants"
    "gmcc/pkg/binutil"
)
```

### Formatting
- Indent with **tabs** (size 4)
- Line endings: **CRLF** (Windows)
- No trailing whitespace trimming
- Follow standard Go formatting (`go fmt`)

### Naming Conventions
- **Types**: PascalCase (e.g., `AccountConfig`, `VarInt`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase (e.g., `playerID`, `useOfficialAuth`)
- **Constants**: PascalCase or ALL_CAPS for exported constants
- **Interfaces**: Single-method interfaces use `-er` suffix (e.g., `Reader`)
- **Test files**: `*_test.go`
- **Test functions**: `Test<Type>_<Method>` or `Test<Function>`

### Types
- Use explicit types when clarity is needed
- Define custom types for domain concepts (e.g., `type VarInt int32`)
- Use struct tags for YAML/JSON marshaling
- Keep struct fields exported if they need to be accessed outside the package

### Error Handling
- Always check errors and handle them explicitly
- Wrap errors with context using `fmt.Errorf("...: %w", err)`
- Return errors rather than logging and continuing
- Use sentinel errors for known error conditions
- Test error cases in unit tests

```go
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}
```

### Comments
- Use Chinese comments for business logic (existing convention)
- Use GoDoc format for exported items
- Keep comments concise and meaningful

```go
// ReadVarInt 读取 VarInt 编码的整数
func (r *Reader) ReadVarInt() (int32, error) {
    // 实现...
}
```

### Testing
- Use table-driven tests
- Name test cases descriptively
- Use `t.Run()` for subtests
- Test both success and error cases
- Skip incomplete tests with `t.Skip("reason")`

```go
func TestReader_ReadVarInt(t *testing.T) {
    tests := []struct {
        name     string
        data     []byte
        expected int32
        wantErr  bool
    }{
        {"zero", []byte{0x00}, 0, false},
        {"empty", []byte{}, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Project Structure

```
gmcc/
├── cmd/gmcc/           # Main application entry point
│   └── main.go
├── internal/           # Private application code
│   ├── auth/          # Authentication (Microsoft, Minecraft)
│   ├── config/        # Configuration management
│   ├── constants/     # Application constants
│   ├── headless/      # Headless mode runner
│   ├── logx/          # Logging utilities
│   ├── mcclient/      # Minecraft client implementation
│   └── session/       # Session management
├── pkg/               # Public libraries
│   ├── binutil/       # Binary data utilities (VarInt, etc.)
│   └── httpx/         # HTTP client utilities
└── docs/              # Documentation
```

### Package Guidelines
- `internal/`: Private code, cannot be imported by external packages
- `pkg/`: Reusable public libraries
- Keep packages focused on a single responsibility
- Minimize dependencies between packages

### Configuration
- Use `config.yaml` for runtime configuration
- Support environment variables (e.g., `GMCC_CONFIG`, `GMCC_DISABLE_AUTO_UPDATE`)
- Provide sensible defaults in code

### Logging
- Use `internal/logx` package for all logging
- Use appropriate log levels: `Debug`, `Info`, `Warn`, `Error`
- Include context in log messages

## Pre-commit Checklist

Before committing code:

1. Run `go fmt ./...` - Ensure code is formatted
2. Run `go vet ./...` - No issues detected
3. Run `go test ./...` - All tests pass
4. Verify build succeeds: `go build ./cmd/gmcc`
