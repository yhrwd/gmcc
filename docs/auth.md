# 认证与账号状态文档

本文档面向前端和后端对接，说明 gmcc 当前的账号认证模型、Microsoft 设备码登录流程、账号状态语义，以及实例创建前需要满足的条件。

## 1. 认证模型概览

gmcc 将账号资源和认证凭据分开管理：

- 账号 metadata 存在 `.state/accounts.yaml`
- 认证凭据存在 `auth.vault.path` 指向的加密目录中
- 实例 metadata 存在 `.state/instances.yaml`

前端可以把它理解为三类对象：

- `account`: 业务资源，表示系统中定义了一个可用账号
- `auth session`: 该账号当前是否已登录，以及是否可用于创建实例
- `instance`: 使用某个账号实际连接 Minecraft 服务器的运行单元

## 2. 账号认证状态

`GET /api/accounts` 和 `GET /api/accounts/:id` 会返回以下字段：

- `id`: 账号 ID
- `player_id`: 已登录时返回的 Minecraft Profile ID
- `enabled`: 账号是否启用
- `label`: 前端展示标签
- `note`: 备注
- `auth_status`: 认证状态
- `has_token`: 是否存在可恢复或异常待处理的凭据记录

### `auth_status` 枚举

| 值 | 含义 | 前端建议 |
|---|---|---|
| `not_logged_in` | 账号尚未登录，或没有可用认证记录 | 引导用户发起 Microsoft 登录 |
| `logged_in` | 账号已有可用刷新凭据，可用于创建实例 | 允许创建实例 |
| `auth_invalid` | 存在认证记录，但记录损坏、过期或需要重新登录 | 提示用户重新登录 |

### `has_token` 的含义

`has_token` 不是“当前 access token 是否未过期”的精确表达，而是“系统中是否存在该账号的认证记录”。

当前实现下：

- `logged_in` -> `has_token=true`
- `auth_invalid` -> `has_token=true`
- `not_logged_in` -> `has_token=false`

前端做按钮可用性判断时，应以 `auth_status` 为主，不要只看 `has_token`。

## 3. Microsoft 设备码登录流程

当前公开的认证 API 有两个：

- `POST /api/auth/microsoft/init`
- `POST /api/auth/microsoft/poll`

这是标准的设备码授权模式，适合前端后台分离场景。

### 流程图

```text
前端                    gmcc                     Microsoft / Minecraft
 |                        |                                |
 | POST /auth/init        |                                |
 |----------------------->| 申请 device code               |
 |                        |------------------------------->|
 |                        | 返回 user_code / verify_url    |
 |<-----------------------|                                |
 | 展示验证码和跳转链接    |                                |
 | 用户浏览器完成授权      |                                |
 | POST /auth/poll        | 查询登录状态                    |
 |----------------------->|------------------------------->|
 |<-----------------------| pending / succeeded / failed   |
 |                        |                                |
 | succeeded 后拿到 profile |                                |
```

### 3.1 初始化登录

请求：

```http
POST /api/auth/microsoft/init
Content-Type: application/json

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

字段说明：

- `user_code`: 需要展示给用户输入的设备码
- `verification_uri`: 用户可访问的授权地址
- `verification_uri_complete`: 可直接跳转的完整地址
- `expires_in`: 设备码剩余秒数
- `interval`: 建议轮询间隔，单位秒
- `account_id`: 绑定的账号 ID

前端建议：

- 弹出登录对话框
- 高亮显示 `user_code`
- 提供“复制验证码”和“打开授权页”按钮
- 使用 `interval` 作为轮询周期，避免过快轮询

### 3.2 轮询登录状态

请求：

```http
POST /api/auth/microsoft/poll
Content-Type: application/json

{
  "account_id": "acc-main"
}
```

处理中响应：

```json
{
  "success": true,
  "status": "pending",
  "message": "Waiting for user authorization...",
  "account_id": "acc-main"
}
```

成功响应：

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

失败或结束响应：

```json
{
  "success": false,
  "status": "expired",
  "message": "Device login expired",
  "account_id": "acc-main"
}
```

### `status` 枚举

| 值 | 含义 | 前端行为 |
|---|---|---|
| `pending` | 用户尚未完成授权 | 保持轮询 |
| `succeeded` | 登录成功，凭据已写入 vault | 结束轮询，刷新账号列表 |
| `expired` | 设备码过期 | 提示重新发起登录 |
| `cancelled` | 登录流程已取消 | 关闭弹窗或允许重试 |
| `failed` | 登录失败 | 展示错误并允许重新发起 |

补充说明：当请求体缺失或运行时认证管理器未初始化时，接口还会返回 `status="error"`，并分别配合 `401` 或 `503` 状态码。

## 4. 认证数据存储

当前实现不再使用旧的 `.session/*.json` 明文缓存模型，而是使用每账号一个加密文件的 vault 模型。

特点：

- 每个账号单独一个 `.vault` 文件
- 所有文件共享同一把主密钥
- 主密钥从环境变量读取，不写入 `config.yaml`
- 文件名带账号可读前缀和短 hash，便于排查但不直接暴露全部信息

启动服务前必须设置：

- `GMCC_AUTH_VAULT_KEY`，或配置中 `auth.vault.key_env` 指向的变量名

## 5. 实例创建前置条件

实例创建不是登录接口。当前模型要求：

1. 账号已存在于账号 metadata 中
2. 账号 `enabled=true`
3. 账号认证状态为 `logged_in`

如果不满足条件，实例创建会失败。

前端推荐流程：

1. 先创建账号
2. 再做 Microsoft 登录
3. 登录成功后刷新账号列表
4. 仅对 `auth_status=logged_in` 的账号开放实例创建

## 6. 前端页面建议

### 账号列表页

建议展示：

- 账号 ID
- 标签与备注
- 启用状态
- 认证状态
- 对应操作按钮：登录、重新登录、删除、创建实例

按钮建议：

- `not_logged_in`: 显示“登录”
- `logged_in`: 显示“创建实例”
- `auth_invalid`: 显示“重新登录”

### 登录弹窗

建议展示：

- 账号 ID
- `user_code`
- 授权地址
- 倒计时（基于 `expires_in`）
- 当前轮询状态文案

## 7. 常见对接问题

### 为什么账号已存在，但不能创建实例？

因为账号资源存在不等于账号已登录。账号只是 metadata，实例创建还要求认证状态为 `logged_in`。

### 为什么 `has_token=true` 还需要重新登录？

因为 `auth_invalid` 也可能带有旧记录。前端应以 `auth_status` 为准。

### 登录成功后要不要重新拉账号列表？

要。登录成功后，建议立刻刷新：

- `GET /api/accounts`
- 如果当前页面只看单账号，也可以调用 `GET /api/accounts/:id`

## 8. 相关文档

- `docs/api.md`: 完整 API 接口文档
- `docs/reference.md`: 配置、存储和运行参考
