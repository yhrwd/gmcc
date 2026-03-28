# gmcc 文档目录

本目录包含 gmcc 项目的完整技术文档，按照 [Diátaxis 框架](https://diataxis.fr/) 组织。

## 快速导航

| 文档类型 | 说明 |
|----------|------|
| **[入门教程](#入门教程)** | 面向新用户的逐步指南 |
| **[操作指南](#操作指南)** | 解决特定问题的步骤说明 |
| **[技术参考](#技术参考)** | 协议、API 和数据格式参考 |
| **[设计文档](#设计文档)** | 架构设计和技术决策 |

---

## 入门教程

面向新用户的入门指南，帮助快速上手 gmcc。

| 文档 | 说明 |
|------|------|
| [项目 README](../README.md) | 项目概述、安装和快速开始 |

---

## 操作指南

解决特定问题的步骤说明。

| 文档 | 说明 |
|------|------|
| [command-development.md](command-development.md) | 命令模块开发教程 |

---

## 技术参考

协议规范、API 参考和数据格式。

### 核心系统

| 文档 | 说明 |
|------|------|
| [auth.md](auth.md) | 认证系统（微软/Xbox/Minecraft 三级认证流程） |
| [protocol.md](protocol.md) | 协议实现（包定义、编解码、加密流程） |
| [tui.md](tui.md) | TUI 框架（终端界面架构与组件） |
| [player.md](player.md) | 玩家数据（状态、位置、背包系统） |

### 数据格式

| 文档 | 说明 |
|------|------|
| [formats/nbt_format.md](formats/nbt_format.md) | NBT 二进制格式规范 |
| [formats/snbt_format.md](formats/snbt_format.md) | SNBT 字符串格式规范 |
| [formats/text_component.md](formats/text_component.md) | 文本组件 JSON 格式 |

---

## 设计文档

系统架构设计和技术决策记录。

### 活跃设计

| 文档 | 说明 | 状态 |
|------|------|------|
| [manager_spec.md](superpowers/specs/manager_spec.md) | 集群管理器规格说明 | 计划中 |

### 已完成设计（存档）

| 文档 | 说明 |
|------|------|
| command-module | 命令模块设计（已完成实现） |
| entity-tracking-design | 实体跟踪系统设计（已完成实现） |
| config-auto-update-design | 配置自动更新设计（已完成实现） |
| component-parsing | 物品组件解析设计（已完成实现） |

完整存档见 `superpowers/plans/archive/` 目录。

---

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

---

## 目录结构

```
docs/
├── README.md              # 本文件
├── auth.md                # 认证系统
├── command-development.md # 命令开发教程
├── player.md             # 玩家系统
├── protocol.md           # 协议规范
├── tui.md                # TUI 框架
├── formats/              # 数据格式参考
│   ├── nbt_format.md
│   ├── snbt_format.md
│   └── text_component.md
└── superpowers/          # 设计文档
    ├── specs/            # 规格说明（活跃）
    └── plans/            # 实现计划
        └── archive/      # 已完成的计划存档
```

---

## 相关资源

- [Minecraft Wiki - 协议](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE)
- [node-minecraft-protocol](https://github.com/PrismarineJS/node-minecraft-protocol)
- [minecraft-data](https://github.com/PrismarineJS/minecraft-data)