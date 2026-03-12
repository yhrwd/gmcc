# gmcc

Go 语言实现的 Minecraft Java 版控制台客户端，支持协议版本 774 (1.21.11)。

## 功能特性

- ✅ 微软正版认证（设备码登录 / 刷新令牌）
- ✅ 离线模式支持
- ✅ 登录加密（AES/CFB8）
- ✅ 登录压缩
- ✅ 配置状态处理
- ✅ 游戏状态心跳（保持在线）
- ✅ 聊天消息接收/发送
- ✅ 命令发送
- ✅ CESU-8 编码支持（正确显示中文和 Emoji）
- ✅ NBT 数据解析与路径查询

## 快速开始

### 下载

从 [Releases](https://github.com/yhrwd/gmcc/releases) 页面下载对应平台的二进制文件。

### 运行

```bash
# 首次运行会生成默认配置文件
./gmcc

# 或指定配置文件路径
GMCC_CONFIG=/path/to/config.yaml ./gmcc
```

### 配置

编辑 `config.yaml`：

```yaml
account:
  player_id: "你的游戏ID"
  use_official_auth: true   # true: 正版认证, false: 离线模式

server:
  address: "mc.example.com:25565"

actions:
  delay_ms: 1200            # 入服后动作延迟（毫秒）
  on_join_commands:         # 入服后自动执行的命令
    - "list"
  on_join_messages:         # 入服后自动发送的消息
    - "大家好"

log:
  log_dir: "logs"           # 日志目录
  max_size: 10              # 单个日志文件最大大小（MB）
  debug: false              # 调试模式
  enable_file: true         # 启用文件日志
```

## 认证流程

### 正版认证（`use_official_auth: true`）

1. 首次运行时，程序会输出设备码登录链接
2. 在浏览器中打开链接，输入设备码完成微软登录
3. 认证成功后，令牌会缓存到 `.session/<玩家ID>.json`
4. 后续运行会自动刷新令牌，无需重复登录

### 离线模式（`use_official_auth: false`）

直接使用 `player_id` 作为游戏用户名连接服务器。

## 命令行用法

```bash
# 显示帮助
./gmcc -h

# 显示版本
./gmcc -version
```

## 项目结构

```
cmd/gmcc/
  main.go                    # 程序入口

internal/
  auth/
    microsoft/service.go     # 微软/Xbox/XSTS 认证
    minecraft/service.go     # Minecraft 登录验证
  config/
    config.go                # 配置结构定义
    loader.go                # YAML 加载校验
  logx/
    logx.go                  # 日志模块
  mcclient/
    client.go                # 客户端状态机
    codec.go                 # 协议编解码
    chat.go                  # 聊天/命令处理
    chat_parser.go           # 聊天 JSON 解析
    handlers_*.go            # 各状态包处理
    text_component.go        # 文本组件解析
    protocol_774.go          # 协议常量定义
    utils.go                 # 工具函数
  nbt/
    decode.go                # NBT 解码器
    encode.go                # NBT 编码器
    snbt.go                  # SNBT 解析器
    path.go                  # NBT 路径查询
    nbt.go                   # CESU-8 转换
  player/
    player.go                # 玩家状态管理
    inventory.go             # 背包系统
  session/
    cache.go                 # 会话令牌缓存
  tui/
    tui.go                   # TUI 界面

pkg/
  httpx/                     # HTTP 工具
```

## 编译

```bash
# 本地编译
go build -o gmcc ./cmd/gmcc

# 交叉编译
GOOS=windows GOARCH=amd64 go build -o gmcc.exe ./cmd/gmcc
GOOS=linux GOARCH=amd64 go build -o gmcc ./cmd/gmcc
GOOS=darwin GOARCH=arm64 go build -o gmcc ./cmd/gmcc
```

## 开发文档

详细开发说明请参阅 [docs/development.md](docs/development.md)。

## 许可证

MIT License