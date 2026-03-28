# SNBT 格式规范

## 概述

**SNBT（Stringified NBT，字符串化 NBT）** 是 NBT 的文本表示形式，用于命令、配置文件和数据包中。

- 官方文档：[Minecraft Wiki - SNBT格式](https://zh.minecraft.wiki/w/SNBT%E6%A0%BC%E5%BC%8F)
- 协议版本：774 (1.21.11)

## 基本语法

### 数值类型后缀

| 类型 | 后缀 | 示例 | 说明 |
|------|------|------|------|
| Byte | `b` / `B` | `127b` | 有符号字节 |
| Short | `s` / `S` | `32767s` | 有符号短整型 |
| Int | 无 | `42` | 整型（默认） |
| Long | `l` / `L` | `100L` | 有符号长整型 |
| Float | `f` / `F` | `3.14f` | 单精度浮点 |
| Double | `d` / `D` | `2.718d` | 双精度浮点（默认） |

### 无符号数值

| 类型 | 后缀 | 示例 |
|------|------|------|
| Unsigned Byte | `ub` / `UB` | `255ub` |
| Unsigned Short | `us` / `US` | `65535us` |
| Unsigned Int | `ui` / `UI` | `4294967295ui` |
| Unsigned Long | `ul` / `UL` | `18446744073709551615ul` |

## 数据结构

### 复合标签（Compound）

```
{key1: value1, key2: value2}
```

示例：
```snbt
{display: {Name: '{"text":"钻石剑"}'}, Damage: 0s}
```

### 列表（List）

```
[value1, value2, value3]
```

示例：
```snbt
[1b, 2b, 3b]
["stone", "dirt", "grass"]
```

### 数组（Arrays）

| 类型 | 前缀 | 示例 |
|------|------|------|
| Byte Array | `[B;` | `[B; 1b, 2b, 3b]` |
| Int Array | `[I;` | `[I; 1, 2, 3]` |
| Long Array | `[L;` | `[L; 1l, 2l, 3l]` |

## 字符串规则

### 引号使用

- 标准字符串：双引号 `"text"`
- 可选引号：无空格/特殊字符时可省略
- 转义序列：
  - `\\` - 反斜杠
  - `\"` - 双引号
  - `\n` - 换行
  - `\t` - 制表符

### 示例

```snbt
# 简单字符串
hello
"hello world"

# 带转义的字符串
"Line 1\\nLine 2"
"Say \\"Hello\\""

# 复杂复合标签
{
    id: "minecraft:diamond_sword",
    Count: 1b,
    tag: {
        Damage: 0s,
        Enchantments: [
            {id: "minecraft:sharpness", lvl: 5s}
        ]
    }
}
```

## 布尔值

| 值 | NBT 表示 |
|----|---------|
| `true` | `1b` |
| `false` | `0b` |

## 路径查询语法

用于命令和代码中访问 NBT 数据：

```
compound.child.property
list[index]
list[{matcher: value}].property
```

示例：
```
Inventory[0].id
Inventory[{Slot:0b}].id
Items[].Count
```

## 与 JSON 的区别

| 特性 | SNBT | JSON |
|------|------|------|
| 键引号 | 可选 | 必需（双引号） |
| 尾部逗号 | 允许 | 不允许 |
| 注释 | 不支持 | 不支持 |
| 单引号 | 支持 | 不支持 |
| 数值类型标记 | 支持 | 不支持 |

## 参考

- [SNBT格式 - Minecraft Wiki](https://zh.minecraft.wiki/w/SNBT%E6%A0%BC%E5%BC%8F)
- [NBT路径 - Minecraft Wiki](https://zh.minecraft.wiki/w/NBT%E8%B7%AF%E5%BE%84)
- [教程:SNBT](https://zh.minecraft.wiki/w/Tutorial:SNBT)
