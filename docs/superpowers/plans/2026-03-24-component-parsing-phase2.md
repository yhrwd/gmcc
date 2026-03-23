# 阶段2: 组件解析实现

**计划文档**: [主计划](../2026-03-24-component-parsing-implementation-plan.md)  
**前置阶段**: [阶段1 - 基础架构搭建](./2026-03-24-component-parsing-phase1.md)  
**预计工期**: 3天  
**依赖**: 阶段1完成  
**阻塞**: 阶段3

---

## 1. 阶段目标

实现组件解析框架和丢弃处理器，支持104个组件类型的解析（当前阶段仅丢弃数据），集成容器组件回调机制。

### 成功标准

- [ ] `ComponentParser` 实现完成
- [ ] 104个组件类型映射表完成
- [ ] 容器组件（ID 73）回调机制工作正常
- [ ] SlotData 包含组件列表信息
- [ ] 与现有容器处理逻辑兼容

---

## 2. 任务清单

### 2.1 实现 internal/item/component/parser.go

**目标**: 实现组件解析器主逻辑

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
    // 实现组件列表解析
}

// RegisterHandler 注册自定义处理器
func (p *Parser) RegisterHandler(typeID int32, handler ComponentHandler)
```

**实现步骤**:
1. [ ] 创建 `Parser` 结构体
2. [ ] 实现 `NewParser()` 构造函数
3. [ ] 实现 `Parse()` 方法，解析组件列表
4. [ ] 实现 `RegisterHandler()` 方法
5. [ ] 添加单元测试

**测试要点**:
- 空组件列表解析
- 多个组件连续解析
- 未知组件类型处理
- 错误恢复

---

### 2.2 创建组件处理器映射表

**文件**: `internal/item/component/handlers.go`

**目标**: 定义104个组件类型的处理器映射

```go
package component

// componentHandlers 全局处理器映射表
var componentHandlers map[int32]ComponentHandler

// defaultHandlers 返回默认处理器映射
func defaultHandlers() map[int32]ComponentHandler {
    handlers := make(map[int32]ComponentHandler)
    
    // ID 范围 0-103
    for typeID := int32(0); typeID <= 103; typeID++ {
        if typeID == 73 {
            // 容器组件特殊处理
            handlers[typeID] = ContainerComponentHandler
        } else {
            handlers[typeID] = makeDiscardHandler(typeID)
        }
    }
    
    return handlers
}
```

**组件ID列表** (从 component_skipping.go 迁移):

| ID | 名称 | 处理器类型 |
|----|------|-----------|
| 0 | custom_data | SkipNBT |
| 1 | max_stack_size | SkipVarInt |
| 2 | max_damage | SkipVarInt |
| 3 | damage | SkipVarInt |
| 4 | unbreakable | SkipNothing |
| 5 | use_effects | SkipUseEffects |
| 6 | custom_name | SkipNBT |
| ... | ... | ... |
| 73 | container | ContainerComponentHandler |
| ... | ... | ... |
| 103 | color | SkipVarInt |

**完整ID映射**: 参见 `.knowledge/1.21.11/types/components.json`

**实现步骤**:
1. [ ] 创建 handlers.go 文件
2. [ ] 实现 `defaultHandlers()` 函数
3. [ ] 注册所有104个组件的丢弃处理器
4. [ ] 特殊处理容器组件（ID 73）

---

### 2.3 实现丢弃处理器

**文件**: `internal/item/component/discard.go`

**目标**: 从 `component_skipping.go` 迁移丢弃逻辑

**迁移清单**:

| 原函数 | 新位置 | 说明 |
|--------|--------|------|
| SkipNBT | binutil.SkipNBT | 移至pkg |
| SkipVarInt | binutil.SkipVarInt | 移至pkg |
| SkipString | binutil.SkipString | 移至pkg |
| SkipBool | binutil.SkipBool | 移至pkg |
| SkipNothing | discard.go | 保留 |
| SkipUseEffects | discard.go | 保留 |
| SkipEnchantments | discard.go | 保留 |
| ... | ... | ... |
| SkipContainer | discard.go | 特殊处理 |

**实现步骤**:

1. **基础跳过函数迁移到 pkg/binutil/**
   ```go
   // pkg/binutil/skip.go
   func SkipNBT(r *bytes.Reader) error
   func SkipVarInt(r *bytes.Reader) error
   func SkipString(r *bytes.Reader) error
   func SkipBool(r *bytes.Reader) error
   ```

2. **组件特定跳过函数保留在 discard.go**
   ```go
   // internal/item/component/discard.go
   func SkipNothing(r *bytes.Reader) error
   func SkipUseEffects(r *bytes.Reader) error
   func SkipEnchantments(r *bytes.Reader) error
   // ... 其他复杂跳过函数
   ```

3. **创建 makeDiscardHandler 工厂函数**
   ```go
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
   ```

---

### 2.4 实现容器组件特殊处理

**目标**: 保留容器组件回调机制，支持后续容器内容处理

**文件**: `internal/item/component/discard.go`

```go
// containerCallback 全局回调变量
var containerCallback func(size int32) error

// SetContainerCallback 注册容器回调
func SetContainerCallback(callback func(size int32) error) {
    containerCallback = callback
}

// ContainerComponentHandler 容器组件特殊处理器
func ContainerComponentHandler(r *bytes.Reader) (*ComponentResult, error) {
    // 1. 读取容器大小
    size, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    // 2. 触发回调（如已注册）
    if containerCallback != nil {
        if err := containerCallback(size); err != nil {
            return nil, err
        }
    }
    
    // 3. 继续读取容器内容
    length, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    for i := int32(0); i < length; i++ {
        if err := item.SkipSlotData(r); err != nil {
            return nil, err
        }
    }
    
    return &ComponentResult{TypeID: 73}, nil
}
```

**回调注册** (在 handlers_container.go 中):

```go
// internal/mcclient/handlers_container.go
func init() {
    component.SetContainerCallback(func(size int32) error {
        logx.Debugf("container component parsed: size=%d", size)
        // 预留：后续实现容器内容处理
        return nil
    })
}
```

**实现步骤**:
1. [ ] 在 discard.go 中定义回调变量
2. [ ] 实现 `SetContainerCallback()` 函数
3. [ ] 实现 `ContainerComponentHandler()`
4. [ ] 在 handlers_container.go 中注册回调
5. [ ] 测试回调触发

---

### 2.5 更新 SlotData 读取逻辑

**文件**: `internal/item/slot.go`

**目标**: 集成组件解析到 SlotData 读取流程

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
    Components []*component.ComponentResult // 新增字段
}

// ReadSlotData 从 Reader 读取物品槽数据
func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
    // 1. 读取数量
    count, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    if count <= 0 {
        return nil, nil // 空槽
    }
    
    // 2. 读取物品ID
    itemID, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    // 3. 读取添加的组件数量
    numAdd, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    // 4. 读取移除的组件数量
    numRemove, err := binutil.ReadVarInt(r)
    if err != nil {
        return nil, err
    }
    
    // 5. 解析组件
    parser := component.NewParser()
    components := make([]*component.ComponentResult, 0, numAdd)
    
    for i := int32(0); i < numAdd; i++ {
        typeID, err := binutil.ReadVarInt(r)
        if err != nil {
            return nil, err
        }
        
        result, err := parser.ParseComponent(typeID, r)
        if err != nil {
            return nil, fmt.Errorf("component %d: %w", typeID, err)
        }
        
        components = append(components, result)
    }
    
    // 6. 跳过移除的组件类型ID
    for i := int32(0); i < numRemove; i++ {
        if _, err := binutil.ReadVarInt(r); err != nil {
            return nil, err
        }
    }
    
    return &SlotData{
        ID:         itemID,
        Count:      count,
        Components: components,
    }, nil
}

// SkipSlotData 跳过物品槽数据
func SkipSlotData(r *bytes.Reader) error {
    _, err := ReadSlotData(r)
    return err
}
```

**实现步骤**:
1. [ ] 更新 SlotData 结构体，添加 Components 字段
2. [ ] 实现完整的 ReadSlotData 函数
3. [ ] 更新 SkipSlotData 使用新的 ReadSlotData
4. [ ] 更新单元测试

---

### 2.6 更新 packet/readers.go

**文件**: `internal/mcclient/packet/readers.go`

**目标**: 使用新的 item 包实现

```go
package packet

import (
    "gmcc/internal/item"
)

// ReadSlotData 使用 internal/item 的实现
func ReadSlotData(r *bytes.Reader) (*item.SlotData, error) {
    return item.ReadSlotData(r)
}

// 保持其他函数不变或迁移到 pkg/binutil
```

**实现步骤**:
1. [ ] 更新导入路径
2. [ ] 修改 ReadSlotData 调用新的实现
3. [ ] 移除重复的 SlotData 定义
4. [ ] 验证编译通过

---

## 3. 实施步骤

### Day 1

**上午（4小时）**:
1. [ ] 创建 `internal/item/component/parser.go`
2. [ ] 实现 Parser 结构体和基本方法
3. [ ] 编写 Parser 单元测试
4. [ ] 创建 `internal/item/component/handlers.go`
5. [ ] 实现 defaultHandlers() 函数

**下午（4小时）**:
1. [ ] 创建 `internal/item/component/discard.go`
2. [ ] 迁移基础跳过函数到 pkg/binutil/skip.go
3. [ ] 实现组件特定的丢弃函数
4. [ ] 实现 makeDiscardHandler 工厂函数
5. [ ] 编写丢弃处理器测试

### Day 2

**上午（4小时）**:
1. [ ] 实现容器组件特殊处理
2. [ ] 实现 SetContainerCallback 函数
3. [ ] 在 handlers_container.go 中注册回调
4. [ ] 测试容器回调触发

**下午（4小时）**:
1. [ ] 更新 internal/item/slot.go
2. [ ] 添加 Components 字段到 SlotData
3. [ ] 更新 ReadSlotData 函数
4. [ ] 编写 SlotData 单元测试

### Day 3

**上午（4小时）**:
1. [ ] 更新 packet/readers.go
2. [ ] 修改导入路径和函数调用
3. [ ] 修复编译错误
4. [ ] 运行单元测试

**下午（4小时）**:
1. [ ] 集成测试
2. [ ] 验证容器回调工作正常
3. [ ] 性能测试对比
4. [ ] 提交代码并打标签

---

## 4. 验证清单

### 编译验证

```bash
# 1. 编译所有包
go build ./...

# 2. 运行组件相关测试
go test ./internal/item/component/... -v
go test ./internal/item/... -v

# 3. 运行 packet 测试
go test ./internal/mcclient/packet/... -v
```

### 功能验证

```go
// 测试组件解析
func TestComponentParsing(t *testing.T) {
    // 创建测试数据
    data := createTestSlotData()
    r := bytes.NewReader(data)
    
    // 解析 SlotData
    slot, err := item.ReadSlotData(r)
    if err != nil {
        t.Fatalf("ReadSlotData failed: %v", err)
    }
    
    // 验证组件列表
    if len(slot.Components) == 0 {
        t.Error("Expected components, got none")
    }
    
    // 验证容器回调
    // ...
}
```

---

## 5. 风险和对策

| 风险 | 概率 | 影响 | 对策 |
|------|------|------|------|
| 组件解析错误 | 中 | 高 | 保持原有 Skip 逻辑，仅添加结果收集 |
| 性能下降 | 中 | 中 | 预分配切片容量，减少内存分配 |
| 回调未触发 | 低 | 高 | 添加详细的调试日志 |
| 向后不兼容 | 低 | 高 | SlotData 添加字段不影响现有代码 |

---

## 6. 提交信息模板

```bash
# Day 1 提交
git add internal/item/component/
git commit -m "feat: implement component parser and handlers

- Add Parser struct with Parse method
- Implement defaultHandlers() for 104 component types
- Create discard handlers for all component types
- Add unit tests for parser and handlers

Refs: phase2"

# Day 2 提交
git commit -m "feat: add container callback mechanism

- Implement ContainerComponentHandler with callback
- Add SetContainerCallback registration function
- Register callback in handlers_container.go
- Update SlotData with Components field

Refs: phase2"

# Day 3 提交
git commit -m "feat: integrate component parsing into SlotData

- Update ReadSlotData to use component parser
- Migrate packet/readers.go to use new item package
- Add comprehensive integration tests
- Update documentation

Refs: phase2"

# 阶段完成标签
git tag -a v0.3.0-phase2 -m "Phase 2 complete: Component parsing implementation"
```

---

## 7. 下一步

完成本阶段后，进入 [阶段3 - 优化和重构](./2026-03-24-component-parsing-phase3.md)
