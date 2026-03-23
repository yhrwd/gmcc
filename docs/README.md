# gmcc 文档目录

## 概述

gmcc 是一个 Go 语言实现的 Minecraft Java 版控制台客户端，支持协议版本 774 (1.21.11)。

## 文档列表

| 文档 | 说明 |
|------|------|
| [development.md](development.md) | 开发指南（项目架构、模块说明、运行流程） |
| [protocol.md](protocol.md) | 协议实现（包定义、编解码、加密流程） |
| [auth.md](auth.md) | 认证系统（微软/Xbox/Minecraft 认证流程） |
| [tui.md](tui.md) | TUI 框架（终端界面架构与组件） |
| [player.md](player.md) | 玩家数据（状态、位置、背包系统） |

## 格式规范

| 文档 | 说明 |
|------|------|
| [nbt_format.md](nbt_format.md) | NBT 二进制格式规范 |
| [snbt_format.md](snbt_format.md) | SNBT 字符串格式规范 |
| [text_component.md](text_component.md) | 文本组件 JSON 格式 |
| [data_components_1.21.11.md](data_components_1.21.11.md) | 1.21.11 数据组件规范 |

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

log:
  log_dir: "logs"
  max_size: 10
  debug: false
  enable_file: true
```

### 运行

```bash
./gmcc
```

## 核心功能

- 微软正版认证（设备码登录/刷新令牌）
- 离线模式支持
- 登录加密（AES/CFB8）
- 登录压缩
- 配置状态处理
- 游戏状态心跳
- 聊天消息收发
- 命令发送
- CESU-8 编码支持
- NBT 数据解析
- TUI 界面
- 玩家状态获取
- 玩家位置获取
- 背包信息获取

## 相关资源

- [Minecraft Wiki - 协议](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE)
- [node-minecraft-protocol](https://github.com/PrismarineJS/node-minecraft-protocol)
- [minecraft-data](https://github.com/PrismarineJS/minecraft-data)
