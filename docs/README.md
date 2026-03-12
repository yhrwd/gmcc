# gmcc 文档目录

## 概述

gmcc 是一个 Go 语言实现的 Minecraft Java 版控制台客户端，支持协议版本 774 (1.21.11)。本文档目录包含了完整的项目文档。

## 文档列表

### [开发指南](development.md)
完整的开发文档，包含：
- 项目架构和设计理念
- 各模块详细说明
- 协议实现细节
- 扩展开发指南

### [协议实现](protocol.md)
Minecraft 协议的实现细节：
- 协议版本 774 的包定义
- 包编解码流程
- 状态机设计
- 特性开关

### [认证系统](auth.md)
Minecraft 认证流程：
- 微软设备码登录
- Xbox/XSTS 认证
- Minecraft 会话认证
- 令牌缓存机制

### [TUI框架](tui.md)
终端用户界面框架：
- TUI 架构设计
- 组件系统
- 事件处理
- 渲染系统

### [玩家数据](player.md)
玩家信息获取：
- 玩家状态
- 位置坐标
- 背包系统

## 快速开始

### 安装

```bash
# 从源码编译
go build -o gmcc ./cmd/gmcc

# 或从 Releases 下载
```

### 配置

创建 `config.yaml` 配置文件：

```yaml
account:
  player_id: "你的游戏ID"
  use_official_auth: true  # 正版认证

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

## 项目结构

```
gmcc/
├── cmd/gmcc/              # 程序入口
├── internal/              # 内部模块
│   ├── auth/              # 认证相关
│   │   ├── microsoft/     # 微软/Xbox/XSTS 认证
│   │   └── minecraft/     # Minecraft 认证
│   ├── config/            # 配置加载
│   ├── logx/              # 日志系统
│   ├── mcclient/          # 客户端核心
│   │   ├── client.go      # 状态机和主循环
│   │   ├── codec.go       # 协议编解码
│   │   ├── chat.go         # 聊天/命令
│   │   ├── handlers_*.go   # 各状态包处理
│   │   ├── text_component.go # 文本组件解析
│   │   └── protocol_774.go # 协议常量
│   ├── nbt/               # NBT 数据处理
│   ├── player/            # 玩家信息管理
│   ├── session/           # 会话缓存
│   └── tui/               # TUI 框架
├── pkg/                   # 公共工具
│   └── httpx/             # HTTP 工具
├── docs/                  # 文档
└── .knowledge/            # 协议知识库
```

## 核心功能

- ✅ 微软正版认证（设备码登录/刷新令牌）
- ✅ 离线模式支持
- ✅ 登录加密（AES/CFB8）
- ✅ 登录压缩
- ✅ 配置状态处理
- ✅ 游戏状态心跳
- ✅ 聊天消息接收/发送
- ✅ 命令发送
- ✅ CESU-8 编码支持
- ✅ NBT 数据解析
- ✅ TUI 界面
- ✅ 玩家状态获取
- ✅ 玩家位置获取
- ✅ 背包信息获取

## 相关资源

- [Minecraft Wiki - 协议](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE)
- [node-minecraft-protocol](https://github.com/PrismarineJS/node-minecraft-protocol)
- [minecraft-data](https://github.com/PrismarineJS/minecraft-data)