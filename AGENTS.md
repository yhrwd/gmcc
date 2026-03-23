# AGENTS.md - gmcc Agent Guidelines

This file provides guidelines for agentic coding agents working on the gmcc project.

## Project Overview

gmcc is a Minecraft Java Edition console client supporting protocol version 774 (1.21.11). Written in Go 1.25.1.

## Build Commands

```bash
# Build the binary
go build -o gmcc ./cmd/gmcc

# Build with version info
go build -ldflags "-s -w -X main.Version=1.0.0" -o gmcc ./cmd/gmcc

# Run the application
./gmcc
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test by name
go test -run TestNBTDecoder_String ./internal/mcclient/...

# Run tests in a specific package
go test -v ./internal/nbt/...

# Run tests with coverage
go test -cover ./...

# Run tests matching a pattern
go test -run 'TestNBT|TestMarshal' ./internal/nbt/...
```

## Linting

```bash
# Run go vet
go vet ./...

# Run golangci-lint if installed
golangci-lint run
```

## Code Style

### Formatting

- Use 4 spaces for indentation (configured in .editorconfig)
- Use CRLF line endings (Windows)
- Keep lines under 100 characters when practical

### Imports

Organize imports in three groups with blank lines between:
1. Standard library
2. External packages
3. Internal packages

```go
import (
    "bytes"
    "encoding/json"
    "fmt"

    "gmcc/internal/mcclient/packet"
    "gmcc/internal/nbt"
)
```

### Naming Conventions

- **Types**: PascalCase (e.g., `Client`, `Config`, `Decoder`)
- **Functions/Variables**: camelCase (e.g., `newDecoder`, `parseSNBT`)
- **Constants**: PascalCase or CamelCase with prefix (e.g., `MaxPacketSize`, `sessionDir`)
- **Packages**: lowercase, short names (e.g., `nbt`, `logx`, `mcclient`)
- **Files**: lowercase with underscores for multi-word names (e.g., `chat_parser.go`)

### Types

- Use `int`/`int64` for general integers, explicit widths (`int32`, `int16`) for protocol fields
- Use `any` instead of `interface{}`
- Use custom types for protocol constants (e.g., `type PacketID int`)

### Error Handling

- Return errors with context using `fmt.Errorf("description: %w", err)`
- Use sentinel errors for known failure cases
- Check errors explicitly, don't ignore with `_`

```go
// Good
if err := dec.Decode(&result); err != nil {
    return fmt.Errorf("decode failed: %w", err)
}

// Bad
data, _ := os.ReadFile(path)
```

### Comments

- Use Chinese comments for user-facing documentation (config fields, public APIs)
- Use English for internal implementation comments
- Comment public APIs; internal functions can be self-explanatory

### Testing

- Test files named `*_test.go` in the same package
- Use table-driven tests with `t.Run`:

```go
tests := []struct {
    name  string
    value any
}{
    {"byte", int8(127)},
    {"short", int16(32767)},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test code
    })
}
```

- Use `t.Fatalf` for setup errors, `t.Errorf` for assertion failures

### Project Structure

```
cmd/gmcc/          # Entry point
internal/          # Core modules (not exported)
  auth/            # Authentication (microsoft, minecraft)
  config/          # Configuration loading
  logx/            # Logging
  mcclient/        # Minecraft client core
  nbt/             # NBT data handling
  player/          # Player state
  session/         # Token caching
  tui/             # Terminal UI
pkg/               # Public utilities
docs/              # Documentation
```

### Key Conventions

- Protocol constants in `internal/mcclient/protocol_774.go`
- Use NBT path queries via `nbt.QueryPath()` for player data
- Chat message handling via `client.SetChatHandler()`
- Token caching in `.session/` directory
