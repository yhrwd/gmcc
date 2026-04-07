# gmcc

gmcc 是一个面向 API 的无界面 Minecraft 客户端服务，重点提供账号管理、实例编排和前端友好的 HTTP 接口。

当前运行时模型拆分为三层：

- account metadata：保存在 `.state/accounts.yaml` 的非敏感账号定义
- instance metadata：保存在 `.state/instances.yaml` 的非敏感实例定义
- auth vault：保存在 `auth.vault.path` 目录下的每账号加密认证记录

这种结构适合作为自定义管理面板、自动化平台或集群控制台的后端服务。

## 提供能力

- 提供账号、实例、集群状态、Microsoft 设备码登录和操作日志的 REST API
- 账号与实例解耦，一个账号可以支撑多个运行时实例
- 认证信息使用加密 auth vault 存储，统一由环境变量主密钥保护
- 重启后可通过 metadata 恢复资源，不把敏感信息写入明文 YAML
- 提供带自动重连策略的集群运行时和状态查询能力

## 快速开始

### 1. 准备环境变量

启动前先设置 auth vault 主密钥：

```bash
# Windows PowerShell
$env:GMCC_AUTH_VAULT_KEY = "replace-with-a-strong-secret"

# Windows CMD
set GMCC_AUTH_VAULT_KEY=replace-with-a-strong-secret
```

### 2. 构建

```bash
go build -o gmcc.exe ./cmd/gmcc
```

### 3. 启动

```bash
./gmcc.exe
```

默认 API 监听地址为 `0.0.0.0:8080`。

## 配置说明

默认从 `config.yaml` 加载配置，也可以通过 `GMCC_CONFIG` 指定其他路径。

示例：

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

## 运行时存储布局

这里的“运行时基目录”默认指当前活动配置文件所在目录；如果通过 `GMCC_CONFIG` 指定了配置文件，则以该文件所在目录为准。

- `config.yaml`: 服务配置
- `.authvault/`: 每账号一个加密认证文件
- `.state/accounts.yaml`: 账号 metadata
- `.state/instances.yaml`: 实例 metadata
- `logs/`: 运行日志（相对运行时基目录）
- `logs/audit/`: JSONL 格式的操作审计日志（相对运行时基目录）

## API 概览

基础路径：`/api`

读接口：

- `GET /api/status`
- `GET /api/accounts`
- `GET /api/accounts/:id`
- `GET /api/instances`
- `GET /api/instances/:id`
- `GET /api/logs/operations`

写接口：

- `POST /api/accounts`
- `DELETE /api/accounts/:id`
- `POST /api/instances`
- `POST /api/instances/:id/start`
- `POST /api/instances/:id/stop`
- `POST /api/instances/:id/restart`
- `DELETE /api/instances/:id`

认证接口：

- `POST /api/auth/microsoft/init`
- `POST /api/auth/microsoft/poll`

详细前端对接文档：

- `docs/api.md`
- `docs/auth.md`
- `docs/reference.md`

## 推荐前端接入流程

1. 使用 `POST /api/accounts` 创建账号 metadata
2. 使用 `POST /api/auth/microsoft/init` 发起 Microsoft 设备码登录
3. 使用 `POST /api/auth/microsoft/poll` 轮询登录状态
4. 当账号状态变为 `logged_in` 后，调用 `POST /api/instances` 创建实例
5. 使用实例控制接口启动、停止或重启实例
6. 使用 `GET /api/status`、`GET /api/instances`、`GET /api/logs/operations` 更新面板数据

前端约束：

- 当前内置 Vue 控制台与后端同源协作，所有请求都应使用 `/api` 下的相对路径，不要在生产构建里写死独立 API 域名。
- Microsoft 设备码登录第一版只支持单会话模型：同一时刻只维护一个活跃的设备码会话，新的发起动作应视为替换上一轮会话。

## 可选内嵌前端工作流

- 默认直接执行 `go build -o gmcc.exe ./cmd/gmcc` 时，gmcc 会以 API-only 方式构建；仓库内只保留 `internal/webui/dist/.keep` 占位文件，不会默认内嵌任何前端资源。
- 如果存在已构建好的前端产物，可先放到 `frontend/dist/`，再运行打包工具把允许内嵌的资源复制到 `internal/webui/dist/` 并完成最终构建。
- 打包工具入口为 `go run ./tools/packager`，默认会读取 `frontend/dist`，并在 `build/` 下生成 Linux / Windows / macOS 的多架构交叉编译产物。
- 允许进入内嵌目录的文件目前包括根目录下的 `index.html`、`favicon.ico`、`favicon.svg`、`manifest.webmanifest`、`robots.txt`，以及 `assets/` 下的非隐藏资源文件；`.map` 文件和隐藏文件/目录会被忽略。
- 如果 `frontend/dist` 不存在，或过滤后没有得到 `index.html`，打包工具会清理 `internal/webui/dist/` 回到仅保留 `.keep`，并继续产出 API-only 二进制。

示例：

```bash
go run ./tools/packager
```

### 前端开发与构建

前端工程位于 `frontend/`，基于 Vite + Vue 3 + TypeScript。

开发模式：

```bash
cd frontend
npm install
npm run dev
```

默认 Vite 开发服务器运行在 `http://localhost:5173`。仓库示例配置已经为该地址放开 CORS，便于本地联调后端 `http://localhost:8080`。

前端测试与构建：

```bash
cd frontend
npm run test
npm run build
```

`npm run build` 会输出 Vite 风格产物到 `frontend/dist/`，通常包括：

- 根目录 `index.html`
- `assets/` 下带内容哈希的 JS / CSS 等静态资源
- 可选白名单根文件，如 `favicon.ico`、`favicon.svg`、`manifest.webmanifest`、`robots.txt`

### Packager 用法

当你希望把前端一起嵌入 gmcc 二进制时，在仓库根目录执行：

```bash
go run ./tools/packager
```

常见工作流：

```bash
cd frontend
npm run build

cd ..
go run ./tools/packager
```

packager 默认读取 `frontend/dist/`，筛选可嵌入文件后复制到 `internal/webui/dist/`，再执行交叉编译并输出到 `build/`。

默认产物包括：

- `build/gmcc-linux-amd64`
- `build/gmcc-linux-arm64`
- `build/gmcc-windows-amd64.exe`
- `build/gmcc-windows-arm64.exe`
- `build/gmcc-darwin-amd64`
- `build/gmcc-darwin-arm64`

如果只想构建特定目标，可通过 `-targets` 指定，例如：

```bash
go run ./tools/packager -targets linux/amd64,darwin/arm64
```

如果仍然需要单文件旧模式输出，可通过 `-output` 保持兼容，例如：

```bash
go run ./tools/packager -targets windows/amd64 -output build/gmcc.exe
```

当前 whitelist 规则为：

- 根目录只允许 `index.html`、`favicon.ico`、`favicon.svg`、`manifest.webmanifest`、`robots.txt`
- `assets/` 目录允许普通静态资源，包括 Vite 默认生成的哈希文件名
- 隐藏文件、隐藏目录以及 `.map` sourcemap 文件不会进入内嵌目录
- 如果最终没有得到 `index.html`，packager 会回退到 API-only 构建

运行时行为：

- 所有 `/api` 路径始终由 API 处理。
- 当内嵌前端可用时，非 `/api` 的 `GET` / `HEAD` 请求会优先返回静态资源，不存在的前端路由会回退到 `index.html`。
- 当内嵌前端不可用时，非 `/api` 的 `GET` / `HEAD` 请求返回 `503 Frontend unavailable` 页面，其他非 `/api` 方法返回 `404`。

## 构建与测试

```bash
go fmt ./...
go test ./...
go build -o gmcc.exe ./cmd/gmcc
```

## 说明

- 当前服务默认以 API-only 方式构建和运行；只有在显式准备并内嵌前端产物后，才会对外提供前端静态资源
- `POST /api/instances` 要求目标账号已存在、已启用且认证状态可用
- `POST /api/instances` 中 `enabled` 省略时默认 `true`；若显式传 `enabled=false`，则不能同时传 `auto_start=true`
- 删除账号前必须先处理引用该账号的实例，否则会失败
- auth vault 中是敏感数据，不应提交 `.authvault/` 或泄露主密钥环境变量
