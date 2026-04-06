# gmcc 运行与集成参考

本文档提供当前 gmcc API 服务的配置、目录结构、运行方式和数据模型参考，方便前端、后端和运维统一理解系统边界。

## 1. 项目定位

gmcc 当前形态是一个 API-first 的 Minecraft 集群运行时，核心职责包括：

- 管理账号 metadata
- 管理实例 metadata
- 保存加密认证凭据
- 启动、停止、重启实例
- 提供可供前端消费的运行状态和审计日志

## 2. 核心目录结构

```text
gmcc/
├── cmd/gmcc/                  # 程序入口
├── internal/auth/             # 认证状态、provider、vault
├── internal/cluster/          # 实例生命周期和集群编排
├── internal/config/           # 统一配置定义与读写
├── internal/resource/         # 账号/实例 metadata 编排层
├── internal/state/            # metadata repository
├── internal/web/              # HTTP API 层
├── internal/webtypes/         # Web DTO
├── docs/                      # 对接文档
├── .authvault/                # 认证加密文件目录（运行时生成）
├── .state/                    # metadata 状态目录（运行时生成）
└── config.yaml                # 主配置文件
```

## 3. 关键模块边界

### `internal/auth/session`

负责账号认证状态管理，包括：

- 账号认证状态判定
- Microsoft 设备码登录流程
- 登录成功后写入 vault
- 查询账号 profile 和登录状态

关键状态：

- `logged_in`
- `not_logged_in`
- `auth_invalid`

### `internal/auth/vault`

负责加密认证记录持久化：

- 每账号一个 vault 文件
- 原子写入
- 按账号独立锁保护
- 启动时依赖环境变量提供主密钥

### `internal/state`

负责 metadata 文件读写：

- `AccountRepository` -> `.state/accounts.yaml`
- `InstanceRepository` -> `.state/instances.yaml`

特点：

- YAML 持久化
- 文件缺失时按空列表处理
- 保存前校验唯一性
- 原子覆盖写入

### `internal/resource`

负责跨 repository 的资源编排与引用校验：

- 创建账号
- 查询账号
- 删除账号
- 恢复资源
- 校验实例创建前的账号状态

它是“账号 metadata / 实例 metadata / auth status”之间的业务粘合层。

### `internal/cluster`

负责实例运行生命周期：

- 创建实例
- 启动/停止/重启/删除实例
- 汇总集群状态
- 自动重连
- 维护实例运行态指标

### `internal/web`

负责 Gin API 路由；当前服务以 API 为主，但如果存在嵌入式前端资源，也会对非 `/api` 路径提供静态资源或入口页回退。

## 4. 配置参考

默认配置结构如下：

```yaml
auth:
  vault:
    path: ".authvault"
    key_env: "GMCC_AUTH_VAULT_KEY"
    scrypt_n: 1048576
    scrypt_r: 8
    scrypt_p: 1
    salt_len: 32

cluster:
  global:
    max_instances: 10
    reconnect_policy:
      enabled: true
      max_retries: 0
      base_delay: 2s
      max_delay: 2m
      multiplier: 1.8
  accounts: []

web:
  bind: "0.0.0.0:8080"
  auth:
    audit_log_retention_days: 30
  cors:
    enabled: true
    origins:
      - "http://localhost:5173"
      - "http://localhost:3000"

log:
  log_dir: "logs"
  max_size: 512
  debug: false
  enable_file: true
```

### 关键配置项说明

| 字段 | 说明 |
|---|---|
| `auth.vault.path` | 认证 vault 目录 |
| `auth.vault.key_env` | 主密钥环境变量名 |
| `cluster.global.max_instances` | 实例数上限，`0` 表示不限制 |
| `cluster.global.reconnect_policy.*` | 自动重连参数 |
| `web.bind` | HTTP API 监听地址 |
| `web.cors.enabled` | 是否启用 CORS |
| `web.cors.origins` | 允许的前端来源 |
| `log.*` | 运行日志输出设置 |

## 5. 环境变量

| 变量名 | 说明 |
|---|---|
| `GMCC_CONFIG` | 自定义配置文件路径 |
| `GMCC_AUTH_VAULT_KEY` | vault 主密钥 |

如果 `auth.vault.key_env` 被改成别的名字，就需要提供对应的新环境变量。

## 6. 数据文件布局

### `.state/accounts.yaml`

保存账号非敏感元数据。

示例：

```yaml
- account_id: acc-main
  enabled: true
  label: Main account
  note: Primary production account
```

### `.state/instances.yaml`

保存实例非敏感元数据。

示例：

```yaml
- instance_id: bot-1
  account_id: acc-main
  server_address: mc.example.com
  enabled: true
```

### `.authvault/*.vault`

保存账号认证记录的加密文件，不适合人工编辑。

## 7. HTTP API 资源模型

### AccountView

```json
{
  "id": "acc-main",
  "player_id": "player-uuid",
  "enabled": true,
  "label": "Main account",
  "note": "Primary",
  "auth_status": "logged_in",
  "has_token": true
}
```

### InstanceView

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

### ClusterStatus

```json
{
  "cluster_status": "partial",
  "total_instances": 3,
  "running_instances": 2,
  "uptime": 60000000000
}
```

说明：`uptime` 是 Go `time.Duration` 的 JSON 表现，当前为纳秒数值。

补充：`GET /api/accounts/:id` 与 `GET /api/instances/:id` 都是直接返回单个对象，而不是 `{ "account": ... }` 或 `{ "instance": ... }` 结构。

## 8. 实例状态参考

| 值 | 含义 |
|---|---|
| `pending` | 已创建，未启动 |
| `starting` | 启动中 |
| `running` | 运行中 |
| `reconnecting` | 网络断开后自动重连中 |
| `stopped` | 已停止 |
| `error` | 运行异常 |

## 9. 审计日志

审计日志通过 `GET /api/logs/operations` 读取，磁盘上按天写入 `logs/audit/YYYY-MM-DD.jsonl`；若通过 `GMCC_CONFIG` 指定配置文件，则日志目录以该配置文件所在目录为基准，否则相对当前工作目录解析。

常见 `action`：

- `account_create`
- `account_delete`
- `instance_create`
- `instance_start`
- `instance_stop`
- `instance_restart`
- `instance_delete`
- `microsoft_auth_init`
- `microsoft_auth_poll`

## 10. 运行与验证命令

### 构建

```bash
go build -o gmcc.exe ./cmd/gmcc
```

### 测试

```bash
go test ./...
```

### 格式化

```bash
go fmt ./...
```

## 11. 对接建议

- 前端初始化时优先拉取 `GET /api/status`、`GET /api/accounts`、`GET /api/instances`
- 创建实例前先校验账号 `auth_status`
- Microsoft 登录成功后立即刷新账号列表
- 删除账号前先确认没有关联实例
- 若需要时间轴页面，可将操作日志作为后台管理时间线数据源

## 12. 相关文档

- `README.md`
- `docs/auth.md`
- `docs/api.md`
