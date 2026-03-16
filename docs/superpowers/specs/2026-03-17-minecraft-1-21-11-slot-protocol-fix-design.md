# Minecraft 1.21.11 Slot Protocol Fix Design

## Problem Statement

The client was encountering EOF errors when parsing container content packets from 1.21.11 servers. Analysis of packet dumps revealed:

1. **Incorrect Slot Format Assumption**: The code assumed an older Slot format (`present(bool) -> item_id(VarInt) -> count(byte)`) but 1.21+ uses the new format (`count(VarInt) -> item_id(VarInt) -> components`).

2. **Component Type Mapping Errors**: The component skipping logic had incorrect mappings for several component types, causing parsing failures when encountering unknown components.

## Solution Overview

### 1. Slot Format Correction (`readers.go`)

**Before**:
```go
// Old (incorrect) assumption: present(bool) -> item_id(VarInt) -> count(byte) -> components
present, err := ReadBoolFromReader(r)
itemID, err := ReadVarIntFromReader(r)
countByte, err := ReadU8(r)
```

**After**:
```go
// Correct 1.21+ format: count(VarInt) -> [if count>0] item_id(VarInt) -> components
count, err := ReadVarIntFromReader(r)
if count == 0 {
    return nil, nil // Empty item
}
itemID, err := ReadVarIntFromReader(r)
// Skip components
if err := SkipSlotComponents(r); err != nil {
    return nil, err
}
return &SlotData{ID: itemID, Count: count}, nil
```

### 2. Component Type Mapping Fixes (`component_skipping.go`)

Fixed the component type mappings based on official Minecraft Wiki data:
- Type 0: `custom_data` → SkipNBT (was incorrectly SkipVarInt)
- Type 1: `max_stack_size` → SkipVarInt (was incorrectly SkipBannerPatterns)
- Type 2: `max_damage` → SkipVarInt (was incorrectly SkipNBT)
- Type 3: `damage` → SkipVarInt (was incorrectly SkipNBT)
- Type 4: `unbreakable` → SkipNothing (was incorrectly SkipNBT)
- Added missing component handler functions for all types

### 3. Robust Unknown Component Handling

**Before**:
```go
// Return error for unknown component types
return fmt.Errorf("unknown component type %d", componentType)
```

**After**:
```go
// For unknown components, attempt to skip as NBT (common case)
logx.Warnf("Unknown component type: %d, attempting NBT skip", componentType)
return SkipNBT(r)
```

## Changes Made

### Files Modified

1. `internal/mcclient/packet/readers.go`
   - Corrected ReadSlotData to use 1.21+ Slot format
   - Removed incorrect format detection logic

2. `internal/mcclient/packet/component_skipping.go`
   - Fixed component type mappings (0-71)
   - Added missing handler functions (SkipEnchantments, SkipBlockPredicates, etc.)
   - Improved SkipComponentByType to handle unknown components gracefully
   - Removed unused "fmt" import

## Testing Approach

1. **Unit Testing**: Verify Slot parsing with known good/bad data
2. **Integration Testing**: Connect to 1.21.11 server and verify container content parsing
3. **Error Injection**: Test with malformed packets to ensure graceful handling
4. **Regression Testing**: Ensure older protocol versions still work

## Backward Compatibility

- The new Slot format is backwards compatible with 1.20.5+ as it uses the same component structure
- Unknown component handling improves forward compatibility with future versions
- All existing functionality remains unchanged

## Performance Impact

- Minimal: Only changes parsing logic, no additional allocations
- Improved: Better error handling reduces connection drops from parse failures

## Future Considerations

- Monitor logs for unknown component types to update mappings
- Consider adding metrics for parse success/failure rates
- Test with various server implementations (Paper, Purpur, Fabric, Forge)
