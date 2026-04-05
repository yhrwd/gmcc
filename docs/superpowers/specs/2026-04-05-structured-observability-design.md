# 结构化可观测性设计（子项目C）

## 1. 背景与目标

Phase A/B 已经为集群实例生命周期与账号级认证打下基础，但当前日志体系仍存在明显不足：

- `internal/logx` 仍以全局文本日志为主，字段化能力弱。
- 控制台日志偏吵，关键摘要与过程细节混在一起。
- 虽然已有 `internal/web/audit`，但尚未形成统一的关键事件模型。
- 日志文件滚动能力只覆盖普通文本日志，事件维度与文件大小控制不足。

本设计聚焦“子项目C：结构化可观测性”，目标是在不引入远程日志平台和复杂查询系统的前提下，建立统一的结构化事件日志、控制台摘要降噪和可控的本地文件滚动体系。

## 2. 设计范围（In Scope / Out of Scope）

### In Scope

- 在 `internal/logx` 内新增结构化事件输出能力。
- 将控制台文本日志收敛为摘要级输出。
- 为文本日志与 JSONL 事件日志引入可控滚动策略。
- 为 Phase A/B 关键节点补充统一事件模型。
- 保持与 `internal/web/audit` 后续对接兼容。

### Out of Scope

- 远程日志采集、ELK、OpenTelemetry、分布式 tracing。
- Web 端日志检索平台与高级筛选 UI。
- 大范围重构所有历史日志调用点为字段式 API。

## 3. 约束与确认决策

基于本轮澄清，固定以下约束：

- 第一阶段优先做结构化事件日志。
- 默认输出落地为：本地 JSONL 文件 + 控制台摘要。
- 结构化事件的核心关联键为：`instance_id + account_id`。
- 控制台默认级别为摘要级，不保留过多中间过程。
- 单个日志文件必须受控，避免文件过大难以打开。

## 4. 架构与职责边界

### 4.1 logx 作为统一入口

- 保留 `internal/logx` 作为唯一日志出口。
- 在其内部拆分为两个逻辑通道：
  - `summary logger`：控制台 + 摘要文本文件。
  - `event logger`：JSONL 结构化事件文件。

### 4.2 调用侧职责

- `cluster`、`auth/session`、`mcclient` 仅负责声明“发生了什么”。
- 业务模块不得自行拼接 JSON 或直接写事件文件。
- 日志格式、落盘、滚动、降噪规则由 `logx` 统一实现。

### 4.3 与 audit 的关系

- 本阶段不把 `internal/web/audit` 作为中心依赖。
- 但事件字段设计需兼容后续审计模型，尤其是“操作意图 + 结果”类事件。

## 5. 日志分层策略

### 5.1 Summary（默认控制台可见）

仅保留需要人即时感知的摘要事件，例如：

- 实例启动成功 / 停止 / 删除失败
- 进入重连 / 重连成功 / 重连失败
- 认证成功 / 认证失败 / 需要重新登录
- 系统级错误

首期摘要保留判定规则：

- 会改变实例或账号可用性的事件，保留摘要。
- 需要人工介入的事件，保留摘要。
- 高频重复的过程事件，不保留摘要。
- 已有结构化事件可完整表达、但对值班者无需即时可见的过程事件，下沉到 event/debug。

### 5.2 Debug（默认关闭）

用于记录中间决策和排障细节，例如：

- token 选择路径
- provider 分支结果
- 设备码轮询状态
- 连接阶段细节和协议过程节点

### 5.3 Event（JSONL 结构化事件）

用于承载关键行为的结构化记录，便于后续检索、审计与统计。

## 6. 事件模型

每条 JSONL 事件最小字段：

- `ts`
- `level`
- `event_type`
- `action`
- `message`
- `instance_id`（系统级事件可省略）
- `account_id`（系统级事件可省略）

可选扩展字段：

- `player_id`
- `from_status`
- `to_status`
- `reason`
- `attempt`
- `backoff_ms`
- `auth_error`
- `result`

字段语义约束：

- `event_type`：固定事件类别，如 `instance.lifecycle`、`instance.reconnect`、`auth.session`。
- `action`：该类别下的具体动作，如 `start`、`ready`、`scheduled`、`auth_failed`。
- `message`：给人读的简短摘要，不承载机器判定语义。
- `result`：仅在需要表达结果时使用，如 `success`、`failed`、`cancelled`。

字段规则：

- 使用固定、小写下划线命名。
- 缺失值可省略，不强制写空字符串。
- 禁止按业务动态生成字段名。
- 若事件属于 `system.error` 且无实例/账号上下文，允许省略 `instance_id` 与 `account_id`。

示例：

```json
{"ts":"2026-04-05T08:30:00Z","level":"info","event_type":"instance.lifecycle","action":"ready","message":"instance ready","instance_id":"bot-1","account_id":"acc-main","player_id":"Steve"}
{"ts":"2026-04-05T08:31:00Z","level":"warn","event_type":"instance.reconnect","action":"scheduled","message":"instance reconnect scheduled","instance_id":"bot-1","account_id":"acc-main","attempt":2,"backoff_ms":3600,"reason":"network_disconnect"}
{"ts":"2026-04-05T08:32:00Z","level":"error","event_type":"auth.session","action":"auth_failed","message":"account authentication failed","instance_id":"bot-1","account_id":"acc-main","auth_error":"device_login_required","result":"failed"}
```

## 7. 首期事件类型

### 7.1 生命周期事件

- `instance.lifecycle`
  - `action=create`
  - `action=start`
  - `action=ready`
  - `action=stop`
  - `action=delete`
  - `action=transition`

### 7.2 重连事件

- `instance.reconnect`
  - `action=scheduled`
  - `action=waiting`
  - `action=success`
  - `action=failed`

### 7.3 认证事件

- `auth.session`
  - `action=cache_hit`
  - `action=refresh`
  - `action=device_login_begin`
  - `action=device_login_status`
  - `action=auth_failed`
  - `action=clear`

### 7.4 系统错误事件

- `system.error`
  - `action=fatal`

## 8. 文件与滚动策略

### 8.1 文件布局

- 文本摘要日志：`logs/gmcc.log`
- 结构化事件日志：`logs/gmcc-events.jsonl`

### 8.2 滚动策略

- 首期按大小滚动，不按日期切片。
- 当文件超过阈值时，将活动文件重命名归档并创建新文件。
- 事件日志阈值应小于文本日志阈值，保证 JSONL 文件更易打开。

默认阈值与保留策略：

- 文本摘要日志默认阈值：`10MB`
- 事件 JSONL 日志默认阈值：`5MB`
- 默认最多保留 `5` 个归档文件，超出后删除最旧文件。
- 首期允许通过 `logx` 初始化参数或集中配置覆盖默认值；若未配置，则使用上述固定默认值。

归档命名约定：

- 文本日志：`gmcc-YYYYMMDD-HHMMSS.log`
- 事件日志：`gmcc-events-YYYYMMDD-HHMMSS.jsonl`

### 8.3 写入失败降级

- 日志写入失败不得影响主业务。
- 控制台摘要优先保证输出。
- 事件文件写入失败时，仅输出降级告警，不阻塞实例运行。

文件系统错误路径要求：

- 初始化失败（目录创建失败、文件打开失败）：
  - 控制台日志继续可用。
  - 文件日志能力降级并输出一次启动告警。
- 滚动失败（rename 失败、新文件创建失败）：
  - 保留当前可用输出能力。
  - 输出单次告警，避免无限刷屏。
- 关闭失败：
  - 仅记录告警，不影响进程退出流程。
- 并发写入期间的写失败：
  - 当前事件丢弃并输出降级告警。
  - 不阻塞业务 goroutine。

## 9. 并发与安全规则

- `logx` 继续作为统一并发保护点，避免多 goroutine 并发写入半行 JSONL。
- 文本与事件写入可以共享统一锁，或各自独立锁，但必须保证单条事件原子写入。
- 事件文件与文本文件都应支持安全关闭与重新初始化。
- 单个日志目录应避免宽松权限；文件权限与目录权限遵循最小可用原则。

敏感信息边界：

- 事件日志和摘要日志禁止记录 access token、refresh token、device code 原文、XSTS token、完整 provider 原始响应。
- `auth_error` 只允许记录归一化错误码，不允许落原始敏感 payload。
- 若底层错误文本可能包含敏感内容，必须先归一化或裁剪后再写事件。

## 10. 旧日志清理策略

- 清理那些只为临时排障存在、且会重复轰炸控制台的日志点。
- 保留真正有运营价值的摘要日志。
- 将重复性的“中间步骤”日志下沉到 `Debug` 或结构化事件。
- 不做全仓库日志点一次性重写，只优先清理 Phase A/B 相关路径。

首期治理清单仅限：

- `internal/cluster`：实例启动、停止、重连、状态迁移相关日志。
- `internal/auth/session`：cache hit、refresh、device login、auth_failed、clear 相关日志。
- `internal/mcclient`：在线认证流程中的 token 选择与中间过程日志。
- `internal/web/audit`：保留其现有职责，不纳入本阶段大规模改造。

迁移规则：

- 关键摘要继续保留在控制台。
- 细节分支转为 `Debugf` 或事件日志。
- 无法归类且短期仍有价值的日志先保留，避免一次性误删。

## 11. 测试策略

### 11.1 单元测试

- 事件序列化字段完整性。
- JSONL 单行写入与换行格式。
- 摘要级降噪规则。
- 文件大小达到阈值时的滚动切换。
- 目录不可写、初始化失败、rename 失败、关闭失败等降级路径。
- 多 goroutine 并发写事件时的原子性与无半行写入。

### 11.2 集成测试

- 实例启动 -> 就绪：产生摘要 + `instance.lifecycle` 事件。
- 认证失败：控制台仅输出摘要，同时产出 `auth.session` 事件。
- 断线重连：产出 `instance.reconnect` 事件，并控制控制台输出量。

### 11.3 回归测试

- 默认控制台输出只允许出现摘要节点，不允许出现 provider 分支、设备码轮询、token 选择等中间细节。
- 关键节点仍然可见且不丢失。
- `go test ./...`、关键 `-race`、`go build` 通过。

## 12. 验收标准

- Phase A/B 的关键流程都能生成统一结构化事件。
- 控制台默认输出收敛为摘要级；以“实例启动 -> 认证 -> 断线重连”流程为例，只允许出现启动成功、认证结果、重连计划、重连结果、严重错误这些摘要节点。
- 事件日志和文本日志均具备单文件大小控制。
- 日志写入失败不会阻塞实例运行。
- 新增测试覆盖事件模型、滚动行为与关键流程联动。

首期必须接入的调用点清单：

- `cluster.Instance.Start/Stop/Restart` 相关生命周期节点。
- `cluster.Manager.handleReconnect` 的 scheduled/success/failed 节点。
- `auth/session.AuthManager.GetSession/Refresh/BeginDeviceLogin/Clear` 相关关键节点。
- `mcclient.prepareOnlineSession` 中认证成功/失败摘要节点。

## 13. 与后续阶段衔接

- 子项目C 完成后，后续可以基于 JSONL 事件模型接入 `internal/web/audit` 查询与展示。
- 首期不引入远程平台，遵循 YAGNI，只建立本地结构化观测基础。

## 14. 实施结果记录（Phase C）

- `logx` 已支持摘要日志与结构化事件双通道。
- Phase A/B 关键流程已接入 `instance.lifecycle`、`instance.reconnect`、`auth.session` 事件。
- 默认控制台日志已收敛为摘要级，中间过程已转为 debug 或事件日志。
- 文本日志与 JSONL 事件日志均已支持受控滚动与归档数量限制。
- 与结构化事件重复的旧日志路径已在目标范围内清理。
