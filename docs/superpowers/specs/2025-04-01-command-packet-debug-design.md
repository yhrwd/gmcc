# Command Packet Debug Enhancement Design

**Date**: 2025-04-01  
**Topic**: Command Packet Debug Logging  
**Status**: Approved

## Summary

Enhance debugging output for command sending to help diagnose signature-related issues like `chat.disabled.invalid_signature` and `chat.disabled.chain_broken`.

## Motivation

Currently when sending commands, only basic log messages are shown:
- `已发送命令: /msg YHRWD Hello`
- `命令包内容 (hex): 0d6d736720...`

This provides insufficient information to debug why signed commands fail on servers with secure chat enabled. We need visibility into the full packet structure including timestamps, salts, and signatures.

## Design

### Overview

Add a structured logging function `logCommandPacket` that outputs detailed command packet information before sending. The debug output is only active when `debug: true` is configured.

### Changes

**File**: `internal/mcclient/chat.go`

1. Add `logCommandPacket` function to format and output command packet details
2. Enhance `sendSignedCommand` to log structured packet information
3. Keep `sendUnsignedCommand` with existing hex dump (simpler structure)

### Debug Output Format

**Signed Command**:
```
[DEBUG] 发送签名命令包:
  命令: msg YHRWD Hello
  时间戳: 1743495600123
  盐值: 9876543210
  参数签名: 1个
    - message: [256 bytes]
  消息计数: 0
  确认位: 000000
  校验和: 1
  完整包体 (hex): 0d6d736720...
```

**Unsigned Command**:
```
[DEBUG] 发送无签名命令:
  命令: msg YHRWD Hello
  包体长度: 13 bytes
  完整包体 (hex): 0d6d736720...
```

### Implementation Details

The `logCommandPacket` function receives:
- `isSigned bool` - whether this is a signed or unsigned command
- `cmd string` - the command text
- `timestamp int64` - Unix timestamp in milliseconds
- `salt int64` - random salt value
- `signatures []commandArgumentSignature` - list of argument signatures
- `payload []byte` - the complete serialized packet

For signed commands, it logs each signature's argument name and byte length. For unsigned commands, it only logs the command text and hex dump.

Hex output is truncated to first 64 bytes to prevent log spam from large payloads.

### Error Handling

- Debug logging failures are silently ignored (using `_ = logx.Debugf` pattern)
- This ensures a logging error doesn't break command sending

## Testing

Test cases to verify:
1. Debug output appears when `debug: true` in config
2. Debug output is suppressed when `debug: false`
3. Signed command shows all fields correctly
4. Unsigned command shows simplified output
5. Hex output is properly truncated for large payloads

## Future Considerations

- Could extend to log chat message packets if needed
- Could add packet receipt debugging in handlers
- Could support configurable hex dump length

---

**Approved by**: User  
**Implementation**: Ready for writing-plans skill
