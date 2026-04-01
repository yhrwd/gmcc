# Chat State Synchronization Fix Design

**Date**: 2025-04-01  
**Topic**: Secure Chat State Machine Implementation  
**Status**: Approved

## Summary

Fix `chat.disabled.chain_broken` and `chat.disabled.invalid_signature` errors by implementing the complete chat state synchronization system, including `lastSeenMessages` tracking, `acknowledged` bitset, and proper message chain maintenance.

## Background

Minecraft 1.19+ requires clients to maintain a chat state machine for secure chat. The client must:
1. Track received signed messages (up to 20)
2. Acknowledge these messages when sending new messages
3. Include acknowledged message signatures in the signing context

Currently, the client sends `messageCount=0` and `acknowledged=[0,0,0]`, causing servers to reject messages.

## Problem Analysis

**Current state** (`chat.go`):
- `sendSignedCommand` always sends `messageCount=0`
- `acknowledged` is always `[0x00, 0x00, 0x00]`
- `buildChatSignableBody` receives empty acknowledgements
- `player_chat` handler does not store received messages

**Root cause**: Missing `lastSeenMessages` tracking causes the message chain to break.

## Design

### Data Structures

```go
// lastSeenMessage stores a signed message identifier
type lastSeenMessage struct {
    signature  [256]byte  // Message signature (256 bytes for RSA)
    senderUUID [16]byte   // Sender UUID
    index      int32     // Message index in sender's chain
}

// lastSeenMessageBuffer is a fixed-size ring buffer (max 20 messages)
type lastSeenMessageBuffer struct {
    messages [20]lastSeenMessage
    head     int        // Next write position
    count    int        // Current message count
    mu       sync.Mutex // Thread safety
}
```

**Methods**:
- `Add(msg lastSeenMessage)` - Add message, evict oldest if full
- `GetAll() []lastSeenMessage` - Return all messages (oldest to newest)
- `Len() int` - Return current count

### Client State Extension

Add to `Client` struct:
```go
lastSeenBuf  lastSeenMessageBuffer  // Track received signed messages
acknowledged [3]byte               // 24-bit bitset for acknowledged messages
```

### Message Receiving Flow

**Location**: `handlers_chat.go` - `handlePlayerChatPacket`

When receiving `player_chat`:
1. Parse `signature`, `senderUUID`, and `index`
2. Lock `lastSeenBuf`
3. Call `lastSeenBuf.Add(lastSeenMessage{...})`
4. Unlock

### Message Sending Flow

**Location**: `chat.go` - `sendSignedCommand` and `SendMessage`

Before building the signable body:
1. Lock `lastSeenBuf`
2. Get `messages := lastSeenBuf.GetAll()`
3. `messageCount := len(messages)`
4. Build `acknowledged` bitset (all 1s for now, or based on actual acknowledgements)
5. Extract signatures: `ackSignatures := make([][]byte, len(messages))`
6. For each message: `ackSignatures[i] = msg.signature[:]`
7. Unlock

**Signed Command Payload** (correct format):
```
[Command String]
[Timestamp (int64)]
[Salt (int64)]
[Argument Signatures Count (VarInt)]
  [Arg Name 1 (String)]
  [Signature 1 (256 bytes)]
[Message Count (VarInt)]
[Acknowledged (3 bytes bitset)]
[Checksum (1 byte)]
```

**Signing Context** (for `buildChatSignableBody`):
```go
signable := [
    0x01 (version),                              // 4 bytes
    playerUUID,                                  // 16 bytes
    sessionID,                                 // 16 bytes
    messageIndex,                                // 4 bytes
    salt,                                        // 8 bytes
    timestamp,                                   // 8 bytes
    content_length (int32),                      // 4 bytes
    content (UTF-8 bytes),                       // variable
    acknowledgement_count (int32),               // 4 bytes
    acknowledged[0].signature (256 bytes),        // 256 bytes each
    acknowledged[1].signature (256 bytes),
    ...
]
```

### Acknowledged Bitset Format

The `acknowledged` field is a 24-bit bitset packed into 3 bytes:
- Bit 0 (LSB): Acknowledged status of `lastSeenMessages[0]`
- Bit 1: Acknowledged status of `lastSeenMessages[1]`
- ...
- Bit 23: Acknowledged status of `lastSeenMessages[23]` (max 20 used)

For initial implementation, all bits can be set to 1 (acknowledging all tracked messages).

## Implementation Files

### Modified Files

1. **`internal/mcclient/chat.go`**
   - Add `lastSeenMessage` and `lastSeenMessageBuffer` types
   - Add `lastSeenBuf` and `acknowledged` to `Client` struct
   - Implement ring buffer methods
   - Update `sendSignedCommand` to use tracked messages
   - Update `SendMessage` to use tracked messages

2. **`internal/mcclient/handlers_chat.go`**
   - Update `handlePlayerChatPacket` to extract and store message info

### New Functions

```go
// chat.go
func (b *lastSeenMessageBuffer) Add(msg lastSeenMessage)
func (b *lastSeenMessageBuffer) GetAll() []lastSeenMessage
func (b *lastSeenMessageBuffer) Len() int
func (c *Client) trackReceivedMessage(signature [256]byte, senderUUID [16]byte, index int32)
func (c *Client) getAcknowledgements() (count int32, bitset [3]byte, signatures [][]byte)
```

## Testing

**Unit Tests**:
1. `lastSeenMessageBuffer.Add` - Verify FIFO eviction after 20 messages
2. `lastSeenMessageBuffer.GetAll` - Verify correct order
3. `getAcknowledgements` - Verify bitset and signatures

**Integration Tests**:
1. Connect to server with secure chat enabled
2. Send command after receiving player_chat
3. Verify no `chain_broken` error
4. Verify message is accepted by server

## Migration Notes

This is a breaking change for the chat subsystem. Existing code will continue to work (just without proper acknowledgements), but will be fixed once this implementation is complete.

## References

- Minecraft Protocol Wiki: Player Chat packet
- Secure Chat protocol documentation
- Current implementation: `internal/mcclient/chat.go`, `internal/mcclient/handlers_chat.go`

---

**Approved by**: User  
**Implementation**: Ready for writing-plans skill
