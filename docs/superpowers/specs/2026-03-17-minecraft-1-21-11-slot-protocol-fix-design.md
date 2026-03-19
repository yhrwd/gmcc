# Minecraft 1.21.11 Reader And Play Protocol Alignment Design

## Context

The current `v774` implementation is only partially aligned with Minecraft `1.21.11`. The repository already moved `ItemStack` parsing to the `1.21+` shape, but `.knowledge/1.21.11` shows that the remaining protocol drift is broader than a Slot-only fix:

- `types/components.json` defines `104` component raw IDs, while `internal/mcclient/packet/component_skipping.go` only covers `0-71`.
- The existing component mapping is based on an older layout. From raw ID `5` onward, many handlers are shifted and no longer match `1.21.11`.
- Several active play handlers still use older packet IDs or older field order assumptions, which can desynchronize parsing even when `ReadSlotData` is correct.
- Container parsing is currently disabled behind `DEBUG_DUMP_CONTAINER_PACKETS`, which confirms that the reader layer is not yet reliable enough for live inventory updates.

This spec supersedes the earlier narrow "slot protocol fix" framing. The real task is to align the reader layer, protocol constants, and the currently used play handlers with `.knowledge/1.21.11`.

## Goals

1. Make the packet reader layer consume `1.21.11` ItemStack/component payloads without drifting.
2. Align `internal/mcclient/protocol/v774.go` with the packet IDs currently used by the client.
3. Align active play handlers with the `1.21.11` field order described in `.knowledge/1.21.11`.
4. Re-enable container parsing after the underlying reader path is trustworthy.

## Non-Goals

- Full implementation of all `139` play clientbound packets.
- A generic multi-version runtime schema engine.
- Rewriting the chat/text stack if the current `readAnonymousNBTJSON` path still matches the packets we already parse successfully.
- Touching unrelated user edits outside the protocol-alignment scope.

## Options Considered

### Option A: Static alignment from `.knowledge/1.21.11` (Recommended)

Use `.knowledge/1.21.11` as the authoritative source and manually align:

- `protocol/v774.go`
- `packet/readers.go`
- `packet/component_skipping.go`
- active handlers in `handlers_container.go`, `handlers_player.go`, and `handlers_play.go`

Pros:

- Smallest change set that fixes the current client.
- Easy to validate against packet dumps and current runtime behavior.
- Keeps the existing code structure intact.

Cons:

- Still version-specific.
- Requires manual upkeep when the next protocol changes.

### Option B: Schema-driven component and packet skipping

Build a generic skipper/reader from `.knowledge/1.21.11` structures.

Pros:

- More future-proof.
- Less manual mapping logic once complete.

Cons:

- Much larger scope.
- Harder to verify in one pass.
- Not necessary to restore the current client.

### Option C: Keep partial static mapping and rely on loose fallback

Continue the current approach of mapping some IDs and falling back to `SkipNBT` for unknown components.

Pros:

- Lowest immediate effort.

Cons:

- Unsafe. Many `1.21.11` components are not NBT-backed.
- Likely to keep inventory parsing flaky and hard to debug.

Option A is the recommended path because it fixes the actual `1.21.11` client without turning this task into a protocol framework rewrite.

## Source Of Truth

The implementation should treat the following files as authoritative for this scope:

- `.knowledge/1.21.11/packets/play_clientbound.json`
- `.knowledge/1.21.11/types/components.json`
- `.knowledge/1.21.11/summary.json`

The generated dataset is more trustworthy than the existing handwritten constants whenever they disagree.

## Design

### 1. Reader Layer Alignment

`internal/mcclient/packet/readers.go` remains the central entry point for packet-level primitive reads and `ItemStack` parsing.

The `ItemStack` wire format should remain:

`count -> [if count > 0] item_id -> merged components`

What changes in this phase is not the overall `ItemStack` shape, but the correctness of the component payload reader/skipper that follows it.

Planned reader-layer adjustments:

- Keep `ReadSlotData` focused on extracting `count`, `itemID`, and advancing the reader over merged components.
- Keep helper readers small and reusable. Add helper skip/read functions only when a `1.21.11` component shape genuinely needs one.
- Continue reusing `ReadAnonymousNBTJSON` for packets that are already known to carry text components in the same form as current chat packets.

### 2. Component Skipper Realignment

`internal/mcclient/packet/component_skipping.go` must be realigned to the `rawId` ordering from `.knowledge/1.21.11/types/components.json`.

This is the highest-priority change because every container-related packet depends on it.

The implementation should:

- Replace the current outdated `rawId -> skipper` table with the `1.21.11` order.
- Add explicit handlers for the new or shifted component groups that the current table does not represent correctly.
- Cover all `rawId` values from `0` through `103`.

Important examples already confirmed from `.knowledge/1.21.11`:

- `5` is `minecraft:use_effects`, no longer `custom_name`
- `30` is `minecraft:attack_range`
- `38` is `minecraft:piercing_weapon`
- `39` is `minecraft:kinetic_weapon`
- `40` is `minecraft:swing_animation`
- `72` is `minecraft:pot_decorations`
- `73` is `minecraft:container`
- `74` is `minecraft:block_state`
- `75` is `minecraft:bees`
- `78` is `minecraft:break_sound`
- `79-103` are mostly registry-backed or enum-backed variant/color components

Fallback behavior also needs to change:

- Do not assume that an unknown component can be skipped as NBT.
- Unknown `rawId` values should surface a clear parsing error with enough context to identify the slot and component ID.
- Debug dumping can still be used as an investigation aid, but not as the steady-state parsing strategy.

### 3. Protocol Constant Audit

`internal/mcclient/protocol/v774.go` should be audited against `.knowledge/1.21.11/packets/play_clientbound.json` for every clientbound play packet the client actively switches on today.

Confirmed mismatches include:

- `PlayClientOpenScreen`: current `0x0D`, expected `0x39`
- `PlayClientGameEvent`: current `0x22`, expected `0x26`
- `PlayClientEntityData`: current `0x5D`, expected `0x61`
- `PlayClientPlayerInfoUpdate`: current `0x42`, expected `0x44`
- `PlayClientPlayerInfoRemove`: current `0x3D`, expected `0x43`
- `PlayClientSetHeldSlot`: current `0x2F`, expected `0x67`

The packet-name maps in the same file should be updated together with the numeric constants so debug logging stays trustworthy.

### 4. Active Handler Alignment

The following handlers are already wired into `handlePlayPacket` and should be aligned with `1.21.11` packet order and shape.

#### `handleSetHealthPacket`

`.knowledge` shows the logical order as:

`food -> health -> saturation`

The implementation should follow the `1.21.11` field order instead of the current older assumption.

#### `handleSetExperiencePacket`

`.knowledge` shows:

`barProgress -> experience -> experienceLevel`

The current code reads `level` before `totalExp`, so this handler should be reordered to match the dataset.

#### `handlePlayerAbilitiesPacket`

`.knowledge` shows a boolean-heavy structure:

- `allowFlying`
- `creativeMode`
- `flySpeed`
- `flying`
- `invulnerable`
- `walkSpeed`

This should replace the current old-style bit-flag parser.

#### `handlePlayLoginPacket`

The current login parser is close in overall shape, but the spec for `minecraft:login` should be rechecked against the nested `commonPlayerSpawnInfo` structure from `.knowledge`.

This phase should:

- keep the current staged parsing style
- align the field order with `commonPlayerSpawnInfo`
- avoid changing unrelated player-state behavior

#### `handleOpenScreenPacket`

`.knowledge` shows the packet fields as:

- `name` (`text_component`)
- `screenHandlerId`
- `syncId`

The current implementation still assumes the older order:

- `windowID`
- `windowType`
- `title`

This handler must be updated, otherwise the reader will desynchronize before any container state is processed.

#### `handleContainerClosePacket`

`.knowledge` exposes a single `syncId`. The handler should be aligned to that structure and kept simple.

#### `handleContainerSetDataPacket`

`.knowledge` shows:

- `propertyId`
- `syncId`
- `value`

The current implementation reads `windowID` first, so this handler must be reordered.

#### `handleContainerContentPacket` and `handleContainerSlotPacket`

These handlers should stay thin wrappers around `ReadSlotData`, but once the component skipper is aligned:

- the debug-only dump shortcut should be removable
- state ID and slot updates should parse without packet drift

### 5. Text Component Handling

For packets already handled successfully as text components, the implementation should keep using the existing `readAnonymousNBTJSON` path unless packet dumps or tests show that a specific packet uses a different wire encoding.

That means:

- `system_chat`
- `action_bar`
- `profileless_chat`
- `open_screen` title/name

should use one consistent text decoding path instead of mixing string, NBT, and ad-hoc logic per handler.

### 6. Error Handling And Diagnostics

Critical packet readers should prefer explicit parse errors over silent drift.

Design rules:

- If a component or field cannot be skipped/read reliably, return an error instead of guessing.
- Packet logs should include packet name, slot index, and component raw ID where applicable.
- Keep dump-to-file support as an opt-in debugging aid, not the default runtime path.
- Do not overwrite or rollback unrelated user changes while aligning these files.

## Files In Scope

- `internal/mcclient/packet/readers.go`
- `internal/mcclient/packet/component_skipping.go`
- `internal/mcclient/protocol/v774.go`
- `internal/mcclient/handlers_container.go`
- `internal/mcclient/handlers_player.go`
- `internal/mcclient/handlers_play.go`
- `internal/mcclient/packet/packet_test.go`
- any new packet-reader tests added for `1.21.11`

## Testing Strategy

### Unit Tests

- Add focused tests for `ReadSlotData` with representative `1.21.11` component layouts.
- Add tests for newly introduced component skipper helpers.
- Add regression tests for packet field order where synthetic payloads are easy to build.

### Packet-Dump Validation

Use real dumps from `container_set_content` and `container_set_slot` to confirm:

- no unexpected EOF
- reader alignment reaches the end of the packet
- carried slot and inventory slots decode consistently

### Integration Validation

Connect to a `1.21.11` server and verify:

- play login completes
- inventory opens without dump fallback
- hotbar selection updates correctly
- health / experience / abilities / game mode updates still reach `Player`

## Risks

1. `.knowledge` describes packet field order reliably, but some `int`-typed fields may still need confirmation at the wire-codec level when choosing between fixed-width and VarInt readers.
2. Some registry-backed component payloads are represented abstractly in the dataset and still require careful low-level skipping logic.
3. `handlers_player.go` already has user edits in the worktree, so changes there must be merged carefully instead of overwritten.

## Success Criteria

This work is complete when all of the following are true:

- `component_skipping.go` matches the `1.21.11` component table through `rawId 103`
- `v774.go` uses the correct packet IDs for the play handlers currently wired in the client
- `open_screen`, container, health, experience, abilities, and game-event parsing follow the `1.21.11` layout
- container parsing can run without the permanent dump shortcut
- tests or packet-dump validation show that packet parsing no longer drifts on `1.21.11`
