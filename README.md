# gmcc

一个基于 Go 开发的 Minecraft 客户端应用程序，提供无界面/自动化游戏能力。

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 功能特性

- **无界面模式** - 无需图形界面即可运行
- **自动认证** - 支持 Microsoft 正版认证
- **配置热更新** - 支持运行时自动重载配置
- **多平台** - 支持 Windows、Linux、macOS
- **结构化日志** - 支持文件和控制台双输出

---

## 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/yourusername/gmcc.git
cd gmcc
```

### 2. 配置环境

确保已安装 Go 1.25.1 或更高版本：

```bash
go version
```

### 3. 配置账户

编辑 `config.yaml` 文件：

```yaml
account:
    player_id: "YourPlayerID"          # 你的玩家 ID
    use_official_auth: false            # 是否使用正版认证

server:
    address: "your.server.com:25565"  # 目标服务器地址

log:
    log_dir: "logs"
    max_size: 512
    debug: false
    enable_file: true
```

### 4. 构建运行

```bash
# 构建
go build -o gmcc.exe ./cmd/gmcc

# 运行
./gmcc.exe
```

---

## 项目结构

```
gmcc/
├── cmd/gmcc/          # 主程序入口
├── internal/          # 内部业务逻辑
│   ├── auth/         # 认证模块
│   ├── config/       # 配置管理
│   ├── headless/     # 无界面运行器
│   ├── logx/         # 日志工具
│   ├── mcclient/     # Minecraft 客户端
│   └── session/      # 会话管理
├── pkg/              # 公共库
│   ├── binutil/      # 二进制数据工具
│   └── httpx/        # HTTP 客户端
├── docs/             # 文档
└── config.yaml       # 默认配置文件
```

---

## 使用方法

### 基本运行

```bash
# 使用默认配置文件
./gmcc.exe

# 指定配置文件
set GMCC_CONFIG=production.yaml
./gmcc.exe
```

### 构建选项

```bash
# 开发构建
go build -o gmcc.exe ./cmd/gmcc

# 生产构建（优化体积）
go build -ldflags="-s -w" -o gmcc.exe ./cmd/gmcc

# 带版本信息
go build -ldflags="-s -w -X main.Version=v1.0.0" -o gmcc.exe ./cmd/gmcc
```

### 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o gmcc-linux ./cmd/gmcc

# macOS
GOOS=darwin GOARCH=amd64 go build -o gmcc-darwin ./cmd/gmcc

# Windows
GOOS=windows GOARCH=amd64 go build -o gmcc.exe ./cmd/gmcc
```

---

## 配置说明

### 配置文件 (config.yaml)

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `account.player_id` | 玩家显示名称 | "your_player_id_here" |
| `account.use_official_auth` | 使用正版认证 | false |
| `server.address` | 服务器地址 | "127.0.0.1:25565" |
| `log.log_dir` | 日志目录 | "logs" |
| `log.max_size` | 日志大小上限 (KB) | 512 |
| `log.debug` | 启用调试日志 | false |
| `log.enable_file` | 启用文件日志 | true |

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `GMCC_CONFIG` | 配置文件路径 | "config.yaml" |
| `GMCC_DISABLE_AUTO_UPDATE` | 禁用配置热更新 | false |

---

## 开发指南

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -v ./pkg/binutil -run TestReader_ReadVarInt

# 覆盖率测试
go test -cover ./...
```

### 代码规范

```bash
# 格式化代码
go fmt ./...

# 检查代码
go vet ./...

# 整理导入
goimports -w .
```

---

## 文档

- [项目参考文档](docs/reference.md) - 详细配置和 API 参考
- [项目架构](docs/architecture.md) - 系统架构说明
- [API 文档](docs/api.md) - 包 API 文档

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

请确保：
- 代码通过 `go fmt` 和 `go vet` 检查
- 所有测试通过
- 提交信息清晰明了

---

## 许可证

本项目采用 [MIT](LICENSE) 许可证。

---

## 感谢

- 感谢所有贡献者
- 使用 [PrismarineJS](https://github.com/PrismarineJS) 作为参考
- 感谢 Go 社区提供的优秀工具

---

<div align="center">

**Made with ❤️ by the gmcc team**

</div>
