# 物品组件内容解析实施计划

**日期**: 2026-03-25  
**对应设计**: [2026-03-24-component-parsing-design.md](../specs/2026-03-24-component-parsing-design.md)  
**状态**: 进行中  
**当前阶段**: 阶段1 - 常用显示组件 ✅ 已完成  
**预计工期**: 7-10天（全阶段）

---

## 1. 计划概述

本计划在 v1.0（丢弃处理器）的基础上，实现组件内容的完整解析。采用分阶段策略，从简单到复杂逐步实现104种组件类型。

### 阶段划分

| 阶段 | 内容 | 预计工期 | 状态 |
|------|------|---------|------|
| 阶段0 | P0 基础类型（11个组件） | 1-2天 | ✅ 已完成 |
| 阶段1 | P1 常用显示组件 | 2-3天 | ✅ 已完成 |
| 阶段2 | P2 复杂结构组件 | 3-4天 | ⏳ 待开始 |
| 阶段3 | 全面集成优化 | 1-2天 | ⏳ 待开始 |

---

## 2. 阶段0: 基础类型实现

### 2.1 目标

实现11个基础数据类型的组件解析器，替换当前的丢弃处理器。

### 2.2 组件列表

| ID | 名称 | 数据类型 | 文件位置 |
|----|------|---------|----------|
| 1 | max_stack_size | VarInt | parsers_varint.go |
| 2 | max_damage | VarInt | parsers_varint.go |
| 3 | damage | VarInt | parsers_varint.go |
| 4 | unbreakable | 无数据 | parsers_bool.go |
| 17 | custom_model_data | Int32 | parsers_int32.go |
| 19 | repair_cost | VarInt | parsers_varint.go |
| 21 | enchantment_glint_override | Bool | parsers_bool.go |
| 31 | enchantable | VarInt | parsers_varint.go |
| 42 | dyed_color | Int32 | parsers_int32.go |
| 43 | map_color | Int32 | parsers_int32.go |
| 44 | map_id | VarInt | parsers_varint.go |

### 2.3 文件结构

```
internal/item/component/
├── parsers_varint.go       # VarInt 解析器
├── parsers_int32.go        # Int32 解析器
├── parsers_bool.go         # Bool 解析器
└── parsers_varint_test.go  # 测试
```

### 2.4 任务清单

- [x] 创建 `parsers_varint.go`，实现6个VarInt组件解析器
- [x] 创建 `parsers_int32.go`，实现3个Int32组件解析器
- [x] 创建 `parsers_bool.go`，实现2个Bool组件解析器
- [x] 更新 `handlers.go`，注册新的解析器
- [x] 编写单元测试
- [x] 运行测试验证

### 2.5 API示例

```go
internal/item/component/parsers_varint.go:
func ParseMaxStackSize(r *bytes.Reader) (*ComponentResult, error) {
    value, err := packet.ReadVarIntFromReader(r)
    if err != nil {
        return nil, err
    }
    return &ComponentResult{
        TypeID: MaxStackSize,
        Data:   value,
    }, nil
}
```

---

## 3. 阶段1: 常用显示组件

### 3.1 目标

实现最常用的显示组件，支持物品名称、描述、附魔等信息的解析。

### 3.2 组件列表

| ID | 名称 | 数据类型 | 依赖 |
|----|------|---------|------|
| 6 | custom_name | Text | 聊天系统 |
| 9 | item_name | Text | 聊天系统 |
| 11 | lore | List\<Text\> | 聊天系统 |
| 12 | rarity | Enum | 新实现 |
| 13 | enchantments | Map | 新实现 |

### 3.3 任务清单

- [x] 实现文本组件解析（custom_name, item_name）
- [x] 实现列表文本解析
- [x] 实现 rarity 枚举解析
- [x] 实现附魔列表解析
- [x] 集成到 handlers
- [x] 测试验证

---

## 4. 阶段2: 复杂结构组件

### 4.1 目标

实现需要复杂结构解析的组件。

### 4.2 组件列表

| ID | 名称 | 数据结构 |
|----|------|----------|
| 5 | use_effects | Record |
| 50 | potion_contents | Complex |
| 54 | trim | Complex |

### 4.3 任务清单

- [ ] 实现 use_effects record 结构
- [ ] 实现 potion_contents 解析
- [ ] 实现 trim 纹饰结构
- [ ] 优化容器组件嵌套解析

---

## 5. 阶段3: 全面集成优化

- [ ] 性能优化和基准测试
- [ ] 文档更新
- [ ] 完整测试覆盖率
- [ ] 代码审查

---

## 6. 验证标准

### 功能验证

- [ ] P0 基础类型组件解析正确
- [ ] P1 显示组件在UI中显示正确
- [ ] 测试覆盖率 ≥ 80%
- [ ] 无性能退化（解析时间 ≤ 原始110%）

### 测试策略

```go
// 示例测试
func TestParseMaxStackSize(t *testing.T) {
    data := []byte{0x40} // 64
    r := bytes.NewReader(data)
    result, err := ParseMaxStackSize(r)
    if err != nil {
        t.Fatal(err)
    }
    if result.Data.(int32) != 64 {
        t.Errorf("got %v, want 64", result.Data)
    }
}
```

---

## 7. 实施优先级

| 优先级 | 组件类型 | 说明 |
|--------|---------|------|
| P0 | 基础类型 | 当前实施阶段 |
| P1 | 常用显示 | 下一个阶段 |
| P2 | 复杂结构 | 后续扩展 |
| P3 | 特殊组件 | 保持现有 |

---

**下一步**: 开始阶段0实现