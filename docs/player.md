# 玩家数据文档

## 概述

本文档描述如何从服务器获取和解码玩家状态、位置和背包信息。

## 玩家状态获取

### 游戏登录包 (0x30 play_client_login)

登录成功后，服务器发送 `login` 包，包含玩家初始状态。

#### 字段结构

```
int32        entity_id
bool         is_hardcore
byte         gamemode
byte         previous_gamemode
array        world_names
nbt         registry_codec
string       max_world_size
string       world_name
long         hashed_seed
int32        max_players
int32        view_distance
int32        simulation_distance
bool         reduced_debug_info
bool         respawn_screen
bool         is_debug
bool         is_flat
byte         last_death_location_exists
optional:{
    string    dimension_name
    string    biome
    position  location
}
```

### 生命值和饥饿值 (0x54 play_client_set_health)

```
float32      health
int32        food
float32      saturation
```

### 经验值 (0x4D play_client_set_experience)

```
float32      experience_bar
int32        level
int32        total_experience
```

### 游戏状态 (0x21 play_client_game_event)

```
int8         event_type
float32      value

event_type:
  0 = INVALID_BED
  1 = END_RAIN
  2 = END_RAIN
  3 = CHANGE_GAMEMODE
  4 = ENTER_CREDITS
  5 = DEMO_MESSAGE
  6 = ARROW_HIT_PLAYER
  7 = RAIN_LEVEL_CHANGE
  8 = THUNDER_LEVEL_CHANGE
  9 = PUFFERFISH_STING
  10 = ELDER_GUARDIAN_EFFECT
  11 = ENABLE_RESPAWN_SCREEN
```

### 玩家能力 (0x36 play_client_player_abilities)

```
int8          flags
float32       flying_speed
float32       field_of_view_modifier

flags:
  0x01 = invulnerable
  0x02 = flying
  0x04 = can_fly
  0x08 = instant_break
```

## 位置获取

### 玩家位置和传送 (0x46 play_client_player_position)

服务器传送玩家时发送：

```
double       x
double       y
double       z
float32      yaw
float32      pitch
int8         relative_arguments
int32        teleport_id
```

客户端必须回复 `accept_teleportation` (0x00):

```
int32        teleport_id
```

### 移动同步

客户端发送位置更新：

- `move_player_pos` (0x16): 只有位置变化
- `move_player_pos_rot` (0x18): 位置和角度变化
- `move_player_rot` (0x17): 只有角度变化
- `move_player_status_only` (0x20): 仅用于 AFK 检测

### 坐标系统

```
X: 东西方向（东正西负）
Y: 垂直方向（上正下负）
Z: 南北方向（南正北负）

Yaw:   水平角度（0=南，90=西，180=北，270=东）
Pitch: 俯仰角度（-90=上，90=下）
```

## 背包系统

### 设置物品槽 (0x2F play_client_set_held_slot)

```
varint       slot  // 0-8, 快捷栏索引
```

### 容器内容 (0x12 play_client_container_set_content)

```
varint       window_id
varint       state_id
varint       count         // 物品槽位数组长度
slot[]       slots         // 所有槽位的物品数据
slot         carried_item  // 光标物品
```

### 容器槽位变化 (0x16 play_client_container_set_slot)

```
varint       window_id
varint       state_id
int16        slot          // 槽位索引
slot         item_data     // 物品数据
```

### 光标物品 (0x5E play_client_set_cursor_item)

```
slot         carried_item  // 光标物品数据
```

### 容器类型

| Window ID | 容器类型 |
|-----------|---------|
| 0 | 玩家背包 |
| 1-99 | 动态容器 ID |

### 玩家背包槽位 (Window ID = 0)

```
槽位 0-8:   快捷栏 (Hotbar)
槽位 9-35:  主背包 (Main Inventory)
槽位 36-39: 盔甲 (Armor: 头盔→靴子)
槽位 40:    左手/副手 (Offhand)
```

注意: 玩家背包总共有 46 个槽位 (索引 0-45)，但窗口大小因版本和容器类型而异。

### Slot 数据格式 (普通格式)

普通格式用于 `container_set_content` 和 `container_set_slot` 包:

```
varint       item_count    // 物品数量，如果 ≤ 0 表示空槽位
if count > 0:
    varint   item_id       // 物品 ID (注册表中的 ID)
    varint   num_add       // 要添加的组件数量
    for each component:
        varint   component_type
        ...      component_data (取决于类型)
    varint   num_remove    // 要移除的组件数量
    for each component:
        varint   component_type
```

### Hashed Slot 数据格式

哈希格式仅用于 `click_container` 包，组件数据以 CRC32C 校验和形式发送:

```
bool         has_item      // 是否有物品
if has_item:
    varint   item_id
    varint   item_count
    ...      components (哈希格式)
```

### 物品 ID

物品 ID 通过注册表获取，格式：

```
minecraft:item_name
```

例如：
- `minecraft:diamond_sword`
- `minecraft:diamond`
- `minecraft:apple`

## 实体数据

### 实体元数据 (0x5D play_client_entity_data)

实体属性变化时发送：

```
int32          entity_id
array         metadata_entries

entry:
    int8       index
    int8       type_id
    value      (根据类型)
```

#### 元数据类型

| ID | 类型 |
|----|------|
| 0 | byte |
| 1 | int16 |
| 2 | int32 |
| 3 | float32 |
| 4 | string |
| 5 | item |
| 6 | boolean |
| 7 | rotation |
| 8 | position |
| 9 | optional_position |
| 10 | direction |
| 11 | optional_uuid |
| 12 | block_state |
| 13 | compound_tag (NBT) |
| 14 | particle |
| 15 | villager_data |
| 16 | optional_int |
| 17 | pose |
| 18 | cat_variant |
| 19 | frog_variant |
| 20 | optional_block_state |
| 21 | vector3 |
| 22 | quaternion |

### 玩家相关元数据索引

| Index | 类型 | 说明 |
|-------|------|------|
| 0 | byte | 旗帜位 |
| 1 | int32 | 空气值 |
| 2 | string | 自定义名称 |
| 3 | boolean | 是否显示自定义名称 |
| 4 | boolean | 是否静音 |
| 5 | optional_position | 绑定位置 |
| 6 | int16 | 绑定实体 ID |
| 9 | boolean | 离地状态 |
| 10 | compound_tag | 离地位置 |
| 14 | int32 | 鞘翅滑翔时间 |
| 15 | int64 | 落地时间 |

## 玩家列表信息

### 玩家信息更新 (0x42 play_client_player_info_update)

```
int8         action_bitset
array        players

player:
    uuid        [16]byte
    
    if action_bitset & 0x01:  // ADD_PLAYER
        string     name
        array      properties
        
    if action_bitset & 0x02:  // INITIALIZE_CHAT
        optional:
            [256]byte  chat_session_uuid
            [32]byte   public_key_expiry
            [32]byte   public_key
            [256]byte  signature
    
    if action_bitset & 0x04:  // UPDATE_GAME_MODE
        int32      gamemode
    
    if action_bitset & 0x08:  // UPDATE_LISTED
        bool       listed
    
    if action_bitset & 0x10:  // UPDATE_LATENCY
        int32      latency (毫秒)
    
    if action_bitset & 0x20:  // UPDATE_DISPLAY_NAME
        optional:
            nbt        display_name
```

### 玩家信息移除 (0x3D play_client_player_info_remove)

```
array        uuids  [16]byte[]
```

## 实现代码结构

### 玩家状态管理器

```go
package player

type Player struct {
    mu          sync.RWMutex
    
    // 基本信息
    EntityID    int32
    UUID        [16]byte
    Name        string
    GameMode    byte
    
    // 位置
    X, Y, Z     float64
    Yaw, Pitch  float32
    Dimension   string
    
    // 状态
    Health      float32
    MaxHealth   float32
    Food        int32
    Saturation  float32
    Level       int32
    Experience  float32
    
    // 能力
    Invulnerable bool
    Flying       bool
    CanFly       bool
    InstantBreak bool
    
    // 背包
    Inventory   map[int8]*Item
    HeldSlot    int8
    
    // 更新回调
    OnHealthChange  func(health float32)
    OnPositionChange func(x, y, z float64)
    OnInventoryChange func(slot int8, item *Item)
}

type Item struct {
    ID          string
    Count       int32
    Damage      int32
    DisplayName string
    Lore        []string
    Enchantments map[string]int32
    NBT         map[string]any
}
```

### 数据更新接口

```go
func (p *Player) UpdateHealth(health float32, food int32, saturation float32) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    oldHealth := p.Health
    p.Health = health
    p.Food = food
    p.Saturation = saturation
    
    if oldHealth != health && p.OnHealthChange != nil {
        go p.OnHealthChange(health)
    }
}

func (p *Player) UpdatePosition(x, y, z float64, yaw, pitch float32, relative int8) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if relative&0x01 != 0 { x += p.X }
    if relative&0x02 != 0 { y += p.Y }
    if relative&0x04 != 0 { z += p.Z }
    if relative&0x08 != 0 { yaw += p.Yaw }
    if relative&0x10 != 0 { pitch += p.Pitch }
    
    p.X, p.Y, p.Z = x, y, z
    p.Yaw, p.Pitch = yaw, pitch
    
    if p.OnPositionChange != nil {
        go p.OnPositionChange(x, y, z)
    }
}

func (p *Player) UpdateInventorySlot(windowID int8, stateID int32, slot int8, item *Item) {
    if windowID != 0 {
        return  // 只处理玩家背包
    }
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if item == nil || item.Count == 0 {
        delete(p.Inventory, slot)
    } else {
        p.Inventory[slot] = item
    }
    
    if p.OnInventoryChange != nil {
        go p.OnInventoryChange(slot, item)
    }
}
```

## 与客户端集成

在 mcclient 包中添加处理函数：

```go
func (c *Client) handleSetHealthPacket(data []byte) error {
    r := bytes.NewReader(data)
    
    health, _ := readFloat32(r)
    food, _ := readVarInt(r)
    saturation, _ := readFloat32(r)
    
    if c.player != nil {
        c.player.UpdateHealth(health, int32(food), saturation)
    }
    
    return nil
}

func (c *Client) handlePlayerPositionPacket(data []byte) error {
    r := bytes.NewReader(data)
    
    x, _ := readFloat64(r)
    y, _ := readFloat64(r)
    z, _ := readFloat64(r)
    yaw, _ := readFloat32(r)
    pitch, _ := readFloat32(r)
    relative, _ := readU8(r)
    teleportID, _ := readVarInt(r)
    
    if c.player != nil {
        c.player.UpdatePosition(x, y, z, yaw, pitch, int8(relative))
    }
    
    // 回复传送确认
    c.conn.WritePacket(playServerAcceptTeleport, encodeVarInt(teleportID))
    
    return nil
}

func (c *Client) handleContainerSetContentPacket(data []byte) error {
    r := bytes.NewReader(data)
    
    windowID, _ := readVarInt(r)
    stateID, _ := readVarInt(r)
    count, _ := readVarInt(r)
    
    items := make([]*SlotData, count)
    for i := int32(0); i < count; i++ {
        items[i], _ = readSlot(r)
    }
    
    carriedItem, _ := readSlot(r)
    
    if c.player != nil {
        c.player.UpdateInventory(windowID, items, carriedItem)
    }
    
    return nil
}

func (c *Client) handleContainerSetSlotPacket(data []byte) error {
    r := bytes.NewReader(data)
    
    windowID, _ := readVarInt(r)
    stateID, _ := readVarInt(r)
    
    var slot int16
    binary.Read(r, binary.BigEndian, &slot)
    
    item, _ := readSlot(r)
    
    if c.player != nil {
        c.player.UpdateInventorySlot(int8(windowID), stateID, int8(slot), item)
    }
    
    return nil
}

func readSlot(r *bytes.Reader) (*SlotData, error) {
    count, err := readVarIntFromReader(r)
    if err != nil {
        return nil, err
    }
    if count <= 0 {
        return nil, nil
    }
    
    itemID, err := readVarIntFromReader(r)
    if err != nil {
        return nil, nil
    }
    
    // 跳过数据组件
    if err := skipSlotComponents(r); err != nil {
        return nil, err
    }
    
    return &SlotData{ID: itemID, Count: count}, nil
}

func skipSlotComponents(r *bytes.Reader) error {
    numAdd, err := readVarIntFromReader(r)
    if err != nil {
        return err
    }
    for i := int32(0); i < numAdd; i++ {
        if err := skipComponentData(r); err != nil {
            return err
        }
    }
    
    numRemove, err := readVarIntFromReader(r)
    if err != nil {
        return err
    }
    for i := int32(0); i < numRemove; i++ {
        _, err := readVarIntFromReader(r)
        if err != nil {
            return err
        }
    }
    return nil
}
```

## 数据组件解析

1.21+ 使用数据组件代替 NBT：

```go
func readSlot(r *bytes.Reader) (*Item, error) {
    hasItem, _ := readBool(r)
    if !hasItem {
        return nil, nil
    }
    
    item := &Item{}
    itemID, _ := readVarInt(r)
    item.Count, _ = readVarInt(r)
    
    // 读取数据组件
    components, _ := readItemComponents(r)
    item.parseComponents(components)
    
    return item, nil
}

func readItemComponents(r *bytes.Reader) (map[int32]any, error) {
    components := make(map[int32]any)
    
    mask, _ := readU8(r)
    hasComponents := mask&0x01 != 0
    hasRemoved := mask&0x02 != 0
    
    if hasComponents {
        count, _ := readVarInt(r)
        for i := int32(0); i < count; i++ {
            typeID, _ := readVarInt(r)
            value, _ := readComponentValue(r, typeID)
            components[typeID] = value
        }
    }
    
    if hasRemoved {
        count, _ := readVarInt(r)
        for i := int32(0); i < count; i++ {
            _ = readVarInt(r)  // 移除的组件 ID
        }
    }
    
    return components, nil
}
```

## 常用物品组件 ID

| ID | 组件名称 | 说明 |
|----|---------|------|
| 4 | custom_data | 自定义 NBT 数据 |
| 5 | max_stack_size | 最大堆叠数 |
| 6 | max_damage | 最大耐久度 |
| 7 | damage | 当前耐久度 |
| 11 | rarity | 稀有度 |
| 12 | enchantments | 附魔 |
| 13 | stored_enchantments | 附魔书存储的附魔 |
| 21 | chargeable | 可充电物品 |
| 22 | fire_resistant | 防火 |
| 25 | food | 食物属性 |
| 26 | tool | 工具属性 |
| 28 | weapon | 武器属性 |
| 32 | damage_resistant | 抗损伤 |
| 34 | death_protection | 死亡保护 |