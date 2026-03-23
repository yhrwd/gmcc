# 物品组件解析系统重构设计

**日期**: 2026-03-24  
**作者**: gmcc agent  
**状态**: 设计中  
**协议版本**: 774 (1.21.11)

---

## 1. 概述

### 1.1 背景

当前系统使用 `component_skipping.go` 来跳过物品槽中的组件数据。这种方式虽然能避免解析错误，但丢弃了所有组件信息，无法支持后续功能（如显示物品自定义名称、附魔、容器内容等）。

### 1.2 目标

将组件处理从"跳过"改为"可扩展的解析"：
- **默认行为**: 所有组件类型使用 `DiscardComponent` 处理器（读取并丢弃数据）
- **可扩展**: 通过注册表模式，后续可为特定组件类型添加实际解析逻辑
- **容器组件**: 预留特殊处理机制

### 1.3 非目标

- 不一次性实现所有组件类型的解析
- 不保留未解析组件的原始字节
- 不修改容器处理逻辑的主体架构

---

## 2. 组件类型映射表

参考官方 wiki (Data Components)，1.21.11 版本共有 120+ 种组件类型：

| ID | 名称 | 处理器状态 |
|----|------|-----------|
| 0 | custom_data | DiscardComponent |
| 1 | max_stack_size | DiscardComponent |
| 2 | max_damage | DiscardComponent |
| 3 | damage | DiscardComponent |
| 4 | unbreakable | DiscardComponent |
| 5 | use_effects | DiscardComponent |
| 6 | custom_name | DiscardComponent |
| ... | ... | ... |
| 73 | container | DiscardComponent + 预留回调 |
| ... | ... | ... |

完整列表见 `component_handlers.go` 中的 `componentHandlers` 映射表。

---

## 3. 设计方案

### 3.1 文件结构变更

```
internal/mcclient/packet/
├── component_skipping.go      # 删除
├── component_handlers.go      # 新增 - 处理器注册表和默认实现
├── component_parser.go        # 新增 - 解析器主逻辑
└── readers.go                 # 修改 - SlotData 读取使用新解析器
```

### 3.2 核心数据结构

```go
// component_handlers.go

// ComponentResult 组件解析结果
type ComponentResult struct {
    TypeID int32
    // Data 字段：已解析的组件数据（当前阶段通常为 nil）
    // 后续扩展时根据组件类型存储具体数据
    Data any
}

// ComponentHandler 组件处理函数签名
type ComponentHandler func(r *bytes.Reader) (*ComponentResult, error)

// 全局处理器映射表 - 初始化时所有组件指向 DiscardComponent
var componentHandlers map[int32]ComponentHandler

// ContainerCallback 容器组件特殊回调
var containerCallback func(size int32) error
```

### 3.3 默认处理器实现

```go
// component_handlers.go

// makeDiscardHandler 创建指定类型的丢弃处理器
// 使用闭包捕获 typeID，避免参数传递问题
func makeDiscardHandler(typeID int32) ComponentHandler {
    return func(r *bytes.Reader) (*ComponentResult, error) {
        // 根据组件类型获取对应的丢弃函数
        skipper, exists := componentDiscards[typeID]
        if !exists {
            // 未知组件：尝试作为 NBT 丢弃
            if err := SkipNBT(r); err != nil {
                return nil, err
            }
            return &ComponentResult{TypeID: typeID}, nil
        }
        
        if err := skipper(r); err != nil {
            return nil, err
        }
        return &ComponentResult{TypeID: typeID}, nil
    }
}

// componentDiscards 组件类型到丢弃函数的映射
// 从现有 component_skipping.go 提取，包含 0-103 所有已知组件
var componentDiscards = map[int32]func(*bytes.Reader) error{
    0:  SkipNBT,           // custom_data
    1:  SkipVarInt,        // max_stack_size
    2:  SkipVarInt,        // max_damage
    // ... 完整映射表，见实现时提取
    73: SkipContainerData, // container
}

// SkipContainerData 容器组件丢弃函数
func SkipContainerData(r *bytes.Reader) error {
    // 读取容器大小（预留回调参数）
    size, err := ReadVarIntFromReader(r)
    if err != nil {
        return err
    }
    
    // 触发预留回调
    if containerCallback != nil {
        if err := containerCallback(size); err != nil {
            return err
        }
    }
    
    // 读取并丢弃槽位数据
    length, err := ReadVarIntFromReader(r)
    if err != nil {
        return err
    }
    for i := int32(0); i < length; i++ {
        if err := SkipSlotData(r); err != nil {
            return err
        }
    }
    
    return nil
}

// 容器组件处理器（包装丢弃函数，添加结果返回）
func ContainerComponentHandler(r *bytes.Reader) (*ComponentResult, error) {
    if err := SkipContainerData(r); err != nil {
        return nil, err
    }
    return &ComponentResult{TypeID: 73}, nil
}

// 注册/获取接口
func RegisterComponentHandler(typeID int32, handler ComponentHandler)
func SetContainerCallback(callback func(size int32) error)
```

### 3.4 组件解析主逻辑

```go
// component_parser.go

// ParseComponents 解析物品槽中的组件列表
func ParseComponents(r *bytes.Reader) ([]*ComponentResult, error) {
    // 读取组件数量（VarInt）
    count, err := ReadVarIntFromReader(r)
    if err != nil {
        return nil, err
    }
    
    results := make([]*ComponentResult, 0, count)
    
    for i := int32(0); i < count; i++ {
        // 读取组件类型ID
        typeID, err := ReadVarIntFromReader(r)
        if err != nil {
            return nil, fmt.Errorf("读取组件类型ID失败: %w", err)
        }
        
        // 查找并执行处理器
        handler, ok := componentHandlers[typeID]
        if !ok {
            // 未知组件：尝试作为 NBT 丢弃
            logx.Warnf("未知组件类型: %d, 尝试丢弃", typeID)
            if err := SkipNBT(r); err != nil {
                return nil, err
            }
            continue
        }
        
        result, err := handler(r)
        if err != nil {
            return nil, fmt.Errorf("处理组件 %d 失败: %w", typeID, err)
        }
        
        results = append(results, result)
    }
    
    return results, nil
}
```

### 3.5 初始化与注册

```go
// component_handlers.go

var componentHandlers map[int32]ComponentHandler

func init() {
    componentHandlers = make(map[int32]ComponentHandler)
    
    // 从现有 component_skipping.go 的组件列表提取
    // 实际 ID 范围为 0-103（见原 component_skipping.go）
    componentIDs := []int32{
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
        20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
        40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
        60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
        80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
        100, 101, 102, 103,
    }
    
    // 为所有已知组件类型注册默认丢弃处理器
    for _, typeID := range componentIDs {
        // 容器组件(73)使用特殊处理器
        if typeID == 73 {
            componentHandlers[typeID] = ContainerComponentHandler
        } else {
            componentHandlers[typeID] = makeDiscardHandler(typeID)
        }
    }
}
```

### 3.6 容器组件回调注册

在 `internal/mcclient/handlers_container.go` 中注册容器回调：

```go
package mcclient

import (
    "gmcc/internal/logx"
    "gmcc/internal/mcclient/packet"
)

func init() {
    packet.SetContainerCallback(func(size int32) error {
        // 预留：后续实现容器内容处理
        logx.Debugf("container component parsed: size=%d", size)
        return nil
    })
}
```

---

## 4. 与现有系统集成

### 4.1 SlotData 读取修改

```go
// readers.go 中的 ReadSlotData 函数

// SlotData 扩展组件字段
type SlotData struct {
    ID         int32
    Count      int32
    // 新增：组件数据（当前阶段仅存储解析结果，数据字段为 nil）
    Components []*ComponentResult
}

func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
    // 原有逻辑读取 ItemID 等基础字段...
    slot := &SlotData{ID: itemID, Count: count}
    
    // 修改：使用新的组件解析
    components, err := ParseComponents(r)
    if err != nil {
        return nil, err
    }
    slot.Components = components
    
    // 当前：组件数据已存储在 SlotData.Components 中
    // 后续：可以遍历 components 提取具体数据（如 custom_name, damage 等）
    
    return slot, nil
}
```

### 4.2 向后兼容

- 默认行为与之前相同：所有组件数据被读取并丢弃
- 添加组件解析不会影响现有功能
- 仅当日志级别为 Debug 时输出组件解析信息

---

## 5. 后续扩展指南

### 5.1 添加新组件解析

```go
// 在合适的位置（如 item/components.go）实现具体解析

func ParseCustomNameComponent(r *bytes.Reader) (*packet.ComponentResult, error) {
    // 解析 NBT 格式的文本组件
    nbt, err := ReadNBT(r)
    if err != nil {
        return nil, err
    }
    
    return &packet.ComponentResult{
        TypeID: 6,
        Data:   nbt, // 或提取为具体类型
    }, nil
}

// 在初始化时注册
func init() {
    packet.RegisterComponentHandler(6, ParseCustomNameComponent)
}
```

### 5.2 容器组件完整实现

```go
// handlers_container.go

func init() {
    packet.SetContainerCallback(func(size int32) error {
        // 存储容器大小到 Player 状态
        // 后续实现容器槽位同步
        return nil
    })
}
```

---

## 6. 测试计划

1. **单元测试**: 测试各种组件类型的读取
2. **集成测试**: 验证实际服务器返回的物品数据解析
3. **回归测试**: 确保与现有容器处理逻辑兼容

---

## 7. 风险评估

| 风险 | 缓解措施 |
|------|---------|
| 新解析逻辑导致协议错误 | 保持与原有 Skip 逻辑一致，仅添加结果收集 |
| 性能开销 | 结果切片预分配，避免重复内存分配 |
| 未知组件类型 | 保留原有未知组件处理：尝试作为 NBT 丢弃 |

---

## 8. 总结

本设计将组件处理从"跳过"重构为"可扩展的解析框架"：
- 所有组件默认使用 `DiscardComponent` 读取并丢弃
- 通过全局映射表支持后续注册具体解析器
- 容器组件预留特殊回调接口
- 代码结构清晰，便于后续逐步添加组件解析功能
