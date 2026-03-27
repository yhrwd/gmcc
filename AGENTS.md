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
cmd/gmcc/          # 程序入口
internal/          # 核心模块（不导出）
  auth/            # 认证 (microsoft, minecraft)
  commands/        # 命令系统 (core, handlers, tracker, auth)
  components/      # 数据组件解析框架
  config/          # 配置加载
  constants/       # 常量定义
  entity/          # 实体跟踪系统
  headless/        # 无头模式运行器
  i18n/            # 国际化 (Minecraft 语言数据)
  item/            # 物品系统
    component/     # 物品组件解析器
  logx/            # 日志系统
  mcclient/        # Minecraft 客户端核心
    chat/          # 聊天消息处理
    crypto/        # 加密/解密
    handlers/      # 数据包处理器
    packet/        # 数据包定义
    protocol/      # 协议定义
  nbt/             # NBT 数据处理
  player/          # 玩家状态
  registry/        # 物品注册表 (Minecraft ID -> 物品信息)
  session/         # Token 缓存
  tui/             # 终端 UI
pkg/               # 公共工具
  binutil/         # 二进制工具
  httpx/           # HTTP 工具
docs/              # 文档
```

### Key Conventions

- Protocol constants in `internal/mcclient/protocol/` directory
- Use NBT path queries via `nbt.QueryPath()` for player data
- Chat message handling via `client.SetChatHandler()`
- Entity tracking via `entity.Tracker`
- Item registry via `registry.GetItemRegistry()`
- Internationalization via `i18n.GetI18n()`
- Token caching in `.session/` directory
- Item component parsing via `item/component` package
- Command system in `internal/commands/` package
