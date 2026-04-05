# gmcc API 文档

本文档面向前端对接，覆盖当前公开 API、请求响应结构、推荐调用流程以及常见错误语义。

## 1. 基本约定

- Base URL: `http://<host>:8080`
- Base Path: `/api`
- 请求体格式: `application/json`
- 响应格式: `application/json`
- 当前接口默认不使用 Web 密码认证

## 2. 推荐接入流程

```text
加载首页
  -> GET /api/status
  -> GET /api/accounts
  -> GET /api/instances

新增账号
  -> POST /api/accounts

账号登录
  -> POST /api/auth/microsoft/init
  -> POST /api/auth/microsoft/poll (循环直到 succeeded)
  -> GET /api/accounts 刷新状态

创建实例
  -> POST /api/instances

控制实例
  -> POST /api/instances/:id/start
  -> POST /api/instances/:id/stop
  -> POST /api/instances/:id/restart

查看日志
  -> GET /api/logs/operations
```

## 3. 状态查询接口

### 3.1 获取集群状态

`GET /api/status`

响应示例：

```json
{
  "cluster_status": "running",
  "total_instances": 2,
  "running_instances": 2,
  "uptime": 120000000000
}
```

字段：

- `cluster_status`: `running` / `partial` / `stopped`
- `total_instances`: 当前实例总数
- `running_instances`: 当前运行中的实例数
- `uptime`: 服务启动到现在的持续时间，单位纳秒

### 3.2 获取账号列表

`GET /api/accounts`

响应示例：

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

### 3.3 获取单个账号

`GET /api/accounts/:id`

响应结构与账号列表中的单项一致。

错误：

- `404 {"error":"account not found"}`

### 3.4 获取实例列表

`GET /api/instances`

响应示例：

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
      "food": 20,
      "position": {
        "x": 100.5,
        "y": 65,
        "z": -24.25
      }
    }
  ]
}
```

### 3.5 获取单个实例

`GET /api/instances/:id`

响应结构与实例列表中的单项一致。

错误：

- `404 {"error":"instance not found"}`

## 4. 账号管理接口

### 4.1 创建账号

`POST /api/accounts`

请求：

```json
{
  "id": "acc-main",
  "label": "Main account",
  "note": "Primary operator account"
}
```

成功响应：

```json
{
  "success": true,
  "message": "Account created"
}
```

失败场景：

- 请求体非法 -> `400`
- 账号已存在 -> `400`
- 资源管理器未初始化 -> `503`

前端注意：创建账号只会写 metadata，不会自动登录。

### 4.2 删除账号

`DELETE /api/accounts/:id`

成功响应：

```json
{
  "success": true,
  "message": "Account deleted"
}
```

典型失败：

```json
{
  "success": false,
  "message": "",
  "error": "account in use"
}
```

前端建议：删除前先检查该账号下是否还有实例。

## 5. Microsoft 登录接口

### 5.1 初始化设备码登录

`POST /api/auth/microsoft/init`

请求：

```json
{
  "account_id": "acc-main"
}
```

成功响应：

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

失败：

- 参数缺失或 JSON 非法 -> `401`
- 服务未初始化 -> `503`
- 上游 provider 错误 -> `500`

### 5.2 轮询设备码登录状态

`POST /api/auth/microsoft/poll`

请求：

```json
{
  "account_id": "acc-main"
}
```

响应示例：

```json
{
  "success": true,
  "status": "pending",
  "message": "Waiting for user authorization...",
  "account_id": "acc-main"
}
```

成功完成：

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

终态失败：

```json
{
  "success": false,
  "status": "failed",
  "message": "Microsoft authentication failed",
  "account_id": "acc-main"
}
```

轮询建议：

- 使用 `init` 返回的 `interval` 作为轮询周期
- 收到 `succeeded` / `expired` / `cancelled` / `failed` 后停止轮询
- `succeeded` 后刷新账号列表

## 6. 实例管理接口

### 6.1 创建实例

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

说明：

- `id`: 实例 ID
- `account_id`: 关联账号 ID
- `server_address`: Minecraft 服务器地址
- `enabled`: 可选；省略时默认 `true`，显式传 `false` 会创建为禁用实例
- `auto_start`: 创建后是否立即启动；当 `enabled=false` 时不能同时为 `true`

成功响应：

```json
{
  "success": true,
  "message": "Instance created"
}
```

典型失败：

```json
{
  "success": false,
  "message": "",
  "error": "account not logged in"
}
```

可能错误语义：

- `account not found`
- `account disabled`
- `account not logged in`
- `account auth invalid`
- `disabled instance cannot auto start`
- `instance already exists`

### 6.2 启动实例

`POST /api/instances/:id/start`

成功响应：

```json
{
  "success": true,
  "message": "Instance started"
}
```

### 6.3 停止实例

`POST /api/instances/:id/stop`

成功响应：

```json
{
  "success": true,
  "message": "Instance stopped"
}
```

### 6.4 重启实例

`POST /api/instances/:id/restart`

成功响应：

```json
{
  "success": true,
  "message": "Instance restarted"
}
```

### 6.5 删除实例

`DELETE /api/instances/:id`

成功响应：

```json
{
  "success": true,
  "message": "Instance deleted"
}
```

## 7. 操作日志接口

### 获取操作日志

`GET /api/logs/operations?start=<RFC3339>&end=<RFC3339>`

审计日志会按天持久化到活动配置文件所在目录下的 `logs/audit/YYYY-MM-DD.jsonl`；如果通过 `GMCC_CONFIG` 指定了配置文件，则以该文件所在目录为基准。

如果不传：

- `start` 默认最近 7 天
- `end` 默认当前时间

请求示例：

```http
GET /api/logs/operations?start=2026-04-01T00:00:00Z&end=2026-04-05T23:59:59Z
```

响应示例：

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
      "client_ip": "127.0.0.1",
      "user_agent": "Mozilla/5.0"
    }
  ]
}
```

## 8. HTTP 状态码约定

| 状态码 | 含义 |
|---|---|
| `200` | 请求成功，或轮询接口返回业务终态 |
| `400` | 业务参数错误、资源冲突、前置条件不满足 |
| `401` | 某些认证请求体缺失或非法 |
| `404` | 资源不存在 |
| `500` | 服务器内部错误或上游 provider 错误 |
| `503` | 核心依赖未初始化 |

## 9. 前端页面拆分建议

### 首页 Dashboard

初始化拉取：

- `GET /api/status`
- `GET /api/accounts`
- `GET /api/instances`

### 账号管理页

操作流：

- 创建账号 -> 登录账号 -> 刷新列表

### 实例管理页

操作流：

- 过滤 `logged_in` 账号 -> 创建实例 -> 启动/停止/重启 -> 刷新实例列表

### 审计页

操作流：

- 选择时间范围 -> 请求 `GET /api/logs/operations`

## 10. 对接注意事项

- 当前 `uptime` 字段是 duration 数值，不是字符串
- 当前接口没有服务端分页，前端应自行控制刷新频率
- 登录流程完成后不要假设实例自动创建，需要显式调用创建实例接口
- 删除账号前必须先处理关联实例，否则会返回 `account in use`

## 11. 相关文档

- `README.md`
- `docs/auth.md`
- `docs/reference.md`
