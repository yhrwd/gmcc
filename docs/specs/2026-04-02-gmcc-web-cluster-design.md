# GMCC Web 集群管理面板设计文档

**文档版本**: v1.0  
**创建日期**: 2026-04-02  
**状态**: Design Complete  

---

## 1. 概述

### 1.1 设计目标

为 GMCC 添加 Web 管理界面，支持：
- **多账号可视化监控**：在网页端查看所有机器人状态
- **运行时操作控制**：启动/停止/删除账号实例
- **安全的 Token 存储**：本地加密存储 Minecraft 认证令牌
- **操作审计追踪**：记录所有敏感操作的执行者

### 1.2 核心特性

| 特性 | 说明 |
|------|------|
| **免登录浏览** | 进入网页即可查看状态，无需认证 |
| **操作时密码验证** | 敏感操作需输入预设密码 |
| **多密码支持** | 支持多个操作密码，日志追踪哪个密码执行的操作 |
| **环境绑定加密** | Token 使用机器码+密码派生密钥加密 |
| **前后端分离** | Vue3 + Go，前端可独立开发部署 |

### 1.3 技术栈

- **后端**: Go 1.25 + Gin 框架
- **前端**: Vue3 + TypeScript + TailwindCSS v4.2（预留目录）
- **认证**: 密码验证 + JWT Token
- **加密**: scrypt 密钥派生 + AES-256-GCM
- **通信**: RESTful API + WebSocket（可选实时状态）

---

## 2. 架构设计

### 2.1 整体架构

```
┌────────────────────────────────────────────────────────────────┐
│                     gmcc web start                              │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   HTTP Server (Gin)                      │   │
│  │                                                          │   │
│  │   ┌────────────┐  ┌────────────┐  ┌──────────────────┐   │   │
│  │   │ Public API │  │ Auth API   │  │ Protected API    │   │   │
│  │   │ (无需认证) │  │ (密码验证) │  │ (需Token/密码)   │   │   │
│  │   │            │  │            │  │                  │   │   │
│  │   │ GET /api/  │  │ POST /auth │  │ POST /api/in-  │   │   │
│  │   │ status     │  │ /verify    │  │ stances/:id/    │   │   │
│  │   │            │  │            │  │ start           │   │   │
│  │   └────────────┘  └────────────┘  └──────────────────┘   │   │
│  │                                                          │   │
│  │   ┌───────────────────────────────────────────────────┐  │   │
│  │   │              Middleware                            │  │   │
│  │   │   - CORS                                           │  │   │
│  │   │   - Static Files (web/*)                           │  │   │
│  │   │   - Optional Auth (password verification)        │  │   │
│  │   └───────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Business Layer                         │   │
│  │                                                          │   │
│  │  ┌─────────────┐  ┌───────────────┐  ┌──────────────┐  │   │
│  │  │ AuthManager │  │ ClusterManager│  │ TokenVault   │  │   │
│  │  │             │  │               │  │              │  │   │
│  │  │ -密码验证   │  │ -实例生命周期 │  │ -密钥派生    │  │   │
│  │  │ -JWT签发    │  │ -状态监控     │  │ -AES加密     │  │   │
│  │  │ -操作日志   │  │ -并发控制     │  │ -Token存取   │  │   │
│  │  └─────────────┘  └───────────────┘  └──────────────┘  │   │
│  │                                                          │   │
│  │  ┌─────────────┐  ┌───────────────┐                    │   │
│  │  │ AuditLogger │  │ KeyManager    │                    │   │
│  │  │             │  │               │                    │   │
│  │  │ -记录操作   │  │ -机器码提取   │                    │   │
│  │  │ -日志轮转   │  │ -密钥派生     │                    │   │
│  │  └─────────────┘  └───────────────┘                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Web Assets (Static)                    │   │
│  │                                                          │   │
│  │   /                     → web/index.html               │   │
│  │   /assets/*             → web/assets/*                 │   │
│  │   /api/*                → API routes                   │   │
│  └─────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────┘
```

### 2.2 组件职责

| 组件 | 职责 |
|------|------|
| **HTTP Server** | 路由分发、中间件、静态文件服务 |
| **AuthManager** | 密码验证、JWT 签发、操作日志记录 |
| **ClusterManager** | 多账号实例管理、状态聚合 |
| **TokenVault** | Token 加密存储、解密读取 |
| **KeyManager** | 机器码提取、密钥派生 |
| **AuditLogger** | 操作审计日志、日志轮转 |

---

## 3. 模块详细设计

### 3.1 认证与操作日志 (AuthManager)

#### 3.1.1 数据模型

```go
// WebAuthConfig Web认证配置
type WebAuthConfig struct {
	// 多个密码条目
	Passwords []PasswordEntry `yaml:"passwords"`
	
	// JWT Token过期时间
	TokenExpiry time.Duration `yaml:"token_expiry"`
	
	// 操作日志保留天数
	AuditLogRetentionDays int `yaml:"audit_log_retention_days"`
}

type PasswordEntry struct {
	// 密码ID（用于日志标识）
	ID string `yaml:"id"`
	
	// bcrypt哈希存储
	Hash string `yaml:"hash"`
	
	// 是否启用
	Enabled bool `yaml:"enabled"`
	
	// 创建时间
	CreatedAt time.Time `yaml:"created_at"`
	
	// 可选：备注
	Note string `yaml:"note,omitempty"`
}

// OperationLog 操作日志
type OperationLog struct {
	ID        string    `json:"id" yaml:"id"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	
	// 哪个密码执行的操作
	PasswordID string `json:"password_id" yaml:"password_id"`
	
	// 操作类型
	Action string `json:"action" yaml:"action"`
	// 可选值: "instance_start", "instance_stop", "instance_delete", 
	//         "account_create", "account_delete", "token_store"
	
	// 操作目标
	TargetInstanceID string `json:"target_instance_id,omitempty" yaml:"target_instance_id,omitempty"`
	TargetAccountID  string `json:"target_account_id,omitempty" yaml:"target_account_id,omitempty"`
	
	// 操作详情（JSON字符串）
	Details string `json:"details,omitempty" yaml:"details,omitempty"`
	
	// 操作结果
	Success  bool   `json:"success" yaml:"success"`
	ErrorMsg string `json:"error_msg,omitempty" yaml:"error_msg,omitempty"`
	
	// 客户端信息
	ClientIP  string `json:"client_ip" yaml:"client_ip"`
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
}
```

#### 3.1.2 认证流程

**密码验证**（操作时）：

```
用户输入密码
       │
       ▼
┌──────────────┐
│ POST /auth/  │  包含: {password, action, target}
│ verify       │
└──────────────┘
       │
       ▼
┌──────────────┐
│ AuthManager. │  1. 遍历所有启用的密码
│ VerifyPass-  │  2. bcrypt.CompareHashAndPassword
│ word()       │  3. 找到匹配的密码条目
└──────────────┘
       │
       ▼
┌──────────────┐
│ 验证成功     │  返回: {success: true, password_id}
│ 返回Token    │
└──────────────┘
       │
       ▼
┌──────────────┐
│ AuditLogger. │  记录: VERIFY_SUCCESS
│ Log()        │  password_id, client_ip
└──────────────┘
```

**操作审计日志**：

```go
// AuditLogger 审计日志管理器
type AuditLogger struct {
	logDir    string
	retention time.Duration
	mu        sync.Mutex
}

// Log 记录操作日志
func (l *AuditLogger) Log(log *OperationLog) error {
	// 写入 JSON Lines 格式日志文件
	// 文件: logs/audit/YYYY-MM-DD.jsonl
}

// Rotate 日志轮转（保留 N 天）
func (l *AuditLogger) Rotate() error {
	// 删除超过 retentionDays 的旧日志
}

// Query 查询日志
func (l *AuditLogger) Query(start, end time.Time, passwordID string) ([]OperationLog, error)
```

#### 3.1.3 API 设计

```go
// POST /api/auth/verify
// 验证操作密码
{
	"password": "user-input-password",
	"action": "instance_start",      // 可选，用于日志
	"target": "instance-001"         // 可选，用于日志
}

// 成功响应
{
	"success": true,
	"token": "jwt-token-here",       // 短时Token（如5分钟）
	"password_id": "admin_001",      // 哪个密码通过验证
	"expires_at": "2026-04-02T10:00:00Z"
}

// 失败响应
{
	"success": false,
	"error": "invalid_password"
}
```

### 3.2 Token Vault（加密存储）

#### 3.2.1 核心设计

Token 使用**环境绑定加密**：
- **派生密钥** = scrypt(机器码 + 操作密码, salt)
- **加密** = AES-256-GCM(Token, 派生密钥, 随机nonce)
- **存储** = JSON 文件 (权限 0600)

```go
// TokenVault Token加密存储
type TokenVault struct {
	// 存储目录
	storagePath string
	
	// scrypt 参数
	scryptN int
	scryptR int
	scryptP int
	saltLen int
}

// EncryptedToken 加密后的Token结构
type EncryptedToken struct {
	// 版本号（用于未来升级）
	Version int `json:"version"`
	
	// 加密算法
	Algorithm string `json:"algorithm"` // "aes-256-gcm"
	
	// KDF参数
	KDF     string `json:"kdf"`      // "scrypt"
	ScryptN int    `json:"scrypt_n"`
	ScryptR int    `json:"scrypt_r"`
	ScryptP int    `json:"scrypt_p"`
	Salt    []byte `json:"salt"`     // base64编码
	
	// 加密数据
	Nonce      []byte `json:"nonce"`       // base64编码
	Ciphertext []byte `json:"ciphertext"`  // base64编码
	
	// 元数据
	PlayerID  string    `json:"player_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 核心接口
type TokenVault interface {
	// Store 加密并存储Token
	// playerID: 玩家ID（用于文件名）
	// token: Minecraft Token
	// password: 操作密码（用于派生密钥）
	Store(playerID string, token *MinecraftToken, password string) error
	
	// Retrieve 解密并读取Token
	Retrieve(playerID string, password string) (*MinecraftToken, error)
	
	// Delete 删除存储的Token
	Delete(playerID string) error
	
	// Exists 检查Token是否存在
	Exists(playerID string) bool
	
	// List 列出所有存储的Token
	List() ([]string, error)
}
```

#### 3.2.2 密钥派生流程

```
操作密码 + 机器码指纹
       │
       ▼
┌──────────────┐
│ KeyManager.  │  1. 提取机器码
│ GetMachine   │     - CPU ID
│ Fingerprint  │     - 主板序列号
│ ()           │     - 磁盘UUID
└──────────────┘     - 组合哈希
       │
       ▼
┌──────────────┐
│ scrypt()     │  2. 密钥派生
│              │     N=2^20, r=8, p=1
│              │     salt = 随机32字节
│              │     key = scrypt(pwd+machine, salt, 32)
└──────────────┘
       │
       ▼
┌──────────────┐
│ AES-256-GCM  │  3. 加密Token
│              │     nonce = 随机12字节
│              │     ciphertext = AES_GCM_encrypt(token, key, nonce)
└──────────────┘
       │
       ▼
┌──────────────┐
│ 存储文件     │  4. 写入文件
│ .tokens/     │     {player_id}.enc
│ {player_id}. │     权限: 0600
│ enc          │
└──────────────┘
```

#### 3.2.3 安全特性

| 特性 | 实现 |
|------|------|
| **环境绑定** | 机器码 + 密码双重派生，文件在其他机器无法解密 |
| **随机盐值** | 每个账号独立 salt，防止彩虹表攻击 |
| **前向保密** | 即使密码泄露，无环境密钥也无法解密 |
| **文件权限** | 存储文件 0600，仅限所有者访问 |
| **内存安全** | 敏感数据使用后立即清零 |

### 3.3 Microsoft OAuth Device Flow 适配

#### 3.3.1 流程设计

Microsoft OAuth 设备码流程适配网页端：

```
┌─────────┐        ┌─────────────────┐        ┌─────────────────┐
│   Web   │        │  gmcc backend   │        │   Microsoft     │
│   UI    │        │                 │        │   OAuth         │
└────┬────┘        └────────┬────────┘        └────────┬────────┘
     │                      │                          │
     │  1. POST /api/auth/  │                          │
     │     microsoft/init   │                          │
     │──────────────────────>                          │
     │                      │                          │
     │                      │  2. 调用 Microsoft       │
     │                      │     Device Code API      │
     │                      │─────────────────────────>│
     │                      │                          │
     │                      │  3. 返回 Device Code     │
     │                      │<─────────────────────────│
     │                      │                          │
     │  4. 返回用户码和链接  │                          │
     │<──────────────────────                          │
     │                      │                          │
     │  5. 显示用户码，      │                          │
     │     引导用户打开链接  │                          │
     │                      │                          │
     │                      │                          │
     │  6. POST /api/auth/  │                          │
     │     microsoft/poll   │                          │
     │──────────────────────>                          │
     │                      │  7. 轮询 Token API       │
     │                      │─────────────────────────>│
     │                      │                          │
     │                      │  8. 等待用户授权...       │
     │                      │<─────────────────────────│
     │                      │  （用户已在浏览器授权）   │
     │                      │                          │
     │                      │  9. 获取 Token           │
     │                      │<─────────────────────────│
     │                      │                          │
     │                      │  10. 获取 Minecraft      │
     │                      │      Profile             │
     │                      │                          │
     │                      │  11. 加密存储 Token       │
     │                      │      到 Vault            │
     │                      │                          │
     │  12. 返回成功，       │                          │
     │     显示玩家信息      │                          │
     │<──────────────────────                          │
     │                      │                          │
```

#### 3.3.2 API 详情

**初始化认证**:
```yaml
POST /api/auth/microsoft/init
Body:
  {
    "password": "user-input-password"
  }
Response:
  {
    "success": true,
    "device_code": "ABC123...",
    "user_code": "ABCD-EFGH",
    "verification_uri": "https://www.microsoft.com/link",
    "verification_uri_complete": "https://www.microsoft.com/link?code=ABCD-EFGH",
    "expires_in": 900,
    "interval": 5
  }
```

**轮询 Token**:
```yaml
POST /api/auth/microsoft/poll
Body:
  {
    "password": "user-input-password",
    "device_code": "ABC123..."
  }
Response:
  {
    "success": true,
    "status": "pending",  # pending, success, expired, error
    "message": "等待用户授权...",
    
    # 当 status = "success" 时返回：
    "minecraft_profile": {
      "id": "...",
      "name": "PlayerName"
    },
    "account_id": "bot_001"  # 自动创建的账号ID
  }
```

#### 3.3.3 前端交互示例

```javascript
// 发起认证
async function startMicrosoftAuth(password) {
  const res = await fetch('/api/auth/microsoft/init', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password })
  });
  
  const data = await res.json();
  
  // 显示用户码和链接
  showAuthModal({
    userCode: data.user_code,
    verificationUrl: data.verification_uri_complete,
    deviceCode: data.device_code,
    expiresIn: data.expires_in,
    interval: data.interval
  });
  
  // 开始轮询
  startPolling(data.device_code, password, data.interval);
}

// 轮询检查
async function startPolling(deviceCode, password, interval) {
  const check = async () => {
    const res = await fetch('/api/auth/microsoft/poll', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password, device_code: deviceCode })
    });
    
    const data = await res.json();
    
    switch (data.status) {
      case 'pending':
        updateStatus('等待用户授权...');
        setTimeout(check, interval * 1000);
        break;
      case 'success':
        closeAuthModal();
        showSuccess(`认证成功！玩家: ${data.minecraft_profile.name}`);
        refreshAccountList();
        break;
      case 'expired':
        closeAuthModal();
        showError('认证已过期，请重试');
        break;
      case 'error':
        closeAuthModal();
        showError(data.message);
        break;
    }
  };
  
  check();
}
```

#### 3.3.4 错误处理

Microsoft OAuth 可能出现的错误：

| XErr 代码 | 含义 | 处理方式 |
|-----------|------|----------|
| 2148916227 | 账号已被 Xbox 封禁 | 提示用户检查账号状态 |
| 2148916233 | 账号未注册 Xbox | 提示访问 minecraft.net 创建 |
| 2148916235 | 地区不支持 | 提示使用支持的地区账号 |
| 2148916236/7 | 韩国成人验证 | 提示在 Xbox 页面完成验证 |
| 2148916238 | 未成年账户 | 提示添加家庭组成人账户 |
| 2148916262 | 未知错误 | 记录日志，提示重试 |

### 3.4 Cluster Manager（多账号管理）

复用现有 `docs/spec-cluster-manager.md` 中的设计，但扩展为支持 Web API 调用。

---

## 4. API 设计

### 4.1 公开 API（无需认证）

```yaml
# 获取集群整体状态
GET /api/status
Response:
  {
    "cluster_status": "running",    # running, stopped, error
    "total_instances": 10,
    "running_instances": 5,
    "uptime": "3h24m"
  }

# 获取账号列表（不含敏感信息）
GET /api/accounts
Response:
  {
    "accounts": [
      {
        "id": "bot_001",
        "player_id": "BotPlayer001",
        "server_address": "127.0.0.1:25565",
        "status": "running",         # pending, running, stopped, error
        "online_duration": "2h30m",
        "last_seen": "2026-04-02T09:30:00Z",
        "has_token": true             # 是否已存储Token
      }
    ]
  }

# 获取单个账号详情
GET /api/accounts/{id}
Response:
  {
    "id": "bot_001",
    "player_id": "BotPlayer001",
    "server_address": "127.0.0.1:25565",
    "status": "running",
    "online_duration": "2h30m",
    "last_seen": "2026-04-02T09:30:00Z",
    "health": 20,
    "food": 20,
    "position": {"x": 100, "y": 64, "z": -50}
  }
```

### 4.2 认证 API

```yaml
# 验证操作密码
POST /api/auth/verify
Body:
  {
    "password": "user-input-password",
    "action": "instance_start",     # 可选，用于日志
    "target": "bot_001"               # 可选，用于日志
  }
Response:
  {
    "success": true,
    "token": "jwt-token-here",
    "password_id": "admin_001",
    "expires_at": "2026-04-02T10:00:00Z"
  }
```

### 4.3 受保护 API（需要密码验证）

```yaml
# 启动实例
POST /api/instances/{id}/start
Headers:
  Authorization: Bearer {token}       # 或通过Body传password
Body:
  {
    "password": "user-input-password"  # 备选方式
  }
Response:
  {
    "success": true,
    "message": "Instance started",
    "operation_id": "op-uuid-here"
  }

# 停止实例
POST /api/instances/{id}/stop
Body: { "password": "..." }

# 重启实例
POST /api/instances/{id}/restart
Body: { "password": "..." }

# 删除实例
DELETE /api/instances/{id}
Body: { "password": "..." }

# 添加新账号（离线模式 - 需要密码验证）
POST /api/accounts
Body:
  {
    "password": "user-input-password",    # 必须，验证操作权限
    "id": "bot_new",
    "player_id": "NewBot",
    "server_address": "127.0.0.1:25565",
    "use_official_auth": false
  }

# 添加新账号（正版模式 - 需要密码验证）
# 注意：正版认证包含两个步骤
# 1. POST /api/auth/microsoft/init - 获取用户码
# 2. POST /api/auth/microsoft/poll   - 轮询授权结果
# 两个步骤都需要 password 参数

# 发起正版认证流程（Microsoft OAuth Device Flow）
POST /api/auth/microsoft/init
Body:
  {
    "password": "user-input-password"
  }
Response:
  {
    "success": true,
    "device_code": "ABC123...",
    "user_code": "ABCD-EFGH",
    "verification_uri": "https://www.microsoft.com/link",
    "verification_uri_complete": "https://www.microsoft.com/link?code=ABCD-EFGH",
    "expires_in": 900,
    "interval": 5
  }

# 轮询获取微软 Token（前端调用）
POST /api/auth/microsoft/poll
Body:
  {
    "password": "user-input-password",
    "device_code": "ABC123..."
  }
Response:
  {
    "success": true,
    "status": "pending",  # pending, success, expired, error
    "message": "等待用户授权...",
    
    # 当 status = "success" 时返回：
    "minecraft_profile": {
      "id": "...",
      "name": "PlayerName"
    },
    "account_id": "bot_001"  # 自动创建的账号ID
  }

# 存储 Token 到加密 Vault（后端自动完成，前端无需调用）
POST /api/auth/microsoft/store
Body:
  {
    "password": "user-input-password",
    "microsoft_token": "...",
    "minecraft_token": "...",
    "profile_id": "...",
    "profile_name": "PlayerName"
  }
```

### 4.4 受保护操作汇总

以下操作都需要密码验证（`password` 参数）：

| 操作 | API | 操作类型 |
|------|-----|----------|
| 启动实例 | `POST /api/instances/{id}/start` | instance_start |
| 停止实例 | `POST /api/instances/{id}/stop` | instance_stop |
| 重启实例 | `POST /api/instances/{id}/restart` | instance_restart |
| 删除实例 | `DELETE /api/instances/{id}` | instance_delete |
| **添加离线账号** | `POST /api/accounts` | **account_create** |
| **添加正版账号** | `POST /api/auth/microsoft/init` | **account_create** |
| 删除账号 | `DELETE /api/accounts/{id}` | account_delete |

**注意**：
- 添加账号（无论是离线还是正版）都需要密码验证
- 正版认证流程中，`/auth/microsoft/init` 和 `/auth/microsoft/poll` 都需要密码
- 密码用于：1) 验证操作权限 2) 派生加密密钥存储 Token

### 4.5 日志 API（可选）

```yaml
# 查询操作日志
GET /api/logs/operations?start=2026-04-01&end=2026-04-02&password_id=admin_001
Response:
  {
    "logs": [
      {
        "id": "log-001",
        "timestamp": "2026-04-02T09:30:00Z",
        "password_id": "admin_001",
        "action": "instance_start",
        "target_instance_id": "bot_001",
        "success": true,
        "client_ip": "192.168.1.100"
      }
    ]
  }
```

---

## 5. 前端架构

### 5.1 目录结构

```
web/                              # 前端目录
├── index.html                    # Phase 1: 原生HTML示例（立即可用）
├── assets/                       # 静态资源
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   └── app.js               # Phase 1: 原生JS实现
│   └── img/
├── src/                          # Phase 2: Vue3源码（预留）
│   ├── components/               # 可复用组件
│   │   ├── AccountCard.vue
│   │   ├── StatusBadge.vue
│   │   ├── PasswordModal.vue
│   │   └── LogViewer.vue
│   ├── views/                    # 页面视图
│   │   ├── Dashboard.vue        # 主面板
│   │   ├── AccountDetail.vue    # 账号详情
│   │   └── Settings.vue         # 设置页面
│   ├── api/                      # API 封装
│   │   └── client.ts
│   ├── stores/                   # Pinia 状态管理
│   │   ├── auth.ts
│   │   └── cluster.ts
│   ├── types/                    # TypeScript 类型
│   │   └── index.ts
│   ├── App.vue
│   ├── main.ts
│   └── vite-env.d.ts
├── index.html                    # Vue 入口（Phase 2）
├── package.json                  # Vue 依赖（Phase 2）
├── vite.config.ts               # Vite 配置（Phase 2）
├── tsconfig.json                # TypeScript 配置（Phase 2）
├── tailwind.config.ts           # Tailwind 配置（Phase 2）
└── README.md                     # 前端开发指南
```

### 5.2 Phase 1: 原生 HTML 实现

```html
<!-- web/index.html -->
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GMCC 集群管理面板</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#3b82f6',
                        success: '#10b981',
                        warning: '#f59e0b',
                        danger: '#ef4444',
                    }
                }
            }
        }
    </script>
</head>
<body class="bg-gray-50 min-h-screen">
    <!-- 导航栏 -->
    <nav class="bg-white shadow-sm border-b">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex items-center">
                    <h1 class="text-xl font-bold text-gray-900">GMCC 集群管理</h1>
                </div>
                <div class="flex items-center space-x-4">
                    <span id="cluster-status-badge" class="px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                        运行中
                    </span>
                </div>
            </div>
        </div>
    </nav>

    <!-- 主内容 -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- 统计卡片 -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div class="bg-white rounded-lg shadow p-6">
                <div class="text-sm font-medium text-gray-500">总实例</div>
                <div id="stat-total" class="mt-2 text-3xl font-bold text-gray-900">-</div>
            </div>
            <div class="bg-white rounded-lg shadow p-6">
                <div class="text-sm font-medium text-gray-500">运行中</div>
                <div id="stat-running" class="mt-2 text-3xl font-bold text-green-600">-</div>
            </div>
            <div class="bg-white rounded-lg shadow p-6">
                <div class="text-sm font-medium text-gray-500">已停止</div>
                <div id="stat-stopped" class="mt-2 text-3xl font-bold text-gray-600">-</div>
            </div>
            <div class="bg-white rounded-lg shadow p-6">
                <div class="text-sm font-medium text-gray-500">运行时长</div>
                <div id="stat-uptime" class="mt-2 text-3xl font-bold text-primary">-</div>
            </div>
        </div>

        <!-- 账号列表 -->
        <div class="bg-white rounded-lg shadow">
            <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-medium text-gray-900">账号列表</h2>
            </div>
            <div id="accounts-list" class="divide-y divide-gray-200">
                <!-- 动态加载 -->
            </div>
        </div>
    </main>

    <!-- 密码输入 Modal -->
    <div id="password-modal" class="fixed inset-0 bg-gray-500 bg-opacity-75 hidden">
        <div class="flex items-center justify-center min-h-screen">
            <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
                <div class="px-6 py-4 border-b">
                    <h3 class="text-lg font-medium text-gray-900">验证操作密码</h3>
                </div>
                <div class="px-6 py-4">
                    <p class="text-sm text-gray-600 mb-4">请输入操作密码以执行敏感操作</p>
                    <input type="password" id="password-input" 
                           class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                           placeholder="输入密码">
                    <p id="password-error" class="mt-2 text-sm text-red-600 hidden"></p>
                </div>
                <div class="px-6 py-4 border-t flex justify-end space-x-3">
                    <button onclick="closePasswordModal()" class="px-4 py-2 text-gray-700 hover:text-gray-900">
                        取消
                    </button>
                    <button onclick="submitPassword()" class="px-4 py-2 bg-primary text-white rounded-md hover:bg-blue-600">
                        确认
                    </button>
                </div>
            </div>
        </div>
    </div>

    <script src="assets/js/app.js"></script>
</body>
</html>
```

```javascript
// web/assets/js/app.js - 核心逻辑
const API_BASE = '/api';

// 状态管理
let currentAction = null;
let currentInstance = null;
let authToken = null;

// 初始化
async function init() {
    await loadStatus();
    await loadAccounts();
    setInterval(() => {
        loadStatus();
        loadAccounts();
    }, 5000); // 5秒刷新
}

// 加载集群状态
async function loadStatus() {
    try {
        const res = await fetch(`${API_BASE}/status`);
        const data = await res.json();
        updateStatusUI(data);
    } catch (err) {
        console.error('加载状态失败:', err);
    }
}

// 加载账号列表
async function loadAccounts() {
    try {
        const res = await fetch(`${API_BASE}/accounts`);
        const data = await res.json();
        renderAccounts(data.accounts);
    } catch (err) {
        console.error('加载账号失败:', err);
    }
}

// 渲染账号列表
function renderAccounts(accounts) {
    const container = document.getElementById('accounts-list');
    container.innerHTML = accounts.map(account => `
        <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50">
            <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                    <div class="h-10 w-10 rounded-full ${getStatusColor(account.status)} flex items-center justify-center">
                        <span class="text-white font-bold">${account.player_id[0]}</span>
                    </div>
                </div>
                <div>
                    <div class="text-sm font-medium text-gray-900">${account.player_id}</div>
                    <div class="text-sm text-gray-500">${account.server_address}</div>
                </div>
            </div>
            <div class="flex items-center space-x-4">
                <span class="px-2 py-1 text-xs font-medium rounded-full ${getStatusBadgeClass(account.status)}">
                    ${getStatusText(account.status)}
                </span>
                ${account.status === 'stopped' ? `
                    <button onclick="startInstance('${account.id}')" 
                            class="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700">
                        启动
                    </button>
                ` : `
                    <button onclick="stopInstance('${account.id}')" 
                            class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700">
                        停止
                    </button>
                `}
            </div>
        </div>
    `).join('');
}

// 启动实例
async function startInstance(instanceId) {
    currentAction = 'start';
    currentInstance = instanceId;
    showPasswordModal();
}

// 停止实例
async function stopInstance(instanceId) {
    currentAction = 'stop';
    currentInstance = instanceId;
    showPasswordModal();
}

// 显示密码 Modal
function showPasswordModal() {
    document.getElementById('password-modal').classList.remove('hidden');
    document.getElementById('password-input').focus();
}

// 关闭密码 Modal
function closePasswordModal() {
    document.getElementById('password-modal').classList.add('hidden');
    document.getElementById('password-input').value = '';
    document.getElementById('password-error').classList.add('hidden');
}

// 提交密码
async function submitPassword() {
    const password = document.getElementById('password-input').value;
    if (!password) return;
    
    try {
        // 验证密码并执行操作
        const res = await fetch(`${API_BASE}/instances/${currentInstance}/${currentAction}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password })
        });
        
        const data = await res.json();
        
        if (data.success) {
            closePasswordModal();
            loadAccounts(); // 刷新列表
        } else {
            document.getElementById('password-error').textContent = data.error || '操作失败';
            document.getElementById('password-error').classList.remove('hidden');
        }
    } catch (err) {
        document.getElementById('password-error').textContent = '网络错误';
        document.getElementById('password-error').classList.remove('hidden');
    }
}

// 工具函数
function getStatusColor(status) {
    const colors = {
        running: 'bg-green-500',
        stopped: 'bg-gray-400',
        error: 'bg-red-500',
        pending: 'bg-yellow-500'
    };
    return colors[status] || 'bg-gray-400';
}

function getStatusBadgeClass(status) {
    const classes = {
        running: 'bg-green-100 text-green-800',
        stopped: 'bg-gray-100 text-gray-800',
        error: 'bg-red-100 text-red-800',
        pending: 'bg-yellow-100 text-yellow-800'
    };
    return classes[status] || 'bg-gray-100 text-gray-800';
}

function getStatusText(status) {
    const texts = {
        running: '运行中',
        stopped: '已停止',
        error: '错误',
        pending: '启动中'
    };
    return texts[status] || status;
}

// 初始化
init();
```

---

## 6. 配置示例

### 6.1 主配置文件

```yaml
# web.yaml - Web 面板配置
web:
  # 监听地址
  bind: "0.0.0.0:8080"
  
  # 静态文件目录（相对路径或绝对路径）
  static_path: "./web"
  
  # 认证配置
  auth:
    # 多个操作密码
    passwords:
      - id: "admin_001"
        hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"  # bcrypt of "password123"
        enabled: true
        note: "管理员主密码"
        created_at: "2026-04-02T00:00:00Z"
      
      - id: "admin_002"
        hash: "$2a$10$..."
        enabled: true
        note: "备用密码"
        created_at: "2026-04-02T00:00:00Z"
    
    # Token 过期时间（验证后的短时Token）
    token_expiry: "5m"
    
    # 操作日志保留天数
    audit_log_retention_days: 30
  
  # Token Vault 配置
  token_vault:
    # 存储路径
    storage_path: ".tokens"
    
    # scrypt 参数（默认值，通常不需要修改）
    scrypt_n: 1048576  # 2^20
    scrypt_r: 8
    scrypt_p: 1
    salt_len: 32
  
  # CORS 配置（开发用）
  cors:
    enabled: true
    origins:
      - "http://localhost:5173"  # Vue dev server
      - "http://localhost:3000"
  
  # API 限流（可选）
  rate_limit:
    enabled: false
    requests_per_minute: 60

# 继承现有的集群配置
cluster:
  global:
    max_instances: 10
    reconnect_policy:
      enabled: true
      max_retries: 5
      base_delay: 5s
      max_delay: 300s
      multiplier: 2.0
      
  accounts:
    - id: bot_001
      player_id: "BotPlayer001"
      use_official_auth: false
      server_address: "127.0.0.1:25565"
      enabled: true
      
  log:
    log_dir: "logs"
    max_size: 512
    debug: false
    enable_file: true
```

### 6.2 命令行使用

```bash
# 启动 Web 面板（默认配置 web.yaml）
gmcc web

# 指定端口
gmcc web --port 8080

# 指定配置文件
gmcc web --config web.yaml

# 指定静态文件目录（用于开发）
gmcc web --static ./web-dev

# 生成密码哈希（用于配置）
gmcc web hash-password
```

---

## 7. 安全考虑

### 7.1 Token 加密安全

- **环境绑定**：Token 文件绑定到特定机器，无法在其他机器解密
- **密码派生**：即使 Token 文件泄露，无正确密码也无法解密
- **内存安全**：敏感数据使用后立即清零，避免内存泄漏
- **文件权限**：存储目录和文件权限设置为 0700/0600

### 7.2 Web 安全

- **密码验证**：所有敏感操作需要密码验证
- **操作审计**：完整记录操作日志，包括操作者、时间、IP
- **HTTPS 支持**：建议生产环境使用 HTTPS（可配置证书）
- **CORS 限制**：可配置允许的源地址

### 7.3 风险缓解

| 风险 | 缓解措施 |
|------|----------|
| 密码泄露 | bcrypt 哈希存储，暴力破解困难 |
| Token 文件被盗 | 环境绑定加密，其他机器无法使用 |
| 中间人攻击 | 建议使用 HTTPS |
| 未授权访问 | 敏感操作强制密码验证 |

---

## 8. 实现阶段

### Phase 1: 核心功能（MVP）

1. **Backend**
   - [ ] HTTP Server (Gin) 基础框架
   - [ ] Static Files 服务
   - [ ] 公开 API: /api/status, /api/accounts
   - [ ] AuthManager: 密码验证、bcrypt
   - [ ] Protected API: 启动/停止实例
   - [ ] TokenVault: scrypt + AES-GCM 加密
   - [ ] KeyManager: 机器码提取
   - [ ] AuditLogger: 操作日志记录

2. **Frontend**
   - [ ] 原生 HTML/CSS/JS 实现
   - [ ] 状态监控页面
   - [ ] 账号列表展示
   - [ ] 密码输入 Modal
   - [ ] 启动/停止操作

### Phase 2: Vue3 前端（后续开发）

1. **Setup**
   - [ ] Vue3 + Vite 初始化
   - [ ] TypeScript 配置
   - [ ] TailwindCSS v4.2 配置
   - [ ] Pinia 状态管理

2. **Components**
   - [ ] Dashboard 面板
   - [ ] AccountCard 组件
   - [ ] Real-time 状态更新 (WebSocket)
   - [ ] LogViewer 日志查看

3. **Features**
   - [ ] 添加新账号向导
   - [ ] Token 管理页面
   - [ ] 操作日志查询
   - [ ] 设置页面

### Phase 3: 增强功能（可选）

- [ ] WebSocket 实时推送状态
- [ ] 操作日志可视化
- [ ] 多语言支持
- [ ] 移动端适配优化
- [ ] API 文档自动生成

---

## 9. 数据结构完整定义

```go
// internal/web/config.go

package web

import "time"

// WebConfig Web面板配置
type WebConfig struct {
	Bind       string          `yaml:"bind"`
	StaticPath string          `yaml:"static_path"`
	Auth       AuthConfig      `yaml:"auth"`
	TokenVault TokenVaultConfig `yaml:"token_vault"`
	CORS       CORSConfig      `yaml:"cors"`
	RateLimit  RateLimitConfig `yaml:"rate_limit"`
}

type AuthConfig struct {
	Passwords             []PasswordEntry `yaml:"passwords"`
	TokenExpiry           time.Duration   `yaml:"token_expiry"`
	AuditLogRetentionDays int             `yaml:"audit_log_retention_days"`
}

type PasswordEntry struct {
	ID        string    `yaml:"id"`
	Hash      string    `yaml:"hash"`
	Enabled   bool      `yaml:"enabled"`
	CreatedAt time.Time `yaml:"created_at"`
	Note      string    `yaml:"note,omitempty"`
}

type TokenVaultConfig struct {
	StoragePath string `yaml:"storage_path"`
	ScryptN     int    `yaml:"scrypt_n"`
	ScryptR     int    `yaml:"scrypt_r"`
	ScryptP     int    `yaml:"scrypt_p"`
	SaltLen     int    `yaml:"salt_len"`
}

type CORSConfig struct {
	Enabled bool     `yaml:"enabled"`
	Origins []string `yaml:"origins"`
}

type RateLimitConfig struct {
	Enabled             bool `yaml:"enabled"`
	RequestsPerMinute   int  `yaml:"requests_per_minute"`
}

// internal/web/models.go

// AccountView 账号展示模型（公开）
type AccountView struct {
	ID               string    `json:"id"`
	PlayerID         string    `json:"player_id"`
	ServerAddress    string    `json:"server_address"`
	Status           string    `json:"status"`
	OnlineDuration   string    `json:"online_duration"`
	LastSeen         time.Time `json:"last_seen"`
	HasToken         bool      `json:"has_token"`
	Health           float32   `json:"health,omitempty"`
	Food             int32     `json:"food,omitempty"`
	Position         *Position `json:"position,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// ClusterStatus 集群状态
type ClusterStatus struct {
	Status           string        `json:"cluster_status"`
	TotalInstances   int           `json:"total_instances"`
	RunningInstances int           `json:"running_instances"`
	Uptime           time.Duration `json:"uptime"`
}

// AuthVerifyRequest 密码验证请求
type AuthVerifyRequest struct {
	Password string `json:"password" binding:"required"`
	Action   string `json:"action,omitempty"`   // 可选
	Target   string `json:"target,omitempty"`     // 可选
}

// AuthVerifyResponse 密码验证响应
type AuthVerifyResponse struct {
	Success    bool      `json:"success"`
	Token      string    `json:"token,omitempty"`
	PasswordID string    `json:"password_id,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	Error      string    `json:"error,omitempty"`
}

// OperationRequest 操作请求（受保护API）
type OperationRequest struct {
	Password string `json:"password" binding:"required"`
}

// OperationResponse 操作响应
type OperationResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	OperationID string `json:"operation_id,omitempty"`
	Error       string `json:"error,omitempty"`
}
```

---

## 10. 附录

### A. 生成密码哈希

```go
// 用于生成配置文件中的密码哈希
func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
```

### B. 机器码提取

```go
// 提取机器唯一标识用于密钥派生
func GetMachineFingerprint() ([]byte, error) {
	// 组合多个硬件标识
	var components []string
	
	// CPU 信息
	if cpuInfo, err := getCPUInfo(); err == nil {
		components = append(components, cpuInfo)
	}
	
	// 磁盘 UUID
	if diskUUID, err := getDiskUUID(); err == nil {
		components = append(components, diskUUID)
	}
	
	// 主板序列号
	if boardSerial, err := getBoardSerial(); err == nil {
		components = append(components, boardSerial)
	}
	
	// 组合并哈希
	combined := strings.Join(components, "|")
	hash := sha256.Sum256([]byte(combined))
	return hash[:], nil
}
```

---

**文档结束**

**审批状态**: 待批准  
**下一步**: 用户审查 → 创建实现计划