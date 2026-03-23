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

### 3.1 包结构重构

新的包结构，将物品相关代码独立管理，工具函数上移：

```
internal/
├── item/                          # 新增 - 物品系统独立包
│   ├── component/                 # 组件解析
│   │   ├── handlers.go            # 处理器注册表（原 component_skipping.go 重构）
│   │   ├── parser.go              # 解析器主逻辑
│   │   ├── discard.go             # 默认丢弃处理器
│   │   └── types.go               # 组件类型定义和常量
│   └── slot.go                    # SlotData 定义（从 packet/readers.go 迁移）
├── mcclient/
│   └── packet/
│       ├── readers.go             # 修改 - 移除 SlotData，使用 internal/item
│       └── utils.go               # 修改 - 移除通用工具函数
└── handlers_container.go          # 修改 - 使用新的组件回调注册

pkg/
├── binutil/                       # 新增 - 二进制读取工具
│   ├── reader.go                  # VarInt, Bool, String 等基础读取
│   ├── writer.go                  # 对应写入函数
│   └── types.go                   # 通用类型定义
└── nbtutil/                       # 新增 - NBT 工具（如适用）
    └── snbt.go                    # SNBT 解析等
```

**迁移说明：**
- `component_skipping.go` 逻辑拆分到 `internal/item/component/` 目录
- `packet/readers.go` 中的 `SlotData` 和组件读取逻辑移到 `internal/item/`
- 通用二进制工具函数（VarInt, String, Bool）移到 `pkg/binutil/`
- `internal/mcclient/packet/` 保留协议特定的读取（如 SlotData 的组合读取）

### 3.2 核心数据结构

**internal/item/component/types.go:**

```go
package component

// ComponentResult 组件解析结果
type ComponentResult struct {
    TypeID int32
    // Data 字段：已解析的组件数据（当前阶段通常为 nil）
    // 后续扩展时根据组件类型存储具体数据
    Data any
}

// ComponentHandler 组件处理函数签名
type ComponentHandler func(r *bytes.Reader) (*ComponentResult, error)
```

**internal/item/component/handlers.go:**

```go
package component

// 全局处理器映射表 - 初始化时所有组件指向 DiscardComponent
var componentHandlers map[int32]ComponentHandler

// ContainerCallback 容器组件特殊回调
var containerCallback func(size int32) error

// RegisterComponentHandler 注册组件处理器
func RegisterComponentHandler(typeID int32, handler ComponentHandler)

// SetContainerCallback 设置容器组件回调
func SetContainerCallback(callback func(size int32) error)

// ParseComponent 解析单个组件
func ParseComponent(typeID int32, r *bytes.Reader) (*ComponentResult, error)
```

**internal/item/slot.go:**

```go
package item

import "gmcc/internal/item/component"

// SlotData 物品槽数据
type SlotData struct {
    ID         int32
    Count      int32
    Components []*component.ComponentResult
}
```

### 3.3 默认处理器实现

**internal/item/component/discard.go:**

```go
package component

import (
    "bytes"
    "gmcc/pkg/binutil"
)

// makeDiscardHandler 创建指定类型的丢弃处理器
func makeDiscardHandler(typeID int32) ComponentHandler {
    return func(r *bytes.Reader) (*ComponentResult, error) {
        skipper, exists := discardFunctions[typeID]
        if !exists {
            // 未知组件：尝试作为 NBT 丢弃
            if err := binutil.SkipNBT(r); err != nil {
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

// discardFunctions 组件类型到丢弃函数的映射
var discardFunctions = map[int32]func(*bytes.Reader) error{
    0:  binutil.SkipNBT,      // custom_data
    1:  binutil.SkipVarInt,   // max_stack_size
    2:  binutil.SkipVarInt,   // max_damage
    // ... 完整列表见实现
    73: discardContainerData, // container
}

// discardContainerData 容器组件丢弃（含回调）
func discardContainerData(r *bytes.Reader) error {
    size, err := binutil.ReadVarInt(r)
    if err != nil {
        return err
    }
    
    if containerCallback != nil {
        if err := containerCallback(size); err != nil {
            return err
        }
    }
    
    length, err := binutil.ReadVarInt(r)
    if err != nil {
        return err
    }
    for i := int32(0); i < length; i++ {
        if err := item.SkipSlotData(r); err != nil {
            return err
        }
    }
    return nil
}
```

**internal/item/component/handlers.go:**

```go
package component

import "bytes"

// ContainerComponentHandler 容器组件特殊处理器
func ContainerComponentHandler(r *bytes.Reader) (*ComponentResult, error) {
    if err := discardContainerData(r); err != nil {
        return nil, err
    }
    return &ComponentResult{TypeID: 73}, nil
}

// containerCallback 全局回调变量
var containerCallback func(size int32) error

// SetContainerCallback 注册容器回调
func SetContainerCallback(callback func(size int32) error) {
    containerCallback = callback
}
```

### 3.4 组件解析主逻辑

**internal/item/component/parser.go:**

```go
package component

import (
    "bytes"
    "fmt"
    
    "gmcc/internal/logx"
    "gmcc/pkg/binutil"
)

// Parser 组件解析器
type Parser struct {
    handlers map[int32]ComponentHandler
}

// NewParser 创建默认解析器
func NewParser() *Parser {
    return &Parser{
        handlers: defaultHandlers(),
    }
}

// Parse 解析组件列表
func (p *Parser) Parse(r *bytes.Reader) ([]*ComponentResult, error) {
    count, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    results := make([]*ComponentResult, 0, count)
    
    for i := int32(0); i < count; i++ {
        typeID, err := binutil.ReadVarInt(r)
        if err != nil {
            return nil, fmt.Errorf("读取组件类型ID失败: %w", err)
        }
        
        result, err := p.parseComponent(typeID, r)
        if err != nil {
            return nil, fmt.Errorf("处理组件 %d 失败: %w", typeID, err)
        }
        
        results = append(results, result)
    }
    
    return results, nil
}

// parseComponent 解析单个组件
func (p *Parser) parseComponent(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
    handler, ok := p.handlers[typeID]
    if !ok {
        logx.Warnf("未知组件类型: %d, 尝试作为 NBT 丢弃", typeID)
        if err := binutil.SkipNBT(r); err != nil {
            return nil, err
        }
        return &ComponentResult{TypeID: typeID}, nil
    }
    
    return handler(r)
}

// RegisterHandler 注册自定义处理器
func (p *Parser) RegisterHandler(typeID int32, handler ComponentHandler) {
    p.handlers[typeID] = handler
}
```

### 3.5 初始化与注册

**internal/item/component/handlers.go:**

```go
package component

var componentHandlers map[int32]ComponentHandler

// defaultHandlers 返回默认处理器映射
func defaultHandlers() map[int32]ComponentHandler {
    handlers := make(map[int32]ComponentHandler)
    
    // ID 范围 0-103（从原 component_skipping.go 提取）
    componentIDs := []int32{
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
        20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
        40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
        60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
        80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
        100, 101, 102, 103,
    }
    
    for _, typeID := range componentIDs {
        if typeID == 73 {
            handlers[typeID] = ContainerComponentHandler
        } else {
            handlers[typeID] = makeDiscardHandler(typeID)
        }
    }
    
    return handlers
}
```

### 3.6 容器组件回调注册

在 `internal/mcclient/handlers_container.go` 中注册容器回调：

```go
package mcclient

import (
    "gmcc/internal/item/component"
    "gmcc/internal/logx"
)

func init() {
    component.SetContainerCallback(func(size int32) error {
        // 预留：后续实现容器内容处理
        logx.Debugf("container component parsed: size=%d", size)
        return nil
    })
}
```

---

## 4. 与现有系统集成

### 4.1 SlotData 读取修改

**internal/item/slot.go:**

```go
package item

import (
    "bytes"
    
    "gmcc/internal/item/component"
    "gmcc/pkg/binutil"
)

// SlotData 物品槽数据
type SlotData struct {
    ID         int32
    Count      int32
    Components []*component.ComponentResult
}

// ReadSlotData 从 Reader 读取物品槽数据
func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
    // 读取基础字段
    count, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    if count <= 0 {
        return nil, nil // 空槽
    }
    
    itemID, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    // 创建解析器并解析组件
    parser := component.NewParser()
    components, err := parser.Parse(r)
    if err != nil {
        return nil, err
    }
    
    return &SlotData{
        ID:         itemID,
        Count:      count,
        Components: components,
    }, nil
}

// SkipSlotData 跳过物品槽数据（不存储）
func SkipSlotData(r *bytes.Reader) error {
    slot, err := ReadSlotData(r)
    if err != nil {
        return err
    }
    // slot 被读取后丢弃
    _ = slot
    return nil
}
```

**internal/mcclient/packet/readers.go:**

```go
package packet

import (
    "gmcc/internal/item"
)

// ReadSlotData 使用 internal/item 的实现
func ReadSlotData(r *bytes.Reader) (*item.SlotData, error) {
    return item.ReadSlotData(r)
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
// internal/item/component/custom.go

package component

import (
    "bytes"
    "gmcc/pkg/binutil"
)

// ParseCustomNameComponent 解析自定义名称组件 (ID: 6)
func ParseCustomNameComponent(r *bytes.Reader) (*ComponentResult, error) {
    // 解析 NBT 格式的文本组件
    nbt, err := binutil.ReadNBT(r)
    if err != nil {
        return nil, err
    }
    
    return &ComponentResult{
        TypeID: 6,
        Data:   nbt, // 或提取为具体结构
    }, nil
}

// 在包初始化时注册
func init() {
    // 注意：这会覆盖默认的丢弃处理器
    RegisterComponentHandler(6, ParseCustomNameComponent)
}
```

### 5.2 容器组件完整实现

```go
// internal/mcclient/handlers_container.go

import "gmcc/internal/item/component"

func init() {
    component.SetContainerCallback(func(size int32) error {
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

**架构改进：**
- 物品系统独立到 `internal/item/` 包，职责更清晰
- 通用二进制工具移至 `pkg/binutil/`，可被其他模块复用
- 组件解析采用注册表模式，支持运行时扩展

**核心机制：**
- 所有组件默认使用丢弃处理器（读取并丢弃）
- 通过 `Parser.RegisterHandler()` 可注册具体解析器
- 容器组件（ID 73）预留特殊回调接口

**后续扩展：**
- 实现具体组件解析器时，在 `internal/item/component/` 添加新文件
- 使用 `binutil` 中的基础读取函数处理二进制数据
- 通过 `component.SetContainerCallback()` 实现容器内容处理
