# 集群会话生命周期统一设计（子项目A）

## 1. 背景与目标

当前 `gmcc` 的单实例能力已可用，但集群模式在会话生命周期上仍存在不一致：

- 实例状态切换缺少统一约束，容易出现并发下状态覆盖。
- 断线后重连机制存在基础实现，但缺少严格的单实例重连治理。
- 实例创建与移除在并发场景下可能出现“删除后回写”或“删除后复活”风险。

本设计聚焦“子项目A：集群会话生命周期统一”，目标是先建立稳定、可预测、可观测的运行底座。

## 2. 设计范围（In Scope / Out of Scope）

### In Scope

- 统一 `Manager` 与 `Instance` 的生命周期职责边界。
- 定义并收敛实例状态机与状态迁移规则。
- 仅在网络层明确断开时触发自动重连。
- 实现无限重试 + 指数退避策略。
- 完善创建/移除流程，规避并发竞态。
- 为认证与连服相关路径提供可本地执行的模拟测试方案（不依赖真实微软账号与真实服务器）。
- 在生命周期改造范围内清理无用旧代码（仅限与本设计直接相关模块）。

### Out of Scope

- 微软认证流程重构与 token 结构升级（子项目B）。
- 完整结构化日志与审计模型改造（子项目C）。
- 高级调度（优先级队列、多租户编排、熔断矩阵）（子项目D）。

## 3. 约束与确认决策

基于本轮澄清，固定以下约束：

- 自动重连策略：无限重试 + 指数退避。
- 默认退避参数：`base=2s`、`multiplier=1.8`、`max=2m`。
- 进程重启后不自动拉起实例，仅手动/API 启动。
- 断线判定为保守型：仅 TCP 断开、读写错误、连接关闭触发重连。

## 4. 架构与职责边界

### 4.1 Manager（编排与监督）

- 持有实例注册表与全局策略配置。
- 接收外部意图命令（start/stop/restart/create/delete）。
- 监听实例停止事件并按策略驱动重连。
- 保障每个实例在任意时刻最多一个重连循环在运行。

### 4.2 Instance（单实例状态机）

- 作为单实例生命周期唯一状态写入者。
- 维护状态、重连计数、最后活动时间、错误归因。
- 托管 runner 运行、取消与就绪判定。
- 对外提供幂等的 `Start/Stop/Restart` 语义。

### 4.3 Runner/Client（执行层）

- 负责连接建立、协议收发、业务执行。
- 返回原始错误，不做集群级重连策略决策。
- 提供明确就绪信号（`Runner.IsReady()==true` 且已进入 `play` 状态）。

### 4.4 接口契约（Manager ↔ Instance）

以下为本阶段必须收敛的契约语义（接口名可按代码风格命名）：

- `Manager` 对 `Instance` 命令接口：
  - `Start(ctx, trigger) error`：
    - 在 `pending/stopped/error` 可进入启动流程。
    - 在 `starting/running` 返回 `ErrInstanceRunningLike`（幂等保护）。
    - 在 `reconnecting`：
      - `trigger=auto_reconnect` 时允许进入 `starting`（用于重连监督循环）。
      - 其他 trigger 返回 `ErrInstanceRunningLike`。
    - `trigger` 取值固定为：`manual_start | manual_restart | auto_reconnect`。
  - `Stop(ctx) error`：
    - 任意状态可调用。
    - 若已 `stopped`，返回 `nil`（幂等）。
    - 必须中断运行协程与重连等待。
  - `Restart(ctx) error`：等价 `Stop(ctx)` 成功后再 `Start(ctx, manual_restart)`。
- `Instance` 对 `Manager` 事件接口：
  - `OnInstanceExit(instanceID, version, category, err)`：执行层退出后上报。
  - `OnInstanceReady(instanceID, version)`：进入就绪后上报。
- 并发约束：
  - 同一 `instanceID` 同时最多一个“启动路径”与一个“重连监督路径”。
  - 事件处理按 `instanceID` 串行化，采用 `map[instanceID]*sync.Mutex`（或等价每实例队列）实现，防止跨线程乱序覆盖状态。

命令仲裁优先级（高 -> 低）：`delete > stop > restart > auto_reconnect`。

重连触发职责：

- `reconnecting -> starting` 仅由 `Manager` 的重连监督循环触发。
- 监督循环调用公开 `Start(ctx, auto_reconnect)` 进入启动路径，不允许调用绕过状态机的内部启动入口。
- `Instance` 只负责状态迁移执行与幂等校验，不自行发起重连。

外部入口映射：

- API `StartInstance` -> `Start(ctx, manual_start)`。
- API `RestartInstance` -> `Restart(ctx)`（内部转 `manual_restart`）。
- 自动重连监督循环 -> `Start(ctx, auto_reconnect)`。

## 5. 实例状态机设计

### 5.1 状态定义

- `pending`：已注册，未启动。
- `starting`：启动中，等待进入可用状态。
- `running`：运行中，可接收命令。
- `reconnecting`：断线后按策略重连中。
- `stopped`：已停止（手动或流程终止）。
- `error`：启动失败或不可恢复错误（本阶段作为可见终态，等待人工 `Start` 或 `Stop`）。

### 5.2 合法迁移

- `pending -> starting`
- `starting -> running | error | stopped`
- `running -> reconnecting | stopped | error`
- `reconnecting -> starting | stopped | error`
- `error -> starting | stopped`
- `stopped -> starting`

所有迁移通过单入口方法执行，携带 `from/to/reason/at`，拒绝非法跳转。

重连成功的规范路径固定为：`running -> reconnecting -> starting -> running`。

## 6. 生命周期数据流

### 6.1 启动

1. API/命令调用 `StartInstance(id)`。
2. `Manager` 获取实例并下发启动意图。
3. `Instance` 进入 `starting`，创建 runner 与运行上下文。
4. runner 成功进入就绪后，`Instance` 切换为 `running`。

### 6.2 断线重连

1. 执行层返回网络断开类错误。
2. `Instance` 上报停止事件与错误归因。
3. `Manager` 进入该实例单飞重连循环。
4. 按退避参数等待并重启实例。
5. 成功按规范路径回到 `running`；失败继续下一轮直到被手动停止/删除或进程退出。

### 6.3 停止

1. API/命令调用 `StopInstance(id)`。
2. `Instance` 取消上下文，终止运行与待重连等待。
3. 状态落为 `stopped`，并阻断后续自动重连。

## 7. 创建/移除一致性设计

### 7.1 创建（Create）

- `CreateInstance` 只注册实例并初始化为 `pending`，不自动启动。
- 同 ID 冲突立即返回错误，避免隐式覆盖。

### 7.2 移除（Delete）

- 删除前必须完成实例终止（运行协程与重连协程都退出）。
- 删除后从实例表移除并使版本失效。
- 后台协程状态回写按“实例存在且版本匹配”校验，不满足即丢弃。
- `Delete` 终止等待超时设为 `delete_timeout = 10s`。
- 若在超时内仍未完成终止，`Delete` 返回 `ErrDeleteTimeout`，实例保留在注册表且状态置为 `error`（不执行移除）。

版本生命周期约定：

- `CreateInstance` 时初始化 `version=1`。
- 每次进入 `starting` 前（即发起一次启动尝试时）执行 `version++`，并绑定 `runVersion` 给本次运行协程与事件。
- `DeleteInstance` 时执行 `version++` 后从注册表移除。
- 任意异步回写必须同时满足：`instance存在` 且 `event.version==instance.version`。
- 若 `instance不存在`，事件按 stale event 直接丢弃并记录 debug 日志。
- 不满足条件的事件直接丢弃并记录 debug 日志。

`Delete` 超时后的处理：

- 不改变当前 `version` 绑定关系，不做移除。
- 允许后续再次 `Stop/Delete` 重试。
- 若之后收到旧运行协程事件，继续按“实例存在且版本匹配”规则校验。

### 7.3 并发安全规则

- `manager.mu` 仅保护实例表与管理器级共享数据。
- `instance.mu` 仅保护实例内部状态。
- 严格持锁顺序：先 `manager` 后 `instance`，禁止反向持锁。

## 8. 错误处理与归因

执行层错误上抛后，实例层归因为以下类别：

判定优先级（高 -> 低）：`manual_stop > network_disconnect > startup_timeout > auth_failed > unknown`。

说明：若手动停止与底层 EOF 同时出现，按 `manual_stop` 归因，确保不会误触发自动重连。

- `network_disconnect`：连接关闭、EOF、读写错误（本阶段触发重连）。
- `startup_timeout`：启动窗口内未 ready（本阶段记录为错误，不自动重连）。
- `auth_failed`：认证失败（记录但策略细化留给子项目B）。
- `manual_stop`：用户主动停止（禁止自动重连）。
- `unknown`：无法归类错误（本阶段记录为错误，不自动重连）。

重连触发策略表（子项目A最终版）：

- `network_disconnect -> reconnect=true`
- `startup_timeout -> reconnect=false`
- `auth_failed -> reconnect=false`
- `manual_stop -> reconnect=false`
- `unknown -> reconnect=false`

说明：本阶段严格遵循“保守型断线判定”，仅网络层明确断开触发重连。

## 9. 关键阈值（单一事实源）

- `startup_timeout = 30s`（从进入 `starting` 开始计时）。
- `ready_criteria = Runner.IsReady()==true 且 Client 已进入 play 状态`。
- `backoff_jitter = off`（本阶段关闭抖动，便于行为可预测与测试稳定）。
- `reconnect_backoff = 2s * 1.8^attempt, capped at 2m`。
- `attempt` 从 1 递增，无上限。
- `attempt_reset` 规则：
  - 重连成功并进入 `running` 后重置为 `1`。
  - 手动 `Restart` 触发新启动周期时重置为 `1`。
  - 手动 `Stop/Delete` 后不再保留历史 attempt。

## 10. 可观测性最小要求（本阶段）

本阶段不做完整审计系统，但日志必须可还原链路，至少包含：

- `instance_id`
- `from_state`
- `to_state`
- `reason`
- `attempt`
- `backoff`

用于支撑断线到恢复过程排查，后续由子项目C升级为结构化审计。

## 11. 测试与验收标准

### 11.1 单元测试

- 状态机合法/非法迁移。
- `Start/Stop/Restart` 幂等语义。
- 指数退避序列与上限收敛。
- `Stop/Delete` 对重连等待的即时中断。
- 并发 `Create/Delete/Start` 下无竞态。

### 11.2 集成测试

- 手动/API 启动后可进入 `running`。
- 人为断链后自动重连并恢复 `running`。
- 进程重启后默认不自动拉起实例。
- 删除实例后不会被历史协程“复活”。
- 重连等待中触发 `Restart`：仅保留一条启动路径，状态迁移无冲突。
- 重连等待中触发 `Delete`：重连循环被中断，按 `Delete` 语义成功移除或超时返回 `ErrDeleteTimeout`。
- `OnInstanceExit` 与 `Delete` 并发：不会出现已删除实例被回写为 `running/error`。

### 11.3 本地模拟测试策略（新增）

- 认证模拟：
  - 通过可注入的 auth provider/stub 返回固定 token 与可控错误（`auth_failed`、超时、刷新失败），覆盖生命周期与错误归因。
  - 禁止在 CI/本地默认测试中依赖真实微软账号交互。
- 连服模拟：
  - 使用本地 fake server（或可控 `net.Listener` harness）模拟握手成功、延迟、EOF、连接重置、启动期无响应。
  - 可脚本化触发“先 ready 后断链”“持续拒绝连接”“间歇性网络抖动”等场景。
- 稳定性验证：
  - 单实例与多实例场景都需在本地可复现。
  - 核心重连行为测试不依赖公网与第三方服务可用性。
- 真实链路测试：
  - 真实微软登录与真实服务器联调作为可选手工验证项，不作为默认自动化门槛。

### 11.4 验收门槛

- `go test ./...` 通过。
- 关键并发路径在 `go test -race ./...` 下无数据竞争。
- 日志可完整重建一次断线与恢复链路。
- 本地模拟测试可稳定通过，且不要求外部账号/公网依赖。

## 12. 旧代码清理策略（新增）

- 清理原则：仅清理由本次生命周期重构替代且已无调用路径的代码，不做无关大规模重构。
- 清理对象：
  - 重复/失效的状态写入分支。
  - 已被统一状态机替代的分散重连逻辑。
  - 不再使用的生命周期辅助字段或函数。
- 安全措施：
  - 删除前通过引用搜索与测试确认无调用。
  - 清理与功能改造分提交（或分 commit 区块）以便回溯。

## 13. 实施边界与后续衔接

子项目A完成后，进入子项目B与C：

- 子项目B：基于本状态机接入认证失败差异化策略与 token 生命周期治理。
- 子项目C：将本阶段最小日志字段升级为统一结构化审计事件。

本设计遵循 YAGNI：只实现当前稳定运行所需能力，不提前引入高复杂编排系统。
