# gmcc API 文档

本文档面向前端对接与 API 调用方，基于当前服务端实现整理公开接口、请求响应结构、错误语义与推荐调用方式。

## 1. 基本约定

- Base URL: `http://<host>:8080`
- Base Path: `/api`
- 请求体格式: `application/json`
- 响应格式: `application/json`
- 当前公开接口默认不使用 Web 密码认证
- 当前服务以 API 为主；若未构建前端资源，访问非 `/api` 路径会返回前端不可用页面

## 2. 推荐接入流程

```text
页面初始化
  -> GET /api/status
  -> GET /api/resources
  -> GET /api/accounts
  -> GET /api/instances

创建账号
  -> POST /api/accounts

发起 Microsoft 登录
  -> POST /api/auth/microsoft/init
  -> 引导用户打开 verification_uri / verification_uri_complete
  -> POST /api/auth/microsoft/poll (按 interval 轮询直到终态)
  -> GET /api/accounts 刷新账号状态

创建实例
  -> POST /api/instances

控制实例
  -> POST /api/instances/:id/start
  -> POST /api/instances/:id/stop
  -> POST /api/instances/:id/restart
  -> DELETE /api/instances/:id

审计查询
  -> GET /api/logs/operations
```

## 3. 通用响应规则

### 3.1 读取类接口

读取类接口通常直接返回业务对象或对象列表，例如：

- `GET /api/status` 直接返回集群状态对象
- `GET /api/accounts` 返回 `{ "accounts": [...] }`
- `GET /api/instances` 返回 `{ "instances": [...] }`
- `GET /api/resources` 直接返回系统资源快照对象

### 3.2 写操作接口

写操作类接口通常返回统一结构：

```json
{
  "success": true,
  "message": "Instance started"
}
```

失败时通常返回：

```json
{
  "success": false,
  "message": "",
  "error": "instance not found"
}
```

字段说明：

- `success`: 操作是否成功
- `message`: 成功时的人类可读提示；当前大多数失败响应不会填写该字段
- `error`: 失败原因
- `operation_id`: 结构中预留字段，当前实现未实际返回

## 4. 数据模型

### 4.1 AccountView

账号公开视图。

```json
{
  "id": "acc-main",
  "player_id": "player-uuid",
  "enabled": true,
  "label": "Main",
  "note": "Primary account",
  "auth_status": "logged_in",
  "has_token": true
}
```

字段：

- `id`: 账号 ID
- `player_id`: Minecraft 玩家档案 ID；未登录时可能省略或为空
- `enabled`: 账号是否启用
- `label`: 可选展示名
- `note`: 可选备注
- `auth_status`: 当前认证状态，可能值为 `logged_in`、`not_logged_in`、`auth_invalid`
- `has_token`: 是否存在可识别的认证缓存；`logged_in` 与 `auth_invalid` 时通常为 `true`

### 4.2 InstanceView

实例公开视图。

```json
{
  "id": "bot-1",
  "account_id": "acc-main",
  "player_id": "player-uuid",
  "server_address": "mc.example.com",
  "status": "running",
  "online_duration": "3m12s",
  "last_seen": "2026-04-05T12:00:00Z",
  "has_token": true,
  "health": 20,
  "food": 20,
  "position": {
    "x": 100.5,
    "y": 65,
    "z": -24.25
  }
}
```

字段：

- `id`: 实例 ID
- `account_id`: 关联账号 ID
- `player_id`: 当前使用的玩家档案 ID；可能为空
- `server_address`: 目标 Minecraft 服务器地址
- `status`: 实例状态，当前可能值为 `pending`、`starting`、`running`、`reconnecting`、`stopped`、`error`
- `online_duration`: 在线时长字符串，由 Go `time.Duration.String()` 生成
- `last_seen`: 最后活动时间，RFC3339 时间戳；字段名固定为 `last_seen`
- `has_token`: 关联账号是否存在认证缓存
- `health`: 角色生命值；未获取到时可能省略
- `food`: 饥饿值；未获取到时可能省略
- `position`: 角色坐标；未获取到时可能省略

### 4.3 ClusterStatus

```json
{
  "cluster_status": "running",
  "total_instances": 2,
  "running_instances": 1,
  "uptime": 120000000000
}
```

字段：

- `cluster_status`: 集群总体状态
- `total_instances`: 当前实例总数
- `running_instances`: 运行中的实例数量
- `uptime`: 服务运行时长，JSON 中表现为纳秒数值

### 4.4 ResourceSnapshotView

```json
{
  "cpu_percent": 12.5,
  "memory": {
    "total_bytes": 17179869184,
    "used_bytes": 8589934592,
    "available_bytes": 7516192768,
    "used_percent": 50
  },
  "collected_at": "2026-04-05T12:00:00Z"
}
```

字段：

- `cpu_percent`: CPU 总体使用率，范围 `0-100`
- `memory.total_bytes`: 总内存字节数
- `memory.used_bytes`: 已使用内存字节数
- `memory.available_bytes`: 可用内存字节数
- `memory.used_percent`: 内存使用率，范围 `0-100`
- `collected_at`: 采集时间，UTC `RFC3339` 时间戳

### 4.5 MicrosoftAuthInitResponse

```json
{
  "success": true,
  "user_code": "ABCD-EFGH",
  "verification_uri": "https://microsoft.com/devicelogin",
  "verification_uri_complete": "https://microsoft.com/devicelogin?code=ABCD-EFGH",
  "expires_in": 900,
  "interval": 5,
  "account_id": "acc-main"
}
```

字段：

- `success`: 初始化是否成功
- `user_code`: 设备码
- `verification_uri`: 用户手动打开的验证地址
- `verification_uri_complete`: 可直接跳转的完整验证地址
- `expires_in`: 设备码剩余有效秒数
- `interval`: 建议轮询间隔，单位秒
- `account_id`: 对应账号 ID

### 4.6 MicrosoftAuthPollResponse

```json
{
  "success": true,
  "status": "succeeded",
  "message": "Microsoft authentication succeeded",
  "minecraft_profile": {
    "id": "player-uuid",
    "name": "PlayerName"
  },
  "account_id": "acc-main"
}
```

字段：

- `success`: 轮询结果是否视为成功；`pending` 与 `succeeded` 会返回 `true`
- `status`: 登录状态，当前可能值为 `pending`、`succeeded`、`expired`、`cancelled`、`failed`、`error`
- `message`: 人类可读状态说明
- `minecraft_profile`: 仅在拿到玩家档案时返回
- `account_id`: 对应账号 ID

## 5. 状态查询接口

### 5.1 获取集群状态

`GET /api/status`

成功响应：返回 `ClusterStatus`。

失败：

- `503 {"error":"cluster manager not initialized"}`

### 5.2 获取账号列表

`GET /api/accounts`

成功响应：

```json
{
  "accounts": [
    {
      "id": "acc-main",
      "player_id": "player-uuid",
      "enabled": true,
      "label": "Main",
      "note": "Primary account",
      "auth_status": "logged_in",
      "has_token": true
    }
  ]
}
```

失败：

- `503 {"error":"resource manager not initialized"}`
- `500 {"error":"failed to load accounts"}`

### 5.3 获取单个账号

`GET /api/accounts/:id`

成功响应：直接返回单个 `AccountView`，不是包在 `account` 字段里。

失败：

- `503 {"error":"resource manager not initialized"}`
- `404 {"error":"account not found"}`

### 5.4 获取实例列表

`GET /api/instances`

成功响应：

```json
{
  "instances": [
    {
      "id": "bot-1",
      "account_id": "acc-main",
      "player_id": "player-uuid",
      "server_address": "mc.example.com",
      "status": "running",
      "online_duration": "3m12s",
      "last_seen": "2026-04-05T12:00:00Z",
      "has_token": true,
      "health": 20,
      "food": 20
    }
  ]
}
```

失败：

- `503 {"error":"cluster manager not initialized"}`

### 5.5 获取单个实例

`GET /api/instances/:id`

成功响应：直接返回单个 `InstanceView`。

失败：

- `503 {"error":"cluster manager not initialized"}`
- `404 {"error":"instance not found"}`

### 5.6 获取宿主机系统资源

`GET /api/resources`

成功响应：返回 `ResourceSnapshotView`。

失败：

- `503 {"error":"resource metrics collector not initialized"}`
- `500 {"error":"failed to collect system resources"}`

## 6. 账号管理接口

### 6.1 创建账号

`POST /api/accounts`

请求：

```json
{
  "id": "acc-main",
  "label": "Main account",
  "note": "Primary operator account"
}
```

说明：

- `id` 必填
- 当前接口创建的账号默认 `enabled=true`
- 当前接口只创建账号 metadata，不会自动完成登录

成功响应：

```json
{
  "success": true,
  "message": "Account created"
}
```

失败：

- `400 {"success":false,"message":"","error":"invalid request"}`
- `400 {"success":false,"message":"","error":"account already exists"}`
- `400 {"success":false,"message":"","error":"account id is required"}`
- `503 {"success":false,"message":"","error":"resource manager not initialized"}`

### 6.2 删除账号

`DELETE /api/accounts/:id`

成功响应：

```json
{
  "success": true,
  "message": "Account deleted"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"account in use"}`
- `400 {"success":false,"message":"","error":"account not found"}`
- `503 {"success":false,"message":"","error":"resource manager not initialized"}`

前端注意：删除账号前应先处理其关联实例，否则大概率收到 `account in use`。

## 7. Microsoft 登录接口

### 7.1 初始化设备码登录

`POST /api/auth/microsoft/init`

请求：

```json
{
  "account_id": "acc-main"
}
```

成功响应：返回 `MicrosoftAuthInitResponse`。

失败语义：

- 请求体缺失、JSON 非法或缺少 `account_id` -> `401 {"success":false,"user_code":"","verification_uri":"","verification_uri_complete":"","expires_in":0,"interval":0}`
- 认证管理器未初始化 -> `503 {"success":false,"user_code":"","verification_uri":"","verification_uri_complete":"","expires_in":0,"interval":0}`
- 上游流程失败 -> `500 {"success":false,"user_code":"","verification_uri":"","verification_uri_complete":"","expires_in":0,"interval":0}`

注意：初始化失败时当前实现不会返回详细 `message` 或 `error` 字段。

### 7.2 轮询设备码登录状态

`POST /api/auth/microsoft/poll`

请求：

```json
{
  "account_id": "acc-main"
}
```

`pending` 响应示例：

```json
{
  "success": true,
  "status": "pending",
  "message": "Waiting for user authorization...",
  "account_id": "acc-main"
}
```

`succeeded` 响应示例：

```json
{
  "success": true,
  "status": "succeeded",
  "message": "Microsoft authentication succeeded",
  "minecraft_profile": {
    "id": "player-uuid",
    "name": "PlayerName"
  },
  "account_id": "acc-main"
}
```

终态失败示例：

```json
{
  "success": false,
  "status": "failed",
  "message": "Microsoft authentication failed",
  "account_id": "acc-main"
}
```

其他终态示例：

```json
{
  "success": false,
  "status": "expired",
  "message": "Device login expired",
  "account_id": "acc-main"
}
```

```json
{
  "success": false,
  "status": "cancelled",
  "message": "Device login cancelled",
  "account_id": "acc-main"
}
```

接口约定：

- 业务轮询结果无论成功、等待中还是多数终态失败，通常都返回 `200`
- 前端必须结合 `success`、`status`、`message` 判断结果
- 收到 `succeeded`、`expired`、`cancelled`、`failed` 后应停止轮询
- 建议按照 `init` 返回的 `interval` 作为轮询周期

失败语义：

- 请求体缺失、JSON 非法或缺少 `account_id` -> `401 {"success":false,"status":"error","message":"account_id required"}`
- 认证管理器未初始化 -> `503 {"success":false,"status":"error","message":"runtime auth not initialized"}`

## 8. 实例管理接口

### 8.1 创建实例

`POST /api/instances`

请求：

```json
{
  "id": "bot-1",
  "account_id": "acc-main",
  "server_address": "mc.example.com",
  "enabled": true,
  "auto_start": false
}
```

字段说明：

- `id`: 实例 ID，必填
- `account_id`: 关联账号 ID，必填
- `server_address`: Minecraft 服务器地址，必填
- `enabled`: 可选；省略时默认 `true`
- `auto_start`: 创建后是否立即启动，默认 `false`

约束：

- 当 `enabled=false` 时，`auto_start` 不能为 `true`
- 创建实例本身不会自动修正账号状态；账号未登录、账号无效等问题会直接返回错误

成功响应：

```json
{
  "success": true,
  "message": "Instance created"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"invalid request"}`
- `400 {"success":false,"message":"","error":"disabled instance cannot auto start"}`
- `400 {"success":false,"message":"","error":"instance already exists"}`
- `400 {"success":false,"message":"","error":"instance id is required"}`
- `400 {"success":false,"message":"","error":"max instances reached"}`
- `400 {"success":false,"message":"","error":"account not found"}`
- `400 {"success":false,"message":"","error":"account disabled"}`
- `400 {"success":false,"message":"","error":"account not logged in"}`
- `400 {"success":false,"message":"","error":"account auth invalid"}`
- `400 {"success":false,"message":"","error":"server address is required"}`
- `503 {"success":false,"message":"","error":"cluster manager not initialized"}`

注意：当 `auto_start=true` 且启动阶段失败时，接口仍会直接返回创建流程中的错误，不会区分“已创建但启动失败”的独立响应结构。

### 8.2 启动实例

`POST /api/instances/:id/start`

成功响应：

```json
{
  "success": true,
  "message": "Instance started"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"instance not found"}`
- `400 {"success":false,"message":"","error":"instance is already running or starting"}`
- `400 {"success":false,"message":"","error":"instance <id> already deleted"}` 或其他底层启动错误
- `503 {"success":false,"message":"","error":"cluster manager not initialized"}`

### 8.3 停止实例

`POST /api/instances/:id/stop`

成功响应：

```json
{
  "success": true,
  "message": "Instance stopped"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"instance not found"}`
- `503 {"success":false,"message":"","error":"cluster manager not initialized"}`

说明：实例已处于停止态时，当前实现会直接视为成功，不返回错误。

### 8.4 重启实例

`POST /api/instances/:id/restart`

成功响应：

```json
{
  "success": true,
  "message": "Instance restarted"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"instance not found"}`
- `400 {"success":false,"message":"","error":"instance is already running or starting"}`
- `400 {"success":false,"message":"","error":"instance <id> already deleted"}` 或其他底层重启错误
- `503 {"success":false,"message":"","error":"cluster manager not initialized"}`

### 8.5 删除实例

`DELETE /api/instances/:id`

成功响应：

```json
{
  "success": true,
  "message": "Instance deleted"
}
```

常见失败：

- `400 {"success":false,"message":"","error":"instance not found"}`
- `400 {"success":false,"message":"","error":"delete timed out waiting for instance shutdown"}`
- `503 {"success":false,"message":"","error":"cluster manager not initialized"}`

## 9. 审计日志接口

### 9.1 获取操作日志

`GET /api/logs/operations?start=<RFC3339>&end=<RFC3339>`

审计日志按天写入活动配置文件所在目录下的 `logs/audit/YYYY-MM-DD.jsonl`；如果通过 `GMCC_CONFIG` 指定配置文件，则以该文件所在目录为基准。

默认行为：

- `start` 未传时，默认最近 7 天
- `end` 未传时，默认当前时间

请求示例：

```http
GET /api/logs/operations?start=2026-04-01T00:00:00Z&end=2026-04-05T23:59:59Z
```

成功响应：

```json
{
  "logs": [
    {
      "id": "log-123",
      "timestamp": "2026-04-05T10:00:00Z",
      "action": "instance_start",
      "target_instance_id": "bot-1",
      "target_account_id": "acc-main",
      "success": true,
      "error_msg": "",
      "client_ip": "127.0.0.1",
      "user_agent": "Mozilla/5.0"
    }
  ]
}
```

字段说明：

- `id`: 日志 ID
- `timestamp`: 记录时间
- `action`: 操作类型，如 `instance_start`、`instance_create`、`account_create`、`microsoft_auth_init`
- `target_instance_id`: 目标实例 ID，可选
- `target_account_id`: 目标账号 ID，可选
- `details`: 附加说明，可选；当前多数接口未填写
- `success`: 操作是否成功
- `error_msg`: 失败原因，可选
- `client_ip`: 请求来源 IP
- `user_agent`: 请求 User-Agent，可选

失败：

- `400 {"error":"invalid start date"}`
- `400 {"error":"invalid end date"}`
- `500 {"error":"failed to query logs"}`

## 10. HTTP 状态码约定

| 状态码 | 含义 |
| --- | --- |
| `200` | 请求成功；Microsoft 轮询接口进入业务终态失败时通常也返回 `200` |
| `400` | 请求体非法、业务参数错误、资源冲突、前置条件不满足 |
| `401` | 当前仅用于 Microsoft 登录相关请求体缺失或非法 |
| `404` | 读取类接口请求的资源不存在，或 API 路径不存在 |
| `500` | 服务器内部错误或上游认证流程错误 |
| `503` | 核心依赖未初始化 |

访问不存在的 API 路径时，服务返回：

```json
{
  "success": false,
  "error": "API endpoint not found"
}
```

## 11. 对接注意事项

- `GET /api/accounts/:id` 与 `GET /api/instances/:id` 都是直接返回对象，不包裹外层字段
- `last_seen` 字段实际来自实例内部 `last_active`，前端应以返回字段名为准
- `uptime` 是数值型 duration，单位为纳秒，不是字符串
- 当前接口没有服务端分页，前端应自行控制刷新频率
- Microsoft 登录轮询不能只看 HTTP 状态码，必须读取 `success`、`status`、`message`
- `AccountView.has_token` 并不等同于 token 一定可直接使用；`auth_invalid` 状态下该值仍可能为 `true`
- 创建账号后不会自动登录，也不会自动创建实例
- 删除账号前必须先清理关联实例，否则通常会返回 `account in use`

## 12. 相关文档

- `README.md`
- `docs/auth.md`
- `docs/reference.md`
