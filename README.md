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

log:
  log_dir: "logs"
  max_size: 10              # 单个日志文件最大大小（MB）
  debug: false
  enable_file: true
```

### 运行

```bash
./gmcc
```

## 项目结构

```
cmd/gmcc/          # 程序入口
internal/          # 核心模块
  auth/            # 认证 (microsoft, minecraft)
  config/          # 配置加载
  item/            # 物品系统 (新增)
    component/     # 组件解析框架
  logx/            # 日志系统
  mcclient/        # 客户端核心
  nbt/             # NBT 数据处理
  player/          # 玩家状态
  session/         # 令牌缓存
  tui/             # 终端 UI
pkg/               # 公共工具
  binutil/         # 二进制工具 (新增)
```

## 文档

详细文档请查看 [docs/README.md](docs/README.md)

## 协议版本

当前支持：**协议 774** (Minecraft Java 1.21.11)

## 许可证

MIT License
