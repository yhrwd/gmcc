# gmcc

Go 语言实现的 Minecraft Java 版控制台客户端，支持协议版本 774 (1.21.11)。

## 功能特性

- 微软正版认证（设备码登录 / 刷新令牌）
- 离线模式支持
- 登录加密（AES/CFB8）
- 登录压缩
- 配置状态处理
- 游戏状态心跳（保持在线）
- 聊天消息接收/发送
- 命令发送
- CESU-8 编码支持（正确显示中文和 Emoji）
- NBT 数据解析与路径查询
- **物品组件解析系统** - 支持 104 种数据组件类型
- **实体跟踪系统** - 实时跟踪玩家和其他实体位置
- **国际化系统** - 支持 Minecraft 语言数据本地化
- **物品注册表** - Minecraft ID 到物品信息的映射
- **无头模式运行器** - 支持自动化脚本
- **命令系统** - 完整的命令框架，支持自定义命令
- TUI 终端用户界面
- 玩家状态和背包系统

## 快速开始

### 安装

```bash
# 从源码编译
go build -o gmcc ./cmd/gmcc

# 或从 Releases 下载
```

### 配置

创建 `config.yaml`：

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
  sign_commands: false      # 入服命令是否签名
  default_sign_commands: true  # 默认命令签名行为

commands:
  enabled: false            # 启用机器人命令系统
  prefix: "!"               # 命令前缀
  allow_all: false          # 允许所有人使用
  whitelist: []             # 允许使用的玩家列表

log:
  log_dir: "logs"
  max_size: 512             # 单个日志文件最大大小（KB）
  debug: false
  enable_file: true

runtime:
  headless: false           # 无界面模式

packets:
  handle_container: true    # 处理容器数据包
```

### 运行

```bash
./gmcc
```

## 项目结构

```
cmd/gmcc/          # 程序入口
internal/          # 核心模块（不导出）
  auth/            # 认证 (microsoft, minecraft)
  commands/        # 命令系统
    adapter/       # Bot 适配器
    core/          # 命令核心（路由、解析、状态）
    handlers/      # 命令处理器
    tracker/       # 命令状态跟踪
    auth/          # 权限管理
    modules/       # 命令模块
    parser/        # 消息解析器
  config/          # 配置加载、热重载、原子更新
  constants/       # 常量定义
  entity/          # 实体跟踪系统
  headless/        # 无头模式运行器
  i18n/            # 国际化 (Minecraft 语言数据)
  item/            # 物品系统
    component/     # 物品组件解析器
  logx/            # 日志系统
  mcclient/        # Minecraft 客户端核心
    chat/          # 聊天消息处理、文本组件解析
    crypto/        # 加密/解密 (CFB8)
    handlers/      # 数据包处理器
    packet/        # 数据包定义、编解码
    protocol/      # 协议常量
  nbt/             # NBT 数据处理（解码、编码、路径查询）
  player/          # 玩家状态（位置、背包、附近玩家）
  registry/        # 物品注册表 (Minecraft ID -> 物品信息)
  session/         # Token 缓存
  tui/             # 终端 UI
pkg/               # 公共工具
  binutil/         # 二进制工具（VarInt、读写器）
  httpx/           # HTTP 工具
docs/              # 文档
  formats/         # 数据格式参考 (NBT, SNBT, 文本组件)
  superpowers/     # 设计文档
    specs/         # 规格说明
    plans/          # 实现计划
      archive/      # 已完成的计划存档
```

## 文档

详细文档请查看 [docs/README.md](docs/README.md)

## 协议版本

当前支持：**协议 774** (Minecraft Java 1.21.11)

## 许可证

MIT License