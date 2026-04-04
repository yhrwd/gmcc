# 认证系统文档

## 概述

gmcc 支持 Minecraft Java 版的两种认证模式：
- **正版认证**：Microsoft/Xbox/Minecraft 三级认证链
- **离线模式**：直接使用玩家名称

## 认证流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                       正版认证流程                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. 检查缓存                                                    │
│     ├─ 有效 Minecraft Token → 直接使用                          │
│     ├─ 有效 MS Access Token → 换取 XSTS                         │
│     └─ 使用 Refresh Token → 刷新 MS Token                       │
│                                                                 │
│  2. 设备码登录（无缓存时）                                       │
│     ├─ 请求 Microsoft 设备码                                    │
│     ├─ 用户访问 URL 输入设备码                                  │
│     ├─ 轮询等待授权完成                                         │
│     └─ 获取 Microsoft Access Token                              │
│                                                                 │
│  3. Xbox Live 认证                                              │
│     └─ MS Token → Xbox Token → XSTS Token                       │
│                                                                 │
│  4. Minecraft 认证                                              │
│     ├─ XSTS Token → Minecraft Access Token                      │
│     ├─ 验证游戏所有权                                           │
│     └─ 获取玩家 Profile                                         │
│                                                                 │
│  5. 会话加入                                                    │
│     └─ JoinServer(serverHash) → 服务器验证                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 微软认证服务

### 端点

| 端点 | URL | 用途 |
|------|-----|------|
| Device Code | `https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode` | 获取设备码 |
| Token | `https://login.microsoftonline.com/consumers/oauth2/v2.0/token` | 获取/刷新令牌 |
| Xbox Auth | `https://user.auth.xboxlive.com/user/authenticate` | Xbox 认证 |
| XSTS Auth | `https://xsts.auth.xboxlive.com/xsts/authorize` | XSTS 认证 |

### 设备码登录

```go
// 1. 请求设备码
type deviceCodeResponse struct {
    DeviceCode      string `json:"device_code"`
    UserCode        string `json:"user_code"`
    VerificationURI string `json:"verification_uri"`
    ExpiresIn       int    `json:"expires_in"`
    Interval        int    `json:"interval"`
}

// 2. 轮询等待授权
for {
    time.Sleep(time.Duration(interval) * time.Second)
    resp, err := pollToken(deviceCode)
    if resp.AccessToken != "" {
        return resp, nil
    }
}
```

### Xbox/XSTS 认证

```go
// Xbox 认证
xboxRequest := {
    "Properties": {
        "AuthMethod": "RPS",
        "SiteName": "user.auth.xboxlive.com",
        "RpsTicket": "d=" + msAccessToken
    },
    "RelyingParty": "http://auth.xboxlive.com",
    "TokenType": "JWT"
}

// XSTS 认证
xstsRequest := {
    "Properties": {
        "SandboxId": "RETAIL",
        "UserTokens": [xboxToken]
    },
    "RelyingParty": "rp://api.minecraftservices.com/",
    "TokenType": "JWT"
}
```

## Minecraft 认证服务

### 端点

| 端点 | URL | 用途 |
|------|-----|------|
| Login | `https://api.minecraftservices.com/authentication/login_with_xbox` | 登录 |
| Ownership | `https://api.minecraftservices.com/entitlements/mcstore` | 验证所有权 |
| Profile | `https://api.minecraftservices.com/minecraft/profile` | 获取档案 |
| Certificates | `https://api.minecraftservices.com/player/certificates` | 玩家证书 |
| Join | `https://sessionserver.mojang.com/session/minecraft/join` | 会话加入 |
| HasJoined | `https://sessionserver.mojang.com/session/minecraft/hasJoined` | 验证加入 |

### 获取 Minecraft Token

```go
type minecraftLoginResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
}

func GetMinecraftToken(xsts *XSTSResponse) (*minecraftLoginResponse, error) {
    // 构建请求
    body := map[string]string{
        "identityToken": "XBL3.0 x=" + xsts.DisplayClaims.Xui[0].Uhs + ";" + xsts.Token,
    }
    // POST 到 Minecraft 登录端点
}
```

### 验证游戏所有权

```go
func VerifyGameOwnership(accessToken string) error {
    // GET /entitlements/mcstore
    // 检查响应中是否包含 Minecraft 项目
}
```

### 获取玩家档案

```go
type MinecraftProfile struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func GetProfile(accessToken string) (*MinecraftProfile, error) {
    // GET /minecraft/profile
    // 返回玩家 UUID 和名称
}
```

### 获取玩家证书

用于聊天签名：

```go
type PlayerCertificate struct {
    KeyPair struct {
        PrivateKey string `json:"privateKey"`
        PublicKey  string `json:"publicKey"`
    } `json:"keyPair"`
    PublicKeySignature   string `json:"publicKeySignature"`
    PublicKeySignatureV2 string `json:"publicKeySignatureV2"`
    ExpiresAt            string `json:"expiresAt"`
}
```

## 会话加入

### 流程

```
1. 收到服务器 encryption_request
   ↓
2. 生成 shared_secret (16 bytes)
   ↓
3. 计算 server_hash = sha1(serverID + shared_secret + publicKey)
   ↓  
4. 调用 JoinServer(accessToken, profileID, serverHash)
   ↓
5. 发送 encryption_response
   ↓
6. 启用 AES/CFB8 加密流
```

### Server Hash 计算

Minecraft 使用特殊的 SHA-1 十六进制格式：

```go
func minecraftServerHash(serverID string, sharedSecret, publicKey []byte) string {
    h := sha1.New()
    h.Write([]byte(serverID))
    h.Write(sharedSecret)
    h.Write(publicKey)
    sum := h.Sum(nil)
    
    // 处理负数情况（二进制补码）
    if sum[0]&0x80 != 0 {
        // 取反加一
        for i := range sum {
            sum[i] = ^sum[i]
        }
        for i := len(sum) - 1; i >= 0; i-- {
            sum[i]++
            if sum[i] != 0 {
                break
            }
        }
        n := new(big.Int).SetBytes(sum)
        return "-" + n.Text(16)
    }
    
    n := new(big.Int).SetBytes(sum)
    return n.Text(16)
}
```

## 令牌缓存

### 缓存结构

```go
type TokenCache struct {
    Microsoft struct {
        AccessToken  string    `json:"access_token,omitempty"`
        RefreshToken string    `json:"refresh_token,omitempty"`
        ExpiresAt    time.Time `json:"expires_at,omitempty"`
    } `json:"microsoft"`
    Minecraft struct {
        AccessToken string    `json:"access_token,omitempty"`
        ExpiresAt   time.Time `json:"expires_at,omitempty"`
        ProfileID   string    `json:"profile_id,omitempty"`
        ProfileName string    `json:"profile_name,omitempty"`
    } `json:"minecraft"`
}
```

### 缓存优先级

1. **Minecraft Access Token**: 直接使用，无需重新认证
2. **Microsoft Refresh Token**: 刷新获取新的 MS Token
3. **设备码登录**: 重新进行完整认证流程

### 缓存文件位置

```
.session/<account_id>.json
```

## 集群认证补充

- 正版认证 token 缓存按 `account_id` 持有，不再按 `player_id` 作为官方认证缓存主键。
- 同账号并发实例共享认证会话，刷新流程由 `AuthManager` 单飞协调。
- `refresh_token` 失效或需要设备码重新登录时，实例进入 `auth_failed`，不会触发自动网络重连。
- 真实微软登录仅作为手工联调项；默认自动化验证使用 fake provider。

### 缓存验证

```go
func (c *TokenCache) HasValidMinecraftToken(now time.Time) bool {
    return c.Minecraft.AccessToken != "" &&
           !c.Minecraft.ExpiresAt.IsZero() &&
           c.Minecraft.ExpiresAt.After(now.Add(5*time.Minute))
}

func (c *TokenCache) HasValidMicrosoftAccess(now time.Time) bool {
    return c.Microsoft.AccessToken != "" &&
           !c.Microsoft.ExpiresAt.IsZero() &&
           c.Microsoft.ExpiresAt.After(now.Add(5*time.Minute))
}

func (c *TokenCache) HasMicrosoftRefreshToken() bool {
    return c.Microsoft.RefreshToken != ""
}
```

## 离线模式

离线模式直接使用配置中的 `player_id` 作为用户名：

```go
// 生成离线 UUID
func offlineUUID(name string) [16]byte {
    hash := md5.Sum([]byte("OfflinePlayer:" + name))
    hash[6] = (hash[6] & 0x0F) | 0x30  // 版本 3
    hash[8] = (hash[8] & 0x3F) | 0x80  // 变体 2
    return hash
}
```

### 离线模式限制

- 无法加入正版服务器（online-mode=true）
- 无法使用皮肤
- 无法使用聊天签名功能
- 无法使用某些需要认证的功能

## 安全考虑

### 令牌存储

- 令牌缓存文件存储在本地 `.session/` 目录
- 建议不要将此目录提交到版本控制
- 敏感信息（如 refresh_token）应妥善保管

### SSL/TLS

所有 Minecraft API 请求均使用 HTTPS：

```go
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    },
}
```

### 证书处理

玩家证书存储格式：

```go
// 解析 PEM 格式私钥
func parsePrivateKeyPEM(pemData string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(pemData))
    
    // 尝试 PKCS#1
    if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
        return key, nil
    }
    
    // 尝试 PKCS#8
    if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
        return key.(*rsa.PrivateKey), nil
    }
    
    return nil, fmt.Errorf("不支持的私钥格式")
}
```

## 错误处理

### 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| 设备码超时 | 用户未在有效期内完成登录 | 重新启动认证流程 |
| 令牌过期 | Access Token 已过期 | 使用 Refresh Token 刷新 |
| 游戏所有权验证失败 | 账号未拥有 Minecraft | 购买游戏 |
| 会话加入失败 | 服务器验证失败 | 检查网络连接 |

### 重试策略

```go
// 令牌刷新失败后回退到设备码登录
if cache.HasMicrosoftRefreshToken() {
    msToken, err := RefreshMicrosoftToken(cache.Microsoft.RefreshToken)
    if err != nil {
        logx.Warnf("refresh_token 刷新失败，回退设备码授权: %v", err)
    }
}

// 完整重新认证
msToken, err := GetMicrosoftToken()
if err != nil {
    return fmt.Errorf("Microsoft 登录失败: %w", err)
}
```

## 代码位置

```
internal/
├── auth/
│   ├── microsoft/
│   │   └── service.go    # 微软/Xbox/XSTS 认证
│   └── minecraft/
│       └── service.go    # Minecraft 认证、会话加入、证书获取
└── session/
    └── cache.go          # 令牌缓存读写
```
