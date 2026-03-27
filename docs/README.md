# gmcc 文档目录

## 概述

gmcc 是一个 Go 语言实现的 Minecraft Java 版控制台客户端，支持协议版本 774 (1.21.11)。

## 文档列表

| 文档 | 说明 |
|------|------|
| [auth.md](auth.md) | 认证系统（微软/Xbox/Minecraft 认证流程） |
| [protocol.md](protocol.md) | 协议实现（包定义、编解码、加密流程） |
| [tui.md](tui.md) | TUI 框架（终端界面架构与组件） |
| [player.md](player.md) | 玩家数据（状态、位置、背包系统） |
| [nbt_format.md](nbt_format.md) | NBT 二进制格式规范 |
| [snbt_format.md](snbt_format.md) | SNBT 字符串格式规范 |
| [text_component.md](text_component.md) | 文本组件 JSON 格式 |

### 超级能力文档 (superpowers/)

| 文档 | 说明 |
|------|------|
| [command-module-spec.md](superpowers/specs/command-module-spec.md) | 命令模块规格说明 |
| [manager_spec.md](superpowers/specs/manager_spec.md) | 管理器规格说明 |
| [entity-tracking-design.md](superpowers/specs/2026-03-26-entity-tracking-design.md) | 实体跟踪系统设计 |
| [config-auto-update-design.md](superpowers/specs/2026-03-28-config-auto-update-design.md) | 配置自动更新设计 |
| [command-module-plan.md](superpowers/plans/command-module-plan.md) | 命令模块实现计划 |

## 格式规范

| 文档 | 说明 |
|------|------|
| [nbt_format.md](nbt_format.md) | NBT 二进制格式规范 |
| [snbt_format.md](snbt_format.md) | SNBT 字符串格式规范 |
| [text_component.md](text_component.md) | 文本组件 JSON 格式 |

## 代码位置

| 模块 | 路径 | 说明 |
|------|------|------|
| 认证 | `internal/auth/` | 微软/Xbox/Minecraft 三级认证 |
| 命令系统 | `internal/commands/` | 路由、解析、状态管理 |
| 配置 | `internal/config/` | 配置加载、热重载 |
| 实体跟踪 | `internal/entity/` | 实体状态与位置跟踪 |
| 无头模式 | `internal/headless/` | 无界面运行器 |
| 国际化 | `internal/i18n/` | Minecraft 语言数据 |
| 物品系统 | `internal/item/` | 物品与组件解析 |
| 日志 | `internal/logx/` | 日志记录系统 |
| 客户端 | `internal/mcclient/` | Minecraft 客户端核心 |
| NBT | `internal/nbt/` | NBT 编解码 |
| 玩家 | `internal/player/` | 玩家状态管理 |
| 注册表 | `internal/registry/` | 物品 ID 注册表 |
| 会话 | `internal/session/` | Token 缓存 |
| TUI | `internal/tui/` | 终端界面 |

## 快速开始

### 安装

```bash
# 从源码编译
go build -o gmcc ./cmd/gmcc
```

### 配置

创建 `config.yaml`：

```yaml
account:
  player_id: "你的游戏ID"
  use_official_auth: true

server:
  address: "mc.example.com:25565"

actions:
  delay_ms: 1200
  on_join_commands:
    - "list"
  on_join_messages:
    - "大家好"
  sign_commands: false        # 入服命令是否签名
  default_sign_commands: true # 默认命令签名行为

commands:
  enabled: false
  prefix: "!"
  allow_all: false
  whitelist: []

log:
  log_dir: "logs"
  max_size: 512
  debug: false
  enable_file: true

runtime:
  headless: false

packets:
  handle_container: true
```

### 运行

```bash
./gmcc
```

## 核心功能

### 认证与连接
- 微软正版认证（设备码登录/刷新令牌）
- 离线模式支持
- 登录加密（AES/CFB8）
- 登录压缩
- 配置状态处理

### 游戏交互
- 游戏状态心跳
- 聊天消息收发（支持 CESU-8 编码）
- 命令发送（支持签名/无签名模式）
- 玩家状态获取（生命值、饥饿值、经验）
- 玩家位置跟踪与传送
- 背包系统与容器管理

### 数据处理
- NBT 数据解析与路径查询
- SNBT 字符串格式规范
- 文本组件 JSON 解析与 ANSI 转换
- 物品组件解析系统（104 种数据组件）
- 物品注册表（Minecraft ID 映射）

### 高级功能
- 实体跟踪系统（玩家、怪物位置）
- 国际化系统（Minecraft 语言数据）
- 无头模式运行器（自动化脚本）
- 命令系统框架（可扩展命令路由）
- 配置热重载（自动更新配置）

## 相关资源

- [Minecraft Wiki - 协议](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE)
- [node-minecraft-protocol](https://github.com/PrismarineJS/node-minecraft-protocol)
- [minecraft-data](https://github.com/PrismarineJS/minecraft-data)
