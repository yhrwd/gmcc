---
name: minecraft-dev-assistant
description: Minecraft 协议规范、数据格式、PrismarineJS 库和知识库的开发参考。
tags: [minecraft, protocol, nbt, chat, knowledge]
---

# Minecraft Dev Assistant Skill

提供 Minecraft 协议规范、数据格式、PrismarineJS 相关库和项目知识库的开发参考。

## 核心功能

- **协议查询**: Java/Bedrock 版数据包结构与状态机 (详见 `.knowledge/MC_Protocol_Data`)
- **数据格式**: NBT 结构定义与聊天组件 JSON 规范
- **库参考**: `prismarine-chat` 和 `prismarine-nbt` 的 API 说明
- **游戏数据**: 方块、物品、实体等 ID 与属性参考 (详见 `.knowledge/minecraft-data`)
- **知识库检索**: 查询项目根目录下 `.knowledge/` 目录下的协议数据、版本映射和 API 文档
- **文档参考**: 项目 `docs/` 目录下有 Minecraft 格式规范文档 (NBT/SNBT/文本组件/数据组件)

## 知识库结构

项目根目录: `C:\Users\Yhrza\Desktop\python\gmcc\`

```
.knowledge/                   # 项目知识库 (位于项目根目录)
├── MC_Protocol_Data/          # Java 版协议定义
│   └── java_edition/
│       ├── packets/          # 数据包结构
│       ├── codec/            # 编解码器
│       ├── indexed_data/     # 版本索引数据
│       └── structures/       # 协议结构
│
├── minecraft-data/            # 游戏数据规范
│   └── data/                 # 各版本数据 JSON
│
├── prismarine-nbt/           # NBT 库参考
├── prismarine-chat/          # 聊天组件库参考
├── mineflayer/               # Bot 框架参考
├── protocol_774/             # 协议 774 特定数据
├── README.md                 # 知识库索引
└── links.md                  # 外部链接

docs/                         # 项目文档
├── nbt_format.md             # NBT 格式规范
├── snbt_format.md            # SNBT 格式规范
├── text_component.md        # 文本组件规范
├── data_components_1.21.11.md  # 1.21.11 数据组件规范
└── ...
```

## 使用指南

### 0. 检索建议
- **文档优先**: 搜索时优先查看项目根目录下的 `.md` 文件和 `docs/` 目录
- **深入分析**: 如果需要具体版本的数据包 ID 或字段，请进入 `.knowledge/` 目录查看对应的 JSON 文件
- **本地文档**: `docs/` 目录下有中文的 NBT/SNBT/文本组件/数据组件规范文档

### 1. 协议开发
查询不同版本的协议差异和握手流程。
- 参考文件: `docs/protocol.md`
- 版本映射: `.knowledge/MC_Protocol_Data/java_edition/versions.json`
- 数据包索引: `.knowledge/MC_Protocol_Data/java_edition/packets.csv`

### 2. 游戏数据查询
| 数据类型 | 路径 |
|---------|------|
| 方块 | `.knowledge/minecraft-data/data/{version}/blocks.json` |
| 物品 | `.knowledge/minecraft-data/data/{version}/items.json` |
| 实体 | `.knowledge/minecraft-data/data/{version}/entities.json` |
| 附魔 | `.knowledge/minecraft-data/data/{version}/enchantments.json` |
| 配方 | `.knowledge/minecraft-data/data/{version}/recipes.json` |
| 协议描述 | `.knowledge/minecraft-data/data/{version}/protocol.json` |

### 3. NBT 处理
了解如何解析和构建 NBT 数据。
- 参考文件: `prismarine-nbt.md`
- 常用方法: `nbt.parse()`, `nbt.simplify()`, `nbt.comp()`

### 4. 聊天消息
处理 Minecraft 复杂的文本组件系统。
- 参考文件: `prismarine-chat.md`
- 常用方法: `new ChatMessage()`, `MessageBuilder`

### 5. 快速查询

```bash
# 查找协议版本
grep -r "protocolVersion.*774" .knowledge/

# 查找数据包定义
find .knowledge/MC_Protocol_Data -name "play_client_login.json"

# 查找物品数据
find .knowledge/minecraft-data -name "items.json" | head -5
```

### 6. 本地文档 (docs/)
项目 `docs/` 目录下有中文 Minecraft 格式规范文档：

| 文档 | 说明 |
|------|------|
| `docs/nbt_format.md` | NBT 二进制格式规范 |
| `docs/snbt_format.md` | SNBT 文本格式规范 |
| `docs/text_component.md` | 文本组件 (JSON 格式) |
| `docs/data_components_1.21.11.md` | 1.21.11 数据组件规范 |

## 版本映射

| Minecraft 版本 | 协议版本 |
|---------------|---------|
| 1.21.11       | 774     |
| 1.21.5        | 773     |
| 1.21          | 772     |
| 1.20.4        | 765     |
| 1.20.2        | 764     |
| 1.20.1        | 763     |
| 1.19.4        | 762     |
| 1.19.3        | 761     |
| 1.19.1/2      | 760     |

完整映射见 `.knowledge/MC_Protocol_Data/java_edition/versions.json`

## 数据包状态机

```
┌─────────────┐     Handshake     ┌─────────────┐
│   Handshake │ ──────────────────>│    Status   │
└─────────────┘                   └─────────────┘
       │
       │ intent=2 (login)
       ▼
┌─────────────┐     Success      ┌─────────────┐
│    Login    │ ─────────────────>│    Play    │
└─────────────┘                   └─────────────┘
```

## NBT 标签类型

| ID | Name | Size |
|----|------|------|
| 0 | TAG_End | 0 |
| 1 | TAG_Byte | 1 |
| 2 | TAG_Short | 2 |
| 3 | TAG_Int | 4 |
| 4 | TAG_Long | 8 |
| 5 | TAG_Float | 4 |
| 6 | TAG_Double | 8 |
| 7 | TAG_Byte_Array | var |
| 8 | TAG_String | var |
| 9 | TAG_List | var |
| 10 | TAG_Compound | var |
| 11 | TAG_Int_Array | var |
| 12 | TAG_Long_Array | var |

## 参考资源
- [Minecraft Wiki 协议页面](https://zh.minecraft.wiki/w/Java版网络协议)
- [PrismarineJS GitHub](https://github.com/PrismarineJS)
- [NBT 格式规范](https://zh.minecraft.wiki/w/NBT格式)
- [文本组件规范](https://zh.minecraft.wiki/w/文本组件)

## 开发规则
- **协议优先**: 始终验证目标 Minecraft 版本的协议 ID
- **版本优先**: 查询数据前确认目标协议版本
- **类型安全**: 在构建 NBT 时使用明确的类型转换（如 `nbt.int()`, `nbt.short()`）
- **性能优化**: 处理大体积 NBT 或 Chunk 数据时，避免过度使用 `simplify()`
- **数据源追踪**: 优先使用 `MC_Protocol_Data` 中的精确数据

---
*Generated from project knowledge base*