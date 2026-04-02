# gmcc 项目参考文档

## 项目概览

**gmcc** 是一个基于 Go 开发的 Minecraft 客户端应用程序，提供无界面/自动化游戏能力。

- **模块名称**: `gmcc`
- **Go 版本**: 1.25.1
- **主要功能**: 连接到 Minecraft 服务器并执行自动化操作

---

## 目录结构

```
gmcc/
├── cmd/gmcc/              # 主程序入口点
│   └── main.go           # 应用程序启动入口
│
├── internal/              # 私有应用代码（不可被外部导入）
│   ├── auth/             # 认证模块（Microsoft/Minecraft）
│   ├── config/           # 配置管理
│   ├── constants/        # 应用常量定义
│   ├── headless/         # 无界面运行器
│   ├── logx/             # 日志工具
│   ├── mcclient/         # Minecraft 客户端实现
│   │   ├── crypto/       # 加密相关
│   │   ├── packet/       # 数据包处理
│   │   └── protocol/     # 协议实现
│   └── session/          # 会话管理
│
├── pkg/                   # 公共库（可被外部导入）
│   ├── binutil/          # 二进制数据工具
│   │   ├── reader.go     # 二进制读取器
│   │   ├── writer.go     # 二进制写入器
│   │   └── types.go      # 类型定义（VarInt 等）
│   └── httpx/            # HTTP 客户端工具
│
├── docs/                  # 文档
├── config.yaml           # 默认配置文件
└── go.mod                # Go 模块定义
```

### 关键目录说明

| 目录 | 用途 | 访问权限 |
|------|------|----------|
| `cmd/` | 包含可执行程序入口 | - |
| `internal/` | 私有业务逻辑代码 | 仅限本项目 |
| `pkg/` | 可复用的公共库 | 可外部导入 |
| `docs/` | 项目文档 | - |

---

## 配置文件参考

配置文件采用 YAML 格式，默认文件名为 `config.yaml`。

### 完整配置结构

```yaml
account:
    player_id: "your_player_id_here"        # 玩家 ID
    use_official_auth: false                # 是否使用正版认证

server:
    address: "127.0.0.1:25565"             # 服务器地址（主机:端口）

log:
    log_dir: "logs"                         # 日志目录
    max_size: 512                           # 单个日志文件最大大小（KB）
    debug: false                            # 是否启用调试日志
    enable_file: true                       # 是否启用文件日志
```

### 配置字段说明

#### account

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `player_id` | string | "your_player_id_here" | 玩家显示名称，最大 16 字符 |
| `use_official_auth` | bool | false | 是否使用 Microsoft 正版认证 |

#### server

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `address` | string | "127.0.0.1:25565" | 目标服务器地址，格式为 `host:port` |

#### log

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `log_dir` | string | "logs" | 日志文件存储目录，自动创建 |
| `max_size` | int64 | 512 | 日志文件大小上限（KB） |
| `debug` | bool | false | 启用后输出 DEBUG 级别日志 |
| `enable_file` | bool | true | 是否将日志写入文件 |

---

## 环境变量

所有环境变量均以 `GMCC_` 为前缀。

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `GMCC_CONFIG` | string | "config.yaml" | 指定配置文件路径 |
| `GMCC_DISABLE_AUTO_UPDATE` | bool | false | 设置为 "true" 禁用配置热更新 |

### 使用示例

```bash
# 使用自定义配置文件
set GMCC_CONFIG=/path/to/custom-config.yaml

# 禁用配置热更新
set GMCC_DISABLE_AUTO_UPDATE=true
```

---

## 构建参考

### 基础构建

```bash
# 构建主程序
go build -o gmcc.exe ./cmd/gmcc

# 构建带版本信息（推荐）
go build -ldflags="-s -w -X main.Version=v1.0.0" -o gmcc.exe ./cmd/gmcc

# 生产环境构建（去除符号表和调试信息）
go build -ldflags="-s -w" -o gmcc.exe ./cmd/gmcc
```

### 交叉编译

| 目标平台 | 命令 |
|----------|------|
| Linux (amd64) | `GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gmcc-linux ./cmd/gmcc` |
| Linux (arm64) | `GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gmcc-linux-arm64 ./cmd/gmcc` |
| Windows (amd64) | `GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o gmcc.exe ./cmd/gmcc` |
| Windows (arm64) | `GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o gmcc-arm64.exe ./cmd/gmcc` |
| macOS (amd64) | `GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o gmcc-darwin ./cmd/gmcc` |
| macOS (arm64) | `GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o gmcc-darwin-arm64 ./cmd/gmcc` |

### 构建参数说明

| 参数 | 说明 |
|------|------|
| `-s` | 去除符号表 |
| `-w` | 去除 DWARF 调试信息 |
| `-X main.Version=xxx` | 注入版本号到 main.Version 变量 |

---

## 运行指南

### 基础运行

```bash
# 使用默认配置文件 (config.yaml)
./gmcc.exe

# 使用环境变量指定配置
set GMCC_CONFIG=production.yaml
./gmcc.exe
```

### 程序输出

启动成功时，控制台将显示：

```
[INFO] gmcc 启动 (无界面模式)
[INFO] 玩家: <player_id>
[INFO] 服务器: <server_address>
[INFO] 正在连接...
```

### 退出方式

- **正常退出**: 按下 `Ctrl+C` 发送 SIGINT 信号
- 程序会优雅地断开连接并清理资源

---

## 测试参考

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/binutil

# 运行特定测试函数
go test -v ./pkg/binutil -run TestReader_ReadVarInt

# 运行测试并显示覆盖率
go test -cover ./...

# 运行基准测试
go test -bench=. ./...
```

### 测试标志

| 标志 | 说明 |
|------|------|
| `-v` | 详细输出模式 |
| `-run <pattern>` | 仅运行匹配名称的测试 |
| `-cover` | 显示覆盖率统计 |
| `-bench=<pattern>` | 运行匹配的基准测试 |

---

## 核心包 API 速览

### binutil 包

提供 Minecraft 协议所需的二进制数据读写功能。

**类型:**
- `VarInt` - 可变长度 32 位整数
- `VarLong` - 可变长度 64 位整数
- `Position` - 三维坐标 (X, Y, Z)
- `UUID` - 16 字节 UUID

**Reader 方法:**
- `NewReader(data []byte) *Reader` - 从字节数组创建读取器
- `ReadVarInt() (int32, error)` - 读取 VarInt
- `ReadString() (string, error)` - 读取带长度前缀的字符串
- `ReadBool() (bool, error)` - 读取布尔值
- `ReadInt32() (int32, error)` - 读取大端序 int32
- `ReadPosition() (x, y, z int32, err error)` - 读取位置坐标

**Writer 方法:**
- `NewWriter() *Writer` - 创建写入器
- `WriteVarInt(v int32) error` - 写入 VarInt
- `WriteString(s string) error` - 写入带长度前缀的字符串
- `Bytes() []byte` - 获取写入的字节

### config 包

配置管理功能。

**函数:**
- `LoadWithAutoUpdate(path string, autoUpdate bool) (*Config, error)` - 加载配置，支持热更新

**类型:**
- `Config` - 完整配置结构
- `AccountConfig` - 账户配置
- `ServerConfig` - 服务器配置
- `LogConfig` - 日志配置

### httpx 包

HTTP 客户端工具。

**函数:**
- `PostForm(rawURL string, form url.Values, ptr interface{}) (*HTTPResponse, error)` - 发送表单 POST
- `PostJSON(rawURL string, reqBody interface{}, ptr interface{}) (*HTTPResponse, error)` - 发送 JSON POST

---

## 版本信息

- **当前版本**: 从构建参数注入 (`main.Version`)
- **Go 版本要求**: 1.25.1+

---

## 相关资源

- 项目仓库: `<repository-url>`
- 问题反馈: GitHub Issues
- 构建脚本参考: `.github/workflows/release.yml`
