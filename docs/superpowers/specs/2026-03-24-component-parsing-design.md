# 物品组件解析系统重构设计

**日期**: 2026-03-24  
**作者**: gmcc agent  
**状态**: 设计中  
**协议版本**: 774 (1.21.11)

**参考实现声明**: 本设计参考了 Tnze/go-mc 项目的组件系统设计，遵循 MIT License (Copyright (c) 2019 Tnze)。本设计仅借鉴其接口设计理念，具体实现方式采用处理器函数映射表模式。

---

## 1. 概述

### 1.1 背景

当前系统使用 `component_skipping.go` 来跳过物品槽中的组件数据。这种方式虽然能避免解析错误，但丢弃了所有组件信息，无法支持后续功能（如显示物品自定义名称、附魔、容器内容等）。

同时，代码结构存在以下问题：
- 工具类分散在多个包中，部分可复用工具与Minecraft客户端紧耦合
- 聊天处理流程依赖内部日志和国际化，难以复用和测试
- 缺少统一的数据获取策略

### 1.2 目标

将组件处理从"跳过"改为"可扩展的解析"：
- **默认行为**: 所有组件类型使用 `DiscardComponent` 处理器（读取并丢弃数据）
- **可扩展**: 通过注册表模式，后续可为特定组件类型添加实际解析逻辑
- **容器组件**: 预留特殊处理机制

**架构优化目标：**
- **工具类提取**: 将通用工具移至 `pkg/` 供其他模块复用
- **聊天解耦**: 移除聊天处理对内部包（i18n/logx）的硬依赖
- **数据标准化**: 建立明确的数据获取优先级策略

### 1.3 非目标

- 不一次性实现所有组件类型的解析
- 不保留未解析组件的原始字节
- 不修改容器处理逻辑的主体架构
- 不重构协议常量定义（保持在 `internal/mcclient/protocol/`）

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

---

## 附录A: 数据获取优先级策略

### A.1 原则

**优先级（从高到低）：**

1. **Minecraft Wiki (zh.minecraft.wiki)** - 协议规范的主要来源
   - URL: https://zh.minecraft.wiki/
   - 所有页面索引: https://zh.minecraft.wiki/w/Special:AllPages
   - **注意**: 确认页面版本是否为 **1.21.11**（或兼容版本）
   - 搜索关键词: "Data Components", "Protocol", "Item Format"

2. **go-mc 参考实现** - Tnze/go-mc 项目
   - 位置: `../go-mc/level/component/`
   - 组件接口设计、序列化模式
   - 注意：版本可能略有差异，需与 wiki 核对

3. **官方 Minecraft 源代码/文档** - 权威参考
   - 通过反编译或官方发布的 obfuscation maps
   - 用于验证协议细节和数据类型

4. **第三方数据资源** - 辅助参考
   - PrismarineJS 库的实现
   - 社区维护的 protocol.json 文件
   - 仅作为补充验证，不作为主要依据

5. **实际网络抓包** - 最终验证
   - 使用实际服务器通信数据进行验证
   - 用于确认实现正确性
   - 不作为设计阶段的主要参考

**重要提示**: 原 wiki.vg 已闭站，现使用 https://zh.minecraft.wiki/ 作为主要参考来源。

### A.2 组件类型定义来源

**Data Components 定义：**
- **主要来源**: Minecraft Wiki (zh.minecraft.wiki - Data Components 页面)
- **参考实现**: go-mc (Tnze/go-mc) 的 `level/component/` 包
  - **License**: MIT License (Copyright (c) 2019 Tnze)
  - 提供组件接口设计、ID 映射、序列化方法
  - 组件数量：~40 个已实现（1.20+ 版本）
  - **ID 映射注意**: go-mc 的 ID 顺序可能与 1.21.11 有差异，需对照 wiki
  - **使用方式**: 本设计仅借鉴其接口设计理念，具体实现采用处理器函数映射表模式
- **验证方式**: 对比 `.knowledge/1.21.11/types/components.json`

**组件ID映射表：**
- 使用 wiki 定义的 raw ID 映射
- 版本: 1.21.11 (协议 774)
- ID 范围: 0-103
- **go-mc 参考**: 查看 `go-mc/level/component/components.go` 中的 NewComponent() switch 语句

### A.3 参考实现对比

**go-mc 组件结构 (参考):**

> 参考实现: [Tnze/go-mc](https://github.com/Tnze/go-mc)
> License: MIT License (Copyright (c) 2019 Tnze)
> 本设计仅借鉴接口设计理念，实现方式不同

```go
// go-mc/level/component/components.go
type DataComponent interface {
    pk.Field           // 实现 ReadFrom/WriteTo
    ID() string        // 返回组件名称如 "minecraft:custom_name"
}

// 示例: CustomName 组件
type CustomName struct {
    Name chat.Message
}

func (CustomName) ID() string { return "minecraft:custom_name" }
func (c *CustomName) ReadFrom(r io.Reader) (n int64, err error) {
    return c.Name.ReadFrom(r)
}
```

**本设计差异:**
- go-mc: 每个组件是一个完整类型，实现序列化接口
- 本设计: 使用处理器函数映射表，更轻量，适合"先丢弃后实现"策略
- 两者都支持 ID 到处理器的映射

### A.4 更新策略

当 Minecraft 版本更新时：
1. 首先检查 zh.minecraft.wiki 的 Data Components 页面更新
2. 参考 go-mc 的组件实现作为结构参考
3. 对比新旧版本组件 ID 映射变化
4. 更新 `discardFunctions` 映射表
5. 在测试服务器验证解析正确性

---

## 附录B: 工具类重构建议

### B.1 可移动的工具类清单

#### 1. CFB8 加密 → pkg/crypto/cfb8/

**当前位置**: `internal/mcclient/crypto/cfb8.go`

**移动理由**:
- 纯加密算法实现，无 Minecraft 特定逻辑
- 可独立复用，符合 `pkg/` 包的通用性要求
- 仅依赖标准库 `crypto/cipher`

**新位置**:
```
pkg/crypto/cfb8/
├── cfb8.go          # 原有实现
└── cfb8_test.go     # 对应测试
```

**依赖修改**:
```go
// internal/mcclient/packet/codec.go
import "gmcc/pkg/crypto/cfb8"  // 修改前: internal/mcclient/crypto
```

#### 2. VarInt 工具 → pkg/mcutil/varint.go

**当前位置**: `internal/mcclient/packet/codec.go` (内嵌)

**移动理由**:
- VarInt 是 Minecraft 协议的基础类型
- 多个模块可能需要（packet, protocol, item）
- 独立后便于测试和复用

**新位置**:
```
pkg/mcutil/
├── varint.go        # ReadVarInt, WriteVarInt, VarIntSize
└── varint_test.go
```

#### 3. UUID 工具 → pkg/mcutil/uuid.go

**当前位置**: `internal/mcclient/packet/utils.go`

**提取函数**:
- `ParseUUID()` - 从字符串解析 UUID
- `FormatUUID()` - 格式化 UUID 为字符串
- `OfflineUUID()` - 生成离线玩家 UUID
- `UUIDToBytes()` - UUID 转字节数组

**新位置**:
```
pkg/mcutil/uuid.go
```

#### 4. NBT 处理 → pkg/nbt/

**当前位置**: `internal/nbt/`

**移动理由**:
- NBT 是 Minecraft 通用数据格式
- 完整的编码/解码/SNBT 解析实现
- 可被其他项目复用

**新位置**:
```
pkg/nbt/
├── decode.go
├── encode.go
├── nbt.go
├── path.go
├── raw.go
├── snbt.go
└── *_test.go
```

**注意**: 移动时需移除对 `internal/logx` 的依赖，改用 error 返回。

#### 5. CESU8 工具 → pkg/nbt/cesu8.go

**当前位置**: `internal/nbt/decode.go:739-782`

**提取函数**:
- `CESU8ToUTF8()` - 当前存在，改为公开
- `UTF8ToCESU8()` - 新增（目前缺失）

**用途**: Minecraft 字符串编码（协议中使用）

### B.2 聊天处理解耦

#### 当前问题

**文件**: `internal/mcclient/chat/chat_parser.go:8`
```go
import "gmcc/internal/logx"
```

**紧耦合点**:
1. `ExtractPlainTextFromChatJSON()` 调用 `logx.Debugf()`
2. `text_component.go` 直接调用 `i18n.Translate()`
3. TUI 直接依赖 `mcclient` 和 `chat` 包

#### 解耦方案

**方案1: 使用接口注入**

```go
// internal/mcclient/chat/interfaces.go
package chat

// Logger 可选调试日志接口
type Logger interface {
    Debugf(format string, args ...any)
}

// Translator 本地化翻译接口
type Translator interface {
    Translate(key string, args ...any) string
}

// ParserOptions 解析器选项
type ParserOptions struct {
    Logger     Logger
    Translator Translator
}
```

**方案2: 返回值替代副作用**

```go
// 当前（副作用）:
func ExtractPlainText(json []byte) string {
    // 内部调用 logx.Debugf
}

// 优化（无副作用）:
func ExtractPlainText(json []byte) (string, error) {
    // 返回错误，调用者决定是否记录
}
```

#### 实施步骤

1. 在 `chat` 包定义 `ParserOptions` 结构体
2. 修改 `TextComponent.ToPlain()` 接收可选 `Translator` 参数
3. 移除 `chat_parser.go` 的 `logx` 导入
4. TUI 层注入 translator 和 logger
5. 保持向后兼容（nil 参数时使用默认行为）

### B.3 重构后包结构

```
gmcc/
├── cmd/gmcc/
├── internal/
│   ├── mcclient/
│   │   ├── chat/           # 解耦后的聊天处理
│   │   └── packet/         # 使用pkg工具
│   └── ...其他/
└── pkg/                    # 公共可复用工具
    ├── binutil/            # 二进制读写（本设计新增）
    ├── crypto/
    │   └── cfb8/           # CFB8加密
    ├── mcutil/             # Minecraft通用工具
    │   ├── uuid.go
    │   └── varint.go
    └── nbt/                # NBT处理（从internal移出）
```

### B.4 迁移优先级

| 优先级 | 包/功能 | 理由 |
|--------|---------|------|
| P0 | binutil | 本设计依赖，必须先实现 |
| P1 | crypto/cfb8 | 独立性强，无依赖风险 |
| P2 | mcutil/varint | packet 重构需要 |
| P2 | mcutil/uuid | 多个地方使用 |
| P3 | nbt/ | 依赖较多，需要分阶段迁移 |
| P3 | chat解耦 | 影响 TUI，可后续优化 |

---

## 附录C: 常量提取与优化建议

### C.1 容器类型常量提取

**当前位置**: `internal/mcclient/handlers_container.go:22-32`

```go
const (
    ContainerTypePlayer    int32 = 0
    ContainerTypeChest     int32 = 1
    ContainerTypeCrafting  int32 = 2
    ContainerTypeFurnace   int32 = 3
    ContainerTypeDispenser int32 = 4
    ContainerTypeHopper    int32 = 5
    ContainerTypeAnvil     int32 = 6
    ContainerTypeBeacon    int32 = 7
    ContainerTypeBrewing   int32 = 8
)
```

**建议移动位置**: `internal/mcclient/protocol/container_types.go`

**理由**:
- 容器类型是协议级别的常量，与具体处理器实现无关
- 多处可能使用（packet handlers, player inventory, UI 显示）
- 应该与 v774.go 中的协议常量保持一致的命名风格

### C.2 组件类型 ID 常量提取

**当前问题**: 组件 ID (0-103) 分布在 `component_skipping.go` 的映射表中，但没有明确定义为常量

**建议**: 在 `internal/item/component/constants.go` 中定义：

```go
package component

// Component IDs per Minecraft 1.21.11 Data Components
const (
    CustomData      int32 = 0
    MaxStackSize    int32 = 1
    MaxDamage       int32 = 2
    Damage          int32 = 3
    Unbreakable     int32 = 4
    UseEffects      int32 = 5
    CustomName      int32 = 6
    MinimumAttackCharge int32 = 7
    DamageType      int32 = 8
    // ... 完整的 0-103 列表
    Container       int32 = 73
    BlockState      int32 = 74
    Bees            int32 = 75
    Lock            int32 = 76
    ContainerLoot   int32 = 77
    BreakSound      int32 = 78
    // ... 继续到 103
)
```

**使用方式**:
```go
// 取代魔法数字
componentHandlers[Container] = ContainerComponentHandler

// 在 switch 中使用
case CustomName:
    // 处理自定义名称
```

### C.3 常量命名规范

统一常量命名风格：

| 位置 | 当前命名 | 建议命名 |
|------|---------|---------|
| handlers_container.go | `ContainerTypePlayer` | `ContainerPlayer` 或保持 |
| v774.go | `PlayClientContainerContent` | 保持（协议包 ID） |
| constants.go | `MaxPacketSize` | 保持 |
| component 新增 | - | `ComponentCustomName` |

**推荐风格**:
- 协议常量: `PlayClientXxx`, `PlayServerXxx`
- 组件常量: `ComponentXxx` (带前缀避免冲突)
- 容器常量: `ContainerXxx`
- 通用常量: `MaxXxx`, `DefaultXxx`

### C.4 NBT 解析优化建议

**当前实现分析**:
- 基于反射的解码器 (`internal/nbt/decode.go`)
- 每次解码创建新的 Decoder
- 反射带来性能开销

**数据来源**: 参考 Minecraft NBT 格式规范 (wiki.vg)

**优化方向**:

#### 1. 添加 RawReader 类型

用于直接读取原始 NBT 数据而不使用反射：

```go
// internal/nbt/raw.go

type RawReader struct {
    r   io.Reader
    buf [8]byte // 复用缓冲区
}

func NewRawReader(r io.Reader) *RawReader {
    return &RawReader{r: r}
}

// 直接读取，不通过反射
func (r *RawReader) ReadByte() (int8, error)
func (r *RawReader) ReadShort() (int16, error)
func (r *RawReader) ReadInt() (int32, error)
func (r *RawReader) ReadString() (string, error)
func (r *RawReader) ReadCompound() (map[string]any, error)
// ...
```

#### 2. Decoder 对象池

```go
var decoderPool = sync.Pool{
    New: func() any {
        return &Decoder{}
    },
}

func AcquireDecoder(r io.Reader) *Decoder {
    d := decoderPool.Get().(*Decoder)
    d.r = r
    d.offset = 0
    d.fieldPath = d.fieldPath[:0]
    return d
}

func ReleaseDecoder(d *Decoder) {
    decoderPool.Put(d)
}
```

#### 3. 特定类型的快速路径

```go
// 对于 map[string]any 使用快速路径
func (d *Decoder) decodeCompoundFast() (map[string]any, error) {
    result := make(map[string]any)
    // 直接读取，避免反射开销
    // ...
    return result, nil
}
```

#### 4. 组件解析中的 NBT 优化

在组件解析时，大多数组件只需要特定字段：

```go
// 优化的 custom_name 组件解析
func ParseCustomNameComponent(r *bytes.Reader) (*ComponentResult, error) {
    // 使用 RawReader 直接读取 NBT
    raw := nbt.NewRawReader(r)
    
    // 只读取我们关心的字段，忽略其他
    compound, err := raw.ReadCompound()
    if err != nil {
        return nil, err
    }
    
    name, ok := compound["text"].(string)
    if !ok {
        return nil, fmt.Errorf("custom_name missing text field")
    }
    
    return &ComponentResult{
        TypeID: CustomName,
        Data:   name,
    }, nil
}
```

#### 5. CESU8 字符串编码优化

**当前位置**: `internal/nbt/decode.go:739-782` (readUTF8String)

**优化建议**: 提取为公开函数并优化

```go
// pkg/nbt/cesu8.go

// CESU8ToUTF8 converts CESU-8 (Minecraft's string encoding) to UTF-8
func CESU8ToUTF8(data []byte) (string, error) {
    // 优化版本：预分配容量，减少扩容
    var result []rune
    // ... 实现
}

// UTF8ToCESU8 converts UTF-8 to CESU-8
func UTF8ToCESU8(s string) ([]byte, error) {
    // 用于协议写入
    // ... 实现
}
```

**使用场景**:
- 协议包中的字符串字段（Minecraft 使用 CESU-8 变体）
- NBT 字符串标签
- 聊天消息文本组件
