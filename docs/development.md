# Development Notes

## Goal

This project targets **Java Edition protocol 774** and keeps the runtime small:

- connect
- login/auth
- enter play
- AFK heartbeat

The design follows the same ideas used by `node-minecraft-protocol`:

- split client pipeline into auth / encrypt / compress / play stages
- treat protocol behavior as feature-gated (not hardcoded for all versions)
- keep logs protocol-aware for fast issue triage

## Package Guide

### `cmd/gmcc`

- Program entry.
- Loads config, initializes logger, starts `mcclient.Client`.

### `internal/config`

- `Load(path)` reads and validates YAML.
- `Default()` outputs baseline config.
- Main fields:
  - `account`: account mode and player id.
  - `server`: host and port.
  - `actions`: on-join commands/messages.
  - `log`: runtime logging behavior.

### `internal/auth/microsoft`

- Handles device-code login and refresh token flow.
- Exposes XSTS conversion helpers for reuse in session pipeline.

### `internal/auth/minecraft`

- Exchanges XSTS token for Minecraft access token.
- Verifies ownership and calls session join.

### `internal/session`

- Token cache in `.session/<player>.json`.
- Used in order: cached MC token -> MS refresh -> device code.

### `internal/mcclient`

- `codec.go`: frame read/write, compression, AES/CFB8 stream.
- `protocol_774.go`: packet IDs, protocol feature flags, packet-name map.
- `client.go`: login/config/play state machine and AFK loop.
- `chat.go`: message/command send + chat packet parsing + JSON extraction.

### `internal/logx`

- Unified stdout/file logger.
- Debug mode logs packet-level details.

### `pkg/httpx`

- Thin HTTP helpers (`GET/POST form/POST json`) with status-aware errors.

### `pkg/rwfile`

- Generic file read/write helpers (string/bytes/gob).

### `pkg/cryptox`

- Generic crypto helpers (not protocol-specific login codec).

## Runtime Pipeline

The end-to-end flow is:

1. Load config (`internal/config`).
2. Resolve account mode (`use_official_auth`):
   - offline only, or
   - Microsoft flow with cache -> refresh -> device code.
3. Connect TCP and send handshake/login start.
4. Handle login encryption:
   - call session join
   - send `encryption_begin` response
   - enable AES/CFB8 stream.
5. Handle login compression and login success.
6. Enter configuration state, ack required packets.
7. Enter play state and keep alive/AFK loop.
8. On first play login:
   - execute `actions.on_join_commands`
   - execute `actions.on_join_messages`.

## Protocol Metadata

`internal/mcclient/protocol_774.go` contains:

- packet IDs for Login / Configuration / Play
- protocol feature flags (`features774`)
- packet name maps (`packetName`) for debug logs
- state name helper (`stateName`)

When upgrading protocol, update this file first.

## Encryption Compatibility

Current implementation uses the legacy serverbound key packet layout:

- `sharedSecret`
- `verifyToken`

For protocol 774 in current data source, this is correct.
If a future version enables signature-based encryption response, implement a new branch where `features774.SignatureEncryption=true`.

## Message and Command I/O

### Receive

Current play-state receive handlers:

- `system_chat` (`0x77`)
- `player_chat` (`0x3F`)
- `action_bar` (`0x55`)
- `profileless_chat` (`0x21`)

For chat components carried as `anonymousNbt`, client decodes them with network NBT format and exposes **raw JSON string** (`ChatMessage.RawJSON`) so plugin-formatted messages can be parsed downstream.

`ChatMessage` fields:

- `Type`: packet category (`system`, `player_chat`, `action_bar`, etc.)
- `PlainText`: extracted text fallback
- `RawJSON`: raw chat JSON from NBT component
- `SenderUUID`: sender UUID when available
- `ReceivedAt`: timestamp

Register callback:

```go
client.SetChatHandler(func(msg mcclient.ChatMessage) {
    // parse msg.RawJSON with your plugin/parser logic
})
```

### Send

Available APIs in `internal/mcclient/chat.go`:

- `SendCommand(command string)` -> serverbound `chat_command_signed` (`0x07`, signed-lite), fallback `chat_command` (`0x06`)
- `SendMessage(message string)` -> serverbound `chat_message` (`0x08`, unsigned)
- client auto sends `chat_session_update` (`0x09`) after first play login to provide profile public key for secure-chat state

Config-driven auto send:

```yaml
actions:
  delay_ms: 1200
  on_join_commands:
    - "list"
  on_join_messages:
    - "hello from gmcc"
```

Notes:

- `on_join_commands` should not include leading `/` (code trims it anyway).
- `delay_ms` controls when on-join actions start (default 1200ms) to reduce race with secure-chat session activation.
- Unsigned chat messages may be rejected by servers enforcing strict secure-chat policy.
- Current command signing is **signed-lite** (no argument signatures). Some plugins/servers may still enforce full signatures for specific arguments.

## Logging and Diagnostics

`internal/logx` writes to stdout and `logs/gmcc.log`.

`log.debug=true` enables protocol diagnostics:

- outgoing packet id/length/preview
- incoming frame length/compression details
- packet id + symbolic packet name + state
- encryption request/response details

Console format is time-only (`HH:MM:SS`), while file logs keep full datetime.
File logs rotate by size (`log.max_size`) into chunk files `gmcc-*.log`.

These logs are enough to pinpoint:

- wrong packet structure
- encryption desync
- compression threshold issues
- unexpected state transitions

## Extension Checklist

When adding packet support or upgrading version:

1. Update packet IDs and feature flags in `protocol_774.go`.
2. Validate encryption packet schema against source protocol data.
3. Add packet handling branch in `client.go`.
4. Add debug logging for new packet path.
5. Add tests for codec/protocol edge cases when possible.
6. Run:
   - `go test ./...`
   - `go build ./cmd/gmcc`

## External References

- node-minecraft-protocol:
  - https://github.com/PrismarineJS/node-minecraft-protocol
  - https://raw.githubusercontent.com/PrismarineJS/node-minecraft-protocol/master/src/createClient.js
  - https://raw.githubusercontent.com/PrismarineJS/node-minecraft-protocol/master/src/client/encrypt.js
- protocol schema source used for verification:
  - https://raw.githubusercontent.com/PrismarineJS/minecraft-data/master/data/pc/1.21.11/protocol.json
