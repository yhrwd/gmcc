# 数据组件规范 (1.21.11)

## 概述

**数据组件（Data Component）** 是 Minecraft 1.20.5+ 引入的结构化数据系统，用于定义物品属性，取代传统的 NBT 标签。

- 官方文档：[Minecraft Wiki - 数据组件](https://zh.minecraft.wiki/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6)
- 协议版本：774 (1.21.11)

## 物品堆叠结构

```
ItemStack {
    id: string (命名空间ID)
    count: varint
    components: {
        component_id: value,
        "!component_id": {}  // 移除默认组件
    }
}
```

## 核心组件

### 基础属性

| 组件 | 类型 | 说明 |
|------|------|------|
| `max_stack_size` | int | 最大堆叠数量 |
| `max_damage` | int | 最大耐久度 |
| `damage` | int | 当前损坏值 |
| `unbreakable` | empty | 无法破坏 |
| `rarity` | string | 稀有度 (common/uncommon/rare/epic) |

### 显示相关

| 组件 | 类型 | 说明 |
|------|------|------|
| `custom_name` | text component | 自定义名称 |
| `item_name` | text component | 物品名称 |
| `lore` | list<text component> | 物品描述 |
| `dyed_color` | int | 染色颜色 (RGB) |
| `enchantment_glint_override` | bool | 附魔光效覆盖 |

### 附魔

| 组件 | 类型 | 说明 |
|------|------|------|
| `enchantments` | map<string, int> | 当前附魔 |
| `stored_enchantments` | map<string, int> | 存储的附魔（附魔书） |
| `enchantable` | int | 附魔能力值 |
| `repair_cost` | int | 铁砧修复成本 |

### 功能性

| 组件 | 类型 | 说明 |
|------|------|------|
| `food` | compound | 食物属性 |
| `tool` | compound | 工具属性（挖掘等级、速度） |
| `weapon` | compound | 武器属性 |
| `glider` | empty | 可滑翔 |
| `firework_explosion` | compound | 烟火之星爆炸数据 |
| `fireworks` | compound | 烟花火箭数据 |
| `potion_contents` | compound | 药水效果 |
| `bundle_contents` | list | 收纳袋内容 |
| `container` | list | 容器内容 |

### 限制相关

| 组件 | 类型 | 说明 |
|------|------|------|
| `can_break` | compound | 可破坏的方块 |
| `can_place_on` | compound | 可放置的方块 |
| `lock` | compound | 容器锁 |

### 实体相关

| 组件 | 类型 | 说明 |
|------|------|------|
| `entity_data` | compound | 实体数据（刷怪蛋等） |
| `bucket_entity_data` | compound | 生物桶中的实体 |
| `profile` | compound | 玩家档案（头颅） |

## Slot 数据格式

物品槽位在协议中的编码格式：

```
Slot {
    item_count: varint      // 如果 <= 0 表示空槽
    if item_count > 0:
        item_id: varint     // 物品注册表 ID
        num_components_to_add: varint
        for each:
            component_type: varint
            component_data: ...
        num_components_to_remove: varint
        for each:
            component_type: varint
}
```

### Hashed Slot 格式

用于 `click_container` 包，组件以 CRC32C 校验和形式发送：

```
HashedSlot {
    present: bool
    if present:
        item_id: varint
        item_count: varint
        num_components_to_add: varint
        for each:
            component_type: varint
            crc32c_hash: int  // 组件数据的 CRC32C 校验和
        num_components_to_remove: varint
        for each:
            component_type: varint
}
```

## 常用组件结构

### Tool 组件

```json
{
    "tool": {
        "rules": [
            {
                "blocks": ["minecraft:stone", "minecraft:cobblestone"],
                "speed": 8.0,
                "correct_for_drops": true
            }
        ],
        "default_mining_speed": 1.0,
        "damage_per_block": 1
    }
}
```

### Food 组件

```json
{
    "food": {
        "nutrition": 4,
        "saturation": 0.3,
        "can_always_eat": false,
        "eat_seconds": 1.6,
        "effects": [
            {
                "effect": {
                    "id": "minecraft:regeneration",
                    "duration": 100
                },
                "probability": 1.0
            }
        ]
    }
}
```

### Enchantments 组件

```json
{
    "enchantments": {
        "minecraft:sharpness": 5,
        "minecraft:unbreaking": 3
    }
}
```

## 玩家背包槽位映射

Window ID = 0（玩家背包）：

| 槽位范围 | 说明 |
|----------|------|
| 0-8 | 快捷栏 (Hotbar) |
| 9-35 | 主背包 (Main Inventory) |
| 36-39 | 盔甲 (头盔→靴子) |
| 40 | 副手 (Offhand) |

总计：46 个槽位（索引 0-45）

## 参考

- [数据组件 - Minecraft Wiki](https://zh.minecraft.wiki/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6)
- [Slot Data - Minecraft Wiki](https://minecraft.wiki/w/Java_Edition_protocol/Slot_data)
- [Tutorial:物品堆叠组件](https://zh.minecraft.wiki/w/Tutorial:%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0%E7%BB%84%E4%BB%B6)
