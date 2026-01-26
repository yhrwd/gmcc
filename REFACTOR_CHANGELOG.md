# 项目结构重构说明

## 重构目标

将 Minecraft 客户端工程的目录结构进行优化，将可复用的公开包迁移到 `pkg/` 目录下，使代码结构更清晰、便于维护和复用。

## 主要变更

### 迁移到 `pkg/` 的包（公开可复用）

1. **`pkg/logger/`**
   - 源自 `internal/logger/`
   - 日志管理模块，支持彩色输出、文件输出、UI 回调

2. **`pkg/crypto/`**
   - 源自 `internal/crypto/`
   - Token 持久化存储，使用 GOB 编码格式

3. **`pkg/protocol/`** - 协议栈（Minecraft 通信核心）
   - **`codec/`** - 数据编码/解码
     - `varint.go` - VarInt 编码（Minecraft 变长整数）
     - `varstring.go` - VarString 编码（长度前缀字符串）
     - `compress.go` - 数据压缩（占位实现）
   - **`connection/`** - 连接管理
     - `conn.go` - 封装 bufio 的网络连接
     - `fsm.go` - 状态机实现（Handshake → Status → Login → Config → Play）
     - `writeraw.go` - 原始包发送（处理压缩、加密、长度前缀）
   - **`packet/`** - 数据包解析
     - `packet.go` - 通用数据包结构和读取器
   - **`stage/`** - 协议状态处理
     - `handshake.go` - 握手状态（客户端初始化）
     - `status.go` - 状态查询状态（Ping/Pong）

4. **`pkg/data/`**
   - 源自 `internal/data/`
   - 数据结构定义（ServerStatus JSON 结构体）

5. **`pkg/fileprocess/`**
   - 源自 `internal/fileprocess/`
   - PNG 图像处理（Player Skin 渲染）

### 保留在 `internal/` 的包（内部依赖）

- **`auth/`** - 身份验证（Microsoft、PlayerInfo）
- **`config/`** - 配置管理
- **`tui/`** - 终端用户界面（Bubble Tea）
- **`protocol/addrcheck.go`** - 地址检查
- **`protocol/startmc.go`** - 连接启动入口

## 导入更新

所有内部模块的 import 已从 `gmcc/internal/...` 更新为 `gmcc/pkg/...`：

```diff
// 示例变更
- import "gmcc/internal/protocol/codec"
+ import "gmcc/pkg/protocol/codec"

- import "gmcc/internal/logger"
+ import "gmcc/pkg/logger"
```

## 构建验证

- ✅ `go mod tidy` 通过
- ✅ `go build ./...` 通过
- ✅ 所有包编译成功，无警告

## 文件变更统计

- **创建**: 21 个新文件（`pkg/` 下）
- **删除**: 12 个旧文件（`internal/` 下已迁移的包）
- **修改**: 11 个 import 语句跨多个文件

## 后续建议

1. 为 `pkg/protocol/codec` 补充完整的 Compress 实现
2. 添加单元测试覆盖 `pkg/protocol` 的关键路径
3. 补充 API 文档（pkg 包的使用说明）
4. 考虑将 `internal/tui` 中的通用组件抽取到 `pkg/ui` 以便复用

## 兼容性说明

该重构是**向后不兼容**的改动，仅影响内部 import 路径。任何依赖本工程的外部代码需要更新其 import 语句。
