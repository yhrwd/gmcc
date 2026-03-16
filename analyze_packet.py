#!/usr/bin/env python3
import sys

# 从用户提供的数据中提取实际二进制数据
# 跳过头部注释行

lines = '''00 03 2E 00 00 00 00 01 BA 07 09 00 09 0A 08 00 
05 63 6F 6C 6F 72 00 07 23 39 35 38 44 44 33 08 00 04 74 65 78 74 00 0F E6 97 A3 E9 A3 8E E6 8A A4 E7 9B AE E9 95 9C 00'''.strip().split('\n')

# 合并所有 hex
hex_data = ''.join(lines).replace(' ', '')
data = bytes.fromhex(hex_data)

print("解析 slot 数据:")
print(f"前 {len(data)} 字节")
print()

idx = 0

# 假设格式：
# container_content:
#   window_id: varint
#   state_id: varint
#   slot_count: varint
#   slots[]...

# 查看开头的 varint
def parse_varint(data, idx):
    result = 0
    shift = 0
    while True:
        b = data[idx]
        result |= int(b & 0x7f) << shift
        idx += 1
        if (b & 0x80) == 0:
            break
        shift += 7
    return result, idx

# window_id
val, idx = parse_varint(data, idx)
print(f"window_id: {val} (offset after: {idx})")

# state_id
val, idx = parse_varint(data, idx)
print(f"state_id: {val} (offset after: {idx})")

# slot_count
val, idx = parse_varint(data, idx)
print(f"slot_count: {val} (offset after: {idx})")

print()
print("解析第一个物品 slot[0]:")

# Slot 格式 (pre-1.20.5):
#   如果是空槽: 0x00 (count 为 0 表示空)
#   否则: count (VarInt), item_id (VarInt), NBT
#
# Slot 格式 (1.20.5+):
#   如果是空槽: count (VarInt) <= 0
#   否则: count (VarInt), item_id (VarInt), components

# 尝试读取 count
count, idx = parse_varint(data, idx)
print(f"  count: {count} (offset after: {idx})")

if count > 0:
    # item_id
    item_id, idx = parse_varint(data, idx)
    print(f"  item_id: {item_id} (offset after: {idx})")
    
    # 组件
    num_add, idx = parse_varint(data, idx)
    print(f"  num_add_components: {num_add} (offset after: {idx})")
    
    print(f"  下一个字节: {data[idx]:02x} at offset {idx}")
    
    # 解析组件
    for comp_idx in range(min(num_add, 5)):  # 只解析前几个
        comp_type, idx = parse_varint(data, idx)
        print(f"    component[{comp_idx}] type: {comp_type} (offset: {idx})")
        
        # 根据 type 跳过数据
        # type 0 = custom_data (NBT)
        if comp_type == 0:
            # NBT
            print(f"      下一个字节 (应该是NBT tag): {data[idx]:02x} at offset {idx}")
            break
        break

print()
print("原始十六进制:")
for i in range(min(100, len(data))):
    if i % 16 == 0:
        print(f"{i:04x}:", end=" ")
    print(f"{data[i]:02x}", end=" ")
    if i % 16 == 15:
        print()
print()