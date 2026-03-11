# 协议实现文档

## 概述

本文档详细描述了 gmcc 支持的 Minecraft Java 版协议实现，目标版本为协议 774 (MC 1.21.11)。

## 协议版本

```go
const protocolVersion = 774
const protocolLabel = "Java 1.21.11 / protocol 774"
```

## 连接状态机

客户端连接经历三种状态：

```
Handshake -> Login -> Configuration -> Play
```

### 状态定义

| 状态 | 说明 |
|------|------|
| `stateLogin` | 登录阶段，处理加密、压缩、登录成功 |
| `stateConfiguration` | 配置阶段，处理注册表、资源包、完成配置 |
| `statePlay` | 游戏阶段，处理心跳、聊天、位置等 |

## 包 ID 定义

### Handshake 阶段

| 包名 | ID | 方向 | 说明 |
|------|-----|------|------|
| `handshaking_server_intention` | 0x00 | C→S | 握手意图 |

### Login 阶段（客户端接收）

| 包名 | ID | 说明 |
|------|-----|------|
| `login_disconnect` | 0x00 | 断开连接 |
| `encryption_request` | 0x01 | 加密请求 |
| `login_success` | 0x02 | 登录成功 |
| `set_compression` | 0x03 | 设置压缩 |
| `login_plugin_request` | 0x04 | 插件请求 |
| `cookie_request` | 0x05 | Cookie 请求 |

### Login 阶段（客户端发送）

| 包名 | ID | 说明 |
|------|-----|------|
| `login_server_hello` | 0x00 | 登录开始 |
| `login_server_key` | 0x01 | 加密响应 |
| `login_server_custom_answer` | 0x02 | 插件响应 |
| `login_server_ack` | 0x03 | 确认登录 |
| `login_server_cookie_resp` | 0x04 | Cookie 响应 |

### Configuration 阶段（客户端接收）

| 包名 | ID | 说明 |
|------|-----|------|
| `cookie_request` | 0x00 | Cookie 请求 |
| `custom_payload` | 0x01 | 自定义负载 |
| `disconnect` | 0x02 | 断开连接 |
| `finish_configuration` | 0x03 | 完成配置 |
| `keep_alive` | 0x04 | 心跳 |
| `ping` | 0x05 | Ping |
| `registry_data` | 0x07 | 注册表数据 |
| `resource_pack_pop` | 0x08 | 资源包弹出 |
| `resource_pack_push` | 0x09 | 资源包推送 |
| `select_known_packs` | 0x0E | 选择已知包 |
| `custom_report_details` | 0x13 | 自定义报告详情 |

### Configuration 阶段（客户端发送）

| 包名 | ID | 说明 |
|------|-----|------|
| `client_information` | 0x00 | 客户端信息 |
| `cookie_response` | 0x01 | Cookie 响应 |
| `custom_payload` | 0x02 | 自定义负载 |
| `finish_configuration` | 0x03 | 完成配置 |
| `keep_alive` | 0x04 | 心跳响应 |
| `pong` | 0x05 | Pong |
| `resource_pack` | 0x06 | 资源包响应 |
| `select_known_packs` | 0x07 | 选择已知包 |
| `accept_chat_details` | 0x09 | 接受聊天详情 |

### Play 阶段（客户端接收 - 主要）

| 包名 | ID | 说明 |
|------|-----|------|
| `declare_commands` | 0x10 | 命令树声明 |
| `cookie_request` | 0x15 | Cookie 请求 |
| `disconnect` | 0x20 | 断开连接 |
| `profileless_chat` | 0x21 | 无档案聊天 |
| `keep_alive` | 0x2B | 心跳 |
| `login` | 0x30 | 游戏登录 |
| `ping` | 0x3B | Ping |
| `player_chat` | 0x3F | 玩家聊天 |
| `player_position` | 0x46 | 玩家位置 |
| `resource_pack_pop` | 0x4E | 资源包弹出 |
| `resource_pack_push` | 0x4F | 资源包推送 |
| `action_bar` | 0x55 | 动作栏 |
| `system_chat` | 0x77 | 系统聊天 |
| `set_health` | 0x54 | 设置生命值 |
| `set_experience` | 0x4D | 设置经验值 |
| `player_info_update` | 0x42 | 玩家信息更新 |
| `player_info_remove` | 0x3D | 玩家信息移除 |
| `entity_data` | 0x5D | 实体数据 |

### Play 阶段（客户端发送 - 主要）

| 包名 | ID | 说明 |
|------|-----|------|
| `accept_teleportation` | 0x00 | 接受传送 |
| `chat_ack` | 0x05 | 聊天确认 |
| `chat_command` | 0x06 | 聊天命令（无签名）|
| `chat_command_signed` | 0x07 | 聊天命令（签名）|
| `chat_message` | 0x08 | 聊天消息 |
| `chat_session_update` | 0x09 | 聊天会话更新 |
| `client_information` | 0x0D | 客户端信息 |
| `cookie_response` | 0x14 | Cookie 响应 |
| `keep_alive` | 0x1B | 心跳响应 |
| `move_player_status_only` | 0x20 | 移动状态 |
| `pong` | 0x2C | Pong |
| `resource_pack` | 0x30 | 资源包响应 |
| `move_player_pos` | 0x16 | 玩家位置更新 |
| `move_player_rot` | 0x17 | 玩家旋转更新 |
| `move_player_pos_rot` | 0x18 | 玩家位置旋转更新 |
| `container_close` | 0x0B | 关闭容器 |
| `container_click` | 0x0C | 容器点击 |

## 数据类型编解码

### 基础类型

```go
// VarInt: 可变长度整数
func readVarInt(r io.ByteReader) (int32, error)
func encodeVarInt(v int32) []byte

// String: UTF-8 字符串（带 CESU-8 转换）
func readString(r io.ByteReader, rr io.Reader) (string, error)
func encodeString(s string) []byte

// Bool: 布尔值
func readBool(r io.Reader) (bool, error)
func encodeBool(v bool) []byte

// Int64: 64位整数（大端）
func readInt64(r io.Reader) (int64, error)
func encodeInt64(v int64) []byte

// Int32: 32位整数（大端）
func readInt32(r io.Reader) (int32, error)
func encodeInt32(v int32) []byte

// UUID: 128位唯一标识符
func readUUID(r io.Reader) ([16]byte, error)

// ByteArray: 字节数组
func readByteArray(r io.ByteReader, rr io.Reader) ([]byte, error)
func encodeByteArray(b []byte) []byte
```

### NBT 数据

NBT (Named Binary Tag) 格式用于存储复杂数据结构：

```go
// 网络格式 NBT（匿名根标签）
func readAnonymousNBTJSON(r io.Reader) (string, error)

// NBT 解码器
dec := nbt.NewDecoder(r)
dec.NetworkFormat(true)
var v any
err := dec.Decode(&v)
```

## 包帧结构

### 未压缩帧

```
[长度: VarInt] [包ID: VarInt] [负载数据: byte[]]
```

### 压缩帧（threshold >= 0）

```
[帧长度: VarInt] [未压缩长度: VarInt] [压缩数据: byte[]]
```

如果 `未压缩长度 == 0`，则后续数据未压缩。
否则，后续数据为 zlib 压缩的数据。

## 加密流程

### 1. 收到 encryption_request

```go
type EncryptionRequest struct {
    ServerID          string  // 通常为空
    PublicKey         []byte  // RSA 公钥 DER 编码
    VerifyToken       []byte  // 验证令牌
    ShouldAuthenticate bool    // 是否需要正版认证
}
```

### 2. 生成共享密钥

```go
sharedSecret := make([]byte, 16)
rand.Read(sharedSecret)
```

### 3. 会话认证（如果需要）

```go
serverHash := minecraftServerHash(serverID, sharedSecret, publicKey)
mcauth.JoinServer(accessToken, profileID, serverHash)
```

### 4. 发送 encryption_response

```go
encryptedSecret := rsa.EncryptPKCS1v15(pubKey, sharedSecret)
encryptedToken := rsa.EncryptPKCS1v15(pubKey, verifyToken)
```

### 5. 启用加密流

```go
// AES/CFB8 加密
block, _ := aes.NewCipher(sharedSecret)
decrypter := newCFB8(block, sharedSecret, false)
encrypter := newCFB8(block, sharedSecret, true)
```

### Server Hash 计算

```go
func minecraftServerHash(serverID string, sharedSecret, publicKey []byte) string {
    h := sha1.New()
    h.Write([]byte(serverID))
    h.Write(sharedSecret)
    h.Write(publicKey)
    sum := h.Sum(nil)
    
    // Minecraft 使用特殊的十六进制格式
    // 负数时前面加 "-"，且无前导零
    if sum[0]&0x80 != 0 {
        // 二进制补码
        // ...
    }
    return hexString
}
```

## 聊天签名系统

### Secure Chat 会话初始化

```go
// 1. 获取玩家证书
cert, _ := mcauth.GetPlayerCertificates(accessToken)

// 2. 解析私钥和公钥
privateKey, _ := parsePrivateKeyPEM(cert.KeyPair.PrivateKey)
pubDER, _ := parsePublicKeyPEMToSPKIDER(cert.KeyPair.PublicKey)

// 3. 生成会话 UUID
sessionID, _ := randomUUIDv4()

// 4. 发送 chat_session_update
payload := sessionID + expireMillis + pubDER + signature
```

### 消息签名

```go
signable := encodeInt32(1)           // 版本
signable += uuid                      // 玩家 UUID
signable += sessionID                 // 会话 ID
signable += encodeInt32(messageIndex) // 消息索引
signable += encodeInt64(salt)         // 盐值
signable += encodeInt64(timestamp)    // 时间戳
signable += encodeInt32(len(content)) // 内容长度
signable += content                   // 消息内容
signable += acknowledgements          // 确认

signature := rsa.SignPKCS1v15(privateKey, crypto.SHA256, sha256(signable))
```

### 命令参数签名

命令树中标记为可签名的参数（如 `message`、`msg`）需要签名：

```go
// 解析命令树提取签名规则
targets := extractSignableCommandTargets(nodes, rootIndex)
// 例如: { "msg": { ArgumentName: "message", SliceIndex: 1 } }
```

## 配置阶段处理

### 客户端信息

```go
payload := encodeString("zh_cn")  // 语言
payload += byte(8)                 // 视距
payload += encodeVarInt(0)         // 聊天模式
payload += encodeBool(true)        // 启用聊天颜色
payload += byte(0x7F)              // 显示皮肤部位
payload += encodeVarInt(1)         // 主手（1=右手）
payload += encodeBool(false)       // 启用文本过滤
payload += encodeBool(true)        // 允许服务器列表
payload += encodeVarInt(0)         // 粒子设置
```

### 资源包响应

```go
// 按顺序发送三种状态
actions := []int32{
    resourcePackAccepted,   // 3: 接受
    resourcePackDownloaded, // 4: 已下载
    resourcePackLoaded,     // 0: 已加载
}
```

## Play 阶段心跳

### Keep-Alive

服务器定期发送 keep_alive 包，客户端需要立即回复相同的 ID：

```go
// 接收
id, _ := readInt64(r)

// 发送
WritePacket(playServerKeepAlive, encodeInt64(id))
```

### AFK 心跳

为保持在线状态，客户端定期发送移动状态包：

```go
if time.Since(lastAFKPacket) >= 15*time.Second {
    WritePacket(playServerMoveStatus, []byte{0x01})
    lastAFKPacket = time.Now()
}
```

## 文本组件

### 解析聊天 JSON

聊天消息采用 JSON 文本组件格式：

```json
{
  "text": "Hello ",
  "color": "red",
  "extra": [
    {
      "text": "World",
      "bold": true
    }
  ]
}
```

### 颜色代码

| 颜色 | 代码 | ANSI |
|------|------|------|
| black | §0 | \033[30m |
| dark_blue | §1 | \033[34m |
| dark_green | §2 | \033[32m |
| dark_aqua | §3 | \033[36m |
| dark_red | §4 | \033[31m |
| dark_purple | §5 | \033[35m |
| gold | §6 | \033[33m |
| gray | §7 | \033[37m |
| dark_gray | §8 | \033[90m |
| blue | §9 | \033[94m |
| green | §a | \033[92m |
| aqua | §b | \033[96m |
| red | §c | \033[91m |
| light_purple | §d | \033[95m |
| yellow | §e | \033[93m |
| white | §f | \033[97m |

## CESU-8 编码

Minecraft 使用 CESU-8 编码字符串（特别是中文和 Emoji）：

```go
// CESU-8 转 UTF-8
func CESU8ToUTF8(b []byte) string {
    // 将 CESU-8 编码的字节转为有效的 UTF-8
    // 处理代理对编码的 4 字节序列
}
```

## 协议特性开关

```go
type protocolFeatures struct {
    HasConfigurationState bool  // 是否有配置阶段
    SignatureEncryption   bool  // 是否启用签名加密
}

var features774 = protocolFeatures{
    HasConfigurationState: true,
    SignatureEncryption:   false,
}
```

## 扩展协议版本

升级协议版本时需更新：

1. `protocol_774.go` 中的包 ID 常量
2. `features774` 特性开关
3. 包名映射 `packetName()`
4. 数据类型编解码（如有变化）
5. 新的数据包处理器