# 账号级认证与 Token 集群化设计（子项目B）

## 1. 背景与目标

Phase A 已完成集群实例生命周期统一，但认证链路仍主要以单实例直连方式工作：

- Microsoft / Minecraft 认证调用分散在底层 service 中。
- token 缓存当前偏 `player_id` 语义，缺少账号级归属与并发仲裁。
- 多实例共享同一账号时，可能出现并发刷新、缓存覆盖、失败归因不一致。

本设计聚焦“子项目B：账号级认证与 token 集群化”，目标是在不引入过度复杂度的前提下，建立账号级认证状态、单飞刷新与安全持久化模型，使实例只消费认证结果，不直接管理认证流程。

## 2. 设计范围（In Scope / Out of Scope）

### In Scope

- 引入账号级 `AuthManager` 统一认证编排。
- 引入账号级 `TokenStore` 管理本地 token 缓存。
- 将 token/cache 主键从 `player_id` 升级为 `account_id`。
- 为同账号并发认证请求提供单飞刷新/登录机制。
- 将认证失败归因为实例可消费的明确状态，并阻断自动重连。
- 提供本地 stub/fake 测试方案，不依赖真实微软账号与公网服务。

### Out of Scope

- 修改 Minecraft 连接协议与 session join 算法。
- 引入数据库、分布式锁、远程凭据中心。
- 重构完整日志/审计系统（后续子项目C处理）。

## 3. 约束与确认决策

基于本轮澄清，固定以下约束：

- token/cache 归属模型：按账号独立持有，实例只消费凭据。
- 同账号刷新策略：单飞刷新，同一时刻仅一个真实刷新流程。
- 持久化方式：继续本地文件持久化，不引入数据库。
- refresh token 失效或需要重新设备码登录时：实例进入认证失败状态，停止自动重连，等待人工重新登录。

## 4. 架构与职责边界

### 4.1 AuthManager（账号级认证编排）

- 作为账号级认证状态唯一编排者。
- 负责缓存命中判断、刷新链路决策、设备码登录触发、失败归因。
- 负责同账号并发请求单飞。
- 对外提供统一认证入口，供实例或管理器调用。

### 4.2 TokenStore（账号级缓存存储）

- 负责账号级 token 缓存的 `Load/Save/Delete`。
- 负责本地文件路径映射、原子写入、基本并发保护。
- 不负责认证策略，不直接调用微软或 Minecraft 服务。

### 4.3 Provider 层（远端认证调用）

- `internal/auth/microsoft`：设备码、refresh、Xbox、XSTS。
- `internal/auth/minecraft`：Minecraft token、所有权、profile、证书。
- Provider 层只负责远端交互，不直接持久化 token，不做账号级并发仲裁。

### 4.4 实例/集群层

- 实例启动前向 `AuthManager` 请求账号会话。
- 实例只消费认证结果（成功会话或失败归因）。
- 若认证失败，实例进入 `auth_failed` 并停止自动重连。

## 5. 接口契约

以下为本阶段需收敛的接口语义（名称可按代码风格微调）：

- `AuthManager.GetSession(accountID) (AuthSession, error)`
  - 优先返回可用 Minecraft token + profile。
  - 若本地缓存不可用，则内部决定是否 refresh 或要求设备码登录。
- `AuthManager.Refresh(accountID) (AuthSession, error)`
  - 显式触发刷新链路。
  - 同账号若已有刷新在进行，则等待已有结果。
- `AuthManager.BeginDeviceLogin(accountID) (DeviceLoginInfo, error)`
  - 触发设备码流程并返回用户可执行信息。
- `AuthManager.GetDeviceLoginStatus(accountID) (DeviceLoginStatus, *AuthSession, error)`
  - 返回 `pending / succeeded / expired / cancelled / failed`。
  - 若状态为 `succeeded`，则同时返回最新 `AuthSession`。
- `AuthManager.CancelDeviceLogin(accountID) error`
  - 取消当前账号仍在进行中的设备码流程。
- `AuthManager.Clear(accountID) error`
  - 清理该账号缓存，仅供显式 relogin/logout 使用。

`AuthSession` 最小字段要求：

- `AccountID`
- `MinecraftAccessToken`
- `ProfileID`
- `ProfileName`
- `MicrosoftExpiresAt`
- `MinecraftExpiresAt`
- `Source`（cache / refresh / device_login）

设备码登录契约：

- `BeginDeviceLogin(accountID)` 负责：
  - 请求设备码
  - 返回 `DeviceLoginInfo`（verification URI、user code、expires_at、poll_interval）
  - 在 `AuthManager` 内部登记该账号登录中的流程
- 设备码轮询由 `AuthManager` 内部负责，不要求实例层自行轮询。
- 当设备码登录成功时，由 `AuthManager` 完成 Microsoft -> XSTS -> Minecraft 全链路并写回缓存。
- 当设备码登录未完成、过期、取消或失败时，`AuthManager` 以统一错误状态暴露给调用方。
- 调用方通过 `GetDeviceLoginStatus(accountID)` 观察登录中、成功、失败或过期结果，无需自行实现轮询逻辑。

设备码状态流转：

- `pending`：`BeginDeviceLogin` 创建流程后进入。
- `succeeded`：用户完成授权且全链路 token/profile 获取成功，由 `AuthManager` 内部轮询触发。
- `expired`：设备码超时未完成授权，由 `AuthManager` 内部轮询触发。
- `cancelled`：运维/API 显式取消当前设备码流程时进入；取消动作归 `AuthManager` 所有。
- `failed`：设备码后续链路执行失败（如 XSTS 拒绝、profile 异常）。

归一化规则：

- `expired -> device_login_required`
- `cancelled -> device_login_required`
- `failed -> 进一步按具体 provider 错误归一化`

## 6. 认证数据流

### 6.1 实例启动取会话

1. 实例启动前通过 `account_id` 请求 `AuthManager.GetSession(accountID)`。
2. `AuthManager` 先读取 `TokenStore`。
3. 若缓存中 Minecraft token 有效，直接返回。
4. 若 Minecraft token 无效但 Microsoft access token 或 refresh token 可用，则尝试刷新链路。
5. 刷新成功后持久化并返回最新会话。
6. 若必须重新设备码登录，则返回明确的 `device_login_required` 类错误。

### 6.2 刷新链路

1. 使用 Microsoft refresh token 获取新的 Microsoft access token。
2. 获取 Xbox token 与 XSTS token。
3. 获取 Minecraft access token。
4. 验证游戏所有权并读取 profile。
5. 更新本地缓存并返回 `AuthSession`。

### 6.3 人工重新登录

1. 运维/API 显式触发 `Clear(accountID)` 或 relogin。
2. 清理失效 token。
3. 走设备码登录流程。
4. 登录成功后写回账号级缓存。

## 7. 错误模型与实例联动

实例侧需要可消费的认证归因，最少包括：

- `auth_failed`：总类，实例可见错误类型。
- `refresh_token_invalid`：refresh token 无效、吊销或不可再用。
- `refresh_upstream_failed`：refresh 链路上的临时上游失败。
- `device_login_required`：必须重新设备码登录。
- `ownership_failed`：未购买 Minecraft 或所有权校验失败。
- `profile_invalid`：profile 数据缺失或异常。
- `provider_unavailable`：微软/Minecraft 服务临时不可用。
- `xsts_denied`：XSTS 授权被拒绝。

联动规则：

- 认证错误不会触发网络重连逻辑。
- refresh token 失效、设备码重新登录需求、所有权失败等都会使实例进入认证失败状态。
- 只有人工重新登录或显式刷新成功后，实例才允许重新启动。

错误联动表：

- `refresh_token_invalid -> instance auth_failed -> operator action: relogin`
- `refresh_upstream_failed -> instance auth_failed -> operator action: manual refresh after upstream recovery`
- `device_login_required -> instance auth_failed -> operator action: complete device login`
- `ownership_failed -> instance auth_failed -> operator action: inspect account entitlement`
- `profile_invalid -> instance auth_failed -> operator action: relogin or inspect provider response`
- `provider_unavailable -> instance auth_failed -> operator action: wait upstream recovery, then manual refresh/restart`
- `xsts_denied -> instance auth_failed -> operator action: inspect Xbox/XSTS account state and relogin if needed`

重试策略说明：

- 本阶段实例侧对所有认证错误统一落入 `auth_failed`，不自动重试。
- `provider_unavailable` 虽然语义上是临时外部错误，但为避免与网络重连混淆，也不进入自动重试；运维可在上游恢复后手动 `Refresh` 或重新启动实例。

Provider 错误归一化规则：

- Microsoft token endpoint 暂时不可用 -> `provider_unavailable`
- refresh token 被拒绝/失效 -> `refresh_token_invalid`
- refresh 链路出现临时 HTTP/网络错误 -> `refresh_upstream_failed`
- XSTS 返回授权拒绝 -> `xsts_denied`
- Minecraft entitlement 校验失败 -> `ownership_failed`
- Minecraft profile 缺字段 -> `profile_invalid`

说明：虽然 `provider_unavailable` 是临时错误，但本阶段仍统一收敛为实例侧 `auth_failed`，避免与 Phase A 的网络重连语义混淆；恢复动作由人工刷新/重启显式触发。

## 8. TokenStore 结构与持久化规则

缓存文件从当前 `player_id` 语义升级为账号级语义，例如按 `account_id` 命名。

最小字段：

- `account_id`
- `updated_at`
- `last_auth_error`
- `microsoft.access_token`
- `microsoft.refresh_token`
- `microsoft.expires_at`
- `minecraft.access_token`
- `minecraft.expires_at`
- `minecraft.profile_id`
- `minecraft.profile_name`

其中 `last_auth_error` 仅用于诊断展示与运维排查，不直接驱动运行时自动决策。

持久化要求：

- 继续使用本地文件。
- 保持原子写入（临时文件 + rename）。
- 文件权限维持最小可用范围。
- `Delete` 仅在显式 logout/relogin 时调用，实例停止/删除不得清理账号缓存。

兼容/迁移规则：

- 旧 `player_id` 缓存文件不做自动迁移。
- 新实现优先只读取 `account_id` 语义缓存文件。
- 若用户已有旧缓存，首次升级后需要重新登录一次以建立账号级缓存。
- 后续可选增加迁移工具，但不属于子项目B范围。

## 9. 并发与一致性规则

- 同账号所有认证入口单飞。
- 同账号缓存写入由 `TokenStore` 串行化。
- 任意刷新结果写回前都要校验当前账号流程代次，防止旧刷新覆盖新状态。
- 多实例共享同一账号时，只允许一个真实 refresh/device-login 流程，其余请求等待统一结果。
- 实例不得直接写 token 文件，所有写入必须经过 `AuthManager -> TokenStore`。

有效期判定规则：

- token 若在未来 `60s` 内过期，视为不可用。
- `GetSession` 与 `Refresh` 共用同一过期缓冲规则，避免实例拿到即将过期的 token。
- 所有比较均以当前 UTC 时间为准。

## 10. 本地模拟测试策略

- 使用 fake Microsoft provider 模拟：
  - refresh 成功
  - refresh token 失效
  - 设备码登录必需
  - XSTS 拒绝
  - provider 临时不可用
- 使用 fake Minecraft provider 模拟：
  - token 获取成功
  - 所有权校验失败
  - profile 数据不完整
- 使用并发表驱动测试覆盖：
  - 同账号多实例并发启动时只发生一次 refresh
  - 失败结果可被等待方共享
  - 旧刷新结果不会覆盖新流程状态

真实微软登录仅作为手工联调项，不作为默认自动化门槛。

设备码流程测试至少覆盖：

- `pending -> succeeded`
- `pending -> expired`
- `pending -> cancelled`
- `pending -> failed`

## 11. 旧代码清理策略

- 清理当前散落在调用点上的直接缓存读写。
- 清理重复的 token 有效性判断与刷新分支。
- 保留 provider 层远端调用逻辑，但移除其直接承担的缓存策略职责。
- 仅清理与子项目B直接相关的认证旧逻辑，不做无关重构。

## 12. 验收标准

- 同账号并发启动多个实例时，只发生一次真实 refresh/login。
- refresh token 失效后，实例进入认证失败并停止自动重连。
- 本地 fake provider 测试可稳定通过，不依赖真实微软账号与公网。
- token 缓存从账号级维度持久化，实例删除不影响账号缓存。
- `go test ./...` 通过，新增认证并发测试在 `-race` 下无竞争。

## 13. 与后续阶段衔接

- 子项目B完成后，子项目C可以基于账号级认证状态输出更完整的结构化日志与审计事件。
- 本阶段不引入数据库与分布式协调，遵循 YAGNI，只为当前本地/单节点集群能力提供稳定认证底座。

## 14. 实施结果记录（Phase B）

- 账号级 `AuthManager + TokenStore` 已落地。
- 账号 token 缓存已从 `player_id` 语义切换为 `account_id` 语义，不做自动迁移。
- 单飞刷新与设备码登录状态流转已实现，并支持本地 fake provider 测试。
- 实例认证失败统一映射为 `auth_failed`，不会触发自动网络重连。
