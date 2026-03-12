# 开发文档

## 目标

本项目针对 **Java 版协议 774** (1.21.11)，提供完整的 Minecraft 客户端功能：

- 连接服务器
- 登录/认证（离线/正版）
- 进入游戏状态
- 保持在线心跳
- 聊天消息收发
- 命令执行
- 玩家状态/位置/背包信息获取
- TUI 终端界面

设计理念：

- 模块化：协议、认证、界面分离
- 可扩展：通过包处理器添加新功能
- 协议行为按功能开关控制，而非硬编码版本
- 日志包含协议细节，便于问题定位

## 工具结构

```
gmcc/
├── cmd/gmcc/                 # 程序入口
│   └── main.go               # 启动、信号处理、协调各模块
│
├── internal/                 # 内部模块（不对外暴露）
│   ├── auth/                 # 认证模块
│   │   ├── microsoft/        # 微软/Xbox/XSTS 认证
│   │   │   └── service.go    # 设备码登录、令牌刷新
│   │   └── minecraft/        # Minecraft 认证
│   │       └── service.go    # MC 令牌获取、会话加入、证书
│   │
│   ├── config/               # 配置管理
│   │   ├── config.go         # 配置结构定义
│   │   └── loader.go         # YAML 加载、验证
│   │
│   ├── logx/                 # 日志系统
│   │   └── logx.go           # 文件日志、控制台输出
│   │
│   ├── mcclient/             # Minecraft 客户端核心
│   │   ├── client.go         # 状态机、主循环、连接管理
│   │   ├── codec.go          # 包读写、VarInt、加密、压缩
│   │   ├── protocol_774.go   # 协议常量、包 ID 映射
│   │   ├── chat.go           # 聊天消息、命令发送
│   │   ├── chat_parser.go    # 聊天 JSON 解析
│   │   ├── text_component.go # 文本组件、ANSI 转换
│   │   ├── handlers_preplay.go # 登录/配置阶段包处理
│   │   ├── handlers_play.go    # 游戏阶段包处理
│   │   ├── handlers_player.go  # 玩家相关包处理
│   │   ├── handlers_chat.go    # 聊天包处理
│   │   └── utils.go            # 工具函数
│   │
│   ├── nbt/                  # NBT 数据处理
│   │   ├── nbt.go            # 类型定义
│   │   ├── decode.go         # 解码器
│   │   ├── encode.go         # 编码器
│   │   ├── snbt.go           # SNBT 解析
│   │   ├── path.go           # 路径查询
│   │   └── raw.go            # 原始 NBT
│   │
│   ├── player/               # 玩家信息管理
│   │   ├── player.go         # 玩家状态、位置、背包
│   │   └── inventory.go      # 背包管理
│   │
│   ├── session/              # 会话缓存
│   │   └── cache.go          # 令牌持久化
│   │
│   └── tui/                  # TUI 框架
│       └── tui.go            # 终端界面
│
├── pkg/                      # 公共工具包
│   └── httpx/                # HTTP 客户端
│       └── client.go         # 请求封装、重试
│
├── docs/                     # 文档
│   ├── README.md             # 文档索引
│   ├── development.md        # 开发指南（本文件）
│   ├── protocol.md           # 协议实现
│   ├── auth.md               # 认证系统
│   ├── tui.md                # TUI 框架
│   └── player.md             # 玩家数据
│
└── .knowledge/               # 协议知识库
    ├── README.md             # 索引
    ├── MC_Protocol_Data/     # 协议数据
    ├── minecraft-data/       # 游戏数据
    ├── prismarine-chat/      # 聊天解析参考
    └── prismarine-nbt/       # NBT 解析参考
```

## 包说明

### `cmd/gmcc`

程序入口，加载配置、初始化日志、启动客户端。

### `internal/config`

- `Load(path)` 读取并校验 YAML 配置
- `Default()` 生成默认配置模板

### `internal/auth/microsoft`

处理设备码登录和刷新令牌流程，提供 XSTS 转换工具。

### `internal/auth/minecraft`

用 XSTS 令牌换取 Minecraft 访问令牌，验证游戏所有权，调用 session join。

### `internal/session`

令牌缓存（`.session/<玩家ID>.json`），优先级：缓存的 MC 令牌 → MS 刷新令牌 → 设备码登录。

### `internal/mcclient`

核心客户端实现：

- `codec.go`: 数据帧读写、压缩、AES/CFB8 加密流、字符串编解码
- `protocol_774.go`: 包 ID、协议特性开关、包名映射
- `client.go`: 登录/配置/游戏状态机，心跳循环
- `chat.go`: 消息/命令收发，聊天包解析
- `chat_parser.go`: 聊天 JSON 文本提取

### `internal/nbt`

NBT 数据处理：

- `decode.go`: NBT 解码器，支持 CESU-8
- `encode.go`: NBT 编码器
- `snbt.go`: SNBT（字符串化 NBT）解析
- `path.go`: NBT 路径查询

**NBT 路径查询示例：**

```go
// 解析 NBT 数据
data := map[string]any{
    "Inventory": []any{
        map[string]any{"Slot": int8(0), "id": "diamond"},
        map[string]any{"Slot": int8(1), "id": "iron"},
    },
}

// 查询路径
results, _ := nbt.QueryPath(data, "Inventory[0].id")  // => "diamond"
results, _ := nbt.QueryPath(data, "Inventory[{Slot:1b}].id")  // => "iron"
results, _ := nbt.QueryPath(data, "Inventory[].id")  // => ["diamond", "iron"]
```

支持的路径节点：

| 语法 | 说明 | 示例 |
|------|------|------|
| `{tags}` | 匹配根复合标签 | `{Invisible:1b}` |
| `name` | 访问子标签 | `foo.bar` |
| `name{tags}` | 带模式匹配 | `VillagerData{type:"plains"}` |
| `[]` | 列表所有元素 | `Items[]` |
| `[index]` | 列表索引（支持负数） | `Items[0]`, `Items[-1]` |
| `[{tags}]` | 按模式过滤列表 | `Items[{count:25}]` |

### `internal/logx`

统一的标准输出/文件日志，调试模式可输出协议级详情。

### `pkg/*`

通用工具包：

- `httpx`: HTTP 请求封装

## 运行流程

1. 加载配置（`internal/config`）
2. 确定认证模式（`use_official_auth`）：
   - 离线模式，或
   - 微软正版认证（缓存 → 刷新 → 设备码）
3. 建立 TCP 连接，发送握手/登录包
4. 处理登录加密：
   - 调用 session join
   - 发送 `encryption_begin` 响应
   - 启用 AES/CFB8 加密流
5. 处理登录压缩和登录成功
6. 进入配置状态，响应必要的包
7. 进入游戏状态，开始心跳循环
8. 首次进入游戏时执行：
   - `actions.on_join_commands`
   - `actions.on_join_messages`

## 协议元数据

`internal/mcclient/protocol_774.go` 包含：

- 登录/配置/游戏状态的包 ID
- 协议特性开关（`features774`）
- 包名映射（`packetName`）用于调试日志
- 状态名辅助函数（`stateName`）

升级协议版本时，首先更新此文件。

## 消息与命令收发

### 接收

游戏状态下的聊天包处理器：

- `system_chat` (`0x77`)
- `player_chat` (`0x3F`)
- `action_bar` (`0x55`)
- `profileless_chat` (`0x21`)

对于 NBT 格式的聊天组件，客户端使用网络 NBT 格式解码，提供原始 JSON 字符串（`ChatMessage.RawJSON`）。

`ChatMessage` 字段：

- `Type`: 包类型（`system`, `player_chat`, `action_bar` 等）
- `PlainText`: 提取的纯文本
- `RawJSON`: 原始聊天 JSON
- `SenderUUID`: 发送者 UUID
- `ReceivedAt`: 接收时间

注册回调：

```go
client.SetChatHandler(func(msg mcclient.ChatMessage) {
    // 解析 msg.RawJSON
})
```

### 发送

`internal/mcclient/chat.go` 提供的 API：

- `SendCommand(command string)`: 发送命令（`chat_command_signed`，签名简化版）
- `SendMessage(message string)`: 发送消息（`chat_message`，无签名）
- 入服后自动发送 `chat_session_update`，提供安全聊天所需的公钥

配置驱动的自动发送：

```yaml
actions:
  delay_ms: 1200
  on_join_commands:
    - "list"
  on_join_messages:
    - "大家好"
```

注意：

- `on_join_commands` 不需要前导 `/`
- `delay_ms` 控制入服动作延迟，默认 1200ms
- 无签名的消息可能被强制安全聊天的服务器拒绝

## 日志与诊断

`internal/logx` 同时输出到控制台和 `logs/gmcc.log`。

`log.debug=true` 启用协议诊断：

- 发出包的 ID/长度/预览
- 接收帧的长度/压缩详情
- 包 ID + 符号名 + 状态
- 加密请求/响应详情

控制台格式仅显示时间（`HH:MM:SS`），文件日志包含完整日期时间。
文件日志按大小（`log.max_size`）滚动为 `gmcc-*.log`。

可用于定位：

- 包结构错误
- 加密同步问题
- 压缩阈值问题
- 意外的状态转换

## 扩展清单

添加包支持或升级版本时：

1. 更新 `protocol_774.go` 中的包 ID 和特性开关
2. 根据协议数据验证加密包结构
3. 在 `client.go` 添加包处理分支
4. 添加调试日志
5. 编写编解码/协议边界测试
6. 运行：
   - `go test ./...`
   - `go build ./cmd/gmcc`

## 外部参考

- node-minecraft-protocol: https://github.com/PrismarineJS/node-minecraft-protocol
- 协议数据源: https://github.com/PrismarineJS/minecraft-data