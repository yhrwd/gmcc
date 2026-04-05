# AGENTS.md

本文件面向会在本仓库内工作的 agentic coding agents。
目标是让代理在尽量少试错的前提下，遵循本仓库真实存在的构建方式、测试方式和代码风格。

## 项目概览

- 仓库语言：Go。
- 模块名：`gmcc`。
- 入口程序：`cmd/gmcc`。
- 主要目录：`internal/`、`pkg/`、`cmd/`、`docs/`。
- 运行形态：API-only 的 Minecraft 客户端服务，不内置前端静态资源。
- `frontend/` 目录当前为空，不要假设存在可构建前端。

## 已确认的仓库规则文件

- 未发现 `.cursorrules`。
- 未发现 `.cursor/rules/` 下的规则文件。
- 未发现 `.github/copilot-instructions.md`。
- 因此本文件即为当前仓库内最接近“代理工作说明”的统一入口。

## 构建、测试、格式化命令

以下命令已在仓库根目录验证可用。

### 常用命令

- 格式化全仓库：`go fmt ./...`
- 运行全部测试：`go test ./...`
- 构建主程序：`go build -o gmcc.exe ./cmd/gmcc`
- 直接运行程序：`go run ./cmd/gmcc`

### 单包测试

- 运行某个包的所有测试：`go test ./internal/web`
- 禁用缓存运行某个包：`go test -count=1 ./internal/web`
- 查看详细输出：`go test -v ./internal/web`

### 单个测试

- 运行单个测试函数：`go test ./internal/web -run TestHandleCreateInstanceAllowsDistinctInstanceAndAccountIDs`
- 禁用缓存运行单测：`go test -count=1 ./internal/web -run TestHandleCreateInstanceAllowsDistinctInstanceAndAccountIDs`
- 运行名称匹配的一组测试：`go test ./internal/web -run HandleCreateInstance`

### 其他实用测试方式

- 只测入口程序包：`go test ./cmd/gmcc`
- 只测集群模块：`go test ./internal/cluster`
- 只测认证会话模块：`go test ./internal/auth/session`
- 只测资源管理模块：`go test ./internal/resource`
- 只测状态持久化模块：`go test ./internal/state`

### 发布相关构建

GitHub Actions 中存在跨平台发布构建，见 `.github/workflows/release.yml`。
如需本地模拟，参考这些命令：

- Linux amd64：`GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-linux-amd64 ./cmd/gmcc`
- Linux arm64：`GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-linux-arm64 ./cmd/gmcc`
- Windows amd64：`GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-windows-amd64.exe ./cmd/gmcc`
- Windows arm64：`GOOS=windows GOARCH=arm64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-windows-arm64.exe ./cmd/gmcc`
- macOS amd64：`GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-darwin-amd64 ./cmd/gmcc`
- macOS arm64：`GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=${VERSION}" -o release/gmcc-darwin-arm64 ./cmd/gmcc`

## Lint 现状

- 仓库中未发现 `golangci-lint`、`staticcheck` 或其他专用 lint 配置。
- 实际可依赖的质量门槛主要是：`go fmt ./...`、`go test ./...`、`go build ./cmd/gmcc`。
- 若你进行了中等以上复杂度修改，建议额外执行：`go vet ./...`。
- 但在没有现成配置的情况下，不要擅自引入新的 lint 工具或强制性配置文件，除非用户明确要求。

## 运行与环境注意事项

- 默认配置文件路径是仓库根目录的 `config.yaml`。
- 可通过环境变量 `GMCC_CONFIG` 指定其他配置文件路径。
- 启动前通常需要设置 auth vault 主密钥环境变量，默认键名为 `GMCC_AUTH_VAULT_KEY`。
- 程序会围绕配置文件所在目录解析运行时基目录，而不是一律以仓库根目录为准。
- 运行时会读写 `.state/`、`.authvault/`、`logs/` 等目录。
- `.authvault/` 含敏感认证数据，绝不要建议提交到版本库。

## 代码风格总原则

- 优先遵循现有代码，而不是套用外部“最佳实践”模板。
- 该仓库大量使用中文注释与中文日志；新增内容应与周围文件语言保持一致。
- 结构上偏向“小而清晰”的函数与显式错误返回，而不是隐式魔法。
- 对外部输入做规范化和校验，再进入业务逻辑。
- 错误信息尽量带上下文，必要时使用 `%w` 包装底层错误。

## 格式化与文件风格

- 使用 `go fmt`，不要手工对齐 import 或缩进。
- `.editorconfig` 规定：空格缩进、缩进宽度 4、`utf-8`、行尾 `crlf`。
- `trim_trailing_whitespace = false`，不要专门做无意义的尾空格清理大改。
- `insert_final_newline = false`，不要因为编辑器习惯强行改整个文件结尾样式。
- 除非用户要求，不要做纯风格化的大面积重排。

## Import 约定

- 遵循 Go 默认分组：标准库、第三方、本项目包。
- 本项目导入路径使用模块名前缀，如 `gmcc/internal/config`。
- 当包名冲突、过长或语义不清时，当前仓库接受显式别名，例如：
  - `authsession "gmcc/internal/auth/session"`
  - `authvault "gmcc/internal/auth/vault"`
  - `appconfig "gmcc/internal/config"`
- 不要为无冲突的普通导入随意起别名。

## 命名约定

- 导出标识符使用 Go 常规 PascalCase。
- 未导出标识符使用 camelCase。
- 接收请求的 handler 命名常见形式：`handleGetStatus`、`handleCreateInstance`。
- 构造函数统一使用 `NewXxx`，如 `NewManager`、`NewServer`、`NewAuthManager`。
- 错误变量使用 `ErrXxx` 形式，如 `ErrInstanceNotFound`。
- 测试函数命名使用 `TestXxx`，且倾向写出行为语义，而不是过短缩写。

## 类型与数据建模

- 优先使用明确 struct，而不是 `map[string]any` 传递业务数据；`map[string]any` 主要出现在测试断言 JSON 时。
- 配置、状态、请求响应模型通常放在明确的领域包中，如 `internal/config`、`internal/state`、`internal/webtypes`。
- 数据结构字段通常带 `yaml` 或 `json` tag，新增字段时保持序列化命名风格一致。
- 布尔字段常直接命名为状态，如 `Enabled`、`Success`。
- 可选输入在 HTTP 请求中有时用指针表达“是否显式传入”，例如 `req.Enabled != nil` 这种模式应保持。

## 输入规范化与校验

- 仓库大量使用 `strings.TrimSpace` 清理输入；新增入口参数时优先考虑是否需要 trim。
- 标识符为空通常应尽早返回，例如账号 ID、实例 ID、路径等。
- 创建前先做 normalize，再做 existence / duplicate / state 校验。
- 对外部请求，先 `ShouldBindJSON`，再做领域校验。
- 不要把原始未经清理的用户输入直接写入状态文件或日志。

## 错误处理约定

- 采用显式返回 `error`，不要用 panic 处理正常业务分支。
- 包内错误尽量定义为 sentinel error，跨层判断时使用 `errors.Is`。
- 包装错误优先使用 `fmt.Errorf("context: %w", err)`。
- 错误文案风格并不完全统一；新增代码应优先匹配所在文件语言和表述方式。
- 在 handler 中，业务错误通常转为 HTTP 4xx/5xx JSON 响应，不向外泄露内部细节。
- 在基础设施层，尽量保留底层错误上下文，方便定位。

## 日志约定

- 统一使用 `internal/logx`，不要直接引入自定义第三方日志框架。
- 常见接口：`logx.Infof`、`logx.Warnf`、`logx.Errorf`、`logx.Debugf`。
- 重要状态变更、恢复流程、失败原因应记录日志。
- 日志信息常为中文；若所在文件已明显使用英文错误文本，可就地保持一致。
- 结构化事件通过 `logx.Emit` 及对应 event builder 写入，不要手写随意 JSON 行格式。

## 并发与状态

- 现有代码广泛使用 `sync.Mutex` / `sync.RWMutex` 保护共享状态。
- 修改管理器类对象时，先看现有锁边界，不要在未知线程安全前提下直接读写 map。
- 停止、关闭、监督类流程常结合 `context.Context`、`cancel`、`WaitGroup`。
- 如需新增并发逻辑，优先延续已有模式，而不是引入复杂 channel 架构。

## HTTP 与 Web 层约定

- Web 层基于 Gin。
- handler 通常先做依赖存在性检查，再做参数解析，再执行业务动作。
- JSON 响应多使用 `webtypes.OperationResponse` 或明确的 view model。
- API 路径统一挂在 `/api` 下。
- 新接口应考虑审计日志是否需要通过 `logOperation` 记录。
- 仓库当前是 API-only 模式，不要擅自恢复静态资源托管。

## 持久化与文件修改约定

- 状态元数据保存在 YAML 文件中，见 `internal/state`。
- 写文件时已有原子写入模式 `writeYAMLAtomic`，新增类似逻辑应复用相同思路。
- 仓库对账号 metadata、实例 metadata、auth vault 做了职责拆分，修改时不要把敏感信息混入明文 YAML。
- 涉及配置文件生成时，注意现有代码会尽量保留带注释输出。

## 测试风格约定

- 测试以标准库 `testing` 为主，无 testify 等断言框架。
- 失败断言主要使用 `t.Fatalf` / `t.Fatal`。
- 测试名倾向完整表达行为，例如 `TestHandleCreateInstanceRejectsAutoStartForDisabledInstance`。
- 常见模式是手写 fake / stub，而不是使用 mock 框架。
- Web 测试使用 `httptest` 和 Gin 的 test mode。
- 涉及临时目录或文件的测试，使用 `t.TempDir()`。

## 改动策略建议

- 先看相邻文件怎么做，再决定是否抽象。
- 对已有公开行为，优先写或更新测试，再改实现。
- 不要顺手大规模重命名中文日志、注释或错误文案，除非任务本身要求。
- 没有必要时，不要把中文注释统一翻成英文，也不要反过来。
- 不要引入新的架构层、依赖注入框架或 ORM 风格封装。

## 代理工作清单

在提交较大改动前，推荐至少完成以下检查：

- `go fmt ./...`
- `go test ./...`
- `go build -o gmcc.exe ./cmd/gmcc`

若只改动单个包，最少应运行对应包测试；若改动 HTTP 路由、状态持久化或集群生命周期逻辑，建议回归全量测试。

## 给代理的最后提醒

- 这是一个以稳定性和显式状态管理为主的 Go 服务，不是炫技型项目。
- 优先保持行为可预测、错误可追踪、数据边界清晰。
- 如果你不确定某个风格，直接模仿正在编辑文件周围 50-100 行的既有写法，通常比套模板更正确。
