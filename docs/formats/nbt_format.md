# NBT 格式规范

## 概述

**NBT（Named Binary Tag，二进制命名标签）** 是 Minecraft 用于数据存储和传输的树状数据结构。

- 官方文档：[Minecraft Wiki - NBT格式](https://zh.minecraft.wiki/w/NBT%E6%A0%BC%E5%BC%8F)
- 协议版本：774 (1.21.11)

## 标签类型

| ID | 名称 | 负载大小 | 说明 |
|----|------|----------|------|
| 0 | TAG_End | 0 | 结束标签，无内容 |
| 1 | TAG_Byte | 1 byte | 有符号字节 (-128~127) |
| 2 | TAG_Short | 2 bytes | 有符号短整型 |
| 3 | TAG_Int | 4 bytes | 有符号整型 |
| 4 | TAG_Long | 8 bytes | 有符号长整型 |
| 5 | TAG_Float | 4 bytes | 单精度浮点数 (IEEE 754) |
| 6 | TAG_Double | 8 bytes | 双精度浮点数 (IEEE 754) |
| 7 | TAG_Byte_Array | 4+n | 字节数组 |
| 8 | TAG_String | 2+n | UTF-8 字符串 |
| 9 | TAG_List | 5+n | 同类型元素列表 |
| 10 | TAG_Compound | n | 子标签复合结构 |
| 11 | TAG_Int_Array | 4+n×4 | 整型数组 |
| 12 | TAG_Long_Array | 4+n×8 | 长整型数组 |

## 字节序

Java 版使用**大端序** (Big Endian)。

## 编码细节

### 字符串格式

```
[长度: unsigned short (2 bytes)] [内容: UTF-8 bytes]
```

### 列表格式

```
[元素类型: byte (1 byte)] [长度: int (4 bytes)] [元素负载...]
```

### 复合标签格式

```
[子标签1] [子标签2] ... [TAG_End]
```

### 带名称标签结构

```
[标签ID: byte] [名称长度: unsigned short] [名称: UTF-8] [负载]
```

## 传输格式 vs 存储格式

| 特性 | 存储格式 | 传输格式 (Java) |
|------|----------|-----------------|
| 根标签名称 | 有 | 无（省略） |
| 压缩 | GZip/Zlib 可选 | 无压缩 |
| 用途 | 存档文件、方块实体 | 网络传输 |

## CESU-8 编码

Minecraft 字符串使用 **CESU-8** 而非标准 UTF-8，用于编码代理对（Surrogate Pairs）字符。

转换规则：
- UCS-4 码点 U+10000 ~ U+10FFFF 编码为 6 字节 CESU-8
- 代理对的高位和低位分别编码

## 参考

- [NBT格式 - Minecraft Wiki](https://zh.minecraft.wiki/w/NBT%E6%A0%BC%E5%BC%8F)
- [Minecraft Protocol - NBT](https://minecraft.wiki/w/NBT)
- [DataOutput JavaDoc](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/io/DataOutput.html)
