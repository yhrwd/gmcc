# 阶段1: 基础架构搭建

**计划文档**: [主计划](../2026-03-24-component-parsing-implementation-plan.md)  
**依赖**: 无（首个阶段）  
**阻塞**: 阶段2  
**状态**: 待开始

---

## 1. 阶段目标

建立新的包结构，迁移工具类，确保代码编译通过，所有导入路径正确更新。

### 成功标准

- [ ] 新的包结构创建完成
- [ ] 所有代码编译通过（`go build ./...`）
- [ ] 基础单元测试通过（`go test ./pkg/binutil/...`）
- [ ] 无循环依赖

---

## 2. 任务清单

### 2.1 创建 pkg/binutil/ 包

**目标**: 将通用的二进制读取工具从 `packet/` 移出

#### 文件: pkg/binutil/types.go

```go
package binutil

// 基础类型定义
// VarInt, VarLong, String, Bool 等的类型别名
```

**内容**: 
- 定义 VarInt、VarLong 类型
- 定义 PacketField 接口（如有需要）

#### 文件: pkg/binutil/reader.go

```go
package binutil

// Reader 提供二进制数据读取功能
type Reader struct {
    data []byte
    pos  int
}

func NewReader(data []byte) *Reader
func (r *Reader) ReadVarInt() (int32, error)
func (r *Reader) ReadVarLong() (int64, error)
func (r *Reader) ReadString() (string, error)
func (r *Reader) ReadBool() (bool, error)
func (r *Reader) ReadByte() (byte, error)
func (r *Reader) ReadInt16() (int16, error)
func (r *Reader) ReadInt32() (int32, error)
func (r *Reader) ReadInt64() (int64, error)
func (r *Reader) ReadFloat32() (float32, error)
func (r *Reader) ReadFloat64() (float64, error)
// ... 其他基础类型
```

**实现步骤**:
1. 从 `internal/mcclient/packet/readers.go` 复制读取函数
2. 适配为 `Reader` 类型的方法
3. 保持原有错误处理逻辑
4. 添加单元测试

#### 文件: pkg/binutil/writer.go

```go
package binutil

// Writer 提供二进制数据写入功能
type Writer struct {
    buf []byte
}

func NewWriter() *Writer
func (w *Writer) WriteVarInt(v int32) error
func (w *Writer) WriteString(s string) error
// ... 对应 reader 的写入方法
func (w *Writer) Bytes() []byte
```

**注意**: 如果项目中暂不需要写入功能，可以先实现 reader 部分。

#### 文件: pkg/binutil/utils.go

```go
package binutil

// VarIntSize 返回编码 VarInt 需要的字节数
func VarIntSize(v int32) int

// 其他工具函数
```

#### 文件: pkg/binutil/reader_test.go

```go
package binutil

import "testing"

func TestReader_ReadVarInt(t *testing.T)
func TestReader_ReadString(t *testing.T)
// ... 完整的单元测试
```

**测试要求**:
- 边界值测试（最大值、最小值）
- 错误处理测试（截断数据）
- 与原有实现对比测试

---

### 2.2 创建 internal/item/ 包结构

#### 目录结构

```
internal/item/
├── slot.go                 # SlotData 定义
├── slot_test.go           # SlotData 测试
└── component/             # 组件子包
    ├── types.go          # ComponentResult 等类型
    ├── parser.go         # Parser 实现
    ├── parser_test.go    # Parser 测试
    ├── handlers.go       # 处理器注册
    ├── discard.go        # 丢弃处理器
    ├── constants.go      # 组件ID常量
    └── component_test.go # 集成测试
```

#### 文件: internal/item/slot.go

```go
package item

import "gmcc/internal/item/component"

// SlotData 表示物品槽中的物品数据
type SlotData struct {
    ID         int32                        // 物品ID
    Count      int32                        // 数量
    Components []*component.ComponentResult // 组件列表
}

// ReadSlotData 从字节流读取物品槽数据
func ReadSlotData(data []byte) (*SlotData, error)

// IsEmpty 检查槽位是否为空
func (s *SlotData) IsEmpty() bool
```

**注意**: 
- 从原有的 `packet.SlotData` 迁移
- 添加 `Components` 字段用于存储解析的组件
- 保持向后兼容（添加字段不影响现有使用）

#### 文件: internal/item/component/types.go

```go
package component

// ComponentResult 组件解析结果
type ComponentResult struct {
    TypeID int32  // 组件类型ID
    Data   any    // 解析后的数据（当前阶段可能为nil）
}

// ComponentHandler 组件处理器函数类型
type ComponentHandler func(data []byte) (*ComponentResult, int, error)

// Parser 组件解析器
type Parser struct {
    handlers map[int32]ComponentHandler
}

// NewParser 创建默认解析器
func NewParser() *Parser

// RegisterHandler 注册组件处理器
func (p *Parser) RegisterHandler(typeID int32, handler ComponentHandler)

// Parse 解析组件数据
func (p *Parser) Parse(data []byte) ([]*ComponentResult, error)
```

**设计说明**:
- 保持与原有 `component_skipping.go` 的兼容性
- 使用函数映射表而非接口，更轻量
- 支持运行时注册处理器

---

### 2.3 迁移和更新导入路径

#### 修改的文件清单

| 原文件 | 修改内容 | 新导入 |
|--------|----------|--------|
| internal/mcclient/packet/readers.go | 移除SlotData，使用新包 | `gmcc/internal/item` |
| internal/mcclient/packet/codec.go | 更新工具函数导入 | `gmcc/pkg/binutil` |
| internal/mcclient/handlers_container.go | 使用新的SlotData | `gmcc/internal/item` |
| internal/player/player.go | 更新Item定义 | `gmcc/internal/item` |

#### 导入路径更新示例

**修改前**:
```go
import "gmcc/internal/mcclient/packet"

func handlePacket(data []byte) {
    slot, err := packet.ReadSlotData(data)
    // ...
}
```

**修改后**:
```go
import "gmcc/internal/item"

func handlePacket(data []byte) {
    slot, err := item.ReadSlotData(data)
    // ...
}
```

---

### 2.4 保持向后兼容

#### 兼容层（如需要）

在 `packet/` 包中提供兼容别名（可选）:

```go
// internal/mcclient/packet/compatibility.go
package packet

import "gmcc/internal/item"

// SlotData 别名（向后兼容）
type SlotData = item.SlotData

// ReadSlotData 别名（向后兼容）
func ReadSlotData(data []byte) (*item.SlotData, error) {
    return item.ReadSlotData(data)
}
```

**注意**: 如果项目代码量不大，建议直接更新所有引用，不保留兼容层。

---

## 3. 实施步骤

### Day 1

**上午（4小时）**:
1. [ ] 创建 `pkg/binutil/` 目录结构
2. [ ] 实现 `types.go`
3. [ ] 实现 `reader.go`（基础读取函数）
4. [ ] 运行测试确保编译通过

**下午（4小时）**:
1. [ ] 创建 `internal/item/` 目录结构
2. [ ] 实现 `slot.go`
3. [ ] 实现 `component/types.go`
4. [ ] 编写基础单元测试

### Day 2

**上午（4小时）**:
1. [ ] 更新 `packet/readers.go` 导入
2. [ ] 更新 `packet/codec.go` 导入
3. [ ] 更新 `handlers_container.go` 导入
4. [ ] 运行完整编译检查

**下午（4小时）**:
1. [ ] 修复编译错误
2. [ ] 运行单元测试
3. [ ] 代码审查
4. [ ] 提交代码并打标签: `git tag -a phase1-complete -m "Phase 1: Basic architecture complete"`

---

## 4. 验证清单

### 编译验证

```bash
# 1. 编译所有包
go build ./...

# 2. 运行单元测试
go test ./pkg/binutil/... -v
go test ./internal/item/... -v

# 3. 检查循环依赖
go mod why -m gmcc

# 4. 检查导入路径
go list -f '{{.ImportPath}}: {{.Imports}}' ./pkg/... ./internal/...
```

### 功能验证

```go
// 测试代码示例
func TestBinutilReader(t *testing.T) {
    data := []byte{0x01, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}
    r := binutil.NewReader(data)
    
    v, err := r.ReadVarInt()
    if err != nil || v != 1 {
        t.Errorf("ReadVarInt failed: %v, %v", v, err)
    }
    
    s, err := r.ReadString()
    if err != nil || s != "hello" {
        t.Errorf("ReadString failed: %v, %v", s, err)
    }
}
```

---

## 5. 风险和对策

| 风险 | 概率 | 影响 | 对策 |
|------|------|------|------|
| 导入路径遗漏 | 高 | 中 | 使用IDE全局搜索和替换 |
| 循环依赖 | 中 | 高 | 使用 `go mod why` 检查 |
| 测试失败 | 中 | 中 | 保留旧实现对比验证 |
| 编译性能下降 | 低 | 低 | 监控编译时间 |

---

## 6. 提交信息模板

```bash
# 阶段1提交
git add pkg/binutil/ internal/item/
git commit -m "feat: add binutil package and item package structure

- Create pkg/binutil/ with binary read utilities
- Create internal/item/ for item system
- Move SlotData to internal/item/slot.go
- Update all import paths
- Add unit tests for binutil

Refs: phase1"

# 打标签
git tag -a v0.2.0-phase1 -m "Phase 1 complete: Basic architecture"
```

---

## 7. 下一步

完成本阶段后，进入 [阶段2: 组件解析实现](2026-03-24-component-parsing-phase2.md)
