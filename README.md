# gmcc

Go: Console Minecraft Client

## Structure

```text
cmd/gmcc
  main.go                 # CLI entry

internal/auth/microsoft
  service.go              # Microsoft/Xbox/XSTS token workflow

internal/auth/minecraft
  service.go              # Minecraft login and ownership check

internal/config
  config.go               # Config model + defaults
  loader.go               # YAML loading/validation

internal/mcclient
  protocol_774.go         # Packet IDs for 1.21.11 (protocol 774)
  codec.go                # Packet framing/compression/encryption codec
  client.go               # Login/Configuration/Play state machine + AFK

pkg/httpx
  client.go               # Reusable HTTP helper

pkg/rwfile
  string.go               # String/bytes file helpers
  gob.go                  # Gob encode/decode helpers

pkg/cryptox
  crypto.go               # AES/RSA helpers
```

## Packages

- `internal/config`: 读取/校验 `config.yaml`。
- `internal/auth/microsoft`: Microsoft/Xbox/XSTS 登录链路。
- `internal/auth/minecraft`: Minecraft token、ownership、session join。
- `internal/session`: 本地 token 缓存（`.session/<player>.json`）。
- `internal/mcclient`: 协议 774 的连接、编解码、状态机、聊天收发。
- `internal/logx`: 控制台/文件日志。
- `pkg/httpx`: HTTP 请求封装。
- `pkg/rwfile`: 文件读写工具。
- `pkg/cryptox`: 通用加解密工具。

## Run

```bash
go run ./cmd/gmcc
```

## Log

- 控制台与文件双写日志（可通过 `log.enable_file` 控制）。
- 控制台仅输出时间（`HH:MM:SS`），不显示日期。
- 默认日志文件为 `logs/gmcc.log`，达到 `log.max_size` 后按大小分块滚动为 `logs/gmcc-*.log`。
- `log.debug=true` 时会输出协议级调试信息（包 ID、状态、frame 长度、加密阶段细节）。

## Development

- 开发说明与扩展清单：`docs/development.md`

## Config

```yaml
account:
  player_id: "YourName"
  use_official_auth: true
actions:
  delay_ms: 1200
  on_join_commands: ["list"]
  on_join_messages: []
```

- `use_official_auth=false`: 仅使用离线模式。
- `use_official_auth=true`: 启用 Microsoft 正版认证流程。
- `actions.on_join_commands`: 进入 Play 后自动发送命令（不带 `/`）。
- `actions.on_join_messages`: 进入 Play 后自动发送聊天消息。
- `actions.delay_ms`: 入服后动作延迟，默认 1200ms（让 secure chat 会话先同步）。

## Implemented Flow (1.21.11 / 774)

- Handshake -> Login -> Configuration -> Play
- Login compression + AES/CFB8 encryption
- Online auth workflow (`use_official_auth=true`):
  - Reuse cached Minecraft token from `.session/<player>.json`
  - Fallback to Microsoft `refresh_token` refresh
  - Fallback to device-code login if refresh is unavailable/invalid
  - Complete `sessionserver join` during encrypted login
- Runtime responses for AFK:
  - keep_alive / ping / teleport confirm
  - cookie request / resource-pack status
  - periodic `move_player_status_only` heartbeat
- Chat and command support:
  - receive `system_chat` / `player_chat` / `action_bar`
  - extract raw chat JSON for plugin-formatted message parsing
  - auto send `chat_session_update` (profile public key) after entering play
  - send command via `chat_command_signed` (signed-lite fallback to `chat_command`), send message via `chat_message`
