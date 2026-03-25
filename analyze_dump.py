#!/usr/bin/env python3
import sys

# 读取dump文件
with open('logs/errors/20260325-012206_container_content.bin', 'rb') as f:
    data = f.read()

# 解析VarInt的函数
def read_varint(data, offset):
    result = 0
    for i in range(5):
        if offset >= len(data):
            return None, offset
        b = data[offset]
        offset += 1
        result |= (b & 0x7F) << (i * 7)
        if b & 0x80 == 0:
            break
    return result, offset

# 解析前几个字段
offset = 0
containerId, offset = read_varint(data, offset)
print(f"ContainerId: {containerId} at offset {offset - 1}")

stateId, offset = read_varint(data, offset)
print(f"StateId: {stateId} at offset {offset - 1}")

numItems, offset = read_varint(data, offset)
print(f"NumItems: {numItems} at offset {offset - 1}")

print(f"当前位置偏移: {offset}")

# 尝试解析slot 6
slots_parsed = 0
while slots_parsed < 6 and offset < len(data):
    print(f"\n=== Slot {slots_parsed} (offset {offset}) ===")
    count, offset = read_varint(data, offset)
    if count is None:
        print("读取count失败")
        break
    
    print(f"  Count: {count}")
    
    if count == 0:
        print("  空物品")
        slots_parsed += 1
        continue
        
    itemId, offset = read_varint(data, offset)
    print(f"  ItemId: {itemId}")

    # 跳过组件 - 但这里可能会出错
    try:
        # 读取添加的组件数量
        numAdd, offset = read_varint(data, offset)
        print(f"  添加的组件数量: {numAdd}")
        
        # 读取移除的组件数量
        numRemove, offset = read_varint(data, offset)
        print(f"  移除的组件数量: {numRemove}")
        
        # 这里应该跳过添加的组件数据，但我们暂时不处理，因为问题可能在这之前就发生了
    except Exception as e:
        print(f"  解析组件时出错: {e}")
        break
    
    slots_parsed += 1
    
print(f"\nSlot 6在偏移 {offset} 处，附近的数据:")
# 显示slot 6附近的数据
for i in range(max(0, offset - 10), min(len(data), offset + 10)):
    print(f"  {i:04d}: {data[i]:02x} ({chr(data[i]) if 32 <= data[i] <= 126 else '.'})")